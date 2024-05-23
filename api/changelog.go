package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jonashiltl/openchangelog/apitypes"
)

type Changelog = apitypes.Changelog

func (c *Client) GetChangelog(ctx context.Context, changelogID int64) (Changelog, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/changelogs/%d", changelogID))
	if err != nil {
		return Changelog{}, err
	}

	resp, err := c.rawRequestWithContext(req)
	if err != nil {
		return Changelog{}, fmt.Errorf("error while getting changelog %d: %w", changelogID, err)
	}
	defer resp.Body.Close()

	var cl Changelog
	err = resp.DecodeJSON(&cl)
	return cl, err
}

func (c *Client) ListChangelogs(ctx context.Context) ([]Changelog, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, "/changelogs")
	if err != nil {
		return nil, err
	}

	resp, err := c.rawRequestWithContext(req)
	if err != nil {
		return nil, fmt.Errorf("error while listing changelogs: %w", err)
	}
	defer resp.Body.Close()

	var cls []Changelog
	err = resp.DecodeJSON(&cls)
	return cls, err
}
