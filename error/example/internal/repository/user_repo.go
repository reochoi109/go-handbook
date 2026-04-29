package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/reochoi109/go-handbook/error/example/internal/domain"
)

var errRowNotFound = errors.New("row not found")

type UserRepo interface {
	FindByID(ctx context.Context, id string) (domain.User, error)
}

type InMemoryUserRepo struct {
	users map[string]domain.User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users: map[string]domain.User{
			"user_123": {ID: "user_123", Name: "Reo"},
		},
	}
}

func (r *InMemoryUserRepo) FindByID(ctx context.Context, id string) (domain.User, error) {
	select {
	case <-time.After(30 * time.Millisecond):
	case <-ctx.Done():
		return domain.User{}, ctx.Err()
	}
	u, ok := r.users[id]
	if !ok {
		return domain.User{}, fmt.Errorf("repo find user id=%s: %w", id, errRowNotFound)
	}
	return u, nil
}
