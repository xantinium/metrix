// Пакет agent содержит реализацию агента для сбора метрик.
package agent

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/xantinium/metrix/internal/infrastructure/runtimemetrics"
	"github.com/xantinium/metrix/internal/models"
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
			log.Printf("failed to update metric: %v", err)
		}

		if resp != nil {
			resp.Body.Close()
		}
	}
}

// getHandlerUrl создаёт URL-адрес для запроса на обновление метрик.
func (agent MetrixAgent) getUpdateMetricHandlerURL(metric models.MetricInfo) string {
	b := strings.Builder{}

	b.WriteString("http://")
	b.WriteString(agent.serverAddr)
	b.WriteString("/update/")
	b.WriteString(string(metric.Type()))
	b.WriteString("/")
	b.WriteString(metric.Name())
	b.WriteString("/")

	switch metric.Type() {
	case models.Gauge:
		b.WriteString(floatToStr(metric.GaugeValue()))
	case models.Counter:
		b.WriteString(intToStr(metric.CounterValue()))
	}

	return b.String()
}

func floatToStr(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func intToStr(v int64) string {
	return strconv.FormatInt(v, 10)
}
