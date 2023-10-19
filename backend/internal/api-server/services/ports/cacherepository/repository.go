package cacherepository

import (
	"context"

	"github.com/taaanechka/order-service/internal/api-server/services/ports/ordersrepository"
)

type Repository interface {
	CreateMany(ctx context.Context, orders []ordersrepository.Order) ([]string, error)
	CreateOne(ctx context.Context, order ordersrepository.Order) (string, error)
	FindAllUUIDs(ctx context.Context) ([]string, error)
	FindOne(ctx context.Context, id string) (ordersrepository.Order, error)
}
