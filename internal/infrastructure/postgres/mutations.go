package postgres

import "context"

// SaveMetrics сохраняет текущие значения метрик.
func (client *PostgresClient) SaveMetrics(ctx context.Context) error {
	return nil
}

// UpdateGaugeMetric обновляет текущее значение метрики типа Gauge
// с именем name, перезаписывая его значением value.
//
// Возвращает обновлённое значение метрики.
func (client *PostgresClient) UpdateGaugeMetric(ctx context.Context, name string, value float64) (float64, error) {
	return 0, nil
}

// UpdateCounterMetric обновляет текущее значение метрики типа Counter
// с именем name, добавляя к нему значение value.
//
// Возвращает обновлённое значение метрики.
func (client *PostgresClient) UpdateCounterMetric(ctx context.Context, name string, value int64) (int64, error) {
	return 0, nil
}
