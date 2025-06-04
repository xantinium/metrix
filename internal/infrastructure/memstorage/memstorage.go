// Package memstorage содержит реализацию хранилища метрик,
// все данные которого хранятся в оперативной памяти.
// Дополнительно, имеется возможность записи/чтения метрик из файла.
package memstorage

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/mailru/easyjson"

	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/models"
)

// NewMemStorage создаёт новое хранилище метрик.
// При необходимости, восстанавливает предыдущие знаениченя метрик.
func NewMemStorage(path string, restore bool) (*MemStorage, error) {
	var err error

	storage := &MemStorage{
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

// MemStorage структура, реализующая хранилище метрик.
type MemStorage struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
	fileW          *fileWriter
	mx             sync.RWMutex
}

func (storage *MemStorage) restore(path string) error {
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
			storage.gaugeMetrics[metric.ID] = metric.Value
		case string(models.Counter):
			storage.counterMetrics[metric.ID] = metric.Delta
		}
	}

	return nil
}

// Destroy уничтожает хранилище метрик.
func (storage *MemStorage) Destroy(ctx context.Context) {
	err := storage.SaveMetrics(ctx)
	if err != nil {
		logger.Errorf("failed to save metrics: %v", err)
	}

	storage.fileW.Wait()
}
