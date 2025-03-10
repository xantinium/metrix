package memstorage

import (
	"context"
	"os"
	"sync"

	"github.com/mailru/easyjson"
	"github.com/xantinium/metrix/internal/models"
)

type metricItem struct {
	ID    string  `json:"id"`
	Type  string  `json:"type"`
	Delta int64   `json:"delta"`
	Value float64 `json:"value"`
}

//easyjson:json
type metricsStruct struct {
	Metrics []metricItem `json:"metrics"`
}

// SaveMetrics сохраняет текущие значения метрик в файл.
func (storage *MemStorage) SaveMetrics(ctx context.Context) error {
	metrics, err := storage.GetAllMetrics(ctx)
	if err != nil {
		return err
	}

	metrisToSave := metricsStruct{Metrics: make([]metricItem, len(metrics))}
	for i := range metrics {
		item := metricItem{
			ID:   metrics[i].ID(),
			Type: string(metrics[i].Type()),
		}
		switch metrics[i].Type() {
		case models.Gauge:
			item.Value = metrics[i].GaugeValue()
		case models.Counter:
			item.Delta = metrics[i].CounterValue()
		}

		metrisToSave.Metrics[i] = item
	}

	bytes, err := easyjson.Marshal(metrisToSave)
	if err != nil {
		return err
	}

	return storage.fileW.Write(bytes)
}

type fileWriter struct {
	mx   sync.Mutex
	wg   sync.WaitGroup
	path string
}

func (w *fileWriter) Write(data []byte) error {
	w.wg.Add(1)
	w.mx.Lock()
	defer func() {
		w.mx.Unlock()
		w.wg.Done()
	}()

	file, err := os.OpenFile(w.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func (w *fileWriter) Wait() {
	w.wg.Wait()
}
