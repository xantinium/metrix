// Пакет agent содержит реализацию агента для сбора метрик.
package agent

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

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
	ReportInterval time.Duration
}

// NewMetrixAgent создаёт новый агент метрик.
func NewMetrixAgent(opts MetrixAgentOptions) *MetrixAgent {
	agent := &MetrixAgent{
		serverAddr:    opts.ServerAddr,
		metricsSource: runtimemetrics.NewRuntimeMetricsSource(opts.PollInterval),
		retrier:       tools.DefaulRetrier,
	}

	agent.worker = newMetrixAgentWorker(opts.ReportInterval, agent.UpdateMetrics)

	return agent
}

// MetrixAgent структура, описывающая агент метрик.
type MetrixAgent struct {
	serverAddr    string
	worker        *metrixAgentWorker
	metricsSource *runtimemetrics.RuntimeMetricsSource
	retrier       *tools.Retrier
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
	agent.updateMetricsBatch(agent.metricsSource.GetSnapshot())
}

// updateMetrics обновление метрики через хендлеры первой версии.
//
// Deprecated: метод устарел, следует использовать updateMetricsV2.
func (agent *MetrixAgent) updateMetrics(metric models.MetricInfo) {
	resp, err := http.Post(agent.getUpdateMetricHandlerURL(metric), "text/plain", nil)
	if err != nil {
		logger.Errorf("failed to update metric: %v", err)
	}
	if resp != nil {
		resp.Body.Close()
	}
}

// updateMetricsV2 обновление метрики через хендлеры второй версии.
func (agent *MetrixAgent) updateMetricsV2(metric models.MetricInfo) {
	var (
		err      error
		httpReq  *http.Request
		reqBytes []byte
		resp     *http.Response
	)

	value := metric.GaugeValue()
	delta := metric.CounterValue()

	req := Metrics{
		ID:    metric.ID(),
		MType: string(metric.Type()),
		Delta: &delta,
		Value: &value,
	}

	reqBytes, err = easyjson.Marshal(req)
	if err != nil {
		logger.Errorf("failed to update metric: %v", err)
		return
	}

	reqBytes, err = tools.Compress(reqBytes)
	if err != nil {
		logger.Errorf("failed to update metric: %v", err)
		return
	}

	reqBody := bytes.NewBuffer(reqBytes)
	httpReq, err = http.NewRequest(http.MethodPost, agent.getUpdateMetricV2HandlerURL(), reqBody)
	if err != nil {
		logger.Errorf("failed to update metric: %v", err)
		return
	}

	httpReq.Header.Set("Accept-Encoding", "gzip")
	httpReq.Header.Set("Content-Encoding", "gzip")
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(httpReq)
	if err != nil {
		logger.Errorf("failed to update metric: %v", err)
	}
	if resp != nil {
		resp.Body.Close()
	}
}

// updateMetricsBatch массововое обновление метрик через хендлеры второй версии.
func (agent *MetrixAgent) updateMetricsBatch(metrics []models.MetricInfo) {
	var (
		err      error
		httpReq  *http.Request
		reqBytes []byte
	)

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

	reqBytes, err = easyjson.Marshal(req)
	if err != nil {
		logger.Errorf("failed to update metric: %v", err)
		return
	}

	reqBytes, err = tools.Compress(reqBytes)
	if err != nil {
		logger.Errorf("failed to update metric: %v", err)
		return
	}

	reqBody := bytes.NewBuffer(reqBytes)
	httpReq, err = http.NewRequest(http.MethodPost, agent.getUpdateMetricBatchHandlerURL(), reqBody)
	if err != nil {
		logger.Errorf("failed to update metric: %v", err)
		return
	}

	httpReq.Header.Set("Accept-Encoding", "gzip")
	httpReq.Header.Set("Content-Encoding", "gzip")
	httpReq.Header.Set("Content-Type", "application/json")

	agent.retrier.Exec(func() bool {
		var resp *http.Response
		resp, err = http.DefaultClient.Do(httpReq)
		if err != nil {
			logger.Errorf("failed to update metric: %v", err)
		}
		if resp != nil {
			resp.Body.Close()
		}
		return err != nil
	})
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

	return fmt.Sprintf("http://%s/update/%s/%s/%s", agent.serverAddr, metricTypeStr, metric.ID(), metricValueStr)
}

// getUpdateMetricV2HandlerURL создаёт URL-адрес для запроса на обновление метрик в JSON формате.
func (agent MetrixAgent) getUpdateMetricV2HandlerURL() string {
	return fmt.Sprintf("http://%s/update", agent.serverAddr)
}

// getUpdateMetricBatchHandlerURL создаёт URL-адрес для запроса на массовое обновление метрик в JSON формате.
func (agent MetrixAgent) getUpdateMetricBatchHandlerURL() string {
	return fmt.Sprintf("http://%s/updates", agent.serverAddr)
}

//easyjson:json
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

//easyjson:json
type MetricsBatch []Metrics
