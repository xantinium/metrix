// Пакет runtimemetrics содержит реализацию источника метрик.
// На данный момент, метрики предоставляет пакет runtime.
package runtimemetrics

import (
	"math"
	"math/rand/v2"
	"runtime"
	"sync"
	"time"

	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/models"
)

// NewRuntimeMetricsSource создаёт новый источник метрик.
//
// pollInterval - интервал между обновлениями метрик (сек).
func NewRuntimeMetricsSource(pollInterval int) *RuntimeMetricsSource {
	return &RuntimeMetricsSource{
		stopChan:     make(chan struct{}, 1),
		pollInterval: time.Duration(pollInterval) * time.Second,
	}
}

// RuntimeMetricsSource структура, реализующая источник метрик.
type RuntimeMetricsSource struct {
	stopChan        chan struct{}
	mx              sync.RWMutex
	pollInterval    time.Duration
	snapshotsCount  int64
	metricsSnapshot []models.MetricInfo
}

// Log логирует события источника метрик.
func (source *RuntimeMetricsSource) Log(msg string) {
	logger.Info(
		msg,
		logger.Field{
			Name:  "entity",
			Value: "runtimemetrics",
		},
	)
}

// Run запускает обновления метрик.
func (source *RuntimeMetricsSource) Run() {
	t := time.NewTicker(source.pollInterval)

	go func() {
		for {
			select {
			case <-source.stopChan:
				source.Log("stopping...")
				t.Stop()
				return
			case <-t.C:
				source.DoShapshot()
			}
		}
	}()
}

// Stop прекращает обновления метрик.
func (source *RuntimeMetricsSource) Stop() {
	source.stopChan <- struct{}{}
}

// DoShapshot сканирует метрики и сохраняет их в памяти.
func (source *RuntimeMetricsSource) DoShapshot() {
	source.Log("saving metrics snapshot...")

	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)

	source.mx.Lock()
	defer source.mx.Unlock()

	source.snapshotsCount++
	source.metricsSnapshot = []models.MetricInfo{
		models.NewGaugeMetric("Alloc", float64(stats.Alloc)),
		models.NewGaugeMetric("BuckHashSys", float64(stats.BuckHashSys)),
		models.NewGaugeMetric("Frees", float64(stats.Frees)),
		models.NewGaugeMetric("GCCPUFraction", stats.GCCPUFraction),
		models.NewGaugeMetric("GCSys", float64(stats.GCSys)),
		models.NewGaugeMetric("HeapAlloc", float64(stats.HeapAlloc)),
		models.NewGaugeMetric("HeapIdle", float64(stats.HeapIdle)),
		models.NewGaugeMetric("HeapInuse", float64(stats.HeapInuse)),
		models.NewGaugeMetric("HeapObjects", float64(stats.HeapObjects)),
		models.NewGaugeMetric("HeapReleased", float64(stats.HeapReleased)),
		models.NewGaugeMetric("HeapSys", float64(stats.HeapSys)),
		models.NewGaugeMetric("LastGC", float64(stats.LastGC)),
		models.NewGaugeMetric("Lookups", float64(stats.Lookups)),
		models.NewGaugeMetric("MCacheInuse", float64(stats.MCacheInuse)),
		models.NewGaugeMetric("MCacheSys", float64(stats.MCacheSys)),
		models.NewGaugeMetric("MSpanInuse", float64(stats.MSpanInuse)),
		models.NewGaugeMetric("MSpanSys", float64(stats.MSpanSys)),
		models.NewGaugeMetric("Mallocs", float64(stats.Mallocs)),
		models.NewGaugeMetric("NextGC", float64(stats.NextGC)),
		models.NewGaugeMetric("NumForcedGC", float64(stats.NumForcedGC)),
		models.NewGaugeMetric("NumGC", float64(stats.NumGC)),
		models.NewGaugeMetric("OtherSys", float64(stats.OtherSys)),
		models.NewGaugeMetric("PauseTotalNs", float64(stats.PauseTotalNs)),
		models.NewGaugeMetric("StackInuse", float64(stats.StackInuse)),
		models.NewGaugeMetric("StackSys", float64(stats.StackSys)),
		models.NewGaugeMetric("Sys", float64(stats.Sys)),
		models.NewGaugeMetric("TotalAlloc", float64(stats.TotalAlloc)),
		models.NewCounterMetric("PollCount", source.snapshotsCount),
		models.NewGaugeMetric("RandomValue", randFloat()),
	}
}

// GetSnapshot возвращает сохранённые метрики.
func (source *RuntimeMetricsSource) GetSnapshot() []models.MetricInfo {
	source.mx.RLock()
	defer source.mx.RUnlock()

	return source.metricsSnapshot
}

func randFloat() float64 {
	min := 0.0
	max := math.MaxFloat64

	return min + rand.Float64()*(max-min)
}
