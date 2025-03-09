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

// GetGaugeMetric возвращает метрику типа GAUGE по имени name.
func (repo *MetricsRepository) GetGaugeMetric(name string) (float64, error) {
	return repo.storage.GetGaugeMetric(name)
}

// GetCounterMetric возвращает метрику типа COUNTER по имени name.
func (repo *MetricsRepository) GetCounterMetric(name string) (int64, error) {
	return repo.storage.GetCounterMetric(name)
}

// UpdateGaugeMetric обновляет текущее значение метрики типа GAUGE
// с именем name, перезаписывая его значением value.
func (repo *MetricsRepository) UpdateGaugeMetric(name string, value float64) (float64, error) {
	updatedValue, err := repo.storage.UpdateGaugeMetric(name, value)
	if err != nil {
		return 0, fmt.Errorf("failed to update gauge metric name=%s value=%f: %v", name, value, err)
	}

	repo.onMetricsUpdate()
	return updatedValue, nil
}

// UpdateCounterMetric обновляет текущее значение метрики типа COUNTER
// с именем name, добавляя к нему значение value.
func (repo *MetricsRepository) UpdateCounterMetric(name string, value int64) (int64, error) {
	updatedValue, err := repo.storage.UpdateCounterMetric(name, value)
	if err != nil {
		return 0, fmt.Errorf("failed to update counter metric name=%s value=%d: %v", name, value, err)
	}

	repo.onMetricsUpdate()
	return updatedValue, nil
}

// GetAllMetrics возвращает все существующие метрики.
func (repo *MetricsRepository) GetAllMetrics() ([]models.MetricInfo, error) {
	metrics, err := repo.storage.GetAllMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to get all metrics: %v", err)
	}

	return metrics, nil
}

// CheckDatabase проверяет соединение с БД.
func (repo *MetricsRepository) CheckDatabase(ctx context.Context) error {
	return repo.dbChecker.CheckDatabase(ctx)
}

func (repo *MetricsRepository) onMetricsUpdate() {
	if repo.syncMetrics {
		err := repo.storage.SaveMetrics()
		if err != nil {
			logger.Errorf("failed to sync metrics: %v", err)
		}
	}
}
