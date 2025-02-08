package logging

import (
	"context"
	"log/slog"
	"os"
)

type LevelHandler struct {
	level   slog.Leveler
	handler slog.Handler
}

func NewLevelHandler(level slog.Leveler, h slog.Handler) *LevelHandler {
	if lh, ok := h.(*LevelHandler); ok {
		h = lh.Handler()
	}
	return &LevelHandler{level, h}
}

func (h *LevelHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *LevelHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.handler.Handle(ctx, r)
}

func (h *LevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewLevelHandler(h.level, h.handler.WithAttrs(attrs))
}

func (h *LevelHandler) WithGroup(name string) slog.Handler {
	return NewLevelHandler(h.level, h.handler.WithGroup(name))
}

func (h *LevelHandler) Handler() slog.Handler {
	return h.handler
}

func strToLevel(val string) slog.Level {
	switch val {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "info":
		return slog.LevelInfo
	}
	return slog.LevelWarn
}

func NewDefaultLogger() *slog.Logger {
	level := slog.LevelWarn
	if val, ok := os.LookupEnv("GO_LOG_LEVEL"); ok {
		level = strToLevel(val)
	}
	return NewLogger(level)
}

func NewLogger(level slog.Leveler) *slog.Logger {
	h := slog.Default().Handler()
	return slog.New(NewLevelHandler(level, h))
}
