package source

import (
	"context"
	"errors"

	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/domain"
)

type configRepo struct {
	cfg config.Config
}

func NewConfigRepo(cfg config.Config) Repo {
	return configRepo{
		cfg: cfg,
	}
}

func (c configRepo) CreateGHSource(ctx context.Context, s GHSource) (GHSource, error) {
	return GHSource{}, domain.NewError(domain.ErrBadRequest, errors.New("github source creation not allowed in local config mode"))
}

func (c configRepo) DeleteGHSource(ctx context.Context, workspaceID string, id ID) error {
	return domain.NewError(domain.ErrBadRequest, errors.New("github source deletion not allowed in local config mode"))
}

func (c configRepo) GetGHSource(ctx context.Context, workspaceID string, id ID) (GHSource, error) {
	if c.cfg.Github == nil {
		return GHSource{}, domain.NewError(domain.ErrNotFound, errors.New("github source not found"))
	}
	g := GHSource{
		ID:    NewID(1),
		Owner: c.cfg.Github.Owner,
		Repo:  c.cfg.Github.Repo,
		Path:  c.cfg.Github.Path,
	}
	if c.cfg.Github.Auth != nil {
		g.InstallationID = c.cfg.Github.Auth.AppInstallationId
	}
	return g, nil
}

func (c configRepo) ListGHSource(ctx context.Context, workspaceID string) ([]GHSource, error) {
	g, err := c.GetGHSource(ctx, workspaceID, NewID(1))
	if err != nil {
		return []GHSource{}, err
	}
	return []GHSource{g}, nil
}
