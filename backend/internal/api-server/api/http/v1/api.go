package v1

import (
	"log/slog"

	"github.com/julienschmidt/httprouter"
	"github.com/taaanechka/order-service/internal/api-server/api/http/v1/orders"
	"github.com/taaanechka/order-service/internal/api-server/services/orderservice"
	"github.com/taaanechka/order-service/internal/handlers"
)

type API struct {
	lg      *slog.Logger
	service *orderservice.Service
}

func NewHandler(lg *slog.Logger, service *orderservice.Service) handlers.Handler {
	return &API{
		service: service,
		lg:      lg,
	}
}

func (a *API) Register(router *httprouter.Router) {
	ordersHandler := orders.NewHandler(a.lg, a.service)
	ordersHandler.Register(router)
}
