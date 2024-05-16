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
	ID          int64  `json:"id"`
	WorkspaceID string `json:"workspaceId"`
	Owner       string `json:"owner"`
	Repo        string `json:"repo"`
	Path        string `json:"path,omitempty"`
}

func (g GHSource) Type() SourceType {
	return GitHub
}

func (g GHSource) MarshalJSON() (b []byte, e error) {
	return json.Marshal(struct {
		Type SourceType `json:"type"`
		GHSource
	}{
		Type:     g.Type(),
		GHSource: g,
	})
}
