package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/reochoi109/go-handbook/error/example/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, err error) {
	var app *domain.AppError
	if errors.As(err, &app) {
		switch app.Code {
		case domain.CodeUnauthorized:
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
			return
		case domain.CodeNotFound:
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
			return
		case domain.CodeInvalid:
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid"})
			return
		case domain.CodeTimeout:
			writeJSON(w, http.StatusGatewayTimeout, map[string]any{"error": "timeout"})
			return
		}
	}
	writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal"})
}
