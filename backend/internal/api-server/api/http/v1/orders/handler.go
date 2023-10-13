package orders

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/taaanechka/order-service/internal/api-server/services/orderservice"
	"github.com/taaanechka/order-service/internal/apperror"
	"github.com/taaanechka/order-service/internal/handlers"
)

const (
	ordersURL = "/orders"
	orderURL  = "/orders/:uuid"
)

type handler struct {
	service *orderservice.Service
	lg      *log.Logger
}

func NewHandler(lg *log.Logger, service *orderservice.Service) handlers.Handler {
	return &handler{
		service: service,
		lg:      lg,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, orderURL, apperror.Middleware(h.GetOrderByUUID))
}

func (h *handler) GetOrderByUUID(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("uuid")

	res, err := h.service.GetByUUID(context.Background(), id)
	if err != nil {
		return err
	}

	resBytes, err := json.Marshal(&res)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)

	return nil
}
