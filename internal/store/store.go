package store

import (
	"context"
	"time"

	"github.com/jonashiltl/openchangelog/apitypes"
	_ "github.com/mattn/go-sqlite3"

	"github.com/guregu/null/v5"
)

type Changelog struct {
	WorkspaceID   WorkspaceID
	ID            ChangelogID
	Subdomain     Subdomain
	Domain        Domain
	Title         apitypes.NullString
	Subtitle      apitypes.NullString
	LogoSrc       apitypes.NullString
	LogoLink      apitypes.NullString
	LogoAlt       apitypes.NullString
	LogoHeight    apitypes.NullString
	LogoWidth     apitypes.NullString
	ColorScheme   ColorScheme
	HidePoweredBy bool
	Protected     bool
	PasswordHash  string
	CreatedAt     time.Time
	GHSource      null.Value[GHSource]
	LocalSource   null.Value[LocalSource]
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
	Title         apitypes.NullString
	Subdomain     apitypes.NullString
	Domain        Domain
	Subtitle      apitypes.NullString
	LogoSrc       apitypes.NullString
	LogoLink      apitypes.NullString
	LogoAlt       apitypes.NullString
	LogoHeight    apitypes.NullString
	LogoWidth     apitypes.NullString
	ColorScheme   ColorScheme
	HidePoweredBy *bool
	Protected     *bool
	PasswordHash  apitypes.NullString
}

type Store interface {
	GetChangelog(context.Context, WorkspaceID, ChangelogID) (Changelog, error)
	GetChangelogByDomainOrSubdomain(ctx context.Context, domain Domain, subdomain Subdomain) (Changelog, error)
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
