package middlewares

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/metrix/internal/logger"
)

// LoggerMiddleware мидлварь для логирования запросов.
func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
				Value: ctx.Request.URL.RawPath,
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

		reqBody, err := io.ReadAll(ctx.Copy().Request.Body)
		if err == nil {
			fields = append(fields, logger.Field{
				Name:  "req",
				Value: string(reqBody),
			})
		}

		logger.Info(msg, fields...)
	}
}
