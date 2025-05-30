package agent

import (
	"context"
	"time"

	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/tools"
)

type uploadFuncT = func()

// MetrixAgentWorkerPoolOptions параметры для пула воркеров.
type MetrixAgentWorkerPoolOptions struct {
	PoolSize        int
	ReportInterval  time.Duration // интервал между запросами на выгрузку метрик (сек).
	ReportRateLimit int           // количество одновременных запросов.
	UploadFunc      uploadFuncT
}

// NewMetrixAgentWorkerPool создаёт новый пул воркеров для агента метрик.
func NewMetrixAgentWorkerPool(opts MetrixAgentWorkerPoolOptions) *MetrixAgentWorkerPool {
	return &MetrixAgentWorkerPool{
		sm:             tools.NewSemaphore(opts.ReportRateLimit),
		poolSize:       opts.PoolSize,
		reportInterval: opts.ReportInterval,
		uploadFunc:     opts.UploadFunc,
	}
}

// MetrixAgentWorkerPool структура, описывающая пул воркеров
// для периодической выгрузки метрик на сервер.
type MetrixAgentWorkerPool struct {
	sm             *tools.Semaphore
	poolSize       int
	reportInterval time.Duration
	uploadFunc     uploadFuncT
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
				pool.sm.Acquire()
				pool.Log(logger.InfoLevel, "uploading metrics...")
				pool.uploadFunc()
				pool.sm.Release()
				t.Reset(pool.reportInterval)
			}
		}
	}()
}
