package natsstreaming

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/stan.go"
	"github.com/taaanechka/order-service/internal/api-server/api/natsstreaming/orders"
	"github.com/taaanechka/order-service/internal/api-server/services/orderservice"
)

type API struct {
	lg      *slog.Logger
	handler *orders.Handler
}

func NewHandler(lg *slog.Logger, service *orderservice.Service, sconn stan.Conn) *API {
	ordersHandler := orders.NewHandler(lg, service, sconn)
	return &API{
		lg:      lg,
		handler: ordersHandler,
	}
}

func (a *API) Subscribe() error {
	opts := []stan.SubscriptionOption{stan.SetManualAckMode(), stan.AckWait(time.Second * 30)}
	if err := a.handler.Subscribe(opts); err != nil {
		return fmt.Errorf("NATS API: failed to subscribe: %w", err)
	}
	return nil
}

func (a *API) Drain() {
	a.handler.Drain()
}
