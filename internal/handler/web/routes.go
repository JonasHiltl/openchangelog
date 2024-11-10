package web

import (
	_ "embed"
	"errors"
	"log"
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/analytics"
	"github.com/jonashiltl/openchangelog/internal/analytics/tinybird"
	"github.com/jonashiltl/openchangelog/internal/changelog"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/handler/web/views"
	"github.com/jonashiltl/openchangelog/internal/store"
)

//go:embed static/base.css
var baseCSS string

func RegisterWebHandler(mux *http.ServeMux, e *env) {
	mux.HandleFunc("GET /", serveHTTP(e, index))
	mux.HandleFunc("POST /password", serveHTTP(e, passwordSubmit))
}

func NewEnv(
	cfg config.Config,
	loader *changelog.Loader,
	render Renderer,
) *env {
	return &env{
		cfg:    cfg,
		loader: loader,
		render: render,
	}
}

type env struct {
	loader  *changelog.Loader
	cfg     config.Config
	render  Renderer
	emitter analytics.Emitter
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
			log.Println("Tinybird analytics is enabled, but the 'analytics.tinybird' config section is missing")
			return analytics.NewNoopEmitter()
		}
		return tinybird.New(tinybird.TinybirdOptions{
			AccessToken: cfg.Analytics.Tinybird.AccessToken,
		})
	}

	return analytics.NewNoopEmitter()
}

func serveHTTP(env *env, h func(e *env, w http.ResponseWriter, r *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
				CSS:     baseCSS,
			}

			var domErr errs.Error
			if errors.As(err, &domErr) {
				args.Message = domErr.Msg()
				args.Status = domErr.Status()
			}

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
	}
}
