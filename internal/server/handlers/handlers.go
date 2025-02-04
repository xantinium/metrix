// Пакет handlers содержит хендлеры всех HTTP-запросов.
package handlers

import (
	"fmt"
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
type httpHandler = func(*gin.Context, server) (int, []byte, error)

// RegisterHandler добавляет хендлер handler в качестве обработчика
// паттерна pattern для метода method.
func RegisterHandler(server server, method httpMethod, pattern string, handler httpHandler) {
	server.GetInternalRouter().Handle(method, pattern, func(ctx *gin.Context) {
		statusCode, response, err := handler(ctx, server)
		if err != nil {
			ctx.JSON(statusCode, createErrResp(err))
			return
		}

		ctx.JSON(statusCode, response)
	})
}

func createErrResp(err error) []byte {
	return []byte(fmt.Sprintf("{\"err\":\"%s\"}", err))
}
