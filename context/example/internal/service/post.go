package service

import (
	"context"
	"fmt"
	"time"

	"github.com/reochoi109/go-handbook/context/example/internal/contextmeta"
	"github.com/reochoi109/go-handbook/context/example/internal/repository"
)

func GetUserPosts(ctx context.Context, userID int, delay time.Duration) (string, error) {
	fmt.Printf("[trace_id=%s] service: start\n", contextmeta.TraceIDFrom(ctx))
	return repository.FindPostsByUserID(ctx, userID, delay)
}
