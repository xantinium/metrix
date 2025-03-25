// Пакет runtimemetrics содержит реализацию источника метрик.
// На данный момент, метрики предоставляет пакеты runtime и gopsutil.
package runtimemetrics

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"

	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/models"
)

// NewRuntimeMetricsSource создаёт новый источник метрик.
//
// pollInterval - интервал между обновлениями метрик (сек).
func NewRuntimeMetricsSource(pollInterval int) *RuntimeMetricsSource {
	return &RuntimeMetricsSource{
		pollInterval:    time.Duration(pollInterval) * time.Second,
		metricsSnapshot: map[string]models.MetricInfo{},
	}
}

// RuntimeMetricsSource структура, реализующая источник метрик.
type RuntimeMetricsSource struct {
	mx              sync.RWMutex
	pollInterval    time.Duration
	snapshotsCount  int64 // учитывает только основные метрики
	metricsSnapshot map[string]models.MetricInfo
}

// Log логирует события источника метрик.
func (source *RuntimeMetricsSource) Log(lvl logger.LogLevel, msg string) {
	field := logger.Field{
		Name:  "entity",
		Value: "runtimemetrics",
	}

	switch lvl {
	case logger.InfoLevel:
		logger.Info(msg, field)
	case logger.ErrorLevel:
		logger.Error(msg, field)
	}
}

// Run запускает обновления метрик.
func (source *RuntimeMetricsSource) Run(ctx context.Context) {
	t := time.NewTicker(source.pollInterval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				source.Log(logger.InfoLevel, "stopping main...")
				t.Stop()
				return
			case <-t.C:
				source.DoShapshot()
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				source.Log(logger.InfoLevel, "stopping additional...")
				t.Stop()
				return
			case <-t.C:
				source.DoAdditionalSnapshot(ctx)
			}
		}
	}()
}

// DoShapshot сканирует метрики и сохраняет их в памяти.
func (source *RuntimeMetricsSource) DoShapshot() {
	source.Log(logger.InfoLevel, "saving metrics snapshot...")

	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)

	source.mx.Lock()
	defer source.mx.Unlock()

	source.snapshotsCount++
	source.metricsSnapshot["Alloc"] = models.NewGaugeMetric("Alloc", float64(stats.Alloc))
	source.metricsSnapshot["BuckHashSys"] = models.NewGaugeMetric("BuckHashSys", float64(stats.BuckHashSys))
	source.metricsSnapshot["Frees"] = models.NewGaugeMetric("Frees", float64(stats.Frees))
	source.metricsSnapshot["GCCPUFraction"] = models.NewGaugeMetric("GCCPUFraction", float64(stats.GCCPUFraction))
	source.metricsSnapshot["GCSys"] = models.NewGaugeMetric("GCSys", float64(stats.GCSys))
	source.metricsSnapshot["HeapAlloc"] = models.NewGaugeMetric("HeapAlloc", float64(stats.HeapAlloc))
	source.metricsSnapshot["HeapIdle"] = models.NewGaugeMetric("HeapIdle", float64(stats.HeapIdle))
	source.metricsSnapshot["HeapInuse"] = models.NewGaugeMetric("HeapInuse", float64(stats.HeapInuse))
	source.metricsSnapshot["HeapObjects"] = models.NewGaugeMetric("HeapObjects", float64(stats.HeapObjects))
	source.metricsSnapshot["HeapReleased"] = models.NewGaugeMetric("HeapReleased", float64(stats.HeapReleased))
	source.metricsSnapshot["HeapSys"] = models.NewGaugeMetric("HeapSys", float64(stats.HeapSys))
	source.metricsSnapshot["LastGC"] = models.NewGaugeMetric("LastGC", float64(stats.LastGC))
	source.metricsSnapshot["Lookups"] = models.NewGaugeMetric("Lookups", float64(stats.Lookups))
	source.metricsSnapshot["MCacheInuse"] = models.NewGaugeMetric("MCacheInuse", float64(stats.MCacheInuse))
	source.metricsSnapshot["MCacheSys"] = models.NewGaugeMetric("MCacheSys", float64(stats.MCacheSys))
	source.metricsSnapshot["MSpanInuse"] = models.NewGaugeMetric("MSpanInuse", float64(stats.MSpanInuse))
	source.metricsSnapshot["MSpanSys"] = models.NewGaugeMetric("MSpanSys", float64(stats.MSpanSys))
	source.metricsSnapshot["Mallocs"] = models.NewGaugeMetric("Mallocs", float64(stats.Mallocs))
	source.metricsSnapshot["NextGC"] = models.NewGaugeMetric("NextGC", float64(stats.NextGC))
	source.metricsSnapshot["NumForcedGC"] = models.NewGaugeMetric("NumForcedGC", float64(stats.NumForcedGC))
	source.metricsSnapshot["NumGC"] = models.NewGaugeMetric("NumGC", float64(stats.NumGC))
	source.metricsSnapshot["OtherSys"] = models.NewGaugeMetric("OtherSys", float64(stats.OtherSys))
	source.metricsSnapshot["PauseTotalNs"] = models.NewGaugeMetric("PauseTotalNs", float64(stats.PauseTotalNs))
	source.metricsSnapshot["StackInuse"] = models.NewGaugeMetric("StackInuse", float64(stats.StackInuse))
	source.metricsSnapshot["StackSys"] = models.NewGaugeMetric("StackSys", float64(stats.StackSys))
	source.metricsSnapshot["Sys"] = models.NewGaugeMetric("Sys", float64(stats.Sys))
	source.metricsSnapshot["TotalAlloc"] = models.NewGaugeMetric("TotalAlloc", float64(stats.TotalAlloc))
	source.metricsSnapshot["PollCount"] = models.NewCounterMetric("PollCount", source.snapshotsCount)
	source.metricsSnapshot["RandomValue"] = models.NewGaugeMetric("RandomValue", randFloat())
}

