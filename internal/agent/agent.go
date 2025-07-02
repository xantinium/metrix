// Package agent содержит реализацию агента для сбора метрик.
package agent

import (
	"context"
	_ "net/http/pprof" // Используется для корректной работы профилировщика.
	"time"

	"github.com/xantinium/metrix/internal/agent/rest"
	"github.com/xantinium/metrix/internal/agent/rpc"
	"github.com/xantinium/metrix/internal/config"
	"github.com/xantinium/metrix/internal/infrastructure/runtimemetrics"
	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/models"
	"github.com/xantinium/metrix/internal/tools"
)

const agentWorkerPoolSize = 3

// MetrixAgentOptions параметры агента метрик.
type MetrixAgentOptions struct {
	ServerAddr         string
	PrivateKey         string
	CryptoPublicKey    string
	RequestMethod      config.RequestMethod
	PollInterval       int
	ReportInterval     time.Duration
	ReportRateLimit    int
	IsProfilingEnabled bool
}

// NewMetrixAgent создаёт новый агент метрик.
func NewMetrixAgent(opts MetrixAgentOptions) *MetrixAgent {
	agent := &MetrixAgent{
		isProfilingEnabled: opts.IsProfilingEnabled,
		metricsSource:      runtimemetrics.NewRuntimeMetricsSource(opts.PollInterval),
	}

	agent.workerPool = NewMetrixAgentWorkerPool(MetrixAgentWorkerPoolOptions{
		PoolSize:        agentWorkerPoolSize,
		ReportInterval:  opts.ReportInterval,
		ReportRateLimit: opts.ReportRateLimit,
		UploadFunc:      agent.updateMetrics,
	})

	// Получаем IP-адрес, который будет использоваться
	// запросах. Не уверен, что это правильно, но пока что так.
	addr, err := tools.FindLocalIP()
	if err != nil {
		logger.Errorf("failet to find local ip: %v", err)
	}

	switch opts.RequestMethod {
	case config.RequestMethodRest:
		agent.client = rest.NewClient(rest.ClientOptions{
			ServerAddr:      opts.ServerAddr,
			PrivateKey:      opts.PrivateKey,
			CryptoPublicKey: opts.CryptoPublicKey,
			Addr:            addr,
		})
	case config.RequestMethodRPC:
		agent.client = rpc.NewClient(rpc.ClientOptions{
			ServerAddr:      opts.ServerAddr,
			PrivateKey:      opts.PrivateKey,
			CryptoPublicKey: opts.CryptoPublicKey,
			Addr:            addr,
		})
	default:
		logger.Errorf("unknown request method %q", opts.RequestMethod)
	}

	return agent
}

type Client interface {
	UpdateMetric(ctx context.Context, metric models.MetricInfo) error
	UpdateMetricsBatch(ctx context.Context, metrics []models.MetricInfo) error
}

// MetrixAgent структура, описывающая агент метрик.
type MetrixAgent struct {
	workerPool    *MetrixAgentWorkerPool
	metricsSource *runtimemetrics.RuntimeMetricsSource
	client        Client

	isProfilingEnabled bool
}

// Run запускает агента метрик.
func (agent *MetrixAgent) Run(ctx context.Context) {
	agent.metricsSource.Run(ctx)
	agent.workerPool.Run(ctx)

	if agent.isProfilingEnabled {
		tools.RunProfilingServer()
	}
}

// Disable отключает агент, запрещая выполнять
// новые запросы.
func (agent *MetrixAgent) Disable() {
	agent.metricsSource.Disable()
	agent.workerPool.Disable()
}

func (agent *MetrixAgent) updateMetrics(ctx context.Context) {
	err := agent.client.UpdateMetricsBatch(ctx, agent.metricsSource.GetSnapshot())
	if err != nil {
		logger.Error(err.Error())
	}
}
