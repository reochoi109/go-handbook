package contextmeta

import "context"

type ctxKey string

const (
	traceIDKey ctxKey = "trace_id"
	userIDKey  ctxKey = "user_id"
)

func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceIDKey, id)
}

func TraceIDFrom(ctx context.Context) string {
	v, _ := ctx.Value(traceIDKey).(string)
	if v == "" {
		return "-"
	}
	return v
}

func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func UserIDFrom(ctx context.Context) string {
	v, _ := ctx.Value(userIDKey).(string)
	return v
}
