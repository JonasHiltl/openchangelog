package rest

import (
	"encoding/json"
	"net/http"

	"github.com/jonashiltl/openchangelog/apitypes"
	"github.com/jonashiltl/openchangelog/internal/store"
)

const (
	ghsource_id_param = "id"
)

func ghToApiType(gh store.GHSource) apitypes.GHSource {
	return apitypes.GHSource{
		ID:          gh.ID.String(),
		WorkspaceID: gh.WorkspaceID.String(),
		Owner:       gh.Owner,
		Repo:        gh.Repo,
		Path:        gh.Path,
	}
}

func encodeGHSource(w http.ResponseWriter, gh store.GHSource) error {
	g := ghToApiType(gh)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(g)
}

func listSources(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	// TODO: in future fetch other source types
	gh, err := e.store.ListGHSources(r.Context(), t.WorkspaceID)
	if err != nil {
		return err
	}

	res := make([]apitypes.GHSource, len(gh))
	for i, gh := range gh {
		g := ghToApiType(gh)
		res[i] = g
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

func createGHSource(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}
	var req struct {
		Name           string `json:"name"`
		Owner          string `json:"owner"`
		Repo           string `json:"repo"`
		Path           string `json:"path"`
		InstallationID int64  `json:"installationID"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	gh, err := e.store.CreateGHSource(r.Context(), store.GHSource{
		WorkspaceID:    t.WorkspaceID,
		ID:             store.NewGHID(),
		Owner:          req.Owner,
		Repo:           req.Repo,
		Path:           req.Path,
		InstallationID: req.InstallationID,
	})
	if err != nil {
		return err
	}

	return encodeGHSource(w, gh)
}

func getGHSource(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	ghID, err := store.ParseGHID(r.PathValue(ghsource_id_param))
	if err != nil {
		return err
	}

	gh, err := e.store.GetGHSource(r.Context(), t.WorkspaceID, ghID)
	if err != nil {
		return err
	}
	return encodeGHSource(w, gh)
}

func listGHSources(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	sources, err := e.store.ListGHSources(r.Context(), t.WorkspaceID)
	if err != nil {
		return err
	}
	res := make([]apitypes.GHSource, len(sources))
	for i, gh := range sources {
		g := ghToApiType(gh)
		res[i] = g
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

func deleteGHSources(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	ghID, err := store.ParseGHID(r.PathValue(ghsource_id_param))
	if err != nil {
		return err
	}

	err = e.store.DeleteGHSource(r.Context(), t.WorkspaceID, ghID)
	if err != nil {
		return err
	}
	return nil
}
