// Пакет models содержит описание объектов бизнес-логики.
package models

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")
)

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

// NewGaugeMetric создаёт новую метрику типа Gauge.
func NewGaugeMetric(name string, value float64) MetricInfo {
	return MetricInfo{
		metricName: name,
		metricType: Gauge,
		gaugeValue: value,
	}
}

// NewCounterMetric создаёт новую метрику типа Counter.
func NewCounterMetric(name string, value int64) MetricInfo {
	return MetricInfo{
		metricName:   name,
		metricType:   Counter,
		counterValue: value,
	}
}

// MetricInfo структура, описывающая метрику.
type MetricInfo struct {
	metricName   string
	metricType   MetricType
	gaugeValue   float64
	counterValue int64
}

// Name возвращает имя метрики.
func (info MetricInfo) Name() string {
	return info.metricName
}

// Type возвращает тип метрики.
func (info MetricInfo) Type() MetricType {
	return info.metricType
}

// GaugeValue возвращает значение метрики типа Gauge.
func (info MetricInfo) GaugeValue() float64 {
	return info.gaugeValue
}

// CounterValue возвращает значение метрики типа Counter.
func (info MetricInfo) CounterValue() int64 {
	return info.counterValue
}
