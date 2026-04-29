package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/reochoi109/go-handbook/log/slog/example/internal/order"
	"github.com/reochoi109/go-handbook/log/slog/example/internal/user"
	"github.com/reochoi109/go-handbook/log/slog/example/pkg/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logLevel := slog.LevelInfo
	if os.Getenv("APP_ENV") == "development" {
		logLevel = slog.LevelDebug
	}

	baseLogger := logger.New(logLevel)
	slog.SetDefault(baseLogger)

	baseLogger.InfoContext(ctx, "App Start.",
		"version", "1.0.0",
		"env", os.Getenv("APP_ENV"),
		"log_format", os.Getenv("LOG_FORMAT"),
	)

	// 로거를 각 도메인 서비스에 주입
	orderSvc := order.NewService(baseLogger)
	userSvc := user.NewService(baseLogger)

	// 비즈니스 로직 실행
	userSvc.SignUp(ctx, "Reo")
	orderSvc.Create(ctx, "ORD-2026-0427")

	if os.Getenv("EXIT_AFTER") == "1" {
		baseLogger.InfoContext(ctx, "App Exit (demo).")
		return
	}

	<-ctx.Done()
	baseLogger.InfoContext(context.Background(), "App Shutdown.", "reason", ctx.Err())
}
