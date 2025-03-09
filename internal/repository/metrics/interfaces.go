package metrics

import "github.com/xantinium/metrix/internal/models"

// MetricsStorage интерфейс хранилища метрик.
type MetricsStorage interface {
	GetGaugeMetric(name string) (float64, error)
	GetCounterMetric(name string) (int64, error)
	GetAllMetrics() ([]models.MetricInfo, error)
	UpdateGaugeMetric(name string, value float64) (float64, error)
	UpdateCounterMetric(name string, value int64) (int64, error)
	SaveMetrics() error
}
