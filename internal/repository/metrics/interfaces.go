package metrics

// MetricsStorage интерфейс хранилища метрик.
type MetricsStorage interface {
	GetGaugeMetric(name string) (float64, error)
	GetCounterMetric(name string) (int64, error)
	UpdateGaugeMetric(name string, value float64) error
	UpdateCounterMetric(name string, value int64) error
}
