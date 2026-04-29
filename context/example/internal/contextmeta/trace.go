package contextmeta

import "context"

type key int

const (
	keyTraceID key = iota
	keyUserID
)

func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyTraceID, id)
}

func TraceIDFrom(ctx context.Context) string {
	v, _ := ctx.Value(keyTraceID).(string)
	return v
}

func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyUserID, id)
}

func UserIDFrom(ctx context.Context) string {
	v, _ := ctx.Value(keyUserID).(string)
	return v
}
