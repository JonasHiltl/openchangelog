package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jonashiltl/openchangelog/apitypes"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/store"
)

const (
	workspace_id_param = "wid"
)

func encodeWorkspace(w http.ResponseWriter, ws store.Workspace) error {
	res := apitypes.Workspace{
		ID:    ws.ID.String(),
		Name:  ws.Name,
		Token: ws.Token.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

func createWorkspace(e *env, w http.ResponseWriter, r *http.Request) error {
	var req apitypes.CreateWorkspaceBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	ws, err := e.store.SaveWorkspace(r.Context(), store.Workspace{
		ID:    store.NewWID(),
		Name:  req.Name,
		Token: store.NewToken(),
	})
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
	wId := r.PathValue(workspace_id_param)
	if wId != t.WorkspaceID.String() {
		return errNoWorkspace
	}

	var req struct {
		Name string `json:"name"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	ws, err := e.store.SaveWorkspace(r.Context(), store.Workspace{
		ID:   t.WorkspaceID,
		Name: req.Name,
	})
	if err != nil {
		return err
	}

	return encodeWorkspace(w, ws)
}

var errNoWorkspace = errs.NewError(errs.ErrNotFound, errors.New("workspace not found"))

func getWorkspace(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	wId := r.PathValue(workspace_id_param)
	if wId != t.WorkspaceID.String() {
		return errNoWorkspace
	}

	ws, err := e.store.GetWorkspace(r.Context(), t.WorkspaceID)
	if err != nil {
		return err
	}

	return encodeWorkspace(w, ws)
}

func deleteWorkspace(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	wId := r.PathValue(workspace_id_param)
	if wId != t.WorkspaceID.String() {
		return errNoWorkspace
	}

	err = e.store.DeleteWorkspace(r.Context(), t.WorkspaceID)
	if err != nil {
		return err
	}
	return nil
}
