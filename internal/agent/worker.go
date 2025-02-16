package agent

import (
	"time"

	"github.com/xantinium/metrix/internal/logger"
)

type uploadFuncT = func()

// newMetrixAgentWorker создаёт новый воркер для агента метрик.
//
// reportInterval - интервал между запросами на выгрузку метрик (сек).
func newMetrixAgentWorker(reportInterval int, uploadFunc uploadFuncT) *metrixAgentWorker {
	return &metrixAgentWorker{
		stopChan:       make(chan struct{}, 1),
		reportInterval: time.Duration(reportInterval) * time.Second,
		uploadFunc:     uploadFunc,
	}
}

// metrixAgentWorker структура, описывающая воркер
// для периодической выгрузки метрик на сервер.
type metrixAgentWorker struct {
	stopChan       chan struct{}
	reportInterval time.Duration
	uploadFunc     uploadFuncT
}

// Log логирует события воркера.
func (worker *metrixAgentWorker) Log(msg string) {
	logger.Info(
		msg,
		logger.Field{
			Name:  "entity",
			Value: "worker",
		},
	)
}

// Run запускает воркер.
func (worker *metrixAgentWorker) Run() {
	t := time.NewTicker(worker.reportInterval)

	go func() {
		for {
			select {
			case <-worker.stopChan:
				worker.Log("stopping...")
				t.Stop()
				return
			case <-t.C:
				worker.uploadFunc()
			}
		}
	}()
}

// Stop прекращает работу воркера.
func (worker *metrixAgentWorker) Stop() {
	worker.stopChan <- struct{}{}
}
