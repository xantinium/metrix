package agent

import (
	"context"
	"time"

	"github.com/xantinium/metrix/internal/logger"
)

type uploadFuncT = func()

// newMetrixAgentWorker создаёт новый воркер для агента метрик.
//
// reportInterval  - интервал между запросами на выгрузку метрик (сек).
// reportRateLimit - количество одновременных запросов.
func newMetrixAgentWorker(reportInterval time.Duration, reportRateLimit int, uploadFunc uploadFuncT) *metrixAgentWorker {
	return &metrixAgentWorker{
		reportInterval:  reportInterval,
		reportRateLimit: reportRateLimit,
		uploadFunc:      uploadFunc,
	}
}

// metrixAgentWorker структура, описывающая воркер
// для периодической выгрузки метрик на сервер.
type metrixAgentWorker struct {
	reportInterval  time.Duration
	reportRateLimit int
	uploadFunc      uploadFuncT
}

// Log логирует события воркера.
func (worker *metrixAgentWorker) Log(msg string) {
	logger.Info(
		msg,
		logger.Field{
			Name:  "entity",
			Value: "agent-worker",
		},
	)
}

// Run запускает воркер.
func (worker *metrixAgentWorker) Run(ctx context.Context) {
	t := time.NewTicker(worker.reportInterval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				worker.Log("stopping...")
				t.Stop()
				return
			case <-t.C:
				worker.uploadFunc()
			}
		}
	}()
}
