// Package server содержит реализацю сервера метрик.
//
// Поддерживается два вида серверов:
//   - REST;
//   - RPC.
package server

import (
	"context"
	"net"
	_ "net/http/pprof" // Используется для корректной работы профилировщика.
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/xantinium/metrix/internal/presentation/rest"
	"github.com/xantinium/metrix/internal/presentation/rpc"
	"github.com/xantinium/metrix/internal/repository/metrics"
	"github.com/xantinium/metrix/internal/tools"
)

// MetrixServerBuilder билдер для создания сервера метрик.
type MetrixServerBuilder struct {
	trustedSubnet *net.IPNet
	dbChecker     metrics.DatabaseChecker
	storage       metrics.MetricsStorage

	addr               string
	rpcAddr            string
	privateKey         string
	cryptoPrivateKey   string
	storeInterval      time.Duration
	isProfilingEnabled bool
	isRPCEnabled       bool
}

// NewMetrixServerBuilder создаёт новый билдер сервера метрик.
func NewMetrixServerBuilder() *MetrixServerBuilder {
	return &MetrixServerBuilder{}
}

// SetAddr устанавливает адрес REST-сервера.
func (b *MetrixServerBuilder) SetAddr(addr string) *MetrixServerBuilder {
	b.addr = addr
	return b
}

// SetRPCAddr устанавливает адрес RPC-сервера.
func (b *MetrixServerBuilder) SetRPCAddr(addr string) *MetrixServerBuilder {
	b.rpcAddr = addr
	return b
}

// SetPrivateKey устанавливает приватный ключ,
// используемый в алгоритмах хеширования.
func (b *MetrixServerBuilder) SetPrivateKey(key string) *MetrixServerBuilder {
	b.privateKey = key
	return b
}

// SetCryptoPrivateKey устанавливает приватный ключ,
// используемый в алгоритмах шифрования.
func (b *MetrixServerBuilder) SetCryptoPrivateKey(key string) *MetrixServerBuilder {
	b.cryptoPrivateKey = key
	return b
}

// SetStoreInterval устанавливает интервал между
// сохранениями метрик.
func (b *MetrixServerBuilder) SetStoreInterval(interval time.Duration) *MetrixServerBuilder {
	b.storeInterval = interval
	return b
}

// SetTrustedSubnet устанавливает доверенную подсеть.
// Используется на проверке входящих HTTP-запросов.
func (b *MetrixServerBuilder) SetTrustedSubnet(subnet *net.IPNet) *MetrixServerBuilder {
	b.trustedSubnet = subnet
	return b
}

// SetEnableRPC активирует RPC-сервер.
func (b *MetrixServerBuilder) SetEnableRPC(enable bool) *MetrixServerBuilder {
	b.isRPCEnabled = enable
	return b
}

// EnabledProfiling активирует профилирование.
func (b *MetrixServerBuilder) EnabledProfiling() *MetrixServerBuilder {
	b.isProfilingEnabled = true
	return b
}

// SetStorage устанавливает базу данных для хранения метрик.
// Также, устанавливает сущность для проверки соединения с БД.
func (b *MetrixServerBuilder) SetStorage(storage metrics.MetricsStorage, checker metrics.DatabaseChecker) *MetrixServerBuilder {
	b.storage = storage
	b.dbChecker = checker
	return b
}

// Build завершает создание сервера.
// Возвращает настроенный экземпляр.
func (b *MetrixServerBuilder) Build() *MetrixServer {
	repo := metrics.NewMetricsRepository(metrics.MetricsRepositoryOptions{
		Storage:     b.storage,
		DBChecker:   b.dbChecker,
		SyncMetrics: b.storeInterval == 0,
	})

	return &MetrixServer{
		restServer:         rest.New(b.addr, b.privateKey, b.cryptoPrivateKey, b.trustedSubnet, repo),
		rpcServer:          rpc.New(b.rpcAddr, repo),
		worker:             NewMetrixServerWorker(b.storeInterval, b.storage),
		isProfilingEnabled: b.isProfilingEnabled,
		isRPCEnabled:       b.isRPCEnabled,
	}
}

// MetrixServer структура, описывающая сервер метрик.
type MetrixServer struct {
	restServer         *rest.Server
	rpcServer          *rpc.Server
	worker             *MetrixServerWorker
	isProfilingEnabled bool
	isRPCEnabled       bool
}

// Run запускает сервер метрик.
func (s *MetrixServer) Run() chan error {
	if s.isProfilingEnabled {
		tools.RunProfilingServer()
	}

	errChan := make(chan error, 1)

	go func() {
		errChan <- s.restServer.Run()
	}()

	if s.isRPCEnabled {
		go func() {
			errChan <- s.rpcServer.Run()
		}()
	}

	s.worker.Run()

	return errChan
}

// Stop останавливает сервер метрик.
func (s *MetrixServer) Stop(timeout time.Duration) error {
	defer func() {
		s.worker.Stop()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var group *errgroup.Group
	group, ctx = errgroup.WithContext(ctx)

	group.Go(func() error {
		return s.restServer.Stop(ctx)
	})

	group.Go(func() error {
		return s.rpcServer.Stop(ctx)
	})

	return group.Wait()
}
