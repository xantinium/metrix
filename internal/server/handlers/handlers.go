// Пакет handlers содержит хендлеры всех HTTP-запросов.
package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"

	"github.com/xantinium/metrix/internal/server/interfaces"
)

// httpHandler общий тип для всех хендлеров.
type httpHandler = func(*gin.Context, interfaces.Server) (int, string, error)

// httpV2Handler общий тип для всех JSON-хендлеров.
type httpV2Handler = func(*gin.Context, interfaces.Server) (int, easyjson.Marshaler, error)

// RegisterHandler добавляет хендлер handler в качестве обработчика
// паттерна pattern для метода method.
func RegisterHandler(server interfaces.Server, method string, pattern string, handler httpHandler) {
	server.GetInternalRouter().Handle(method, pattern, func(ctx *gin.Context) {
		statusCode, response, err := handler(ctx, server)
		if err != nil {
			ctx.String(statusCode, err.Error())
			return
		}

		ctx.String(statusCode, response)
	})
}

const baseTemplate = "<html><head><title>Metrix</title></head><body>%s</body></html>"

// RegisterHTMLHandler добавляет хендлер handler в качестве обработчика
// паттерна pattern. Ожидается, что хендлер вернёт валидную HTML-строку.
func RegisterHTMLHandler(server interfaces.Server, pattern string, handler httpHandler) {
	server.GetInternalRouter().Handle(http.MethodGet, pattern, func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "text/html")

		statusCode, response, err := handler(ctx, server)
		if err != nil {
			ctx.String(http.StatusOK, baseTemplate, fmt.Sprintf("status: %d, err: %s", statusCode, err.Error()))
			return
		}

		ctx.String(http.StatusOK, baseTemplate, response)
	})
}

// RegisterV2Handler добавляет хендлер handler в качестве обработчика
// паттерна pattern для метода method.
//
// В отличии от функции RegisterHandler, предназначен
// для хендлеров, работающих с данными в формате JSON.
func RegisterV2Handler(server interfaces.Server, method string, pattern string, handler httpV2Handler) {
	server.GetInternalRouter().Handle(method, pattern, func(ctx *gin.Context) {
		var (
			err           error
			statusCode    int
			response      easyjson.Marshaler
			responseBytes []byte
		)

		statusCode, response, err = handler(ctx, server)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, []byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}

		responseBytes, err = easyjson.Marshal(response)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, []byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
		}

		ctx.JSON(statusCode, responseBytes)
	})
}
