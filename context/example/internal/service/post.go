package service

import (
	"context"
	"fmt"

	"github.com/reochoi109/go-handbook/context/example/internal/contextmeta"
	"github.com/reochoi109/go-handbook/context/example/internal/repository"
)

func GetUserPosts(ctx context.Context, userID int) (string, error) {
	fmt.Printf("[Trace: %s] Service\n", contextmeta.TraceIDFrom(ctx))
	return repository.FindPostsByUserID(ctx, userID)
}
