// Пакет models содержит описание объектов бизнес-логики.
package models

import "fmt"

// MetricType тип метрики.
type MetricType string

const (
	// Gauge перезаписываемая метрика.
	Gauge MetricType = "gauge"
	// Counter суммируемая метрика.
	Counter MetricType = "counter"
)

// ParseStringAsMetricType парсит строку в тип метрики.
func ParseStringAsMetricType(maybeMetricType string) (MetricType, error) {
	switch maybeMetricType {
	case string(Gauge):
		return Gauge, nil
	case string(Counter):
		return Counter, nil
	default:
		return "", fmt.Errorf("unknown metric type")
	}
}
