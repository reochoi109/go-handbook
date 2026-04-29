package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/reochoi109/go-handbook/testing/example/internal/domain"
)

type UserService interface {
	GetUser(ctx context.Context, id string) (domain.User, error)
}

type Handler struct {
	users UserService
}

func NewHandler(users UserService) *Handler {
	return &Handler{users: users}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" || id == "/users" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	u, err := h.users.GetUser(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidID):
			http.Error(w, "invalid id", http.StatusBadRequest)
		case errors.Is(err, domain.ErrUserNotFound):
			http.Error(w, "not found", http.StatusNotFound)
		default:
			http.Error(w, "internal", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(u)
}
