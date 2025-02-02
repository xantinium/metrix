// Пакет handlers содержит хендлеры всех HTTP-запросов.
package handlers

import (
	"net/http"

	"github.com/xantinium/metrix/internal/repository/metrics"
)

// httpMethod тип HTTP-метода.
type httpMethod = string

const (
	// HTTP-метод "GET"
	METHOD_GET httpMethod = http.MethodGet
	// HTTP-метод "POST"
	METHOD_POST httpMethod = http.MethodPost
)

// server интерфейс сервера, доступного в хендлерах.
type server interface {
	GetInternalMux() *http.ServeMux
	GetMetricsRepo() *metrics.MetricsRepository
}

// httpHandler общий тип для всех хендлеров.
type httpHandler = func(server, *http.Request) (int, []byte, error)

// RegisterHandler добавляет хендлер handler в качестве обработчика
// паттерна pattern для метода method.
func RegisterHandler(server server, method httpMethod, pattern string, handler httpHandler) {
	server.GetInternalMux().HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		statusCode, response, err := handler(server, r)
		if err != nil {
			w.WriteHeader(statusCode)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(statusCode)
		w.Write(response)
	})
}
