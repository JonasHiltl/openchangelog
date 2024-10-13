package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jonashiltl/openchangelog/apitypes"
)

type Changelog = apitypes.Changelog
type FullChangelog = apitypes.FullChangelog

func (c *Client) GetChangelog(ctx context.Context, changelogID string) (Changelog, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/changelogs/%s", changelogID), nil)
	if err != nil {
		return Changelog{}, err
	}

	resp, err := c.rawRequestWithContext(req)
	if err != nil {
		return Changelog{}, fmt.Errorf("error while getting changelog %s: %w", changelogID, err)
	}
	defer resp.Body.Close()

	var cl Changelog
	err = resp.DecodeJSON(&cl)
	return cl, err
}

type GetFullChangelogParams struct {
	ChangelogID string
	Page        int
	PageSize    int
}

func (c *Client) GetFullChangelog(ctx context.Context, args GetFullChangelogParams) (FullChangelog, error) {
	q := url.Values{}
	if args.Page != 0 {
		q.Set("page", fmt.Sprint(args.Page))
	}
	if args.PageSize != 0 {
		q.Set("page-size", fmt.Sprint(args.PageSize))
	}

	url := fmt.Sprintf("/changelogs/%s/full?%s", args.ChangelogID, q.Encode())

	req, err := c.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return FullChangelog{}, err
	}

	resp, err := c.rawRequestWithContext(req)
	if err != nil {
		return FullChangelog{}, fmt.Errorf("error while getting full changelog %s: %w", args.ChangelogID, err)
	}
	defer resp.Body.Close()

	var cl FullChangelog
	err = resp.DecodeJSON(&cl)
	return cl, err
}

func (c *Client) ListChangelogs(ctx context.Context) ([]Changelog, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, "/changelogs", nil)
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

func (c *Client) CreateChangelog(ctx context.Context, args apitypes.CreateChangelogBody) (Changelog, error) {
	body, err := json.Marshal(args)
	if err != nil {
		return Changelog{}, err
	}

	req, err := c.NewRequest(ctx, http.MethodPost, "/changelogs", bytes.NewReader(body))
	if err != nil {
		return Changelog{}, err
	}

	resp, err := c.rawRequestWithContext(req)
	if err != nil {
		return Changelog{}, fmt.Errorf("error while creating changelog: %w", err)
	}
	defer resp.Body.Close()

	var cl Changelog
	err = resp.DecodeJSON(&cl)
	return cl, err
}

func (c *Client) UpdateChangelog(ctx context.Context, changelogID string, args apitypes.UpdateChangelogBody) (Changelog, error) {
	body, err := json.Marshal(args)
	if err != nil {
		return Changelog{}, err
	}

	req, err := c.NewRequest(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("/changelogs/%s", changelogID),
		bytes.NewReader(body),
	)
	if err != nil {
		return Changelog{}, err
	}

	resp, err := c.rawRequestWithContext(req)
	if err != nil {
		return Changelog{}, fmt.Errorf("error while updating changelog %s: %w", changelogID, err)
	}
	defer resp.Body.Close()

	var cl Changelog
	err = resp.DecodeJSON(&cl)
	return cl, err
}

func (c *Client) DeleteChangelog(ctx context.Context, changelogID string) error {
	req, err := c.NewRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("/changelogs/%s", changelogID),
		nil,
	)
	if err != nil {
		return err
	}

	_, err = c.rawRequestWithContext(req)
	if err != nil {
		return fmt.Errorf("error while deleting changelog %s: %w", changelogID, err)
	}
	return nil
}

func (c *Client) DeleteChangelogSource(ctx context.Context, changelogID string) error {
	req, err := c.NewRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("/changelogs/%s/source", changelogID),
		nil,
	)
	if err != nil {
		return err
	}

	_, err = c.rawRequestWithContext(req)
	if err != nil {
		return fmt.Errorf("error while deleting changelog %s source: %w", changelogID, err)
	}
	return nil
}

func (c *Client) SetChangelogSource(ctx context.Context, changelogID string, sourceID string) error {
	req, err := c.NewRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/changelogs/%s/source/%s", changelogID, sourceID),
		nil,
	)
	if err != nil {
		return err
	}

	_, err = c.rawRequestWithContext(req)
	if err != nil {
		return fmt.Errorf("error while setting changelog %s source %s: %w", changelogID, sourceID, err)
	}
	return nil
}
