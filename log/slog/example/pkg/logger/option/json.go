package option

import (
	"log/slog"
	"os"
)

func JsonFormat() {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return slog.String("timestamp", a.Value.Time().Format("2006-01-02 15:04:05.000"))
			case slog.LevelKey:
				return slog.Attr{Key: "level", Value: a.Value}
			case slog.MessageKey:
				return slog.Attr{Key: "message", Value: a.Value}
			case slog.SourceKey:
				return slog.Attr{Key: "caller", Value: a.Value}
			}
			return a
		},
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}