// GetSnapshot возвращает сохранённые метрики.
func (source *RuntimeMetricsSource) GetSnapshot() []models.MetricInfo {
	source.mx.RLock()
	defer source.mx.RUnlock()

	i := 0
	metrics := make([]models.MetricInfo, len(source.metricsSnapshot))
	for _, metric := range source.metricsSnapshot {
		metrics[i] = metric
		i++
	}

	return metrics
}

func randFloat() float64 {
	min := 0.0
	max := math.MaxFloat64

	return min + rand.Float64()*(max-min)
}

// DoAdditionalSnapshot сканирует дополнительные метрики и сохраняет их в памяти.
// Отдельный метод нужен для отдельной горутины, требуемой по заданию.
func (source *RuntimeMetricsSource) DoAdditionalSnapshot(ctx context.Context) {
	var (
		err      error
		memStats *mem.VirtualMemoryStat
		cpuStats []float64
	)

	source.Log(logger.InfoLevel, "saving metrics snapshot...")

	memStats, err = mem.VirtualMemory()
	if err != nil {
		source.Log(logger.ErrorLevel, fmt.Sprintf("failed to get additional metrics: %v", err))
		return
	}

	cpuStats, err = cpu.PercentWithContext(ctx, time.Second, true)
	if err != nil {
		source.Log(logger.ErrorLevel, fmt.Sprintf("failed to get additional metrics: %v", err))
		return
	}

	source.mx.Lock()
	defer source.mx.Unlock()

	source.metricsSnapshot["TotalMemory"] = models.NewGaugeMetric("TotalMemory", float64(memStats.Total))
	source.metricsSnapshot["FreeMemory"] = models.NewGaugeMetric("FreeMemory", float64(memStats.Free))

	for i, coreUsage := range cpuStats {
		id := fmt.Sprintf("CPUutilization%d", i)
		source.metricsSnapshot[id] = models.NewGaugeMetric(id, coreUsage)
	}
}
