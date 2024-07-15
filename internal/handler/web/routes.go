package web

import (
	"errors"
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/handler/web/views"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/parse"
	"github.com/jonashiltl/openchangelog/render"
	"github.com/naveensrinivasan/httpcache"
)

func RegisterWebHandler(mux *http.ServeMux, e *env) {
	fs := http.FileServer(http.Dir("./internal/handler/web/public/"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("GET /", serveHTTP(e, index))
}

func NewEnv(
	cfg config.Config,
	store store.Store,
	render render.Renderer,
	parse parse.Parser,
	cache httpcache.Cache,
) *env {
	return &env{
		cfg:    cfg,
		store:  store,
		render: render,
		parse:  parse,
		cache:  cache,
	}
}

type env struct {
	cfg    config.Config
	store  store.Store
	render render.Renderer
	parse  parse.Parser
	cache  httpcache.Cache
}

type handler = func(e *env, w http.ResponseWriter, r *http.Request) error

func serveHTTP(env *env, h handler) func(http.ResponseWriter, *http.Request) {
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
				}
			}

			w.WriteHeader(args.Status)
			err := views.Error(args).Render(r.Context(), w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
