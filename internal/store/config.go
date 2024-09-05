package store

import (
	"context"
	"errors"

	"github.com/guregu/null/v5"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/errs"
)

const (
	CL_DEFAULT_ID = ChangelogID("cl_config")
	GH_DEFAULT_ID = GHSourceID("gh_config")
	WS_DEFAULT_ID = WorkspaceID("ws_config")
)

// Create a new store implementation, backed by the config file
func NewConfigStore(cfg config.Config) Store {
	return &configStore{
		cfg: cfg,
	}
}

type configStore struct {
	cfg config.Config
}

func (s *configStore) CreateChangelog(context.Context, Changelog) (Changelog, error) {
	return Changelog{}, errs.NewError(errs.ErrBadRequest, errors.New("changelog creation not allowed in local config mode"))
}

func (s *configStore) UpdateChangelog(context.Context, WorkspaceID, ChangelogID, UpdateChangelogArgs) (Changelog, error) {
	return Changelog{}, errs.NewError(errs.ErrBadRequest, errors.New("update changelog not allowed in local config mode"))
}

func (s *configStore) GetChangelog(ctx context.Context, wID WorkspaceID, cID ChangelogID) (Changelog, error) {
	cl := Changelog{
		ID: CL_DEFAULT_ID,
	}

	if s.cfg.Page != nil {
		cl.Title = null.NewString(s.cfg.Page.Title, s.cfg.Page.Title != "")
		cl.Subtitle = null.NewString(s.cfg.Page.Subtitle, s.cfg.Page.Subtitle != "")
	}

	if s.cfg.Page.Logo != nil {
		l := s.cfg.Page.Logo
		cl.LogoSrc = null.NewString(l.Src, l.Src != "")
		cl.LogoLink = null.NewString(l.Link, l.Link != "")
		cl.LogoAlt = null.NewString(l.Alt, l.Alt != "")
		cl.LogoHeight = null.NewString(l.Height, l.Height != "")
		cl.LogoWidth = null.NewString(l.Width, l.Width != "")
	}

	// parse local source from config
	if s.cfg.Local != nil {
		cl.LocalSource = null.NewValue(LocalSource{
			Path: s.cfg.Local.FilesPath,
		}, true)
	}

	// parse github source from config
	g, err := s.GetGHSource(ctx, wID, GH_DEFAULT_ID)
	if err == nil {
		cl.GHSource = null.NewValue(g, true)
	}

	return cl, nil
}

func (s *configStore) GetChangelogByDomainOrSubdomain(ctx context.Context, domain Domain, subdomain string) (Changelog, error) {
	return s.GetChangelog(ctx, WS_DEFAULT_ID, CL_DEFAULT_ID)
}

func (s *configStore) ListChangelogs(ctx context.Context, wID WorkspaceID) ([]Changelog, error) {
	cl, err := s.GetChangelog(ctx, wID, CL_DEFAULT_ID)
	if err != nil {
		return []Changelog{}, err
	}
	return []Changelog{cl}, nil
}

func (s *configStore) DeleteChangelog(context.Context, WorkspaceID, ChangelogID) error {
	return errs.NewError(errs.ErrBadRequest, errors.New("changelog deletion not allowed in local config mode"))
}

func (s *configStore) SetChangelogGHSource(context.Context, WorkspaceID, ChangelogID, GHSourceID) error {
	return errs.NewError(errs.ErrBadRequest, errors.New("changeing changelog source not allowed in local config mode"))
}

func (s *configStore) DeleteChangelogSource(context.Context, WorkspaceID, ChangelogID) error {
	return errs.NewError(errs.ErrBadRequest, errors.New("changelog source deletion not allowed in local config mode"))
}

func (s *configStore) CreateGHSource(context.Context, GHSource) (GHSource, error) {
	return GHSource{}, errs.NewError(errs.ErrBadRequest, errors.New("github source creation not allowed in local config mode"))
}

func (s *configStore) DeleteGHSource(context.Context, WorkspaceID, GHSourceID) error {
	return errs.NewError(errs.ErrBadRequest, errors.New("github source deletion not allowed in local config mode"))
}

func (s *configStore) ListGHSources(ctx context.Context, wID WorkspaceID) ([]GHSource, error) {
	g, err := s.GetGHSource(ctx, wID, GH_DEFAULT_ID)
	if err != nil {
		return []GHSource{}, err
	}
	return []GHSource{g}, nil
}

func (s *configStore) GetGHSource(context.Context, WorkspaceID, GHSourceID) (GHSource, error) {
	if s.cfg.Github == nil {
		return GHSource{}, errs.NewError(errs.ErrNotFound, errors.New("github source not found"))
	}
	g := GHSource{
		ID:          GH_DEFAULT_ID,
		Owner:       s.cfg.Github.Owner,
		Repo:        s.cfg.Github.Repo,
		Path:        s.cfg.Github.Path,
		WorkspaceID: WS_DEFAULT_ID,
	}
	if s.cfg.Github.Auth != nil {
		g.InstallationID = s.cfg.Github.Auth.AppInstallationId
	}
	return g, nil
}

func (s *configStore) SaveWorkspace(context.Context, Workspace) (Workspace, error) {
	return Workspace{}, errs.NewError(errs.ErrBadRequest, errors.New("workspace creation not allowed in local config mode"))
}

func (s *configStore) DeleteWorkspace(context.Context, WorkspaceID) error {
	return errs.NewError(errs.ErrBadRequest, errors.New("workspace deletion not allowed in local config mode"))
}

func (s *configStore) GetWorkspace(context.Context, WorkspaceID) (Workspace, error) {
	return Workspace{}, errs.NewError(errs.ErrBadRequest, errors.New("get workspace not allowed in local config mode"))
}

func (s *configStore) GetWorkspaceIDByToken(ctx context.Context, token string) (WorkspaceID, error) {
	return WS_DEFAULT_ID, nil
}
