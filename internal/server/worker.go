package server

import (
	"time"

	"github.com/xantinium/metrix/internal/infrastructure/metricsstorage"
	"github.com/xantinium/metrix/internal/logger"
)

// newMetrixServerWorker создаёт новый воркер для сервера метрик.
//
// storeInterval - интервал между сохранениями метрик (сек).
func newMetrixServerWorker(storeInterval time.Duration, metricsStorage *metricsstorage.MetricsStorage) *metrixServerWorker {
	return &metrixServerWorker{
		stopChan:       make(chan struct{}, 1),
		storeInterval:  storeInterval,
		metricsStorage: metricsStorage,
	}
}

// metrixServerWorker структура, описывающая воркер
// для периодического сохранения метрик.
type metrixServerWorker struct {
	stopChan       chan struct{}
	storeInterval  time.Duration
	metricsStorage *metricsstorage.MetricsStorage
}

// Log логирует события воркера.
func (worker *metrixServerWorker) Log(msg string) {
	logger.Info(
		msg,
		logger.Field{
			Name:  "entity",
			Value: "server-worker",
		},
	)
}

// Run запускает воркер.
func (worker *metrixServerWorker) Run() {
	t := time.NewTicker(worker.storeInterval)

	go func() {
		for {
			select {
			case <-worker.stopChan:
				worker.Log("stopping...")
				t.Stop()
				return
			case <-t.C:
				err := worker.metricsStorage.SaveMetrics()
				if err != nil {
					worker.Log("failed to save metrics")
				}
			}
		}
	}()
}

// Stop прекращает работу воркера.
func (worker *metrixServerWorker) Stop() {
	worker.stopChan <- struct{}{}
}
