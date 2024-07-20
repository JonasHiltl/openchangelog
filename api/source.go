package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jonashiltl/openchangelog/apitypes"
)

type Source = apitypes.Source
type GHSource = apitypes.CreateGHSourceBody
type CreateGHSourceBody = apitypes.CreateGHSourceBody

func (c *Client) CreateGHSource(ctx context.Context, args CreateGHSourceBody) (GHSource, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, "/sources/gh", nil)
	if err != nil {
		return GHSource{}, err
	}

	resp, err := c.rawRequestWithContext(req)
	if err != nil {
		return GHSource{}, fmt.Errorf("error while creating github source: %w", err)
	}
	defer resp.Body.Close()

	var s GHSource
	err = resp.DecodeJSON(&s)
	return s, err
}

func (c *Client) DeleteGHSource(ctx context.Context, sourceID string) error {
	req, err := c.NewRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("/api/sources/gh/%s", sourceID),
		nil,
	)
	if err != nil {
		return err
	}

	_, err = c.rawRequestWithContext(req)
	if err != nil {
		return fmt.Errorf("error while deleting github source %s: %w", sourceID, err)
	}
	return nil
}

func (c *Client) ListSources(ctx context.Context) ([]Source, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, "/sources", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.rawRequestWithContext(req)
	if err != nil {
		return nil, fmt.Errorf("error while listing sources: %w", err)
	}
	defer resp.Body.Close()

	var objs []json.RawMessage
	err = resp.DecodeJSON(&objs)
	if err != nil {
		return nil, err
	}

	res := make([]Source, len(objs))
	for i, obj := range objs {
		res[i] = apitypes.DecodeSource(obj)
	}

	return res, nil
}
