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
	s.log.InfoContext(ctx, "Order created", 
			slog.Group("order", 
				slog.String("id", id),
				slog.String("status", "created"),
		),
    )
}

func (s *Service) CreateSimple(id string) {
	s.log.Info("Order created (simple)", "order_id", id)
}
