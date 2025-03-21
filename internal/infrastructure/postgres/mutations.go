package postgres

import (
	"context"

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
	row := client.db.QueryRowContext(ctx, "INSERT INTO metrics (id, type, gauge_value, counter_value)"+
		" VALUES ($1, $2, $3, 0)"+
		" ON CONFLICT (id, type)"+
		" DO UPDATE SET"+
		" gauge_value = $3"+
		" RETURNING gauge_value;",
		id,
		serializeMetricType(models.Gauge),
		value)

	var gaugeValue float64
	err := row.Scan(&gaugeValue)
	if err != nil {
		return -1, convertError(err)
	}

	return gaugeValue, nil
}

// UpdateCounterMetric обновляет текущее значение метрики типа Counter
// с идентификатором id, добавляя к нему значение value.
//
// Возвращает обновлённое значение метрики.
func (client *PostgresClient) UpdateCounterMetric(ctx context.Context, id string, value int64) (int64, error) {
	row := client.db.QueryRowContext(ctx, "INSERT INTO metrics (id, type, gauge_value, counter_value)"+
		" VALUES ($1, $2, 0, $3)"+
		" ON CONFLICT (id, type)"+
		" DO UPDATE SET"+
		" counter_value = metrics.counter_value + $3"+
		" RETURNING counter_value;",
		id,
		serializeMetricType(models.Gauge),
		value)

	var counterValue int64
	err := row.Scan(&counterValue)
	if err != nil {
		return -1, convertError(err)
	}

	return counterValue, nil
}
