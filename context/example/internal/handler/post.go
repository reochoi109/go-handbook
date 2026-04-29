package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/reochoi109/go-handbook/context/example/internal/contextmeta"
	"github.com/reochoi109/go-handbook/context/example/internal/service"
)

func GetUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*800)
	defer cancel()

	traceID := r.Header.Get("X-Trace-Id")
	if traceID == "" {
		traceID = "REQ-1234"
	}
	ctx = contextmeta.WithTraceID(ctx, traceID)

	posts, err := service.GetUserPosts(ctx, 1)
	if err != nil {
		switch {
		case err == context.DeadlineExceeded:
			http.Error(w, "timeout", http.StatusRequestTimeout)
		case err == context.Canceled:
			http.Error(w, "canceled", http.StatusRequestTimeout)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Trace-Id", contextmeta.TraceIDFrom(ctx))
	_, _ = w.Write([]byte(posts))
}
