// testing/example/internal/web/handler_test.go
package web

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reochoi109/go-handbook/testing/example/internal/domain"
)

type fakeUserService struct {
	get func(ctx context.Context, id string) (domain.User, error)
}

func (f fakeUserService) GetUser(ctx context.Context, id string) (domain.User, error) {
	return f.get(ctx, id)
}

func TestHandler_GetUser(t *testing.T) {
	t.Run("method_not_allowed", func(t *testing.T) {
		h := NewHandler(fakeUserService{
			get: func(ctx context.Context, id string) (domain.User, error) {
				t.Fatalf("service should not be called")
				return domain.User{}, nil
			},
		})

		req := httptest.NewRequest(http.MethodPost, "/users/u1", nil)
		rr := httptest.NewRecorder()

		h.GetUser(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("got %d, want %d", rr.Code, http.StatusMethodNotAllowed)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		h := NewHandler(fakeUserService{
			get: func(ctx context.Context, id string) (domain.User, error) {
				return domain.User{}, domain.ErrUserNotFound
			},
		})

		req := httptest.NewRequest(http.MethodGet, "/users/u404", nil)
		rr := httptest.NewRecorder()

		h.GetUser(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("got %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("ok", func(t *testing.T) {
		h := NewHandler(fakeUserService{
			get: func(ctx context.Context, id string) (domain.User, error) {
				_ = ctx.(context.Context)
				if id != "u1" {
					return domain.User{}, errors.New("unexpected id")
				}
				return domain.User{ID: "u1", Name: "Reo"}, nil
			},
		})

		req := httptest.NewRequest(http.MethodGet, "/users/u1", nil)
		rr := httptest.NewRecorder()

		h.GetUser(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("got %d, want %d", rr.Code, http.StatusOK)
		}
	})
}
