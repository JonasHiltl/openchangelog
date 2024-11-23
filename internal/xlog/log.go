package xlog

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/rs/xid"
)

type contextKey string

func (c contextKey) String() string {
	return "handler context key " + string(c)
}

var (
	ctxKeyRequestID     = contextKey("request-id")
	ctxKeyWorkspaceID   = contextKey("workspace-id")
	ctxKeyRequestPath   = contextKey("request-path")
	ctxKeyRequestMethod = contextKey("request-method")
	ctxKeyRequestStart  = contextKey("request-start")
	ctxKeyRequestHost   = contextKey("request-host")
)

// Attaches logger attributes to the handler.
func AttachLogger(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addRequestID(r)
		addRequestPath(r)
		addRequestMethod(r)
		addRequestStart(r)
		addHost(r)
		fn(w, r)
	}
}

// Attaches a workspace to the request context.
// Useful for logging.
func AddWorkspaceID(r *http.Request, id string) {
	*r = *r.WithContext(context.WithValue(r.Context(), ctxKeyWorkspaceID, id))
}

func addRequestID(r *http.Request) {
	id := r.Header.Get("X-Request-ID")
	if id == "" {
		id = xid.New().String()
	}
	*r = *r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id))
}

func addRequestPath(r *http.Request) {
	*r = *r.WithContext(context.WithValue(r.Context(), ctxKeyRequestPath, r.URL.String()))
}

func addRequestMethod(r *http.Request) {
	*r = *r.WithContext(context.WithValue(r.Context(), ctxKeyRequestMethod, r.Method))
}

func addRequestStart(r *http.Request) {
	*r = *r.WithContext(context.WithValue(r.Context(), ctxKeyRequestStart, time.Now()))
}

func addHost(r *http.Request) {
	host := r.Host
	if r.Header.Get("X-Forwarded-Host") != "" {
		host = r.Header.Get("X-Forwarded-Host")
	}
	*r = *r.WithContext(context.WithValue(r.Context(), ctxKeyRequestHost, host))
}

func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}

func LogRequest(ctx context.Context, status int, msg string) {
	level := slog.LevelError
	if status < 400 {
		level = slog.LevelDebug
	}

	slog.LogAttrs(ctx, level, msg, slog.Int("status", status))
}
