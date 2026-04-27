package order

import (
	"context"
	"log/slog"
)

type Service struct {
	log *slog.Logger
}

func NewService(base *slog.Logger) *Service {
	return &Service{log: base.With("domain", "order")}
}

func (s *Service) Create(ctx context.Context, id string) {
	if ctx == nil {
		ctx = context.Background()
	}
	s.log.InfoContext(ctx, "Order created", "order_id", id)
}

func (s *Service) CreateSimple(id string) {
	s.log.Info("Order created (simple)", "order_id", id)
}
