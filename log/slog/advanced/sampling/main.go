package main

import (
	"context"
	"log/slog"
	"os"
	"sync/atomic"
	"time"
)

type samplingHandler struct {
	next   slog.Handler
	everyN uint64

	minLevel slog.Level
	maxLevel slog.Level

	keepAbove slog.Level
	seq       *atomic.Uint64
}

func main() {
	base := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})

	h := &samplingHandler{
		next:      base,            // 최종 출력을 담당할 핸들러
		everyN:    10,              // 10번째 로그마다 1개씩 출력 (나머지는 드랍)
		minLevel:  slog.LevelInfo,  // 샘플링을 적용할 최소 레벨
		maxLevel:  slog.LevelInfo,  // 샘플링을 적용할 최대 레벨 (여기서는 INFO만 대상)
		keepAbove: slog.LevelError, // 이 레벨(ERROR) 이상은 샘플링 없이 무조건 통과
		seq:       new(atomic.Uint64),
	}

	log := slog.New(h)
	ctx := context.Background()
	for i := 0; i < 35; i++ {
		log.InfoContext(ctx, "heartbeat", "i", i)

		if i%17 == 0 {
			log.ErrorContext(ctx, "db error", "i", i, "code", "E_DB")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func (h *samplingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *samplingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	cp := *h
	cp.next = h.next.WithAttrs(attrs)
	return &cp
}

func (h *samplingHandler) WithGroup(name string) slog.Handler {
	cp := *h
	cp.next = h.next.WithGroup(name)
	return &cp
}

func (h *samplingHandler) Handle(ctx context.Context, r slog.Record) error {
	// 1. 중요한 로그는 샘플링 로직타지 않고 즉시 출력
	if r.Level >= h.keepAbove {
		return h.next.Handle(ctx, r)
	}

	// 2. 샘플링 대상 레벨 범위가 아니면 출력
	if r.Level < h.minLevel || r.Level > h.maxLevel || h.everyN == 0 {
		return h.next.Handle(ctx, r)
	}

	// 3. 카운터 증가
	n := h.seq.Add(1)

	// 4. n번째로그가 아니면 드랍
	if n%h.everyN != 0 {
		return nil
	}
	return h.next.Handle(ctx, r)
}

var _ slog.Handler = (*samplingHandler)(nil)
