package metrics

import (
	"context"

	"github.com/xantinium/metrix/internal/models"
)

// MetricsStorage интерфейс хранилища метрик.
type MetricsStorage interface {
	Destroy(ctx context.Context)
	GetGaugeMetric(ctx context.Context, id string) (float64, error)
	GetCounterMetric(ctx context.Context, id string) (int64, error)
	GetAllMetrics(ctx context.Context) ([]models.MetricInfo, error)
	UpdateGaugeMetric(ctx context.Context, id string, value float64) (float64, error)
	UpdateCounterMetric(ctx context.Context, id string, value int64) (int64, error)
	SaveMetrics(ctx context.Context) error
}

// DatabaseChecker интерфейс для проверки соединения с БД.
type DatabaseChecker interface {
	Ping(ctx context.Context) error
}
