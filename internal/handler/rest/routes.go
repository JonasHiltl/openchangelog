package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/changelog"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/store"
)

func RegisterRestHandler(mux *http.ServeMux, e *env) {
	// Workspace
	mux.HandleFunc("POST /api/workspaces", serveHTTP(e, createWorkspace))
	mux.HandleFunc("GET /api/workspaces/my", serveHTTP(e, getMyWorkspace))
	mux.HandleFunc("GET /api/workspaces/{wid}", serveHTTP(e, getWorkspace))
	mux.HandleFunc("PATCH /api/workspaces/{wid}", serveHTTP(e, updateWorkspace))
	mux.HandleFunc("DELETE /api/workspaces/{wid}", serveHTTP(e, deleteWorkspace))

	// Sources
	mux.HandleFunc("GET /api/sources", serveHTTP(e, listSources))

	// GH sources
	mux.HandleFunc("POST /api/sources/gh", serveHTTP(e, createGHSource))
	mux.HandleFunc("GET /api/sources/gh", serveHTTP(e, listGHSources))
	mux.HandleFunc("GET /api/sources/gh/{id}", serveHTTP(e, getGHSource))
	mux.HandleFunc("DELETE /api/sources/gh/{id}", serveHTTP(e, deleteGHSources))

	// changelog
	mux.HandleFunc("POST /api/changelogs", serveHTTP(e, createChangelog))
	mux.HandleFunc("GET /api/changelogs", serveHTTP(e, listChangelogs))
	mux.HandleFunc("GET /api/changelogs/{cid}", serveHTTP(e, getChangelog))
	mux.HandleFunc("GET /api/changelogs/{cid}/full", serveHTTP(e, getFullChangelog))
	mux.HandleFunc("PATCH /api/changelogs/{cid}", serveHTTP(e, updateChangelog))
	mux.HandleFunc("DELETE /api/changelogs/{cid}", serveHTTP(e, deleteChangelog))
	mux.HandleFunc("PUT /api/changelogs/{cid}/source/{sid}", serveHTTP(e, setChangelogSource))
	mux.HandleFunc("DELETE /api/changelogs/{cid}/source", serveHTTP(e, deleteChangelogSource))
}

func NewEnv(store store.Store, loader *changelog.Loader) *env {
	return &env{
		store:  store,
		loader: loader,
	}
}

type env struct {
	store  store.Store
	loader *changelog.Loader
}

func serveHTTP(env *env, h func(e *env, w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

			res := map[string]any{
				"message": msg,
				"code":    status,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			err := json.NewEncoder(w).Encode(res)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
