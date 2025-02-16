package middlewares

import (
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

		logger.Info(
			msg,
			logger.Field{
				Name:  "status",
				Value: ctx.Writer.Status(),
			},
			logger.Field{
				Name:  "method",
				Value: ctx.Request.Method,
			},
			logger.Field{
				Name:  "url",
				Value: ctx.Request.URL.RawPath,
			},
			logger.Field{
				Name:  "duration",
				Value: duration,
			},
			logger.Field{
				Name:  "size",
				Value: ctx.Writer.Size(),
			},
		)
	}
}
