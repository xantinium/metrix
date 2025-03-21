package memstorage

import (
	"context"

	"github.com/xantinium/metrix/internal/models"
)

// GetGaugeMetric возвращает метрику типа Gauge по идентификатору id.
func (storage *MemStorage) GetGaugeMetric(_ context.Context, id string) (float64, error) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	value, exists := storage.gaugeMetrics[id]
	if !exists {
		return 0, models.ErrNotFound
	}

	return value, nil
}

// GetCounterMetric возвращает метрику типа Counter по идентификатору id.
func (storage *MemStorage) GetCounterMetric(_ context.Context, id string) (int64, error) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	value, exists := storage.counterMetrics[id]
	if !exists {
		return 0, models.ErrNotFound
	}

	return value, nil
}

// GetAllMetrics возвращает все существующие метрики.
func (storage *MemStorage) GetAllMetrics(_ context.Context) ([]models.MetricInfo, error) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	metrics := make([]models.MetricInfo, len(storage.gaugeMetrics)+len(storage.counterMetrics))

	i := 0
	for id, value := range storage.gaugeMetrics {
		metrics[i] = models.NewGaugeMetric(id, value)
		i++
	}
	for id, value := range storage.counterMetrics {
		metrics[i] = models.NewCounterMetric(id, value)
		i++
	}

	return metrics, nil
}
