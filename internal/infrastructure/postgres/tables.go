package postgres

import (
	"context"
	"fmt"

	"github.com/xantinium/metrix/internal/models"
)

type psqlMetricType = uint8

const (
	unknown psqlMetricType = iota
	gauge
	counter
)

func serializeMetricType(metricType models.MetricType) psqlMetricType {
	switch metricType {
	case models.Gauge:
		return gauge
	case models.Counter:
		return counter
	default:
		return unknown
	}
}

func deserializeMetricType(metricType psqlMetricType) (models.MetricType, error) {
	switch metricType {
	case gauge:
		return models.Gauge, nil
	case counter:
		return models.Counter, nil
	default:
		return "", fmt.Errorf("unknown metric type: %q", metricType)
	}
}

func (client *PostgresClient) initTables(ctx context.Context) error {
	err := client.initMetricsTable(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (client *PostgresClient) initMetricsTable(ctx context.Context) error {
	var err error

	client.retrier.Exec(func() bool {
		_, err = client.db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS metrics ("+
			"id VARCHAR(50) NOT NULL,"+
			"type SMALLINT NOT NULL,"+
			"gauge_value DOUBLE PRECISION NOT NULL,"+
			"counter_value BIGINT NOT NULL,"+
			"PRIMARY KEY (id, type)"+
			");")
		return shouldRetry(err)
	})

	return convertError(err)
}
