package postgres

import (
	"context"
	"database/sql"

	"github.com/xantinium/metrix/internal/models"
)

// GetGaugeMetric возвращает метрику типа Gauge по идентификатору id.
func (client *PostgresClient) GetGaugeMetric(ctx context.Context, id string) (float64, error) {
	var (
		err   error
		value float64
	)

	client.retrier.Exec(func() bool {
		row := client.db.QueryRowContext(ctx, "SELECT gauge_value FROM metrics"+
			" WHERE id = $1 AND type = $2;",
			id,
			serializeMetricType(models.Gauge))

		err = row.Scan(&value)
		return shouldRetry(err)
	})

	return value, convertError(err)
}

// GetCounterMetric возвращает метрику типа Counter по идентификатору id.
func (client *PostgresClient) GetCounterMetric(ctx context.Context, id string) (int64, error) {
	var (
		err   error
		value int64
	)

	client.retrier.Exec(func() bool {
		row := client.db.QueryRowContext(ctx, "SELECT counter_value FROM metrics"+
			" WHERE id = $1 AND type = $2;",
			id,
			serializeMetricType(models.Counter))

		err = row.Scan(&value)
		return shouldRetry(err)
	})

	return value, convertError(err)
}

// GetAllMetrics возвращает все существующие метрики.
func (client *PostgresClient) GetAllMetrics(ctx context.Context) ([]models.MetricInfo, error) {
	var (
		err     error
		rows    *sql.Rows
		metrics []models.MetricInfo
	)

	client.retrier.Exec(func() bool {
		rows, err = client.db.QueryContext(ctx, "SELECT id, type, gauge_value, counter_value FROM metrics;")
		if err != nil {
			return shouldRetry(err)
		}
		defer rows.Close()

		metrics = make([]models.MetricInfo, 0)

		for rows.Next() {
			err = rows.Err()
			if err != nil {
				return shouldRetry(err)
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
				return shouldRetry(err)
			}

			metricType, err = deserializeMetricType(maybeMetricType)
			if err != nil {
				return shouldRetry(err)
			}

			switch metricType {
			case models.Gauge:
				metrics = append(metrics, models.NewGaugeMetric(metricID, gaugeValue))
			case models.Counter:
				metrics = append(metrics, models.NewCounterMetric(metricID, counterValue))
			}
		}

		return false
	})

	return metrics, convertError(err)
}
