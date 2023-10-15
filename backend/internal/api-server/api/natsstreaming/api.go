package natsstreaming

import (
	"fmt"
	"log/slog"

	"github.com/nats-io/stan.go"
	"github.com/taaanechka/order-service/internal/api-server/api/natsstreaming/orders"
	"github.com/taaanechka/order-service/internal/api-server/services/orderservice"
)

type API struct {
	lg      *slog.Logger
	sconn   stan.Conn
	service *orderservice.Service
}

func NewHandler(lg *slog.Logger, service *orderservice.Service, sconn stan.Conn) *API {
	return &API{
		service: service,
		sconn:   sconn,
		lg:      lg,
	}
}

func (a *API) Subscribe() error {
	ordersHandler := orders.NewHandler(a.lg, a.service, a.sconn)

	if err := ordersHandler.Subscribe(); err != nil {
		return fmt.Errorf("NATS API: failed to subscribe: %w", err)
	}
	return nil
}
