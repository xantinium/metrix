package postgres

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/xantinium/metrix/internal/models"
)

// NewPostgresClient создаёт новый клиент для работы с PostgreSQL.
func NewPostgresClient(ctx context.Context, connStr string) (*PostgresClient, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	client := &PostgresClient{db: db}

	err = client.initTables(ctx)
	if err != nil {
		db.Close()
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
func (client *PostgresClient) Destroy(_ context.Context) {
	client.db.Close()
}

func convertError(err error) error {
	if err == sql.ErrNoRows {
		return models.ErrNotFound
	}

	return err
}
