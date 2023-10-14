package orders

import (
	"context"
	"encoding/json"
	"log/slog"
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
	lg      *slog.Logger
}

func NewHandler(lg *slog.Logger, service *orderservice.Service) handlers.Handler {
	return &handler{
		service: service,
		lg:      lg,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, orderURL, apperror.Middleware(h.lg, h.GetOrderByUUID))
}

func (h *handler) GetOrderByUUID(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("uuid")

	h.lg.Info("called h.service.GetByUUID: get order by uid")
	res, err := h.service.GetByUUID(context.Background(), id)
	if err != nil {
		h.lg.Error("failed to get order by uid", "err", err)
		return err
	}

	resBytes, err := json.Marshal(&res)
	if err != nil {
		h.lg.Error("failed to marshal order", "err", err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	if _, errWr := w.Write(resBytes); errWr != nil {
		h.lg.Error("failed to write res data in response", "err", errWr)
		return errWr
	}

	return nil
}
