package main

import (
	"context"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

type piiMaskHandler struct {
	next slog.Handler
}

func main(){
	// 1. base logger
	base:= slog.NewJSONHandler(os.Stdout,&slog.HandlerOptions{Level: slog.LevelInfo})
	
	// 2. custom logger
	log := slog.New(&piiMaskHandler{next: base})

	ctx:=context.Background()
	log.InfoContext(ctx, "signup",
		"user_id", "u_123",
		"email", "user@example.com",
		"token", "secret-token-value",
		"note", "email=user@example.com token=abc123",
	)
}

// level check
func (h *piiMaskHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *piiMaskHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &piiMaskHandler{next: h.next.WithAttrs(attrs)}
}

func (h *piiMaskHandler) WithGroup(name string) slog.Handler {
	return &piiMaskHandler{next: h.next.WithGroup(name)}
}

// record method
func (h *piiMaskHandler) Handle(ctx context.Context, r slog.Record)error{
	out := slog.NewRecord(r.Time,r.Level,r.Message,r.PC)
	
	// 모든 속성 순회하며 마스킹 처리 후 새 레코드에 추가
	r.Attrs(func(a slog.Attr) bool {
		out.AddAttrs(maskAttr(a)) // 변조
		return  true
	})
	return h.next.Handle(ctx , out)
}


var _ slog.Handler = (*piiMaskHandler)(nil)

func maskAttr(a slog.Attr) slog.Attr{
	switch strings.ToLower(a.Key) {
	case "email":
		return slog.String(a.Key, maskEmail(a.Value.String()))
	case "token", "password", "secret":
		return slog.String(a.Key, "***")
	case "note":
		return slog.String(a.Key, maskInline(a.Value.String()))
	default:
		return a
	}
}

func maskEmail(s string) string {
	at := strings.IndexByte(s, '@')
	if at <= 1 {
		return "***"
	}
	return s[:1] + "***" + s[at:]
}


var (
	reEmail = regexp.MustCompile(`[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}`)
	reToken = regexp.MustCompile(`(?i)\btoken=([A-Za-z0-9._\-]+)\b`)
)

func maskInline(s string) string {
	s = reEmail.ReplaceAllStringFunc(s, maskEmail) // email pattern
	s = reToken.ReplaceAllString(s, "token=***") // token pattern
	return s
}
