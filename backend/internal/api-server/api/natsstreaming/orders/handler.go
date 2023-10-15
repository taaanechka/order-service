package orders

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	stan "github.com/nats-io/stan.go"
	"github.com/taaanechka/order-service/internal/api-server/services/orderservice"
	"github.com/taaanechka/order-service/internal/api-server/services/ports/ordersrepository"
)

type Handler struct {
	service *orderservice.Service
	sc      stan.Conn
	sub     stan.Subscription
	lg      *slog.Logger
}

func NewHandler(lg *slog.Logger, service *orderservice.Service, sc stan.Conn) *Handler {
	lg.Info("natsHandler: new order handler")
	return &Handler{
		service: service,
		sc:      sc,
		lg:      lg,
	}
}

func (h *Handler) Subscribe() error {
	h.lg.Info("natsHandler: subscribe to the channel")
	var err error
	h.sub, err = h.sc.Subscribe("orders", func(msg *stan.Msg) {
		h.lg.Info("natsHandler", "received message", string(msg.Data))

		var order ordersrepository.Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			h.lg.Error("natsHandler: failed to unmarshal msg.Data", "err", err)
			return
		}
		if _, err := h.service.Create(context.Background(), order); err != nil {
			h.lg.Error("natsHandler: failed to create order", "err", err)
			return
		}
	})
	if err != nil {
		return fmt.Errorf("natsHandler: failed to subscribe")
	}
	return nil
}
