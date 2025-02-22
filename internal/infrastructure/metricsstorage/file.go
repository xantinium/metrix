package metricsstorage

import (
	"github.com/mailru/easyjson"
	"github.com/xantinium/metrix/internal/models"
)

type metricItem struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Delta int64   `json:"delta"`
	Value float64 `json:"value"`
}

//easyjson:json
type metricsStruct struct {
	Metrics []metricItem `json:"metrics"`
}

// SaveMetrics сохраняет текущие значения метрик в файл.
func (storage *MetricsStorage) SaveMetrics() error {
	metrics, err := storage.GetAllMetrics()
	if err != nil {
		return err
	}

	metrisToSave := metricsStruct{Metrics: make([]metricItem, len(metrics))}
	for i := range metrics {
		item := metricItem{
			Name: metrics[i].Name(),
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

	storage.mx.Lock()
	defer storage.mx.Unlock()

	_, err = storage.file.WriteAt(bytes, 0)
	return err
}
