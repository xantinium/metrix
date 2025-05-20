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
	sendHTTPRequest(http.MethodPost, "/value", reqBytes)
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
		sendHTTPRequest(http.MethodPost, "/update", reqBytes)
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
		sendHTTPRequest(http.MethodPost, "/update", reqBytes)
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

	// Выполняем запрос на батчевое обновление метрик.
	sendHTTPRequest(http.MethodPost, "/updates", reqBytes)
}

func newInt64(v int64) *int64 {
	return &v
}

func newFloat64(v float64) *float64 {
	return &v
}

func sendHTTPRequest(method, url string, body []byte) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err == nil {
		resp.Body.Close()
	}
}
