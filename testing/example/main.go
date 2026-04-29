package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/reochoi109/go-handbook/testing/example/internal/domain"
	"github.com/reochoi109/go-handbook/testing/example/internal/service"
	"github.com/reochoi109/go-handbook/testing/example/internal/web"
	"github.com/reochoi109/go-handbook/testing/example/repository"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	repo := repository.NewInMemoryUserRepo([]domain.User{
		{ID: "user_123", Name: "Reo"},
		{ID: "user_456", Name: "Alex"},
	})
	svc := service.NewUserService(repo)
	h := web.NewHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/users/", h.GetUser)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	log.Println("listening on :8080")
	log.Println(`try: curl -i localhost:8080/users/user_123`)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Println("server error:", err)
		os.Exit(1)
	}
}
