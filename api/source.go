package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jonashiltl/openchangelog/apitypes"
)

type Source = apitypes.Source

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
