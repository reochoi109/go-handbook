package user

import (
	"context"
	"log/slog"
)

type Service struct {
	log *slog.Logger
}

func NewService(base *slog.Logger) *Service {
	return &Service{
		log: base.With("domain", "user"),
	}
}

func (s *Service) SignUp(ctx context.Context, username string) {
	s.log.InfoContext(ctx, "User sign-up requested", "username", username)
}

func (s *Service) SignUpSimple(username string) {
	s.log.Info("User sign-up requested (simple log)", "username", username)
}
