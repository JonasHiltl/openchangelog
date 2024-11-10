package lgr

import (
	"context"
	"log/slog"
	"os"

	"github.com/jonashiltl/openchangelog/internal/config"
)

func NewLogger(cfg config.Config) *slog.Logger {
	var sh slog.Handler
	switch cfg.LogStyle {
	case "json":
		sh = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel.ToSlog()})
	default:
		sh = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel.ToSlog()})
	}
	return slog.New(&myHandler{
		Handler: sh,
	})
}

type myHandler struct {
	slog.Handler
}

func (h *myHandler) Handle(ctx context.Context, r slog.Record) error {
	if rID, ok := ctx.Value(ctxKeyRequestID).(string); ok {
		r.AddAttrs(slog.String("request_id", rID))
	}
	if rID, ok := ctx.Value(ctxKeyWorkspaceID).(string); ok {
		r.AddAttrs(slog.String("workspace_id", rID))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *myHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *myHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &myHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *myHandler) WithGroup(name string) slog.Handler {
	return &myHandler{Handler: h.Handler.WithGroup(name)}
}
