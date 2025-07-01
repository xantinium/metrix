package rpc

import (
	"context"

	"github.com/xantinium/metrix/internal/models"
	"github.com/xantinium/metrix/internal/presentation/rpc/gen"
)

func metricTypeToGen(mType models.MetricType) gen.MetricType {
	switch mType {
	case models.Gauge:
		return gen.MetricType_GAUGE
	case models.Counter:
		return gen.MetricType_COUNTER
	default:
		return -1
	}
}

func metricToGen(metric models.MetricInfo) *gen.Metric {
	return &gen.Metric{
		Delta: metric.CounterValue(),
		Value: metric.GaugeValue(),
		Id:    metric.ID(),
		MType: metricTypeToGen(metric.Type()),
	}
}

// UpdateMetric обновление метрики.
func (client *Client) UpdateMetric(ctx context.Context, metric models.MetricInfo) error {
	_, err := client.getMetricsClient().UpdateMetric(ctx, &gen.UpdateMetricRequest{
		Metric: metricToGen(metric),
	})

	return err
}

// UpdateMetricsBatch массововое обновление метрик.
func (client *Client) UpdateMetricsBatch(ctx context.Context, metrics []models.MetricInfo) error {
	req := &gen.UpdateMetricsBatchRequest{
		Metrics: make([]*gen.Metric, len(metrics)),
	}

	for i := range metrics {
		req.Metrics[i] = metricToGen(metrics[i])
	}

	_, err := client.getMetricsClient().UpdateMetricsBatch(ctx, req)

	return err
}
