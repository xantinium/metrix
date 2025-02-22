// Пакет metricsstorage содержит реализацию хранилища метрик.
// На данный момент, все данные хранятся в оперативной памяти.
// Дополнительно, имеется возможность записи/чтения метрик из файла.
package metricsstorage

import (
	"errors"
	"os"
	"sync"

	"github.com/mailru/easyjson"
	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/models"
)

// NewMemStorage создаёт новое хранилище метрик.
// При необходимости, восстанавливает предыдущие знаениченя метрик.
func NewMetricsStorage(path string, restore bool) (*MetricsStorage, error) {
	var err error

	storage := &MetricsStorage{
		fileW:          &fileWriter{path: path},
		gaugeMetrics:   make(map[string]float64),
		counterMetrics: make(map[string]int64),
	}

	if restore {
		err = storage.restore(path)
		if err != nil {
			return nil, err
		}
	}

	return storage, nil
}

// MetricsStorage структура, реализующая хранилище метрик.
type MetricsStorage struct {
	mx             sync.RWMutex
	fileW          *fileWriter
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
}

func (storage *MetricsStorage) restore(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		// Если файл не существует, то просто выходим.
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	rawMetrics, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	metrics := new(metricsStruct)
	err = easyjson.Unmarshal(rawMetrics, metrics)
	if err != nil {
		return err
	}

	storage.mx.Lock()
	defer storage.mx.Unlock()

	for _, metric := range metrics.Metrics {
		switch metric.Type {
		case string(models.Gauge):
			storage.gaugeMetrics[metric.Name] = metric.Value
		case string(models.Counter):
			storage.counterMetrics[metric.Name] = metric.Delta
		}
	}

	return nil
}

// Destroy уничтожает хранилище метрик.
func (storage *MetricsStorage) Destroy() {
	err := storage.SaveMetrics()
	if err != nil {
		logger.Errorf("failed to save metrics: %v", err)
	}

	storage.fileW.Wait()
}
