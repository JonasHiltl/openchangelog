package xlog

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
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
)

// Attaches logger attributes to the handler.
func AttachLogger(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = xid.New().String()
		}
		addRequestID(r, requestID)
		addRequestPath(r, r.URL)
		addRequestMethod(r, r.Method)
		addRequestStart(r)
		fn(w, r)
	}
}

// Attaches a workspace to the request context.
// Useful for logging.
func AddWorkspaceID(r *http.Request, id string) {
	*r = *r.WithContext(context.WithValue(r.Context(), ctxKeyWorkspaceID, id))
}

func addRequestID(r *http.Request, id string) {
	*r = *r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id))
}

func addRequestPath(r *http.Request, url *url.URL) {
	*r = *r.WithContext(context.WithValue(r.Context(), ctxKeyRequestPath, url.String()))
}

func addRequestMethod(r *http.Request, method string) {
	*r = *r.WithContext(context.WithValue(r.Context(), ctxKeyRequestMethod, method))
}

func addRequestStart(r *http.Request) {
	*r = *r.WithContext(context.WithValue(r.Context(), ctxKeyRequestStart, time.Now()))
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
