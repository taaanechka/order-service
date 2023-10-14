package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/taaanechka/order-service/internal/api-server/services/ports/ordersrepository"
)

type Client interface {
	ExecContext(ctx context.Context, sql string, arguments ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, sql string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func NewClient(ctx context.Context, rc ordersrepository.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		rc.Host, rc.Port, rc.Username, rc.Password, rc.Database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return db, nil
}
