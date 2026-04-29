package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/reochoi109/go-handbook/testing/example/internal/domain"
	"github.com/reochoi109/go-handbook/testing/example/repository"
)

type UserService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, id string) (domain.User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.User{}, fmt.Errorf("service get user: %w", domain.ErrInvalidID)
	}

	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.User{}, fmt.Errorf("service get user: %w", domain.ErrUserNotFound)
		}
		return domain.User{}, fmt.Errorf("service get user: %w", err)
	}
	return u, nil
}
