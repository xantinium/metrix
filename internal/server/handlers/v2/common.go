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
