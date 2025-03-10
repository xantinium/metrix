package memstorage

import (
	"context"

	"github.com/xantinium/metrix/internal/models"
)

// GetGaugeMetric возвращает метрику типа Gauge по имени name.
func (storage *MemStorage) GetGaugeMetric(_ context.Context, name string) (float64, error) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	value, exists := storage.gaugeMetrics[name]
	if !exists {
		return 0, models.ErrNotFound
	}

	return value, nil
}

// GetCounterMetric возвращает метрику типа Counter по имени name.
func (storage *MemStorage) GetCounterMetric(_ context.Context, name string) (int64, error) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	value, exists := storage.counterMetrics[name]
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
	for name, value := range storage.gaugeMetrics {
		metrics[i] = models.NewGaugeMetric(name, value)
		i++
	}
	for name, value := range storage.counterMetrics {
		metrics[i] = models.NewCounterMetric(name, value)
		i++
	}

	return metrics, nil
}
