package memstorage

import "context"

// UpdateGaugeMetric обновляет текущее значение метрики типа Gauge
// с идентификатором id, перезаписывая его значением value.
//
// Возвращает обновлённое значение метрики.
func (storage *MemStorage) UpdateGaugeMetric(_ context.Context, id string, value float64) (float64, error) {
	storage.mx.Lock()
	defer storage.mx.Unlock()

	storage.gaugeMetrics[id] = value

	return storage.gaugeMetrics[id], nil
}

// UpdateCounterMetric обновляет текущее значение метрики типа Counter
// с идентификатором id, добавляя к нему значение value.
//
// Возвращает обновлённое значение метрики.
func (storage *MemStorage) UpdateCounterMetric(_ context.Context, id string, value int64) (int64, error) {
	storage.mx.Lock()
	defer storage.mx.Unlock()

	storage.counterMetrics[id] += value

	return storage.counterMetrics[id], nil
}
