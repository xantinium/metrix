package metricsstorage

import "github.com/xantinium/metrix/internal/models"

// GetGaugeMetric возвращает метрику типа GAUGE по имени name.
func (storage *MetricsStorage) GetGaugeMetric(name string) (float64, error) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	value, exists := storage.gaugeMetrics[name]
	if !exists {
		return 0, models.ErrNotFound
	}

	return value, nil
}

// GetCounterMetric возвращает метрику типа COUNTER по имени name.
func (storage *MetricsStorage) GetCounterMetric(name string) (int64, error) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	value, exists := storage.counterMetrics[name]
	if !exists {
		return 0, models.ErrNotFound
	}

	return value, nil
}

// GetAllMetrics возвращает все существующие метрики.
func (storage *MetricsStorage) GetAllMetrics() ([]models.MetricInfo, error) {
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

// UpdateGaugeMetric обновляет текущее значение метрики типа GAUGE
// с именем name, перезаписывая его значением value.
func (storage *MetricsStorage) UpdateGaugeMetric(name string, value float64) (float64, error) {
	storage.mx.Lock()
	defer storage.mx.Unlock()

	storage.gaugeMetrics[name] = value

	return storage.gaugeMetrics[name], nil
}

// UpdateCounterMetric обновляет текущее значение метрики типа COUNTER
// с именем name, добавляя к нему значение value.
func (storage *MetricsStorage) UpdateCounterMetric(name string, value int64) (int64, error) {
	storage.mx.Lock()
	defer storage.mx.Unlock()

	storage.counterMetrics[name] += value

	return storage.counterMetrics[name], nil
}
