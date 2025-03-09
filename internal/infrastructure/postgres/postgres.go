package postgres

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewPostgresClient создаёт новый клиент для работы с PostgreSQL.
func NewPostgresClient(connStr string) (*PostgresClient, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresClient{db: db}, nil
}

// PostgresClient клиент для работы с PostgreSQL.
type PostgresClient struct {
	db *sql.DB
}

// Ping проверка соединения.
func (client *PostgresClient) Ping(ctx context.Context) error {
	return client.db.PingContext(ctx)
}

// Destroy уничтожает клиент.
func (client *PostgresClient) Destroy() {
	client.db.Close()
}
