package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/reochoi109/go-handbook/error/example/internal/repository"
	"github.com/reochoi109/go-handbook/error/example/internal/service"
	"github.com/reochoi109/go-handbook/error/example/internal/web"
)

func main() {
	repo := repository.NewInMemoryUserRepo()
	svc := service.NewUserService(repo)
	h := web.NewHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", h.GetMe)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()

	log.Println("listening on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Println("server error:", err)
		os.Exit(1)
	}
}
