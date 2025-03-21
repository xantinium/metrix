package v2handlers

import (
	"fmt"

	"github.com/xantinium/metrix/internal/models"
)

//easyjson:json
type Metrics struct {
	ID    string   `json:"id"`              // идентификатор метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

//easyjson:json
type MetricsBatch []Metrics

// parseType парсит тип метрики.
func parseType(maybeMetricType string) (models.MetricType, error) {
	switch maybeMetricType {
	case string(models.Gauge):
		return models.Gauge, nil
	case string(models.Counter):
		return models.Counter, nil
	default:
		return "", fmt.Errorf("unknown metric type: %q", maybeMetricType)
	}
}

func parseMetric(rawMetric Metrics) (models.MetricInfo, error) {
	var (
		err        error
		metricID   string
		metricType models.MetricType
	)

	metricID = rawMetric.ID
	if metricID == "" {
		return models.MetricInfo{}, fmt.Errorf("metric id cannot be empty")
	}

	metricType, err = parseType(rawMetric.MType)
	if err != nil {
		return models.MetricInfo{}, err
	}

	switch metricType {
	case models.Gauge:
		if rawMetric.Value == nil {
			return models.MetricInfo{}, fmt.Errorf("value is missing")
		}

		return models.NewGaugeMetric(metricID, *rawMetric.Value), nil
	case models.Counter:
		if rawMetric.Delta == nil {
			return models.MetricInfo{}, fmt.Errorf("value is missing")
		}

		return models.NewCounterMetric(metricID, *rawMetric.Delta), nil
	default:
		return models.MetricInfo{}, fmt.Errorf("unknown metric type: %q", metricType)
	}
}
