package orders

import (
	"encoding/json"
	"fmt"

	"html/template"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/taaanechka/order-service/internal/api-server/services/orderservice"
	"github.com/taaanechka/order-service/internal/apperror"
)

const (
	apiOrdersURL = "/api/v1/orders"
	apiOrderURL  = "/api/v1/orders/:uuid"
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
	// Backend
	router.HandlerFunc(http.MethodGet, apiOrderURL, LogMiddleware(h.lg, apperror.Middleware(h.lg, h.GetOrderByUUID)))
	router.HandlerFunc(http.MethodGet, apiOrdersURL, LogMiddleware(h.lg, apperror.Middleware(h.lg, h.GetAllUUIDs)))

	// Frontend
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.HandlerFunc(http.MethodGet, "/", LogMiddleware(h.lg, apperror.Middleware(h.lg, h.OrderPage)))
}

func (h *Handler) GetAllUUIDs(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	uids, err := h.service.GetAllUUIDs(r.Context())
	if err != nil {
		if len(uids) == 0 {
			h.lg.Error("httpHandler: failed to find orders in cache", "err", err)
			return apperror.ErrNotFound
		}
		h.lg.Error("httpHandler: failed to get orders uids", "err", err)
		return err
	}

	resBytes, err := json.Marshal(&uids)
	if err != nil {
		h.lg.Error("httpHandler: failed to marshal orders uids", "err", err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	if _, errWr := w.Write(resBytes); errWr != nil {
		h.lg.Error("httpHandler: failed to write res data in response", "err", errWr)
		return errWr
	}

	return nil
}

func (h *Handler) GetOrderByUUID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := httprouter.ParamsFromContext(ctx)
	id := params.ByName("uuid")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

func (h *Handler) OrderPage(w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("./templates/order.html")
	if err != nil {
		h.lg.Error("httpHandler: failed to parse files for new template", "err", err)
		fmt.Fprint(w, err.Error())
		return err
	}

	if err := t.Execute(w, nil); err != nil {
		fmt.Fprint(w, err.Error())
		return err
	}
	return nil
}
