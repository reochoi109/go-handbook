package repository

import (
	"context"
	"time"
)

func FindPostsByUserID(ctx context.Context, userID int) (string, error) {
	select {
	case <-time.After(3 * time.Second):
		return "Post List", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
