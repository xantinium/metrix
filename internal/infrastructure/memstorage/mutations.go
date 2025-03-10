package memstorage

import "context"

// UpdateGaugeMetric обновляет текущее значение метрики типа Gauge
// с именем name, перезаписывая его значением value.
//
// Возвращает обновлённое значение метрики.
func (storage *MemStorage) UpdateGaugeMetric(_ context.Context, name string, value float64) (float64, error) {
	storage.mx.Lock()
	defer storage.mx.Unlock()

	storage.gaugeMetrics[name] = value

	return storage.gaugeMetrics[name], nil
}

// UpdateCounterMetric обновляет текущее значение метрики типа Counter
// с именем name, добавляя к нему значение value.
//
// Возвращает обновлённое значение метрики.
func (storage *MemStorage) UpdateCounterMetric(_ context.Context, name string, value int64) (int64, error) {
	storage.mx.Lock()
	defer storage.mx.Unlock()

	storage.counterMetrics[name] += value

	return storage.counterMetrics[name], nil
}
