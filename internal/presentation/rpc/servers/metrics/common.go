package metrics

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xantinium/metrix/internal/models"
	"github.com/xantinium/metrix/internal/presentation/rpc/gen"
)

func parseMetricType(mType gen.MetricType) (models.MetricType, error) {
	switch mType {
	case gen.MetricType_GAUGE:
		return models.Gauge, nil
	case gen.MetricType_COUNTER:
		return models.Counter, nil
	default:
		return "", status.Errorf(codes.InvalidArgument, "unknown metric type: %q", mType)
	}
}

func parseMetric(metric *gen.Metric) (models.MetricInfo, error) {
	if metric == nil {
		return models.MetricInfo{}, status.Error(codes.InvalidArgument, "metric data is missing")
	}

	id := metric.GetId()
	if id == "" {
		return models.MetricInfo{}, status.Error(codes.InvalidArgument, "metric id cannot be empty")
	}

	mType, err := parseMetricType(metric.GetMType())
	if err != nil {
		return models.MetricInfo{}, err
	}

	switch mType {
	case models.Gauge:
		return models.NewGaugeMetric(id, metric.GetValue()), nil
	case models.Counter:
		return models.NewCounterMetric(id, metric.GetDelta()), nil
	default:
		return models.MetricInfo{}, status.Errorf(codes.InvalidArgument, "unknown metric type: %q", mType)
	}
}

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
