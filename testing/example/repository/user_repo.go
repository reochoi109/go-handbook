package repository

import (
	"context"
	"sync"

	"github.com/reochoi109/go-handbook/testing/example/internal/domain"
)

type UserRepo interface {
	FindByID(ctx context.Context, id string) (domain.User, error)
}

type InMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]domain.User
}

func NewInMemoryUserRepo(seed []domain.User) *InMemoryUserRepo {
	m := make(map[string]domain.User, len(seed))
	for _, u := range seed {
		m[u.ID] = u
	}
	return &InMemoryUserRepo{users: m}
}

func (r *InMemoryUserRepo) FindByID(ctx context.Context, id string) (domain.User, error) {
	select {
	case <-ctx.Done():
		return domain.User{}, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}
	return u, nil
}
