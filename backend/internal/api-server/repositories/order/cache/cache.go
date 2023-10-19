package cache

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/taaanechka/order-service/internal/api-server/services/ports/ordersrepository"
	"github.com/taaanechka/order-service/internal/apperror"
)

type Cache struct {
	sync.RWMutex
	orders map[string]string
	lg     *slog.Logger
}

func NewRepository(lg *slog.Logger) (*Cache, error) {
	lg.Info("Cache: new order repository")
	return &Cache{
		orders: make(map[string]string),
		lg:     lg,
	}, nil
}

func (c *Cache) CreateMany(ctx context.Context, orders []ordersrepository.Order) ([]string, error) {
	c.Lock()

	defer c.Unlock()

	uids := make([]string, 0)
	for _, or := range orders {
		var byteData []byte
		byteData, err := json.Marshal(&or)
		if err != nil {
			c.lg.Error("Cache: failed to marshal order data", "err", err)
			return nil, apperror.ErrCreate
		}

		uid := or.OrderUid
		c.orders[uid] = string(byteData)
		uids = append(uids, uid)
	}
	return uids, nil
}

func (c *Cache) CreateOne(ctx context.Context, order ordersrepository.Order) (string, error) {
	c.Lock()

	defer c.Unlock()

	var byteData []byte
	byteData, err := json.Marshal(&order)
	if err != nil {
		c.lg.Error("Cache: failed to marshal order data", "err", err)
		return "", apperror.ErrCreate
	}

	c.orders[order.OrderUid] = string(byteData)
	return order.OrderUid, nil
}

func (c *Cache) FindAllUUIDs(ctx context.Context) ([]string, error) {
	c.RLock()

	defer c.RUnlock()

	uids := make([]string, 0)
	for uid := range c.orders {
		uids = append(uids, uid)
	}
	return uids, nil
}

func (c *Cache) FindOne(ctx context.Context, id string) (ordersrepository.Order, error) {
	c.RLock()

	defer c.RUnlock()

	orStr, found := c.orders[id]
	if !found {
		c.lg.Error("Cache: failed to find order by uid", "err", apperror.ErrNotFound)
		return ordersrepository.Order{}, apperror.ErrNotFound
	}

	var order ordersrepository.Order
	if err := json.Unmarshal([]byte(orStr), &order); err != nil {
		c.lg.Error("Cache: failed to unmarshal order data", "err", err)
		return ordersrepository.Order{}, err
	}
	return order, nil
}
