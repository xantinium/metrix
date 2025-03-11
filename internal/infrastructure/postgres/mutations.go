package postgres

import (
	"context"
	"database/sql"

	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/models"
)

// SaveMetrics сохраняет текущие значения метрик.
func (client *PostgresClient) SaveMetrics(ctx context.Context) error {
	return nil
}

// UpdateGaugeMetric обновляет текущее значение метрики типа Gauge
// с идентификатором id, перезаписывая его значением value.
//
// Возвращает обновлённое значение метрики.
func (client *PostgresClient) UpdateGaugeMetric(ctx context.Context, id string, value float64) (float64, error) {
	var (
		err      error
		newValue float64
	)

	client.retrier.Exec(func() bool {
		row := client.db.QueryRowContext(ctx, "INSERT INTO metrics (id, type, gauge_value, counter_value)"+
			" VALUES ($1, $2, $3, 0)"+
			" ON CONFLICT (id, type)"+
			" DO UPDATE SET"+
			" gauge_value = $3"+
			" RETURNING gauge_value;",
			id,
			serializeMetricType(models.Gauge),
			value)

		err = row.Scan(&newValue)
		return shouldRetry(err)
	})

	return newValue, convertError(err)
}

// UpdateCounterMetric обновляет текущее значение метрики типа Counter
// с идентификатором id, добавляя к нему значение value.
//
// Возвращает обновлённое значение метрики.
func (client *PostgresClient) UpdateCounterMetric(ctx context.Context, id string, value int64) (int64, error) {
	var (
		err      error
		newValue int64
	)

	client.retrier.Exec(func() bool {
		row := client.db.QueryRowContext(ctx, "INSERT INTO metrics (id, type, gauge_value, counter_value)"+
			" VALUES ($1, $2, 0, $3)"+
			" ON CONFLICT (id, type)"+
			" DO UPDATE SET"+
			" counter_value = metrics.counter_value + $3"+
			" RETURNING counter_value;",
			id,
			serializeMetricType(models.Gauge),
			value)

		err = row.Scan(&newValue)
		return shouldRetry(err)
	})

	return newValue, convertError(err)
}

// UpdateMetrics обновляет текущее значение метрик.
// Используется батчевое обновление через транзакцию.
func (client *PostgresClient) UpdateMetrics(ctx context.Context, metrics []models.MetricInfo) error {
	if len(metrics) == 0 {
		return nil
	}

	var (
		err error
		tx  *sql.Tx
	)

	client.retrier.Exec(func() bool {
		tx, err = client.db.BeginTx(ctx, nil)
		if err != nil {
			return shouldRetry(err)
		}

		for _, metric := range metrics {
			switch metric.Type() {
			case models.Gauge:
				_, err = tx.ExecContext(ctx, "INSERT INTO metrics (id, type, gauge_value, counter_value)"+
					" VALUES ($1, $2, $3, 0)"+
					" ON CONFLICT (id, type)"+
					" DO UPDATE SET"+
					" gauge_value = $3;",
					metric.ID(),
					serializeMetricType(models.Gauge),
					metric.GaugeValue())
			case models.Counter:
				_, err = tx.ExecContext(ctx, "INSERT INTO metrics (id, type, gauge_value, counter_value)"+
					" VALUES ($1, $2, 0, $3)"+
					" ON CONFLICT (id, type)"+
					" DO UPDATE SET"+
					" counter_value = metrics.counter_value + $3;",
					metric.ID(),
					serializeMetricType(models.Counter),
					metric.CounterValue())
			default:
				logger.Info("unknown metric type", logger.Field{Name: "type", Value: metric.Type()})
			}
			if err != nil {
				tx.Rollback()
				return shouldRetry(err)
			}
		}

		err = tx.Commit()
		return shouldRetry(err)
	})

	return convertError(err)
}
