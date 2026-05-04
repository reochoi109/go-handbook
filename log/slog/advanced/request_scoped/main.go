package main

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey string

const (
	RequestIDKey ctxKey = "request_id"
	TraceIDKey   ctxKey = "trace_id"
)

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if rid, ok := ctx.Value(RequestIDKey).(string); ok {
		r.AddAttrs(slog.String("request_id", rid))
	}

	if tid, ok := ctx.Value(TraceIDKey).(string); ok {
		r.AddAttrs(slog.String("trace_id", tid))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{Handler: h.Handler.WithGroup(name)}
}

func AppendCtx(parent context.Context, rid, tid string) context.Context {
	ctx := context.WithValue(parent, RequestIDKey, rid)
	return context.WithValue(ctx, TraceIDKey, tid)
}

func main() {
	innerHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(&ContextHandler{innerHandler})

	// 요청 컨텍스트 생성
	ctx := AppendCtx(context.Background(), "REQ-12345", "TRACE-67890")

	logger.InfoContext(ctx, "user login attempt", slog.String("user_email", "reo@example.com"))
	performBusinessLogic(ctx, logger)
}

func performBusinessLogic(ctx context.Context, log *slog.Logger) {
	// 비즈니스 로직 함수는 내부에서 request_id나 trace_id를 전혀 몰라도 된다.
	log.DebugContext(ctx, "database query completed")
}
