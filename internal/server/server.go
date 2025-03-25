// Пакет server содержит реализацю HTTP-сервера, использующий
// http.ServeMux для обработки HTTP-запросов.
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	Addr          string
	PrivateKey    string
	StoreInterval time.Duration
	Storage       metrics.MetricsStorage
	DBChecker     metrics.DatabaseChecker
}

// NewMetrixServer создаёт новый сервер метрик.
func NewMetrixServer(opts MetrixServerOptions) *MetrixServer {
	router := gin.New()
	applyMiddlewares(router, opts.PrivateKey)

	internalServer := &internalMetrixServer{
		router: router,
		metricsRepo: metrics.NewMetricsRepository(metrics.MetricsRepositoryOptions{
			Storage:     opts.Storage,
			SyncMetrics: opts.StoreInterval == 0,
			DBChecker:   opts.DBChecker,
		}),
	}

	handlers.RegisterHTMLHandler(internalServer, "/", handlers.GetAllMetricHandler)
	handlers.RegisterHandler(internalServer, http.MethodGet, "/value/:type/:id", handlers.GetMetricHandler)
	handlers.RegisterHandler(internalServer, http.MethodPost, "/update/:type/:id/:value", handlers.UpdateMetricHandler)
	handlers.RegisterHandler(internalServer, http.MethodGet, "/ping", handlers.PingHandler)
	handlers.RegisterV2Handler(internalServer, http.MethodPost, "/value", v2handlers.GetMetricHandler)
	handlers.RegisterV2Handler(internalServer, http.MethodPost, "/update", v2handlers.UpdateMetricHandler)
	handlers.RegisterV2Handler(internalServer, http.MethodPost, "/updates", v2handlers.UpdateMetricsHandler)

	return &MetrixServer{
		server: &http.Server{
			Addr:    opts.Addr,
			Handler: router,
		},
		internalServer: internalServer,
		worker:         newMetrixServerWorker(opts.StoreInterval, opts.Storage),
	}
}

// MetrixServer структура, описывающая сервер метрик.
type MetrixServer struct {
	server         *http.Server
	internalServer *internalMetrixServer
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
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func applyMiddlewares(router *gin.Engine, privateKey string) {
	mw := []gin.HandlerFunc{gin.Recovery()}

	if privateKey != "" {
		mw = append(mw, middlewares.HashCheckMiddleware(privateKey))
	}
	mw = append(mw, middlewares.CompressMiddleware())
	if privateKey != "" {
		mw = append(mw, middlewares.ResponseHasherMiddleware(privateKey))
	}
	mw = append(mw, middlewares.LoggerMiddleware())

	router.Use(mw...)
}
