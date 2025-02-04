// Пакет handlers содержит хендлеры всех HTTP-запросов.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xantinium/metrix/internal/repository/metrics"
)

// httpMethod тип HTTP-метода.
type httpMethod = string

const (
	// HTTP-метод "GET"
	MethodGet httpMethod = http.MethodGet
	// HTTP-метод "POST"
	MethodPost httpMethod = http.MethodPost
)

// server интерфейс сервера, доступного в хендлерах.
type server interface {
	GetInternalRouter() *gin.Engine
	GetMetricsRepo() *metrics.MetricsRepository
}

// httpHandler общий тип для всех хендлеров.
type httpHandler = func(*gin.Context, server) (int, string, error)

// RegisterHandler добавляет хендлер handler в качестве обработчика
// паттерна pattern для метода method.
func RegisterHandler(server server, method httpMethod, pattern string, handler httpHandler) {
	server.GetInternalRouter().Handle(method, pattern, func(ctx *gin.Context) {
		statusCode, response, err := handler(ctx, server)
		if err != nil {
			ctx.String(statusCode, err.Error())
			return
		}

		ctx.String(statusCode, response)
	})
}
