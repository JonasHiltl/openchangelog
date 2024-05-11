package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/domain/workspace"
)

func encodeWorkspace(w http.ResponseWriter, ws workspace.Workspace) error {
	res := map[string]any{
		"id":    ws.ID,
		"name":  ws.Name,
		"token": ws.Token.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

func createWorkspace(e *env, w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	ws, err := e.workspaceSrv.CreateWorkspace(r.Context(), req.Name)
	if err != nil {
		return err
	}
	return encodeWorkspace(w, ws)
}

func updateWorkspace(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}
	wId := r.PathValue("wid")
	if wId != t.WorkspaceID {
		return errNoWorkspace
	}

	var req struct {
		Name string `json:"name"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	ws, err := e.workspaceSrv.UpdateWorkspace(r.Context(), wId, req.Name)
	if err != nil {
		return RestErrorFromDomain(err)
	}

	return encodeWorkspace(w, ws)
}

var errNoWorkspace = RestError{Code: http.StatusNotFound, Err: errors.New("workspace not found")}

func getWorkspace(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	wId := r.PathValue("wid")
	if wId != t.WorkspaceID {
		return errNoWorkspace
	}

	ws, err := e.workspaceSrv.GetWorkspace(r.Context(), t.WorkspaceID)
	if err != nil {
		return RestErrorFromDomain(err)
	}

	return encodeWorkspace(w, ws)
}

func deleteWorkspace(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	wId := r.PathValue("wid")
	if wId != t.WorkspaceID {
		return errNoWorkspace
	}

	err = e.workspaceSrv.DeleteWorkspace(r.Context(), wId)
	if err != nil {
		return RestErrorFromDomain(err)
	}
	return nil
}
