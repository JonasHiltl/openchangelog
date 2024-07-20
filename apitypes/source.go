package apitypes

import (
	"encoding/json"
)

type SourceType string

const (
	GitHub SourceType = "github"
)

type Source interface {
	Type() SourceType
}

type GHSource struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspaceId"`
	Owner       string `json:"owner"`
	Repo        string `json:"repo"`
	Path        string `json:"path,omitempty"`
}

func (g GHSource) Type() SourceType {
	return GitHub
}

func (g GHSource) MarshalJSON() (b []byte, e error) {
	// needed to bypass recursive marshaling of GHSource
	type Alias GHSource
	return json.Marshal(struct {
		Type SourceType `json:"type"`
		Alias
	}{
		Type:  g.Type(),
		Alias: Alias(g),
	})
}

type CreateGHSourceBody struct {
	Owner          string `json:"owner"`
	Repo           string `json:"repo"`
	Path           string `json:"path"`
	InstallationID int64  `json:"installationID"`
}
