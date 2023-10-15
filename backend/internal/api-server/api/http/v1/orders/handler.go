package orders

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/taaanechka/order-service/internal/api-server/services/orderservice"
	"github.com/taaanechka/order-service/internal/apperror"
)

const (
	ordersURL = "/orders"
	orderURL  = "/orders/:uuid"
)

type Handler struct {
	service *orderservice.Service
	lg      *slog.Logger
}

func NewHandler(lg *slog.Logger, service *orderservice.Service) *Handler {
	lg.Info("httpHandler: new order handler")
	return &Handler{
		service: service,
		lg:      lg,
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LogMiddleware(lg *slog.Logger, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)
		h(lrw, r)
		lg.Info("LogMiddleware", "Method", r.Method, "URL", r.URL, "Status", lrw.statusCode, "Addr", r.RemoteAddr)
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	h.lg.Info("httpHandler: register order handler")
	router.HandlerFunc(http.MethodGet, orderURL, LogMiddleware(h.lg, apperror.Middleware(h.lg, h.GetOrderByUUID)))
}

func (h *Handler) GetOrderByUUID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	id := params.ByName("uuid")

	res, err := h.service.GetByUUID(ctx, id)
	if err != nil {
		h.lg.Error("httpHandler: failed to get order by uid", "err", err)
		return err
	}

	resBytes, err := json.Marshal(&res)
	if err != nil {
		h.lg.Error("httpHandler: failed to marshal order", "err", err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	if _, errWr := w.Write(resBytes); errWr != nil {
		h.lg.Error("httpHandler: failed to write res data in response", "err", errWr)
		return errWr
	}

	return nil
}
