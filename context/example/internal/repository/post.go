package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/reochoi109/go-handbook/context/example/internal/contextmeta"
)

func FindPostsByUserID(ctx context.Context, userID int, delay time.Duration) (string, error) {
	fmt.Printf("[trace_id=%s] repo: query start (delay=%s)\n", contextmeta.TraceIDFrom(ctx), delay)
	select {
	case <-time.After(delay):
		fmt.Printf("[trace_id=%s] repo: query done\n", contextmeta.TraceIDFrom(ctx))
		return "Post List", nil
	case <-ctx.Done():
		fmt.Printf("[trace_id=%s] repo: canceled (%v)\n", contextmeta.TraceIDFrom(ctx), ctx.Err())
		return "", ctx.Err()
	}
}
