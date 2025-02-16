package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"

	"github.com/xantinium/metrix/internal/logger"
)

// LoggerMiddleware мидлварь для логирования запросов.
func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)
		msg := fmt.Sprintf("%d %s %s", ctx.Writer.Status(), ctx.Request.Method, ctx.Request.URL.RawPath)

		logger.Info(
			msg,
			zapcore.Field{
				Key:       "duration",
				Type:      zapcore.DurationType,
				Interface: duration,
			},
			zapcore.Field{
				Key:     "size",
				Integer: int64(ctx.Writer.Size()),
			},
		)
	}
}
