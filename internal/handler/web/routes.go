package web

import (
	_ "embed"
	"errors"
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/changelog"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/handler/web/views"
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
	loader *changelog.Loader
	cfg    config.Config
	render Renderer
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
				args.Message = domErr.AppErr().Error()
				switch domErr.DomainErr() {
				case errs.ErrBadRequest:
					args.Status = http.StatusBadRequest
				case errs.ErrNotFound:
					args.Status = http.StatusNotFound
				case errs.ErrUnauthorized:
					args.Status = http.StatusUnauthorized
				case errs.ErrServiceUnavailable:
					args.Status = http.StatusServiceUnavailable
				}
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
