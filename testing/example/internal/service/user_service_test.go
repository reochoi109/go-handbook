// testing/example/internal/service/user_service_test.go
package service

import (
	"context"
	"errors"
	"testing"

	"github.com/reochoi109/go-handbook/testing/example/internal/domain"
)

type fakeRepo struct {
	find func(ctx context.Context, id string) (domain.User, error)
}

func (f fakeRepo) FindByID(ctx context.Context, id string) (domain.User, error) {
	return f.find(ctx, id)
}

func TestUserService_GetUser(t *testing.T) {
	t.Run("invalid_id", func(t *testing.T) {
		svc := NewUserService(fakeRepo{
			find: func(ctx context.Context, id string) (domain.User, error) {
				t.Fatalf("repo should not be called")
				return domain.User{}, nil
			},
		})

		_, err := svc.GetUser(context.Background(), "   ")
		if !errors.Is(err, domain.ErrInvalidID) {
			t.Fatalf("got err=%v, want ErrInvalidID", err)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		svc := NewUserService(fakeRepo{
			find: func(ctx context.Context, id string) (domain.User, error) {
				return domain.User{}, domain.ErrUserNotFound
			},
		})

		_, err := svc.GetUser(context.Background(), "u404")
		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Fatalf("got err=%v, want ErrUserNotFound", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		svc := NewUserService(fakeRepo{
			find: func(ctx context.Context, id string) (domain.User, error) {
				return domain.User{ID: id, Name: "Reo"}, nil
			},
		})

		u, err := svc.GetUser(context.Background(), "u1")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if u.ID != "u1" {
			t.Fatalf("got id=%q, want %q", u.ID, "u1")
		}
	})
}
