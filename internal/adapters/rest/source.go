package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jonashiltl/openchangelog/apitypes"
	"github.com/jonashiltl/openchangelog/internal/domain/source"
)

func ghToApiType(gh source.GHSource) apitypes.GHSource {
	return apitypes.GHSource{
		ID:          gh.ID.ToInt(),
		WorkspaceID: gh.WorkspaceID,
		Owner:       gh.Owner,
		Repo:        gh.Repo,
		Path:        gh.Path,
	}
}

func encodeGHSource(w http.ResponseWriter, gh source.GHSource) error {
	g := ghToApiType(gh)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(g)
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

	gh, err := e.sourceSrv.CreateGHSource(r.Context(), source.CreateGHSourceArgs{
		WorkspaceID:    t.WorkspaceID,
		Owner:          req.Owner,
		Repo:           req.Repo,
		Path:           req.Path,
		InstallationID: req.InstallationID,
	})
	if err != nil {
		return RestErrorFromDomain(err)
	}

	return encodeGHSource(w, gh)
}

func getGHSource(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	sId, err := strconv.ParseInt(r.PathValue("sid"), 10, 64)
	if err != nil {
		return err
	}

	gh, err := e.sourceSrv.GetGHSource(r.Context(), t.WorkspaceID, sId)
	if err != nil {
		return RestErrorFromDomain(err)
	}
	return encodeGHSource(w, gh)
}

func listGHSources(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	sources, err := e.sourceSrv.ListGHSources(r.Context(), t.WorkspaceID)
	if err != nil {
		return RestErrorFromDomain(err)
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
	sId, err := strconv.ParseInt(r.PathValue("sid"), 10, 64)
	if err != nil {
		return err
	}

	err = e.sourceSrv.DeleteGHSource(r.Context(), t.WorkspaceID, sId)
	if err != nil {
		return RestErrorFromDomain(err)
	}
	return nil
}
