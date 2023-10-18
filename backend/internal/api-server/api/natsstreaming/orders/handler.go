package orders

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	stan "github.com/nats-io/stan.go"
	"github.com/taaanechka/order-service/internal/api-server/services/orderservice"
	"github.com/taaanechka/order-service/internal/api-server/services/ports/ordersrepository"
	"github.com/taaanechka/order-service/internal/apperror"
)

type Handler struct {
	service *orderservice.Service
	sconn   stan.Conn
	sub     stan.Subscription
	lg      *slog.Logger
}

func NewHandler(lg *slog.Logger, service *orderservice.Service, sconn stan.Conn) *Handler {
	lg.Info("natsHandler: new order handler")
	return &Handler{
		service: service,
		sconn:   sconn,
		lg:      lg,
	}
}

func (h *Handler) Subscribe(startOpts []stan.SubscriptionOption) error {
	h.lg.Info("natsHandler: subscribe to the channel")
	var err error
	h.sub, err = h.sconn.Subscribe("orders", h.Handle, startOpts...)
	if err != nil {
		return fmt.Errorf("natsHandler: failed to subscribe")
	}
	return nil
}

func (h *Handler) Handle(msg *stan.Msg) {
	h.lg.Info("natsHandler", "received message", string(msg.Data))

	var order ordersrepository.Order
	bmsg := msg.Data
	if err := json.Unmarshal(bmsg, &order); err != nil {
		h.lg.Error("natsHandler: failed to unmarshal msg.Data", "err", err)
		return
	}
	ctx := context.Background()
	if _, err := h.service.Create(ctx, order); err != nil {
		h.lg.Error("natsHandler: failed to create order", "err", err)
		if errors.Is(err, apperror.ErrValidate) {
			if err := msg.Ack(); err != nil {
				h.lg.Warn("natsHandler: failed to manually acknowledge a message", "err", err)
			}
		}
		return
	}

	if err := msg.Ack(); err != nil {
		h.lg.Warn("natsHandler: failed to manually acknowledge a message", "err", err)
	}
}

func (h *Handler) Drain() {
	h.lg.Info("natsHandler: drain the connection")
	if err := h.sconn.NatsConn().Drain(); err != nil {
        h.lg.Error("natsHandler: failed to drain", "err", err)
    }
}
