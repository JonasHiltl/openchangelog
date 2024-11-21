package web

import (
	"errors"
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/analytics"
	"github.com/jonashiltl/openchangelog/internal/analytics/tinybird"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/handler/web/static"
	"github.com/jonashiltl/openchangelog/internal/handler/web/views"
	"github.com/jonashiltl/openchangelog/internal/load"
	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/search"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/internal/xlog"
	"golang.org/x/exp/slog"
)

func RegisterWebHandler(mux *http.ServeMux, e *env) {
	mux.HandleFunc("GET /", serveHTTP(e, index))
	mux.HandleFunc("POST /password", serveHTTP(e, passwordSubmit))
	mux.HandleFunc("POST /search", serveHTTP(e, searchSubmit))
	mux.HandleFunc("GET /search/tags", serveHTTP(e, searchTags))
}

func NewEnv(
	cfg config.Config,
	loader *load.Loader,
	parser parse.Parser,
	render Renderer,
	searcher search.Searcher,
) *env {
	return &env{
		cfg:      cfg,
		loader:   loader,
		parser:   parser,
		render:   render,
		searcher: searcher,
	}
}

type env struct {
	cfg      config.Config
	render   Renderer
	emitter  analytics.Emitter
	loader   *load.Loader
	parser   parse.Parser
	searcher search.Searcher
}

// Returns the analytics emitter of the changelog.
func (e *env) getAnalyticsEmitter(cl store.Changelog) analytics.Emitter {
	if !cl.Analytics {
		return analytics.NewNoopEmitter()
	}

	// we cache the emitter, since some need global state
	// like the tinybird one for batching
	if e.emitter == nil {
		t := createEmitter(e.cfg)
		e.emitter = t
	}
	return e.emitter
}

func createEmitter(cfg config.Config) analytics.Emitter {
	if cfg.Analytics == nil {
		return analytics.NewNoopEmitter()
	}

	switch cfg.Analytics.Provider {
	case config.Tinybird:
		if cfg.Analytics.Tinybird == nil {
			slog.Warn("Tinybird analytics is enabled, but the 'analytics.tinybird' config section is missing")
			return analytics.NewNoopEmitter()
		}
		return tinybird.New(tinybird.TinybirdOptions{
			AccessToken: cfg.Analytics.Tinybird.AccessToken,
		})
	}

	return analytics.NewNoopEmitter()
}

func serveHTTP(env *env, h func(e *env, w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return xlog.AttachLogger(func(w http.ResponseWriter, r *http.Request) {
		err := h(env, w, r)
		if err != nil {
			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path += "?" + r.URL.RawQuery
			}

			var args views.ErrorArgs = views.ErrorArgs{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
				Path:    path,
				CSS:     static.BaseCSS,
			}

			var domErr errs.Error
			if errors.As(err, &domErr) {
				args.Message = domErr.Msg()
				args.Status = domErr.Status()
			}

			defer xlog.LogRequest(r.Context(), args.Status, args.Message)

			// if requesting widget, don't render html error, just error message
			if _, ok := r.URL.Query()["widget"]; ok {
				http.Error(w, args.Message, args.Status)
				return
			}

			w.WriteHeader(args.Status)
			err := views.Error(args).Render(r.Context(), w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})
}
