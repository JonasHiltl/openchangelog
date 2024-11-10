package admin

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/store"
)

func RegisterAdminHandler(mux *http.ServeMux, e *env) {
	if e.cfg.Admin == nil {
		return
	}

	slog.Info("admin view is enabled at /admin")
	mux.HandleFunc("GET /admin", serveHTTP(e, adminOverview))
	mux.HandleFunc("GET /admin/{wid}", serveHTTP(e, details))
}

func NewEnv(cfg config.Config, st store.Store) *env {
	return &env{
		cfg: cfg,
		st:  st,
	}
}

type env struct {
	cfg config.Config
	st  store.Store
}

func serveHTTP(env *env, h func(e *env, w http.ResponseWriter, r *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(env, w, r)
		if err != nil {
			var domErr errs.Error
			if errors.As(err, &domErr) {
				http.Error(w, domErr.Msg(), domErr.Status())
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
