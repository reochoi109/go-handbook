package order

import "github.com/sirupsen/logrus"

type OrderService struct {
	log *logrus.Entry
}

func NewOrderService(baseLogger *logrus.Logger) *OrderService {
	return &OrderService{
		log: baseLogger.WithField("domain", "order"),
	}
}

func (s *OrderService) CreateOrder() {
	s.log.Info("주문 생성 프로세스 시작")
}
