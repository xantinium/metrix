// Пакет models содержит описание объектов бизнес-логики.
package models

import "fmt"

// MetricType тип метрики.
type MetricType string

const (
	// GAUGE перезаписываемая метрика.
	GAUGE MetricType = "gauge"
	// COUNTER суммируемая метрика.
	COUNTER MetricType = "counter"
)

// ParseStringAsMetricType парсит строку в тип метрики.
func ParseStringAsMetricType(maybeMetricType string) (MetricType, error) {
	switch maybeMetricType {
	case string(GAUGE):
		return GAUGE, nil
	case string(COUNTER):
		return COUNTER, nil
	default:
		return "", fmt.Errorf("unknown metric type")
	}
}
