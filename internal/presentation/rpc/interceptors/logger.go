package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/xantinium/metrix/internal/logger"
)

// LoggerInterceptor перехватчик для логирования.
func LoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)

		duration := time.Since(start)
		msg := "rpc request"

		fields := []logger.Field{
			{
				Name:  "method",
				Value: info.FullMethod,
			},
			{
				Name:  "duration",
				Value: duration,
			},
		}

		if err != nil {
			fields = append(fields, logger.Field{
				Name:  "error",
				Value: err.Error(),
			})
		}

		logger.Info(msg, fields...)

		return resp, err
	}
}
