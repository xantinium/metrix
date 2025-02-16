// Пакет metrics содержит репозиторий для работы с метриками.
package metrics

import (
	"fmt"

	"github.com/xantinium/metrix/internal/models"
)

// NewMetricsRepository создаёт новый репозиторий метрик.
func NewMetricsRepository(storage MetricsStorage) *MetricsRepository {
	return &MetricsRepository{storage: storage}
}

// MetricsRepository структура, описывающая репозиторий метрик.
type MetricsRepository struct {
	storage MetricsStorage
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
func (repo *MetricsRepository) UpdateGaugeMetric(name string, value float64) error {
	err := repo.storage.UpdateGaugeMetric(name, value)
	if err != nil {
		return fmt.Errorf("failed to update gauge metric name=%s value=%f: %v", name, value, err)
	}

	return nil
}

// UpdateCounterMetric обновляет текущее значение метрики типа COUNTER
// с именем name, добавляя к нему значение value.
func (repo *MetricsRepository) UpdateCounterMetric(name string, value int64) error {
	err := repo.storage.UpdateCounterMetric(name, value)
	if err != nil {
		return fmt.Errorf("failed to update counter metric name=%s value=%d: %v", name, value, err)
	}

	return nil
}

// GetAllMetrics возвращает все существующие метрики.
func (repo *MetricsRepository) GetAllMetrics() ([]models.MetricInfo, error) {
	metrics, err := repo.storage.GetAllMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to get all metrics: %v", err)
	}

	return metrics, nil
}
