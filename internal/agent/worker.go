package agent

import (
	"context"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/xantinium/metrix/internal/logger"
)

type uploadFuncT = func()

// MetrixAgentWorkerPoolOptions параметры для пула воркеров.
type MetrixAgentWorkerPoolOptions struct {
	UploadFunc      uploadFuncT
	ReportInterval  time.Duration // интервал между запросами на выгрузку метрик (сек).
	PoolSize        int
	ReportRateLimit int // количество одновременных запросов.
}

// NewMetrixAgentWorkerPool создаёт новый пул воркеров для агента метрик.
func NewMetrixAgentWorkerPool(opts MetrixAgentWorkerPoolOptions) *MetrixAgentWorkerPool {
	return &MetrixAgentWorkerPool{
		sm:             semaphore.NewWeighted(int64(opts.ReportRateLimit)),
		poolSize:       opts.PoolSize,
		reportInterval: opts.ReportInterval,
		uploadFunc:     opts.UploadFunc,
	}
}

// MetrixAgentWorkerPool структура, описывающая пул воркеров
// для периодической выгрузки метрик на сервер.
type MetrixAgentWorkerPool struct {
	uploadFunc     uploadFuncT
	sm             *semaphore.Weighted
	reportInterval time.Duration
	poolSize       int
	disabled       bool
}

// Log логирует события воркеров.
func (pool *MetrixAgentWorkerPool) Log(lvl logger.LogLevel, msg string) {
	field := logger.Field{
		Name:  "entity",
		Value: "agent-worker",
	}

	switch lvl {
	case logger.InfoLevel:
		logger.Info(msg, field)
	case logger.ErrorLevel:
		logger.Error(msg, field)
	}
}

// Run запускает воркеры.
func (pool *MetrixAgentWorkerPool) Run(ctx context.Context) {
	for range pool.poolSize {
		go pool.runWorker(ctx)
	}
}

// Disable отключает воркеры, запрещая
// выполнять uploadFunc.
func (pool *MetrixAgentWorkerPool) Disable() {
	pool.disabled = true
}

func (pool *MetrixAgentWorkerPool) runWorker(ctx context.Context) {
	t := time.NewTimer(pool.reportInterval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				pool.Log(logger.InfoLevel, "stopping...")
				t.Stop()
				return
			case <-t.C:
				if !pool.disabled {
					pool.sm.Acquire(ctx, 1)
					pool.Log(logger.InfoLevel, "uploading metrics...")
					pool.uploadFunc()
					pool.sm.Release(1)
					t.Reset(pool.reportInterval)
				}
			}
		}
	}()
}
