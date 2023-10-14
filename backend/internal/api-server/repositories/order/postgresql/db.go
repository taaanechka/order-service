package postgresql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/taaanechka/order-service/internal/api-server/services/ports/ordersrepository"
	"github.com/taaanechka/order-service/internal/apperror"
	"github.com/taaanechka/order-service/pkg/client/postgresql"
)

type DB struct {
	client postgresql.Client
	lg     *slog.Logger
}

func NewRepository(lg *slog.Logger, cfg ordersrepository.Config) (*DB, error) {
	lg.Info("Postgres: new order repository")
	ctx := context.Background()
	db, err := postgresql.NewClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init client: %w", err)
	}

	return &DB{
		client: db,
		lg:     lg,
	}, nil
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (db *DB) Create(ctx context.Context, order ordersrepository.Order) (string, error) {
	return "", nil
}

func (db *DB) FindAll(ctx context.Context) ([]ordersrepository.Order, error) {
	q := `
		SELECT data FROM orders;
	`
	db.lg.Info("Postgres: FindAll", "query", formatQuery(q))

	rows, err := db.client.QueryContext(ctx, q)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.ErrNotFound
		}
		db.lg.Error("invalid sql query", "err", err)
		return nil, err
	}

	orders := make([]ordersrepository.Order, 0)
	for rows.Next() {
		var jsonData []byte
		if err = rows.Scan(&jsonData); err != nil {
			db.lg.Error("failed to scan rows.next to jsonData", "err", err)
			return nil, err
		}

		var or ordersrepository.Order
		if err = json.Unmarshal(jsonData, &or); err != nil {
			db.lg.Error("invalid order data format", "err", err)
			return nil, apperror.ErrValidate
		}

		orders = append(orders, or)
	}

	if err = rows.Err(); err != nil {
		db.lg.Error("rows error", "err", err)
		return nil, err
	}

	return orders, nil
}

func (db *DB) FindOne(ctx context.Context, id string) (ordersrepository.Order, error) {
	q := `
		SELECT data FROM orders WHERE data->>'order_uid' = $1
	`
	db.lg.Info("Postgres: FindOne", "query", formatQuery(q))

	var jsonData []byte
	err := db.client.QueryRowContext(ctx, q, id).Scan(&jsonData)
	if err != nil {
		if err == sql.ErrNoRows {
			return ordersrepository.Order{}, apperror.ErrNotFound
		}
		db.lg.Error("invalid sql query", "err", err)
		return ordersrepository.Order{}, err
	}

	var order ordersrepository.Order
	if err = json.Unmarshal(jsonData, &order); err != nil {
		db.lg.Error("invalid order data format", "err", err)
		return ordersrepository.Order{}, apperror.ErrValidate
	}

	return order, nil
}
