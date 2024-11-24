package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btvoidx/mint"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/load"
	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/internal/xlog"
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

func NewEnv(store store.Store, loader *load.Loader, parser parse.Parser, e *mint.Emitter) *env {
	return &env{
		store:  store,
		loader: loader,
		parser: parser,
		e:      e,
	}
}

type env struct {
	store  store.Store
	loader *load.Loader
	parser parse.Parser
	e      *mint.Emitter
}

func serveHTTP(env *env, h func(e *env, w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return xlog.AttachLogger(func(w http.ResponseWriter, r *http.Request) {
		err := h(env, w, r)

		if err != nil {
			status := http.StatusInternalServerError
			msg := err.Error()

			var domErr errs.Error
			if errors.As(err, &domErr) {
				msg = domErr.Msg()
				status = domErr.Status()
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

			xlog.LogRequest(r.Context(), status, msg)
		}
	})
}
