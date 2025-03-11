package postgres

import (
	"context"

	"github.com/xantinium/metrix/internal/models"
)

// GetGaugeMetric возвращает метрику типа Gauge по идентификатору id.
func (client *PostgresClient) GetGaugeMetric(ctx context.Context, id string) (float64, error) {
	row := client.db.QueryRowContext(ctx, "SELECT gauge_value FROM metrics"+
		" WHERE id = $1 AND type = $2;",
		id,
		serializeMetricType(models.Gauge))

	var gaugeValue float64
	err := row.Scan(&gaugeValue)
	if err != nil {
		return -1, convertError(err)
	}

	return gaugeValue, nil
}

// GetCounterMetric возвращает метрику типа Counter по идентификатору id.
func (client *PostgresClient) GetCounterMetric(ctx context.Context, id string) (int64, error) {
	row := client.db.QueryRowContext(ctx, "SELECT counter_value FROM metrics"+
		" WHERE id = $1 AND type = $2;",
		id,
		serializeMetricType(models.Gauge))

	var counterValue int64
	err := row.Scan(&counterValue)
	if err != nil {
		return -1, convertError(err)
	}

	return counterValue, nil
}

// GetAllMetrics возвращает все существующие метрики.
func (client *PostgresClient) GetAllMetrics(ctx context.Context) ([]models.MetricInfo, error) {
	rows, err := client.db.QueryContext(ctx, "SELECT id, type, gauge_value, counter_value FROM metrics;")
	if err != nil {
		return nil, convertError(err)
	}
	defer rows.Close()

	metrics := make([]models.MetricInfo, 0)

	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return nil, convertError(err)
		}

		var (
			metricID        string
			metricType      models.MetricType
			maybeMetricType psqlMetricType
			gaugeValue      float64
			counterValue    int64
		)

		err = rows.Scan(&metricID, &maybeMetricType, &gaugeValue, &counterValue)
		if err != nil {
			return nil, convertError(err)
		}

		metricType, err = deserializeMetricType(maybeMetricType)
		if err != nil {
			return nil, convertError(err)
		}

		switch metricType {
		case models.Gauge:
			metrics = append(metrics, models.NewGaugeMetric(metricID, gaugeValue))
		case models.Counter:
			metrics = append(metrics, models.NewCounterMetric(metricID, counterValue))
		}
	}

	return metrics, nil
}
