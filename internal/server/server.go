// Пакет server содержит реализацю HTTP-сервера, использующий
// http.ServeMux для обработки HTTP-запросов.
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/xantinium/metrix/internal/infrastructure/memstorage"
	"github.com/xantinium/metrix/internal/repository/metrics"
	"github.com/xantinium/metrix/internal/server/handlers"
)

// internalMetrixServer внутренняя структура сервера.
// Является реализацией интерфейса сервера, получаемого хендлерами.
type internalMetrixServer struct {
	mux         *http.ServeMux
	metricsRepo *metrics.MetricsRepository
}

// GetInternalMux возвращает используемый http.ServeMux.
func (server *internalMetrixServer) GetInternalMux() *http.ServeMux {
	return server.mux
}

// GetMetricsRepo возвращает репозиторий метрик.
func (server *internalMetrixServer) GetMetricsRepo() *metrics.MetricsRepository {
	return server.metricsRepo
}

// NewMetrixServer создаёт новый сервер метрик.
func NewMetrixServer(port int) *MetrixServer {
	metricsStorage := memstorage.NewMemStorage()

	internalServer := &internalMetrixServer{
		mux:         http.NewServeMux(),
		metricsRepo: metrics.NewMetricsRepository(metricsStorage),
	}

	handlers.RegisterHandler(internalServer, handlers.METHOD_POST, "/update/:type/:name/:value/", handlers.UpdateMetric)

	return &MetrixServer{
		port: port,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: internalServer.mux,
		},
		internalServer: internalServer,
	}
}

// MetrixServer структура, описывающая сервер метрик.
type MetrixServer struct {
	port           int
	server         *http.Server
	internalServer *internalMetrixServer
}

// Run запускает сервер метрик.
func (s *MetrixServer) Run() chan error {
	errChan := make(chan error, 1)

	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	return errChan
}

// Stop останавливает сервер метрик.
func (s *MetrixServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return s.server.Shutdown(ctx)
}
