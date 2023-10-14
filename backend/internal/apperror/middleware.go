package apperror

import (
	"errors"
	"log/slog"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(lg *slog.Logger, h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := h(w, r)
		if err == nil {
			return
		}

		var appErr *AppError
		if !errors.As(err, &appErr) {
			w.WriteHeader(http.StatusTeapot)
			if _, errWr := w.Write(systemError(err).Marshal()); errWr != nil {
				lg.Error("failed to write appErr data in response", "err", errWr)
			}

			return
		}

		switch {
		case errors.Is(appErr, ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
			if _, errWr := w.Write(ErrNotFound.Marshal()); errWr != nil {
				lg.Error("failed to write ErrNotFound data in response", "err", errWr)
			}
		case errors.Is(appErr, ErrValidate):
			w.WriteHeader(http.StatusInternalServerError)
			if _, errWr := w.Write(ErrValidate.Marshal()); errWr != nil {
				lg.Error("failed to write ErrValidate data in response", "err", errWr)
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			if _, errWr := w.Write(appErr.Marshal()); errWr != nil {
				lg.Error("failed to write appErr data in response", "err", errWr)
			}
		}
	}
}
