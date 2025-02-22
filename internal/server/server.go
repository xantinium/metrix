// Пакет server содержит реализацю HTTP-сервера, использующий
// http.ServeMux для обработки HTTP-запросов.
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xantinium/metrix/internal/infrastructure/metricsstorage"
	"github.com/xantinium/metrix/internal/repository/metrics"
	"github.com/xantinium/metrix/internal/server/handlers"
	v2handlers "github.com/xantinium/metrix/internal/server/handlers/v2"
	"github.com/xantinium/metrix/internal/server/middlewares"
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

type MetrixServerOptions struct {
	Addr           string
	StoragePath    string
	StoreInterval  time.Duration
	RestoreMetrics bool
}

// NewMetrixServer создаёт новый сервер метрик.
func NewMetrixServer(opts MetrixServerOptions) *MetrixServer {
	metricsStorage, err := metricsstorage.NewMetricsStorage(opts.StoragePath, opts.RestoreMetrics)
	if err != nil {
		panic(err)
	}

	router := gin.New()
	router.Use(gin.Recovery(), middlewares.CompressMiddleware(), middlewares.LoggerMiddleware())

	internalServer := &internalMetrixServer{
		router:      router,
		metricsRepo: metrics.NewMetricsRepository(metricsStorage),
	}

	handlers.RegisterHTMLHandler(internalServer, "/", handlers.GetAllMetricHandler)
	handlers.RegisterHandler(internalServer, http.MethodGet, "/value/:type/:name", handlers.GetMetricHandler)
	handlers.RegisterHandler(internalServer, http.MethodPost, "/update/:type/:name/:value", handlers.UpdateMetricHandler)
	handlers.RegisterV2Handler(internalServer, http.MethodPost, "/value", v2handlers.GetMetricHandler)
	handlers.RegisterV2Handler(internalServer, http.MethodPost, "/update", v2handlers.UpdateMetricHandler)

	return &MetrixServer{
		server: &http.Server{
			Addr:    opts.Addr,
			Handler: router,
		},
		internalServer: internalServer,
		metricsStorage: metricsStorage,
		worker:         newMetrixServerWorker(opts.StoreInterval, metricsStorage),
	}
}

// MetrixServer структура, описывающая сервер метрик.
type MetrixServer struct {
	server         *http.Server
	internalServer *internalMetrixServer
	metricsStorage *metricsstorage.MetricsStorage
	worker         *metrixServerWorker
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

	s.worker.Run()

	return errChan
}

// Stop останавливает сервер метрик.
func (s *MetrixServer) Stop() error {
	defer func() {
		s.worker.Stop()
		s.metricsStorage.Destroy()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return s.server.Shutdown(ctx)
}
