package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/reochoi109/go-handbook/context/example/internal/contextmeta"
	"github.com/reochoi109/go-handbook/context/example/internal/service"
)

func GetUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	startedAt := time.Now()

	timeout := parseDurationMillis(r, "timeout_ms", 800)
	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	traceID := r.Header.Get("X-Trace-Id")
	if traceID == "" {
		traceID = "REQ-1234"
	}
	ctx = contextmeta.WithTraceID(ctx, traceID)
	w.Header().Set("X-Trace-Id", contextmeta.TraceIDFrom(ctx))

	delay := parseDurationMillis(r, "delay_ms", 3000)
	posts, err := service.GetUserPosts(ctx, 1, delay)
	if err != nil {
		log.Printf("trace_id=%s request failed: err=%v elapsed=%s", contextmeta.TraceIDFrom(ctx), err, time.Since(startedAt))
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			http.Error(w, fmt.Sprintf("timeout (trace_id=%s)", contextmeta.TraceIDFrom(ctx)), http.StatusRequestTimeout)
		case errors.Is(err, context.Canceled):
			http.Error(w, fmt.Sprintf("canceled (trace_id=%s)", contextmeta.TraceIDFrom(ctx)), http.StatusRequestTimeout)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(posts))
}

func parseDurationMillis(r *http.Request, key string, fallback int) time.Duration {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return time.Duration(fallback) * time.Millisecond
	}
	ms, err := strconv.Atoi(raw)
	if err != nil || ms <= 0 {
		return time.Duration(fallback) * time.Millisecond
	}
	return time.Duration(ms) * time.Millisecond
}
