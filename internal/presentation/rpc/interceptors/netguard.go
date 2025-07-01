package interceptors

import (
	"context"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/tools"
)

// NetGuardInterceptor перехватчик для проверки IP-адреса клиента,
// при помощи переданной доверенной подсети.
func NetGuardInterceptor(trustedSubnet *net.IPNet) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		realIP := metadata.ExtractIncoming(ctx).Get(tools.HeaderXRealIP)
		if realIP == "" {
			logger.Infof("%s header is missing, but required", tools.HeaderXRealIP)
			return nil, status.Errorf(codes.Unavailable, "%s header is missing, but required", tools.HeaderXRealIP)
		}

		ip := net.ParseIP(realIP)
		if ip == nil {
			logger.Infof("invalid ip address in %s header", tools.HeaderXRealIP)
			return nil, status.Errorf(codes.Unavailable, "invalid ip address in %s header", tools.HeaderXRealIP)
		}

		if !trustedSubnet.Contains(ip) {
			logger.Infof("ip address %q not in trusted subnet", ip.String())
			return nil, status.Errorf(codes.Unavailable, "ip address %q not in trusted subnet", ip.String())
		}

		return handler(ctx, req)
	}
}
