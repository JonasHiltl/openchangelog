package source

import (
	"context"
)

type Repo interface {
	CreateGHSource(ctx context.Context, s GHSource) (GHSource, error)
	DeleteGHSource(ctx context.Context, workspaceID string, id ID) error
	GetGHSource(ctx context.Context, workspaceID string, id ID) (GHSource, error)
	ListGHSource(ctx context.Context, workspaceID string) ([]GHSource, error)
}
