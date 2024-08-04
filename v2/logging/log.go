package logging

import (
	"context"
	"log/slog"
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

func NewDefaultLogger(level slog.Leveler) *slog.Logger {
	h := slog.Default().Handler()
	return slog.New(NewLevelHandler(level, h))
}
