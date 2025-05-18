// Пакет server содержит реализацю HTTP-сервера, использующий
// http.ServeMux для обработки HTTP-запросов.
package server

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xantinium/metrix/internal/logger"
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

// MetrixServerBuilder билдер для создания сервера метрик.
type MetrixServerBuilder struct {
	addr               string
	privateKey         string
	storeInterval      time.Duration
	isProfilingEnabled bool
	dbChecker          metrics.DatabaseChecker
	storage            metrics.MetricsStorage
}

// NewMetrixServerBuilder создаёт новый билдер сервера метрик.
func NewMetrixServerBuilder() *MetrixServerBuilder {
	return &MetrixServerBuilder{}
}

// SetAddr устанавливает адрес сервера метрик.
func (b *MetrixServerBuilder) SetAddr(addr string) *MetrixServerBuilder {
	b.addr = addr
	return b
}

// SetPrivateKey устанавливает приватный ключ,
// используемый в алгоритмах хеширования.
func (b *MetrixServerBuilder) SetPrivateKey(key string) *MetrixServerBuilder {
	b.privateKey = key
	return b
}

// SetStoreInterval устанавливает интервал между
// сохранениями метрик.
func (b *MetrixServerBuilder) SetStoreInterval(interval time.Duration) *MetrixServerBuilder {
	b.storeInterval = interval
	return b
}

// EnabledProfiling активирует профилирование.
func (b *MetrixServerBuilder) EnabledProfiling() *MetrixServerBuilder {
	b.isProfilingEnabled = true
	return b
}

// SetDatabaseChecker устанавливает сущность для проверки
// состояния базы данных.
func (b *MetrixServerBuilder) SetDatabaseChecker(checker metrics.DatabaseChecker) *MetrixServerBuilder {
	b.dbChecker = checker
	return b
}

// SetStorage устанавливает базу данных для хранения метрик.
func (b *MetrixServerBuilder) SetStorage(storage metrics.MetricsStorage) *MetrixServerBuilder {
	b.storage = storage
	return b
}

func (b *MetrixServerBuilder) Build() *MetrixServer {
	router := gin.New()
	applyMiddlewares(router, b.privateKey)

	internalServer := &internalMetrixServer{
		router: router,
		metricsRepo: metrics.NewMetricsRepository(metrics.MetricsRepositoryOptions{
			Storage:     b.storage,
			SyncMetrics: b.storeInterval == 0,
			DBChecker:   b.dbChecker,
		}),
	}

	handlers.RegisterHTMLHandler(internalServer, "/", handlers.GetAllMetricHandler)
	handlers.RegisterHandler(internalServer, http.MethodGet, "/value/:type/:id", handlers.GetMetricHandler)
	handlers.RegisterHandler(internalServer, http.MethodPost, "/update/:type/:id/:value", handlers.UpdateMetricHandler)
	handlers.RegisterHandler(internalServer, http.MethodGet, "/ping", handlers.PingHandler)
	handlers.RegisterV2Handler(internalServer, http.MethodPost, "/value/", v2handlers.GetMetricHandler)
	handlers.RegisterV2Handler(internalServer, http.MethodPost, "/update/", v2handlers.UpdateMetricHandler)
	handlers.RegisterV2Handler(internalServer, http.MethodPost, "/updates/", v2handlers.UpdateMetricsHandler)

	return &MetrixServer{
		server: &http.Server{
			Addr:    b.addr,
			Handler: router,
		},
		internalServer:     internalServer,
		worker:             NewMetrixServerWorker(b.storeInterval, b.storage),
		isProfilingEnabled: b.isProfilingEnabled,
	}
}

// MetrixServer структура, описывающая сервер метрик.
type MetrixServer struct {
	server             *http.Server
	internalServer     *internalMetrixServer
	worker             *MetrixServerWorker
	isProfilingEnabled bool
}

// Run запускает сервер метрик.
func (s *MetrixServer) Run() chan error {
	if s.isProfilingEnabled {
		s.runProfilingServer()
	}

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

func (s *MetrixServer) runProfilingServer() {
	go func() {
		err := http.ListenAndServe(":9090", nil)
		if err != nil {
			logger.Errorf("failed to start pprof server: %v", err)
		}
	}()
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
