package server

import (
	"context"
	"time"

	"github.com/xantinium/metrix/internal/logger"
)

type MetricsSaver interface {
	SaveMetrics() error
}

// newMetrixServerWorker создаёт новый воркер для сервера метрик.
//
// storeInterval - интервал между сохранениями метрик (сек).
func newMetrixServerWorker(storeInterval time.Duration, metricsSaver MetricsSaver) *metrixServerWorker {
	return &metrixServerWorker{
		stopFunc:      func() {},
		storeInterval: storeInterval,
		metricsSaver:  metricsSaver,
	}
}

// metrixServerWorker структура, описывающая воркер
// для периодического сохранения метрик.
type metrixServerWorker struct {
	stopFunc      context.CancelFunc
	storeInterval time.Duration
	metricsSaver  MetricsSaver
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
	var ctx context.Context
	ctx, worker.stopFunc = context.WithCancel(context.TODO())

	t := time.NewTicker(worker.storeInterval)

	// Периодическая запись работает только при ненулевом storeInterval.
	if worker.storeInterval != 0 {
		go func() {
			for {
				select {
				case <-ctx.Done():
					worker.Log("stopping...")
					t.Stop()
					return
				case <-t.C:
					err := worker.metricsSaver.SaveMetrics()
					if err != nil {
						worker.Log("failed to save metrics")
					}
				}
			}
		}()
	}
}

// Stop прекращает работу воркера.
func (worker *metrixServerWorker) Stop() {
	worker.stopFunc()
}
