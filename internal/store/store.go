package store

import (
	"context"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/guregu/null/v5"
)

type Changelog struct {
	WorkspaceID WorkspaceID
	ID          ChangelogID
	Subdomain   string
	Title       null.String
	Subtitle    null.String
	LogoSrc     null.String
	LogoLink    null.String
	LogoAlt     null.String
	LogoHeight  null.String
	LogoWidth   null.String
	CreatedAt   time.Time
	GHSource    null.Value[GHSource]
	LocalSource null.Value[LocalSource]
}

type Workspace struct {
	ID    WorkspaceID
	Name  string
	Token Token
}

type GHSource struct {
	ID             GHSourceID
	WorkspaceID    WorkspaceID
	Owner          string
	Repo           string
	Path           string
	InstallationID int64
}

type LocalSource struct {
	Path string
}

type UpdateChangelogArgs struct {
	Title      null.String
	Subdomain  null.String
	Subtitle   null.String
	LogoSrc    null.String
	LogoLink   null.String
	LogoAlt    null.String
	LogoHeight null.String
	LogoWidth  null.String
}

type Store interface {
	GetChangelog(context.Context, WorkspaceID, ChangelogID) (Changelog, error)
	GetChangelogBySubdomain(context.Context, string) (Changelog, error)
	ListChangelogs(context.Context, WorkspaceID) ([]Changelog, error)
	CreateChangelog(context.Context, Changelog) (Changelog, error)
	UpdateChangelog(context.Context, WorkspaceID, ChangelogID, UpdateChangelogArgs) (Changelog, error)
	DeleteChangelog(context.Context, WorkspaceID, ChangelogID) error
	SetChangelogGHSource(context.Context, WorkspaceID, ChangelogID, GHSourceID) error
	DeleteChangelogSource(context.Context, WorkspaceID, ChangelogID) error

	// Workspace
	GetWorkspace(context.Context, WorkspaceID) (Workspace, error)
	SaveWorkspace(context.Context, Workspace) (Workspace, error)
	GetWorkspaceIDByToken(ctx context.Context, token string) (WorkspaceID, error)
	DeleteWorkspace(context.Context, WorkspaceID) error

	// Source
	CreateGHSource(context.Context, GHSource) (GHSource, error)
	GetGHSource(context.Context, WorkspaceID, GHSourceID) (GHSource, error)
	ListGHSources(context.Context, WorkspaceID) ([]GHSource, error)
	DeleteGHSource(context.Context, WorkspaceID, GHSourceID) error
}
