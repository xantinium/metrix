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
	row := client.db.QueryRowContext(ctx, "UPDATE metrics"+
		" SET gauge_value = @value"+
		" WHERE id = @id AND type = @type"+
		" RETURNING gauge_value;", id, models.Gauge)

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
	row := client.db.QueryRowContext(ctx, "UPDATE metrics"+
		" SET counter_value = counter_value + @value"+
		" WHERE id = @id AND type = @type"+
		" RETURNING counter_value;", id, models.Counter)

	var counterValue int64
	err := row.Scan(&counterValue)
	if err != nil {
		return -1, convertError(err)
	}

	return counterValue, nil
}
