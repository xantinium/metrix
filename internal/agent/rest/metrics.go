package rest

import (
	"context"
	"fmt"

	"github.com/xantinium/metrix/internal/models"
)

// UpdateMetric обновление метрики.
func (client *Client) UpdateMetric(ctx context.Context, metric models.MetricInfo) error {
	value := metric.GaugeValue()
	delta := metric.CounterValue()

	req := Metrics{
		ID:    metric.ID(),
		MType: string(metric.Type()),
		Delta: &delta,
		Value: &value,
	}

	err := client.sendRequest(ctx, client.getUpdateMetricHandlerURL(), req)
	if err != nil {
		err = fmt.Errorf("failed to update metric: %w", err)
		client.logError(err)
	}

	return err
}

// UpdateMetricsBatch массововое обновление метрик.
func (client *Client) UpdateMetricsBatch(ctx context.Context, metrics []models.MetricInfo) error {
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

	err := client.sendRequest(ctx, client.getUpdateMetricBatchHandlerURL(), req)
	if err != nil {
		err = fmt.Errorf("failed to batch update metrics: %w", err)
		client.logError(err)
	}

	return err
}

// getUpdateMetricHandlerURL создаёт URL-адрес для запроса на обновление метрик в JSON формате.
func (client *Client) getUpdateMetricHandlerURL() string {
	return fmt.Sprintf("http://%s/update/", client.serverAddr)
}

// getUpdateMetricBatchHandlerURL создаёт URL-адрес для запроса на массовое обновление метрик в JSON формате.
func (client *Client) getUpdateMetricBatchHandlerURL() string {
	return fmt.Sprintf("http://%s/updates/", client.serverAddr)
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
