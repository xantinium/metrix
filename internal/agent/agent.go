// Пакет agent содержит реализацию агента для сбора метрик.
package agent

import (
	"fmt"
	"net/http"

	"github.com/xantinium/metrix/internal/infrastructure/runtimemetrics"
	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/models"
	"github.com/xantinium/metrix/internal/tools"
)

// MetrixAgentOptions параметры агента метрик.
type MetrixAgentOptions struct {
	ServerAddr     string
	PollInterval   int
	ReportInterval int
}

// NewMetrixAgent создаёт новый агент метрик.
func NewMetrixAgent(opts MetrixAgentOptions) *MetrixAgent {
	agent := &MetrixAgent{
		serverAddr:    opts.ServerAddr,
		metricsSource: runtimemetrics.NewRuntimeMetricsSource(opts.PollInterval),
	}

	agent.worker = newMetrixAgentWorker(opts.ReportInterval, agent.UpdateMetrics)

	return agent
}

// MetrixAgent структура, описывающая агент метрик.
type MetrixAgent struct {
	serverAddr    string
	worker        *metrixAgentWorker
	metricsSource *runtimemetrics.RuntimeMetricsSource
}

// Run запускает агента метрик.
func (agent *MetrixAgent) Run() {
	agent.metricsSource.Run()
	agent.worker.Run()
}

// Run прекращает работу агента метрик.
func (agent *MetrixAgent) Stop() {
	agent.metricsSource.Stop()
	agent.worker.Stop()
}

// UpdateMetrics обновляет метрики на сервере.
func (agent *MetrixAgent) UpdateMetrics() {
	metrics := agent.metricsSource.GetSnapshot()

	for _, metric := range metrics {
		resp, err := http.Post(agent.getUpdateMetricHandlerURL(metric), "text/plain", nil)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to update metric: %v", err))
		}

		if resp != nil {
			resp.Body.Close()
		}
	}
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

	return fmt.Sprintf("http://%s/update/%s/%s/%s", agent.serverAddr, metricTypeStr, metric.Name(), metricValueStr)
}
