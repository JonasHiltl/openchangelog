package workspace

import "context"

type Repo interface {
	SaveWorkspace(ctx context.Context, w Workspace) error
	DeleteWorkspace(ctx context.Context, id ID) error
	GetWorkspace(ctx context.Context, id ID) (Workspace, error)

	GetWorkspaceIDByToken(ctx context.Context, tkn Token) (ID, error)
}
