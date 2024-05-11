package routes

import (
	"errors"
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/adapters/web/views"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/domain"
	"github.com/jonashiltl/openchangelog/internal/domain/changelog"
	"github.com/jonashiltl/openchangelog/parse"
)

type env struct {
	cfg          config.Config
	parser       parse.Parser
	changelogSrv changelog.Service
}

func NewEnv(cfg config.Config, parser parse.Parser, changelogSrv changelog.Service) *env {
	e := &env{
		cfg:          cfg,
		parser:       parser,
		changelogSrv: changelogSrv,
	}
	return e
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

			var domErr domain.Error
			if errors.As(err, &domErr) {
				args.Message = domErr.AppErr().Error()
				switch domErr.DomainErr() {
				case domain.ErrBadRequest:
					args.Status = http.StatusBadRequest
				case domain.ErrNotFound:
					args.Status = http.StatusNotFound
				}
			}

			err := views.Error(args).Render(r.Context(), w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func Init(mux *http.ServeMux, e *env) {
	fs := http.FileServer(http.Dir("./internal/adapters/web/public/"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("GET /", serveHTTP(e, index))
}
