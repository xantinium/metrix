// Package agent содержит реализацию агента для сбора метрик.
package agent

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" // Используется для корректной работы профилировщика.
	"time"

	"github.com/mailru/easyjson"

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
	PollInterval       int
	ReportInterval     time.Duration
	ReportRateLimit    int
	IsProfilingEnabled bool
}

// NewMetrixAgent создаёт новый агент метрик.
func NewMetrixAgent(opts MetrixAgentOptions) *MetrixAgent {
	agent := &MetrixAgent{
		serverAddr:         opts.ServerAddr,
		privateKey:         opts.PrivateKey,
		cryptoPublicKey:    opts.CryptoPublicKey,
		isProfilingEnabled: opts.IsProfilingEnabled,
		metricsSource:      runtimemetrics.NewRuntimeMetricsSource(opts.PollInterval),
		retrier:            tools.DefaulRetrier,
	}

	agent.workerPool = NewMetrixAgentWorkerPool(MetrixAgentWorkerPoolOptions{
		PoolSize:        agentWorkerPoolSize,
		ReportInterval:  opts.ReportInterval,
		ReportRateLimit: opts.ReportRateLimit,
		UploadFunc:      agent.UpdateMetrics,
	})

	return agent
}

// MetrixAgent структура, описывающая агент метрик.
type MetrixAgent struct {
	workerPool         *MetrixAgentWorkerPool
	metricsSource      *runtimemetrics.RuntimeMetricsSource
	retrier            *tools.Retrier
	serverAddr         string
	privateKey         string
	cryptoPublicKey    string
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

// UpdateMetrics обновляет метрики на сервере.
func (agent *MetrixAgent) UpdateMetrics() {
	agent.updateMetricsBatch(agent.metricsSource.GetSnapshot())
}

// updateMetric обновление метрики через хендлеры первой версии.
//
// Deprecated: метод устарел, следует использовать updateMetricsV2.
func (agent *MetrixAgent) updateMetric(metric models.MetricInfo) {
	resp, err := http.Post(agent.getUpdateMetricHandlerURL(metric), "text/plain", nil)
	if err != nil {
		logger.Errorf("failed to update metric: %v", err)
	}
	if resp != nil {
		resp.Body.Close()
	}
}

// updateMetricV2 обновление метрики через хендлеры второй версии.
func (agent *MetrixAgent) updateMetricV2(metric models.MetricInfo) {
	value := metric.GaugeValue()
	delta := metric.CounterValue()

	req := Metrics{
		ID:    metric.ID(),
		MType: string(metric.Type()),
		Delta: &delta,
		Value: &value,
	}

	err := agent.sendV2Request(agent.getUpdateMetricV2HandlerURL(), req)
	if err != nil {
		logger.Errorf("failed to update metric: %v", err)
	}
}

// updateMetricsBatch массововое обновление метрик через хендлеры второй версии.
func (agent *MetrixAgent) updateMetricsBatch(metrics []models.MetricInfo) {
	req := make(MetricsBatch, len(metrics))
	for i, metric := range metrics {
		value := metric.GaugeValue()
		delta := metric.CounterValue()

		req[i] = Metrics{
			ID:    metric.ID(),
			MType: string(metric.Type()),
			Delta: &delta,
			Value: &value,
		}
	}

	err := agent.sendV2Request(agent.getUpdateMetricBatchHandlerURL(), req)
	if err != nil {
		logger.Errorf("failed to batch update metrics: %v", err)
	}
}

func (agent *MetrixAgent) sendV2Request(url string, req easyjson.Marshaler) error {
	var (
		err      error
		httpReq  *http.Request
		reqBytes []byte
	)

	reqBytes, err = easyjson.Marshal(req)
	if err != nil {
		return err
	}

	reqBytes, err = tools.Compress(reqBytes)
	if err != nil {
		return err
	}

	if agent.cryptoPublicKey != "" {
		reqBytes, err = tools.Encrypt(agent.cryptoPublicKey, reqBytes)
		if err != nil {
			return err
		}
	}

	reqBody := bytes.NewBuffer(reqBytes)
	httpReq, err = http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return err
	}

	httpReq.Header.Set(tools.AcceptEncoding, "gzip")
	httpReq.Header.Set(tools.ContentEncoding, "gzip")
	httpReq.Header.Set(tools.ContentType, "application/json")

	if agent.privateKey != "" {
		var hashedReq string
		hashedReq, err = tools.CalcSHA256(reqBytes, agent.privateKey)
		if err != nil {
			return err
		}

		httpReq.Header.Set(tools.HashSHA256, hashedReq)
	}

	agent.retrier.Exec(func() bool {
		var resp *http.Response
		resp, err = http.DefaultClient.Do(httpReq)
		if resp != nil {
			resp.Body.Close()
		}
		return err != nil
	})

	return err
}

// getHandlerUrl создаёт URL-адрес для запроса на обновление метрик.
func (agent MetrixAgent) getUpdateMetricHandlerURL(metric models.MetricInfo) string {
	metricTypeStr := string(metric.Type())

	var metricValueStr string
	switch metric.Type() {
	case models.Gauge:
		metricValueStr = tools.FloatToStr(metric.GaugeValue())
	case models.Counter:
		metricValueStr = tools.IntToStr(metric.CounterValue())
	}

	return fmt.Sprintf("http://%s/update/%s/%s/%s/", agent.serverAddr, metricTypeStr, metric.ID(), metricValueStr)
}

// getUpdateMetricV2HandlerURL создаёт URL-адрес для запроса на обновление метрик в JSON формате.
func (agent MetrixAgent) getUpdateMetricV2HandlerURL() string {
	return fmt.Sprintf("http://%s/update/", agent.serverAddr)
}

// getUpdateMetricBatchHandlerURL создаёт URL-адрес для запроса на массовое обновление метрик в JSON формате.
func (agent MetrixAgent) getUpdateMetricBatchHandlerURL() string {
	return fmt.Sprintf("http://%s/updates/", agent.serverAddr)
}

//easyjson:json
type Metrics struct {
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
}

//easyjson:json
type MetricsBatch []Metrics
