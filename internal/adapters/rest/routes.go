package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/domain/changelog"
	"github.com/jonashiltl/openchangelog/internal/domain/source"
	"github.com/jonashiltl/openchangelog/internal/domain/workspace"
)

type env struct {
	cfg          config.Config
	sourceSrv    source.Service
	workspaceSrv workspace.Service
	changelogSrv changelog.Service
}

type handler = func(e *env, w http.ResponseWriter, r *http.Request) error

func serveHTTP(env *env, h handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(env, w, r)

		if err != nil {
			var status int
			var msg string
			var statErr RestError
			if errors.As(err, &statErr) {
				status = statErr.Code
				msg = statErr.Error()
			} else {
				status = http.StatusInternalServerError
				msg = err.Error()
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

func initRoutes(mux *http.ServeMux, e *env) {
	// Workspace
	mux.HandleFunc("POST /workspaces", serveHTTP(e, createWorkspace))
	mux.HandleFunc("GET /workspaces/{wid}", serveHTTP(e, getWorkspace))
	mux.HandleFunc("PATCH /workspaces/{wid}", serveHTTP(e, updateWorkspace))
	mux.HandleFunc("DELETE /workspaces/{wid}", serveHTTP(e, deleteWorkspace))

	// GH sources
	mux.HandleFunc("POST /sources/gh", serveHTTP(e, createGHSource))
	mux.HandleFunc("GET /sources/gh", serveHTTP(e, listGHSources))
	mux.HandleFunc("GET /sources/gh/{sid}", serveHTTP(e, getGHSource))
	mux.HandleFunc("DELETE /sources/gh/{sid}", serveHTTP(e, deleteGHSources))

	// changelog
	mux.HandleFunc("POST /changelogs", serveHTTP(e, createChangelog))
	mux.HandleFunc("GET /changelogs", serveHTTP(e, listChangelogs))
	mux.HandleFunc("GET /changelogs/{cid}", serveHTTP(e, getChangelog))
	mux.HandleFunc("PATCH /changelogs/{cid}", serveHTTP(e, updateChangelog))
	mux.HandleFunc("DELETE /changelogs/{cid}", serveHTTP(e, deleteChangelog))
	mux.HandleFunc("PUT /changelogs/{cid}/source", serveHTTP(e, setChangelogSource))
}
