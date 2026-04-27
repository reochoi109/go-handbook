package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/reochoi109/go-handbook/log/slog/example/internal/order"
	"github.com/reochoi109/go-handbook/log/slog/example/internal/user"
	"github.com/reochoi109/go-handbook/log/slog/example/pkg/logger"
)

func main() {
	logLevel := slog.LevelInfo
	if os.Getenv("APP_ENV") == "development" {
		logLevel = slog.LevelDebug
	}

	baseLogger := logger.New(logLevel)
	slog.SetDefault(baseLogger)

	baseLogger.Info("App Start.",
		"version", "1.0.0",
		"env", os.Getenv("APP_ENV"),
	)

	// 로거를 각 도메인 서비스에 주입
	orderSvc := order.NewService(baseLogger)
	userSvc := user.NewService(baseLogger)

	// 비즈니스 로직 실행
	ctx := context.Background()

	userSvc.SignUp(ctx, "Reo")
	orderSvc.Create(ctx, "ORD-2026-0427")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	baseLogger.Info("App Shutdown.")

	time.Sleep(500 * time.Millisecond)
}
