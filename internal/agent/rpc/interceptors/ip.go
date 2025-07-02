package interceptors

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/xantinium/metrix/internal/tools"
)

// XRealIPInterceptor перехватчик для подстановки
// IP-адреса клиента.
func XRealIPInterceptor(addr net.IP) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, tools.HeaderXRealIP, addr.String())

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
