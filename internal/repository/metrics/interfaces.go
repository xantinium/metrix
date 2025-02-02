package metrics

// MetricsStorage интерфейс хранилища метрик.
type MetricsStorage interface {
	UpdateGaugeMetric(name string, value float64) error
	UpdateCounterMetric(name string, value int64) error
}
