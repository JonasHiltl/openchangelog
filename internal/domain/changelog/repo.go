package changelog

import (
	"context"
)

type Repo interface {
	CreateChangelog(ctx context.Context, c Changelog) (Changelog, error)
	UpdateChangelog(ctx context.Context, c Changelog) (Changelog, error)
	GetChangelog(ctx context.Context, workspaceID string, id ID) (Changelog, error)
	DeleteChangelog(ctx context.Context, workspaceID string, id ID) error
	ListChangelogs(ctx context.Context, workspaceID string) ([]Changelog, error)
}
