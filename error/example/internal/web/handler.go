package web

import (
	"context"
	"net/http"
	"time"

	"github.com/reochoi109/go-handbook/error/example/internal/service"
)

type Handler struct {
	users *service.UserService
}

func NewHandler(users *service.UserService) *Handler {
	return &Handler{users: users}
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 800*time.Millisecond)
	defer cancel()

	actor := r.Header.Get("X-User-Id")
	u, err := h.users.Me(ctx, actor)
	if err != nil {
		writeErr(w, err)
		return
	}

	writeJSON(w, http.StatusOK, u)
}
