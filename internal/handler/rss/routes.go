package rss

import (
	"encoding/xml"
	"errors"
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/load"
	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/xlog"
)

type env struct {
	cfg    config.Config
	loader *load.Loader
	parser parse.Parser
}

func NewEnv(cfg config.Config, loader *load.Loader, parser parse.Parser) *env {
	return &env{
		cfg:    cfg,
		loader: loader,
		parser: parser,
	}
}

func RegisterRSSHandler(mux *http.ServeMux, e *env) {
	mux.HandleFunc("GET /feed", serveHTTP(e, feedHandler))
}

func serveHTTP(env *env, h func(e *env, w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return xlog.AttachLogger(func(w http.ResponseWriter, r *http.Request) {
		err := h(env, w, r)
		if err != nil {
			status := http.StatusInternalServerError
			msg := err.Error()

			var domErr errs.Error
			if errors.As(err, &domErr) {
				msg = domErr.AppErr().Error()
				switch domErr.DomainErr() {
				case errs.ErrBadRequest:
					status = http.StatusBadRequest
				case errs.ErrNotFound:
					status = http.StatusNotFound
				case errs.ErrUnauthorized:
					status = http.StatusUnauthorized
				case errs.ErrServiceUnavailable:
					status = http.StatusServiceUnavailable
				}
			}

			type XMLError struct {
				XMLName xml.Name `xml:"xml"`
				Message string   `xml:"string"`
				Code    int      `xml:"code"`
			}

			res := XMLError{
				Message: msg,
				Code:    status,
			}
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(status)
			err := xml.NewEncoder(w).Encode(res)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			xlog.LogRequest(r.Context(), status, msg)
		}
	})
}
