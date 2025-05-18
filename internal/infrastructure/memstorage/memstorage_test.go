package memstorage_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xantinium/metrix/internal/infrastructure/memstorage"
	"github.com/xantinium/metrix/internal/models"
)

func TestMemStorage(t *testing.T) {
	ctx := context.Background()

	storage, err := memstorage.NewMemStorage("metrix.db", false)
	if err != nil {
		t.Fatal(err)
	}

	// С помощью горутин проверяем потокобезопасность.
	for range 5 {
		// Запускаем 5 горутин, каждая устанавливает значение 5.
		// В итоге ожидаем получить: 5.
		go func() {
			for range 20 {
				storage.UpdateGaugeMetric(ctx, "Alloc", 5)
			}
		}()

		// Запускаем 5 горутин, каждая увеличивает значение на 2.
		// В итоге ожидаем получить: 5 * 20 * 2.
		go func() {
			for range 20 {
				storage.UpdateCounterMetric(ctx, "PollCount", 2)
			}
		}()
	}

	time.Sleep(time.Second)

	var (
		gaugeMetric   float64
		counterMetric int64
		metrics       []models.MetricInfo
	)

	gaugeMetric, err = storage.GetGaugeMetric(ctx, "Alloc")
	if err != nil {
		t.Fatal(err)
	}

	counterMetric, err = storage.GetCounterMetric(ctx, "PollCount")
	if err != nil {
		t.Fatal(err)
	}

	metrics, err = storage.GetAllMetrics(ctx)
	if err != nil {
		t.Fatal(err)
	}

	require.Len(t, metrics, 2)
	require.Equal(t, float64(5), gaugeMetric)
	require.Equal(t, int64(200), counterMetric)
}
