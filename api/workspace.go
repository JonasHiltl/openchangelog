package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jonashiltl/openchangelog/apitypes"
)

type Workspace = apitypes.Workspace

func (c *Client) CreateWorkspace(ctx context.Context, args apitypes.CreateWorkspaceBody) (Workspace, error) {
	body, err := json.Marshal(args)
	if err != nil {
		return Workspace{}, err
	}

	req, err := c.NewRequest(
		ctx,
		http.MethodPost,
		"/workspaces",
		bytes.NewReader(body),
	)
	if err != nil {
		return Workspace{}, err
	}

	resp, err := c.rawRequestWithContext(req)
	if err != nil {
		return Workspace{}, fmt.Errorf("error while creating workspace: %w", err)
	}
	defer resp.Body.Close()

	var w Workspace
	err = resp.DecodeJSON(&w)
	return w, err
}

func (c *Client) DeleteWorkspace(ctx context.Context, id string) error {
	req, err := c.NewRequest(
		ctx, http.MethodDelete,
		fmt.Sprintf("/workspaces/%s", id),
		nil,
	)
	if err != nil {
		return err
	}

	_, err = c.rawRequestWithContext(req)
	if err != nil {
		return fmt.Errorf("error while creating workspace: %w", err)
	}
	return nil
}
