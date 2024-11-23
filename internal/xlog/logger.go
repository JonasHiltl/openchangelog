package xlog

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jonashiltl/openchangelog/internal/config"
)

func NewLogger(cfg config.Config) *slog.Logger {
	var sh slog.Handler

	if cfg.Log == nil {
		sh = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	} else {
		switch cfg.Log.Style {
		case "json":
			sh = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.Log.Level.ToSlog()})
		default:
			sh = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.Log.Level.ToSlog()})
		}
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
	if wID, ok := ctx.Value(ctxKeyWorkspaceID).(string); ok {
		r.AddAttrs(slog.String("workspace_id", wID))
	}
	if host, ok := ctx.Value(ctxKeyRequestHost).(string); ok {
		r.AddAttrs(slog.String("host", host))
	}
	if rPath, ok := ctx.Value(ctxKeyRequestPath).(string); ok {
		r.AddAttrs(slog.String("path", rPath))
	}
	if rMth, ok := ctx.Value(ctxKeyRequestMethod).(string); ok {
		r.AddAttrs(slog.String("method", rMth))
	}
	if start, ok := ctx.Value(ctxKeyRequestStart).(time.Time); ok {
		r.AddAttrs(slog.Duration("duration", time.Since(start)))
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
