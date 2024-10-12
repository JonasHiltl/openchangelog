package web

import (
	"crypto/md5"
	"embed"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/changelog"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/handler/web/views"
	"github.com/jonashiltl/openchangelog/render"
)

//go:embed static/*
var staticAssets embed.FS

func RegisterWebHandler(mux *http.ServeMux, e *env) {
	fs := http.FileServer(http.FS(staticAssets))
	mux.Handle("GET /static/*", fs)
	mux.HandleFunc("GET /", serveHTTP(e, index))
	mux.HandleFunc("POST /password", serveHTTP(e, passwordSubmit))
}

func NewEnv(
	cfg config.Config,
	loader *changelog.Loader,
	render render.Renderer,
) *env {
	return &env{
		cfg:            cfg,
		loader:         loader,
		render:         render,
		baseCSSVersion: calculateBaseCSSVersion(),
	}
}

// Calculates the md5 hash of the static/base.css file.
// Needed for cache busting if base.css changes.
func calculateBaseCSSVersion() string {
	fileContent, err := staticAssets.ReadFile("static/base.css")
	if err != nil {
		return ""
	}

	hash := md5.New()
	hash.Write(fileContent)
	return hex.EncodeToString(hash.Sum(nil))
}

type env struct {
	loader         *changelog.Loader
	cfg            config.Config
	render         render.Renderer
	baseCSSVersion string
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

			w.WriteHeader(args.Status)
			err := views.Error(args).Render(r.Context(), w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
