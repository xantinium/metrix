package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/xantinium/metrix/internal/models"
	"github.com/xantinium/metrix/internal/tools"
)

// NewPostgresClient создаёт новый клиент для работы с PostgreSQL.
func NewPostgresClient(ctx context.Context, connStr string) (*PostgresClient, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	client := &PostgresClient{
		db:      db,
		retrier: tools.DefaulRetrier,
	}

	err = client.initTables(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return client, nil
}

// PostgresClient клиент для работы с PostgreSQL.
type PostgresClient struct {
	db      *sql.DB
	retrier *tools.Retrier
}

// Ping проверка соединения.
func (client *PostgresClient) Ping(ctx context.Context) error {
	var err error

	client.retrier.Exec(func() bool {
		err = client.db.PingContext(ctx)
		return shouldRetry(err)
	})

	return convertError(err)
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

func shouldRetry(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
}
