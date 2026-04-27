package order

import (
	"context"

	"github.com/sirupsen/logrus"
)

type OrderService struct {
	log *logrus.Entry
}

func NewOrderService(baseLogger *logrus.Logger) *OrderService {
	return &OrderService{
		log: baseLogger.WithField("domain", "order"),
	}
}

func (s *OrderService) CreateOrderWithContext(ctx context.Context) {
	traceID, _ := ctx.Value("trace_id").(string)

	s.log.WithFields(logrus.Fields{
		"trace_id": traceID,
	}).Info("Order creation process started (with context)")
}

func (s *OrderService) CreateOrder() {
	s.log.Info("Order creation process started")
}
