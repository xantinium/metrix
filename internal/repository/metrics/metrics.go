// Пакет metrics содержит репозиторий для работы с метриками.
package metrics

import (
	"context"
	"fmt"

	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/models"
)

type MetricsRepositoryOptions struct {
	Storage MetricsStorage
	// SyncMetrics нужно ли сохранять метрики после каждой мутации.
	SyncMetrics bool
	DBChecker   DatabaseChecker
}

// NewMetricsRepository создаёт новый репозиторий метрик.
func NewMetricsRepository(opts MetricsRepositoryOptions) *MetricsRepository {
	return &MetricsRepository{
		storage:     opts.Storage,
		syncMetrics: opts.SyncMetrics,
		dbChecker:   opts.DBChecker,
	}
}

// MetricsRepository структура, описывающая репозиторий метрик.
type MetricsRepository struct {
	storage     MetricsStorage
	syncMetrics bool
	dbChecker   DatabaseChecker
}

// GetGaugeMetric возвращает метрику типа Gauge по идентификатору id.
func (repo *MetricsRepository) GetGaugeMetric(ctx context.Context, id string) (float64, error) {
	return repo.storage.GetGaugeMetric(ctx, id)
}

// GetCounterMetric возвращает метрику типа Counter по идентификатору id.
func (repo *MetricsRepository) GetCounterMetric(ctx context.Context, id string) (int64, error) {
	return repo.storage.GetCounterMetric(ctx, id)
}

// UpdateGaugeMetric обновляет текущее значение метрики типа Gauge
// с идентификатором id, перезаписывая его значением value.
func (repo *MetricsRepository) UpdateGaugeMetric(ctx context.Context, id string, value float64) (float64, error) {
	updatedValue, err := repo.storage.UpdateGaugeMetric(ctx, id, value)
	if err != nil {
		return 0, fmt.Errorf("failed to update gauge metric id=%s value=%f: %v", id, value, err)
	}

	repo.onMetricsUpdate(ctx)
	return updatedValue, nil
}

// UpdateCounterMetric обновляет текущее значение метрики типа Counter
// с идентификатором id, добавляя к нему значение value.
func (repo *MetricsRepository) UpdateCounterMetric(ctx context.Context, id string, value int64) (int64, error) {
	updatedValue, err := repo.storage.UpdateCounterMetric(ctx, id, value)
	if err != nil {
		return 0, fmt.Errorf("failed to update counter metric id=%s value=%d: %v", id, value, err)
	}

	repo.onMetricsUpdate(ctx)
	return updatedValue, nil
}

// GetAllMetrics возвращает все существующие метрики.
func (repo *MetricsRepository) GetAllMetrics(ctx context.Context) ([]models.MetricInfo, error) {
	metrics, err := repo.storage.GetAllMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all metrics: %v", err)
	}

	return metrics, nil
}

// CheckDatabase проверяет соединение с БД.
func (repo *MetricsRepository) CheckDatabase(ctx context.Context) error {
	return repo.dbChecker.Ping(ctx)
}

func (repo *MetricsRepository) onMetricsUpdate(ctx context.Context) {
	if repo.syncMetrics {
		err := repo.storage.SaveMetrics(ctx)
		if err != nil {
			logger.Errorf("failed to sync metrics: %v", err)
		}
	}
}
