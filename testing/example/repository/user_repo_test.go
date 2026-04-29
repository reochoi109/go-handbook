package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/reochoi109/go-handbook/testing/example/internal/domain"
)

func TestInMemoryUserRepo_FindByID(t *testing.T) {
	repo := NewInMemoryUserRepo([]domain.User{
		{ID: "u1", Name: "Reo"},
	})

	t.Run("found", func(t *testing.T) {
		u, err := repo.FindByID(context.Background(), "u1")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if u.ID != "u1" {
			t.Fatalf("got id=%q, want %q", u.ID, "u1")
		}
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := repo.FindByID(context.Background(), "u2")
		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Fatalf("got err = %v, want = %v", err, domain.ErrUserNotFound)
		}
	})

}
