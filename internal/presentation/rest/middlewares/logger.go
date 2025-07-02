package middlewares

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/metrix/internal/logger"
)

// LoggerMiddleware мидлварь для логирования запросов.
func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			err       error
			bodyBytes []byte
		)

		if ctx.Request.Body != nil {
			bodyBytes, err = io.ReadAll(ctx.Request.Body)
			if err == nil {
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		start := time.Now()
		ctx.Next()

		duration := time.Since(start)
		msg := "api request"

		fields := []logger.Field{
			{
				Name:  "status",
				Value: ctx.Writer.Status(),
			},
			{
				Name:  "method",
				Value: ctx.Request.Method,
			},
			{
				Name:  "url",
				Value: ctx.Request.URL.String(),
			},
			{
				Name:  "duration",
				Value: duration,
			},
			{
				Name:  "size",
				Value: ctx.Writer.Size(),
			},
		}

		if err == nil {
			fields = append(fields, logger.Field{
				Name:  "req",
				Value: string(bodyBytes),
			})
		}

		logger.Info(msg, fields...)
	}
}
