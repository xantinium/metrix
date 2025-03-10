package postgres

import (
	"context"
	"database/sql"

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
	row := client.db.QueryRowContext(ctx, "INSERT INTO metrix (metric_id, metric_type, gauge_value, counter_value)"+
		" VALUES (@metric_id, @metric_type, @gauge_value, 0)"+
		" ON CONFLICT (metric_id, metric_type)"+
		" DO UPDATE SET"+
		" gauge_value = @gauge_value"+
		" RETURNING gauge_value;",
		sql.Named("metric_id", id),
		sql.Named("metric_type", models.Gauge),
		sql.Named("gauge_value", value))

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
	row := client.db.QueryRowContext(ctx, "INSERT INTO metrix (metric_id, metric_type, gauge_value, counter_value)"+
		" VALUES (@metric_id, @metric_type, 0, @counter_value)"+
		" ON CONFLICT (metric_id, metric_type)"+
		" DO UPDATE SET"+
		" counter_value = counter_value + @counter_value"+
		" RETURNING counter_value;",
		sql.Named("metric_id", id),
		sql.Named("metric_type", models.Gauge),
		sql.Named("counter_value", value))

	var counterValue int64
	err := row.Scan(&counterValue)
	if err != nil {
		return -1, convertError(err)
	}

	return counterValue, nil
}
