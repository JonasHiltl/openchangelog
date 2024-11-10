package lgr

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/rs/xid"
)

type contextKey string

func (c contextKey) String() string {
	return "handler context key " + string(c)
}

var (
	ctxKeyRequestID   = contextKey("request-id")
	ctxKeyWorkspaceID = contextKey("workspace-id")
)

// Attaches logger attributes to the handler.
func AttachLogger(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := xid.New().String()
		AddRequestID(r, requestID)
		fn(w, r)
	}
}

// Attaches a workspace to the request context.
// Useful for logging.
func AddWorkspaceID(r *http.Request, id string) {
	*r = *r.WithContext(context.WithValue(r.Context(), ctxKeyWorkspaceID, id))
}

// Attaches a request id to the request context.
// Useful for logging.
func AddRequestID(r *http.Request, id string) {
	*r = *r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id))
}

func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}
