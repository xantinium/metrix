// Пакет server содержит реализацю HTTP-сервера, использующий
// http.ServeMux для обработки HTTP-запросов.
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xantinium/metrix/internal/infrastructure/memstorage"
	"github.com/xantinium/metrix/internal/repository/metrics"
	"github.com/xantinium/metrix/internal/server/handlers"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// internalMetrixServer внутренняя структура сервера.
// Является реализацией интерфейса сервера, получаемого хендлерами.
type internalMetrixServer struct {
	router      *gin.Engine
	metricsRepo *metrics.MetricsRepository
}

// GetInternalRouter возвращает используемый роутер.
func (server *internalMetrixServer) GetInternalRouter() *gin.Engine {
	return server.router
}

// GetMetricsRepo возвращает репозиторий метрик.
func (server *internalMetrixServer) GetMetricsRepo() *metrics.MetricsRepository {
	return server.metricsRepo
}

// NewMetrixServer создаёт новый сервер метрик.
func NewMetrixServer(port int) *MetrixServer {
	metricsStorage := memstorage.NewMemStorage()

	router := gin.New()
	router.Use(gin.Recovery())

	internalServer := &internalMetrixServer{
		router:      router,
		metricsRepo: metrics.NewMetricsRepository(metricsStorage),
	}

	handlers.RegisterHTMLHandler(internalServer, "/", handlers.GetAllMetricHandler)
	handlers.RegisterHandler(internalServer, handlers.MethodGet, "/value/:type/:name", handlers.GetMetricHandler)
	handlers.RegisterHandler(internalServer, handlers.MethodPost, "/update/:type/:name/:value", handlers.UpdateMetricHandler)

	return &MetrixServer{
		port: port,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
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
