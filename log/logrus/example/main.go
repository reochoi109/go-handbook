package main

import (
	"github.com/reochoi109/go-handbook/log/logrus/example/internal/order"
	"github.com/reochoi109/go-handbook/log/logrus/example/internal/user"
	"github.com/reochoi109/go-handbook/log/logrus/example/pkg/logger"
)

func main() {
	baseLog := logger.New()

	orderSvc := order.NewOrderService(baseLog)
	userSvc := user.NewUserService(baseLog)

	orderSvc.CreateOrder()
	userSvc.SignUp("username")
}
