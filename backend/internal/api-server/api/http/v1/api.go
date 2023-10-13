package v1

import (
	"log"

	"github.com/julienschmidt/httprouter"
	"github.com/taaanechka/order-service/internal/api-server/api/http/v1/orders"
	"github.com/taaanechka/order-service/internal/api-server/services/orderservice"
	"github.com/taaanechka/order-service/internal/handlers"
)

type API struct {
	lg      *log.Logger
	service *orderservice.Service
}

func NewHandler(lg *log.Logger, service *orderservice.Service) handlers.Handler {
	return &API{
		service: service,
		lg:      lg,
	}
}

func (a *API) Register(router *httprouter.Router) {
	users := orders.NewHandler(a.lg, a.service)
	users.Register(router)
}
