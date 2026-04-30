package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"
)

func main() {
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}

	// 목적지 정의
	out := slog.NewJSONHandler(os.Stdout, opts)    // 정상 로그용
	errOut := slog.NewJSONHandler(os.Stderr, opts) // 에러 로그용

	// 핸들러 체이닝 및 라우팅 구성
	h := &levelRouterHandler{
		low: &filterHandler{
			next: out,
			drop: func(_ context.Context, r slog.Record) bool {
				return hasAttr(r, "tag", "noise")
			},
		},
		high: errOut,
		cut:  slog.LevelError,
	}

	log := slog.New(h)

	log.Info("startup", "time", time.Now().Format(time.RFC3339Nano))
	log.Debug("polling", "tag", "noise", "iteration", 1)
	log.Debug("cache miss", "tag", "cache", "key", "u:123")
	log.Error("db error", "code", "E_DB", "detail", "connection refused")

}

type filterHandler struct {
	next slog.Handler
	drop func(ctx context.Context, r slog.Record) bool
}

func (h *filterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *filterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &filterHandler{next: h.next.WithAttrs(attrs), drop: h.drop}
}

func (h *filterHandler) WithGroup(name string) slog.Handler {
	return &filterHandler{next: h.next.WithGroup(name), drop: h.drop}
}

func (h *filterHandler) Handle(ctx context.Context, r slog.Record) error {
	if h.drop != nil && h.drop(ctx, r) {
		return nil // 출력 안함
	}
	return h.next.Handle(ctx, r)
}

type levelRouterHandler struct {
	low  slog.Handler
	high slog.Handler
	cut  slog.Level
}

func (h *levelRouterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if level >= h.cut {
		return h.high.Enabled(ctx, level)
	}
	return h.low.Enabled(ctx, level)
}

func (h *levelRouterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &levelRouterHandler{
		low:  h.low.WithAttrs(attrs),
		high: h.high.WithAttrs(attrs),
		cut:  h.cut,
	}
}

func (h *levelRouterHandler) WithGroup(name string) slog.Handler {
	return &levelRouterHandler{
		low:  h.low.WithGroup(name),
		high: h.high.WithGroup(name),
		cut:  h.cut,
	}
}

func (h *levelRouterHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level >= h.cut {
		return h.high.Handle(ctx, r) // 중요 로그는 high
	}
	return h.low.Handle(ctx, r) // 일반 로그 low
}

func hasAttr(r slog.Record, key string, want string) bool {
	found := false
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == key && strings.Contains(a.Value.String(), want) {
			found = true
			return false // 속성 찾았음으로 반복 종료 (false)
		}
		return true // 계속 반속
	})
	return found
}

var _ slog.Handler = (*filterHandler)(nil)
var _ slog.Handler = (*levelRouterHandler)(nil)
