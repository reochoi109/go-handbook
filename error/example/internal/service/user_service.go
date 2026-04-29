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
		if errors.Is(err, errors.New("row not found")) {
			return domain.User{}, domain.NotFound("user.me", fmt.Errorf("user not found: %w", err))
		}
		return domain.User{}, fmt.Errorf("user.me : %w", err)
	}
	return u, nil
}
