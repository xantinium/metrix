package postgres

import "context"

func (client *PostgresClient) initTables(ctx context.Context) error {
	err := client.initMetricsTable(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (client *PostgresClient) initMetricsTable(ctx context.Context) error {
	_, err := client.db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS metrics ("+
		"name VARCHAR(50) NOT NULL "+
		"type SMALLINT NOT NULL "+
		"gauge_value DOUBLE PRECISION NOT NULL "+
		"counter_value INTEGER NOT NULL "+
		"PRIMARY KEY (name, type)"+
		");")

	return err
}
