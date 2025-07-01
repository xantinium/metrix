// Package rpc содержит RPC-сервер.
package rpc

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/presentation/rpc/gen"
	"github.com/xantinium/metrix/internal/presentation/rpc/interceptors"
	"github.com/xantinium/metrix/internal/presentation/rpc/servers/metrics"
	metricsRepo "github.com/xantinium/metrix/internal/repository/metrics"
)

// New создаёт новый RPC-сервер.
func New(addr string, trustedSubnet *net.IPNet, repo *metricsRepo.MetricsRepository) *Server {
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(prepareInterceptors(trustedSubnet)...))

	gen.RegisterMetricsServer(server, metrics.New(repo))

	return &Server{
		addr:   addr,
		server: server,
	}
}

// Server структура, описывающая RPC-сервер.
type Server struct {
	addr   string
	server *grpc.Server
}

// Run запускает RPC-сервер.
func (s *Server) Run() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	return s.server.Serve(l)
}

// Stop останавливает REST-сервер.
func (s *Server) Stop(ctx context.Context) error {
	stopChan := make(chan struct{})
	defer close(stopChan)

	go func() {
		s.server.GracefulStop()
		stopChan <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to stop")
	case <-stopChan:
		return nil
	}
}

func prepareInterceptors(trustedSubnet *net.IPNet) []grpc.UnaryServerInterceptor {
	result := []grpc.UnaryServerInterceptor{recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
		logger.Error(
			"recovered from panic",
			logger.Field{
				Name:  "panic",
				Value: p,
			},
			logger.Field{
				Name:  "stack",
				Value: debug.Stack(),
			},
		)
		return status.Errorf(codes.Internal, "%s", p)
	}))}

	if trustedSubnet != nil {
		result = append(result, interceptors.NetGuardInterceptor(trustedSubnet))
	}

	// TODO: реализовать перехватчики:
	// 1) HashCheckInterceptor
	// 2) DecryptInterceptor

	result = append(result, interceptors.LoggerInterceptor())

	return result
}
