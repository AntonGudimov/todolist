package logging

import (
	"context"
	"log/slog"
	"todolist/internal/constants"
)

type ContextJsonHandler struct {
	slog.Handler
}

func (h *ContextJsonHandler) Handle(
	ctx context.Context,
	r slog.Record,
) error {
	requestId, _ := ctx.Value(constants.RequestIDKey).(string)
	if requestId != "" {
		r.AddAttrs(slog.String(string(constants.RequestIDKey), requestId))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *ContextJsonHandler) WithAttrs(
	attrs []slog.Attr,
) slog.Handler {
	return &ContextJsonHandler{
		Handler: h.Handler.WithAttrs(attrs),
	}
}

func (h *ContextJsonHandler) WithGroup(
	name string,
) slog.Handler {
	return &ContextJsonHandler{
		Handler: h.Handler.WithGroup(name),
	}
}
