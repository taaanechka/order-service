package postgresql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/taaanechka/order-service/internal/api-server/services/ports/ordersrepository"
	"github.com/taaanechka/order-service/internal/apperror"
	"github.com/taaanechka/order-service/pkg/client/postgresql"
)

type DB struct {
	client postgresql.Client
	lg     *log.Logger
}

func NewRepository(lg *log.Logger, cfg ordersrepository.Config) (*DB, error) {
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

func (db *DB) FindOne(ctx context.Context, id string) (ordersrepository.Order, error) {
	q := `
		SELECT data FROM orders WHERE data->>'order_uid' = $1
	`
	db.lg.Printf(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var jsonData []byte
	err := db.client.QueryRowContext(ctx, q, id).Scan(&jsonData)
	if err != nil {
		if err == sql.ErrNoRows {
			return ordersrepository.Order{}, apperror.ErrNotFound
		}
		return ordersrepository.Order{}, err
	}

	var order ordersrepository.Order
	err = json.Unmarshal(jsonData, &order)
	if err != nil {
		db.lg.Println(err)
		return ordersrepository.Order{}, apperror.ErrValidate
	}

	return order, nil
}
