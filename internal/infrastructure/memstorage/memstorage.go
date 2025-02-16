// Пакет memstorage содержит реализацию хранилища метрик.
// На данный момент, все данные хранятся в оперативной памяти.
package memstorage

import "sync"

// NewMemStorage создаёт новое хранилище метрик.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}
}

// MemStorage структура, реализующая хранилище метрик.
type MemStorage struct {
	mx             sync.RWMutex
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
}

// GetGaugeMetric возвращает метрику типа GAUGE по имени name.
func (storage *MemStorage) GetGaugeMetric(name string) (float64, error) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	return storage.gaugeMetrics[name], nil
}

// GetCounterMetric возвращает метрику типа COUNTER по имени name.
func (storage *MemStorage) GetCounterMetric(name string) (int64, error) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	return storage.counterMetrics[name], nil
}

// UpdateGaugeMetric обновляет текущее значение метрики типа GAUGE
// с именем name, перезаписывая его значением value.
func (storage *MemStorage) UpdateGaugeMetric(name string, value float64) error {
	storage.mx.Lock()
	defer storage.mx.Unlock()

	storage.gaugeMetrics[name] = value

	return nil
}

// UpdateCounterMetric обновляет текущее значение метрики типа COUNTER
// с именем name, добавляя к нему значение value.
func (storage *MemStorage) UpdateCounterMetric(name string, value int64) error {
	storage.mx.Lock()
	defer storage.mx.Unlock()

	storage.counterMetrics[name] += value

	return nil
}
