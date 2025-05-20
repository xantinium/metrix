package v2handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"

	v2handlers "github.com/xantinium/metrix/internal/server/handlers/v2"
)

func ExampleGetMetricHandler() {
	// Составляем запрос на получение метрики
	// с идентификатором "Alloc" и типом "gauge".
	req := v2handlers.Metrics{
		ID:    "Alloc",
		MType: "gauge",
	}

	// Преобразуем в JSON.
	reqBytes, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	// Выполняем запрос на получение метрики.
	http.Post("/value", "application/json", bytes.NewBuffer(reqBytes))
}

func ExampleUpdateMetricHandler() {
	{
		// Составляем запрос на обновление метрики
		// с идентификатором "Alloc" и типом "gauge".
		// В качестве нового значения используем 12.6.
		req := v2handlers.Metrics{
			ID:    "Alloc",
			MType: "gauge",
			Value: newFloat64(12.6),
		}

		// Преобразуем в JSON.
		reqBytes, err := json.Marshal(req)
		if err != nil {
			panic(err)
		}

		// Выполняем запрос на обновление метрики.
		http.Post("/update", "application/json", bytes.NewBuffer(reqBytes))
	}

	{
		// Составляем запрос на обновление метрики
		// с идентификатором "PollCount" и типом "counter".
		// В качестве смещения используем 5.
		req := v2handlers.Metrics{
			ID:    "PollCount",
			MType: "counter",
			Delta: newInt64(5),
		}

		// Преобразуем в JSON.
		reqBytes, err := json.Marshal(req)
		if err != nil {
			panic(err)
		}

		// Выполняем запрос на обновление метрики.
		http.Post("/update", "application/json", bytes.NewBuffer(reqBytes))
	}
}

func ExampleUpdateMetricsHandler() {
	// Составляем запрос на батчевое обновление метрик.
	// Обновляем две метрики разных типов одновременно.
	req := []v2handlers.Metrics{
		{
			ID:    "Alloc",
			MType: "gauge",
			Value: newFloat64(12.6),
		},
		{
			ID:    "PollCount",
			MType: "counter",
			Delta: newInt64(5),
		},
	}

	// Преобразуем в JSON.
	reqBytes, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	// Выполняем запрос на получение метрики.
	http.Post("/updates", "application/json", bytes.NewBuffer(reqBytes))
}

func newInt64(v int64) *int64 {
	return &v
}

func newFloat64(v float64) *float64 {
	return &v
}
