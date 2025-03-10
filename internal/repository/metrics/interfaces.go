package metrics

import (
	"context"

	"github.com/xantinium/metrix/internal/models"
)

// MetricsStorage интерфейс хранилища метрик.
type MetricsStorage interface {
	Destroy(ctx context.Context)
	GetGaugeMetric(ctx context.Context, name string) (float64, error)
	GetCounterMetric(ctx context.Context, name string) (int64, error)
	GetAllMetrics(ctx context.Context) ([]models.MetricInfo, error)
	UpdateGaugeMetric(ctx context.Context, name string, value float64) (float64, error)
	UpdateCounterMetric(ctx context.Context, name string, value int64) (int64, error)
	SaveMetrics(ctx context.Context) error
}

// DatabaseChecker интерфейс для проверки соединения с БД.
type DatabaseChecker interface {
	Ping(ctx context.Context) error
}
