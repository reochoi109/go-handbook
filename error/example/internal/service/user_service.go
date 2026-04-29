package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/reochoi109/go-handbook/error/example/internal/domain"
	"github.com/reochoi109/go-handbook/error/example/internal/repository"
)

type UserService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Me(ctx context.Context, actorUserID string) (domain.User, error) {
	if actorUserID == "" {
		return domain.User{}, domain.Unauthorized("user.me", errors.New("missing actor"))
	}

	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	u, err := s.repo.FindByID(ctx, actorUserID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRowNotFound):
			return domain.User{}, domain.NotFound("user.me", fmt.Errorf("user not found: %w", err))
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return domain.User{}, domain.Timeout("user.me", err)
		}
		return domain.User{}, fmt.Errorf("user.me: %w", err)
	}
	return u, nil
}
