// Package rpc содержит RPC-сервер.
package rpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/xantinium/metrix/internal/presentation/rpc/gen"
	"github.com/xantinium/metrix/internal/presentation/rpc/servers/metrics"
	metricsRepo "github.com/xantinium/metrix/internal/repository/metrics"
)

// New создаёт новый RPC-сервер.
func New(addr string, repo *metricsRepo.MetricsRepository) *Server {
	server := grpc.NewServer()

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
