// Пакет agent содержит реализацию агента для сбора метрик.
package agent

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/mailru/easyjson"

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
	UsingV2        bool
}

// NewMetrixAgent создаёт новый агент метрик.
func NewMetrixAgent(opts MetrixAgentOptions) *MetrixAgent {
	agent := &MetrixAgent{
		serverAddr:    opts.ServerAddr,
		metricsSource: runtimemetrics.NewRuntimeMetricsSource(opts.PollInterval),
		usingV2:       opts.UsingV2,
	}

	agent.worker = newMetrixAgentWorker(opts.ReportInterval, agent.UpdateMetrics)

	return agent
}

// MetrixAgent структура, описывающая агент метрик.
type MetrixAgent struct {
	serverAddr    string
	worker        *metrixAgentWorker
	metricsSource *runtimemetrics.RuntimeMetricsSource
	usingV2       bool
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
		if agent.usingV2 {
			agent.updateMetricsV2(metric)
		} else {
			agent.updateMetrics(metric)
		}
	}
}

func (agent *MetrixAgent) updateMetrics(metric models.MetricInfo) {
	resp, err := http.Post(agent.getUpdateMetricHandlerURL(metric), "text/plain", nil)

	if err != nil {
		logger.Errorf("failed to update metric: %v", err)
	}

	if resp != nil {
		resp.Body.Close()
	}
}

func (agent *MetrixAgent) updateMetricsV2(metric models.MetricInfo) {
	var (
		err      error
		reqBytes []byte
		resp     *http.Response
	)

	value := metric.GaugeValue()
	delta := metric.CounterValue()

	req := Metrics{
		ID:    metric.Name(),
		MType: string(metric.Type()),
		Delta: &delta,
		Value: &value,
	}

	reqBytes, err = easyjson.Marshal(req)
	if err != nil {
		logger.Errorf("failed to update metric: %v", err)
	}

	reqBody := bytes.NewBuffer(reqBytes)
	resp, err = http.Post(agent.getUpdateMetricV2HandlerURL(), "application/json", reqBody)

	if err != nil {
		logger.Errorf("failed to update metric: %v", err)
	}

	if resp != nil {
		resp.Body.Close()
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

// getUpdateMetricV2HandlerURL создаёт URL-адрес для запроса на обновление метрик в JSON формате.
func (agent MetrixAgent) getUpdateMetricV2HandlerURL() string {
	return fmt.Sprintf("http://%s/update", agent.serverAddr)
}

//easyjson:json
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
