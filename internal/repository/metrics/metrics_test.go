package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xantinium/metrix/internal/infrastructure/metricsstorage"
	"github.com/xantinium/metrix/internal/models"
)

func TestMetricsRepository_UpdateGaugeMetric(t *testing.T) {
	storage, err := metricsstorage.NewMetricsStorage("metrix.db", false)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewMetricsRepository(storage, false)

	updateOperations := []struct {
		metricType  models.MetricType
		metricName  string
		metricValue float64
	}{
		{
			metricType:  models.Gauge,
			metricName:  "Alloc",
			metricValue: 123.45,
		},
		{
			metricType:  models.Counter,
			metricName:  "PollCount",
			metricValue: 222,
		},
		{
			metricType:  models.Gauge,
			metricName:  "RandomValue",
			metricValue: 78,
		},
		{
			metricType:  models.Gauge,
			metricName:  "Alloc",
			metricValue: 2.1,
		},
		{
			metricType:  models.Counter,
			metricName:  "PollCount",
			metricValue: 78,
		},
	}

	for _, oper := range updateOperations {
		var err error

		switch oper.metricType {
		case models.Gauge:
			_, err = repo.UpdateGaugeMetric(oper.metricName, oper.metricValue)
		case models.Counter:
			_, err = repo.UpdateCounterMetric(oper.metricName, int64(oper.metricValue))
		}

		require.NoError(t, err)
	}

	var (
		pollCount               int64
		allocValue, randomValue float64
	)

	allocValue, err = repo.GetGaugeMetric("Alloc")
	require.NoError(t, err)
	require.Equal(t, 2.1, allocValue)

	randomValue, err = repo.GetGaugeMetric("RandomValue")
	require.NoError(t, err)
	require.Equal(t, 78.0, randomValue)

	pollCount, err = repo.GetCounterMetric("PollCount")
	require.NoError(t, err)
	require.Equal(t, int64(300), pollCount)
}
