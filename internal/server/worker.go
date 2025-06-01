package server

import (
	"context"
	"time"

	"github.com/xantinium/metrix/internal/logger"
)

type MetricsSaver interface {
	SaveMetrics(ctx context.Context) error
}

// NewMetrixServerWorker создаёт новый воркер для сервера метрик.
//
// storeInterval - интервал между сохранениями метрик (сек).
func NewMetrixServerWorker(storeInterval time.Duration, metricsSaver MetricsSaver) *MetrixServerWorker {
	return &MetrixServerWorker{
		stopFunc:      func() {},
		storeInterval: storeInterval,
		metricsSaver:  metricsSaver,
	}
}

// MetrixServerWorker структура, описывающая воркер
// для периодического сохранения метрик.
type MetrixServerWorker struct {
	metricsSaver  MetricsSaver
	stopFunc      context.CancelFunc
	storeInterval time.Duration
}

// Run запускает воркер.
func (worker *MetrixServerWorker) Run() {
	var ctx context.Context
	ctx, worker.stopFunc = context.WithCancel(context.TODO())

	t := time.NewTicker(worker.storeInterval)

	// Периодическая запись работает только при ненулевом storeInterval.
	if worker.storeInterval != 0 {
		go func() {
			for {
				select {
				case <-ctx.Done():
					worker.log("stopping...")
					t.Stop()
					return
				case <-t.C:
					err := worker.metricsSaver.SaveMetrics(ctx)
					if err != nil {
						worker.log("failed to save metrics")
					}
				}
			}
		}()
	}
}

// Stop прекращает работу воркера.
func (worker *MetrixServerWorker) Stop() {
	worker.stopFunc()
}

// log логирует события воркера.
func (worker *MetrixServerWorker) log(msg string) {
	logger.Info(
		msg,
		logger.Field{
			Name:  "entity",
			Value: "server-worker",
		},
	)
}
