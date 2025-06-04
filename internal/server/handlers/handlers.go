// Package handlers содержит хендлеры всех HTTP-запросов.
package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"

	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/server/interfaces"
	"github.com/xantinium/metrix/internal/tools"
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
			logger.Error(
				"error response",
				logger.Field{Name: "status", Value: statusCode},
				logger.Field{Name: "error", Value: err.Error()},
			)
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
		ctx.Writer.Header().Set(tools.ContentType, "text/html")

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
			logger.Error(
				"error response",
				logger.Field{Name: "status", Value: statusCode},
				logger.Field{Name: "error", Value: err.Error()},
			)
			writeJSON(ctx, statusCode, []byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}

		responseBytes, err = easyjson.Marshal(response)
		if err != nil {
			writeJSON(ctx, http.StatusInternalServerError, []byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
		}

		writeJSON(ctx, statusCode, responseBytes)
	})
}

func writeJSON(ctx *gin.Context, statusCode int, json []byte) {
	ctx.Header(tools.ContentType, "application/json; charset=utf-8")
	ctx.String(statusCode, string(json))
}
