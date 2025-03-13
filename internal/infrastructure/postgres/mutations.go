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
		err    error
		metric models.MetricInfo
	)

	client.retrier.Exec(func() bool {
		metric, err = client.updateMetric(ctx, models.NewGaugeMetric(id, value))
		return shouldRetry(err)
	})

	return metric.GaugeValue(), err
}

// UpdateCounterMetric обновляет текущее значение метрики типа Counter
// с идентификатором id, добавляя к нему значение value.
//
// Возвращает обновлённое значение метрики.
func (client *PostgresClient) UpdateCounterMetric(ctx context.Context, id string, value int64) (int64, error) {
	var (
		err    error
		metric models.MetricInfo
	)

	client.retrier.Exec(func() bool {
		metric, err = client.updateMetric(ctx, models.NewCounterMetric(id, value))
		return shouldRetry(err)
	})

	return metric.CounterValue(), err
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
			_, err = client.updateMetric(ctx, metric)
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

// updateMetric обновляет текущее значение метрики с идентификатором id.
//
// Возвращает обновлённую структуру метрики.
func (client *PostgresClient) updateMetric(ctx context.Context, metric models.MetricInfo) (models.MetricInfo, error) {
	var (
		err       error
		newMetric models.MetricInfo
	)

	getOnConflictExpression := func() string {
		expression := "DO NOTHING"

		switch metric.Type() {
		case models.Gauge:
			expression = " DO UPDATE SET gauge_value = $3"
		case models.Counter:
			expression = " DO UPDATE SET counter_value = metrics.counter_value + $4"
		}

		return expression
	}

	row := client.db.QueryRowContext(ctx, "INSERT INTO metrics (id, type, gauge_value, counter_value)"+
		" VALUES ($1, $2, $3, $4)"+
		" ON CONFLICT (id, type)"+
		getOnConflictExpression()+
		" RETURNING counter_value;",
		metric.ID(),
		serializeMetricType(metric.Type()),
		metric.GaugeValue(),
		metric.CounterValue())

	switch metric.Type() {
	case models.Gauge:
		var newValue float64
		err = row.Scan(&newValue)
		newMetric = models.NewGaugeMetric(metric.ID(), newValue)
	case models.Counter:
		var newValue int64
		err = row.Scan(&newValue)
		newMetric = models.NewCounterMetric(metric.ID(), newValue)
	default:
		logger.Info("unknown metric type", logger.Field{Name: "type", Value: metric.Type()})
	}

	return newMetric, convertError(err)
}
