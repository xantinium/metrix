// Package rest содержит реализацю HTTP-сервера, использующий
// http.ServeMux для обработки HTTP-запросов.
package rest

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/metrix/internal/presentation/rest/handlers"
	v2handlers "github.com/xantinium/metrix/internal/presentation/rest/handlers/v2"
	"github.com/xantinium/metrix/internal/presentation/rest/middlewares"
	"github.com/xantinium/metrix/internal/repository/metrics"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// internalServer внутренняя структура сервера.
// Является реализацией интерфейса сервера, получаемого хендлерами.
type internalServer struct {
	router      *gin.Engine
	metricsRepo *metrics.MetricsRepository
}

// GetInternalRouter возвращает используемый роутер.
func (server *internalServer) GetInternalRouter() *gin.Engine {
	return server.router
}

// GetMetricsRepo возвращает репозиторий метрик.
func (server *internalServer) GetMetricsRepo() *metrics.MetricsRepository {
	return server.metricsRepo
}

// New создаёт новый REST-сервер.
func New(addr, privateKey, cryptoPrivateKey string, trustedSubnet *net.IPNet, repo *metrics.MetricsRepository) *Server {
	router := gin.New()
	applyMiddlewares(router, privateKey, cryptoPrivateKey, trustedSubnet)

	internalServer := &internalServer{
		router:      router,
		metricsRepo: repo,
	}

	handlers.RegisterHTMLHandler(internalServer, "/", handlers.GetAllMetricHandler)
	handlers.RegisterHandler(internalServer, http.MethodGet, "/value/:type/:id", handlers.GetMetricHandler)
	handlers.RegisterHandler(internalServer, http.MethodPost, "/update/:type/:id/:value", handlers.UpdateMetricHandler)
	handlers.RegisterHandler(internalServer, http.MethodGet, "/ping", handlers.PingHandler)
	handlers.RegisterV2Handler(internalServer, http.MethodPost, "/value/", v2handlers.GetMetricHandler)
	handlers.RegisterV2Handler(internalServer, http.MethodPost, "/update/", v2handlers.UpdateMetricHandler)
	handlers.RegisterV2Handler(internalServer, http.MethodPost, "/updates/", v2handlers.UpdateMetricsHandler)

	return &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		internalServer: internalServer,
	}
}

// Server структура, описывающая REST-сервер.
type Server struct {
	server         *http.Server
	internalServer *internalServer
}

// Run запускает REST-сервер.
func (s *Server) Run() error {
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Stop останавливает REST-сервер.
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func applyMiddlewares(router *gin.Engine, privateKey, cryptoPrivateKey string, trustedSubnet *net.IPNet) {
	mw := []gin.HandlerFunc{gin.Recovery()}

	if trustedSubnet != nil {
		mw = append(mw, middlewares.NetGuardMiddleware(trustedSubnet))
	}
	if privateKey != "" {
		mw = append(mw, middlewares.HashCheckMiddleware(privateKey))
	}
	if cryptoPrivateKey != "" {
		mw = append(mw, middlewares.DecryptMiddleware(cryptoPrivateKey))
	}
	mw = append(mw, middlewares.CompressMiddleware())
	if privateKey != "" {
		mw = append(mw, middlewares.ResponseHasherMiddleware(privateKey))
	}
	mw = append(mw, middlewares.LoggerMiddleware())

	router.Use(mw...)
}
