package changelog

import (
	"context"
	"errors"

	"github.com/guregu/null/v5"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/domain"
	"github.com/jonashiltl/openchangelog/internal/domain/source"
)

type configRepo struct {
	cfg config.Config
}

func NewConfigRepo(cfg config.Config) Repo {
	return configRepo{
		cfg: cfg,
	}
}

func (r configRepo) CreateChangelog(ctx context.Context, c Changelog) (Changelog, error) {
	return Changelog{}, domain.NewError(domain.ErrBadRequest, errors.New("changelog creation not allowed in local config mode"))
}

func (r configRepo) UpdateChangelog(ctx context.Context, c Changelog) (Changelog, error) {
	return Changelog{}, domain.NewError(domain.ErrBadRequest, errors.New("update changelog not allowed in local config mode"))
}

func (r configRepo) GetChangelog(ctx context.Context, workspaceID string, id ID) (Changelog, error) {
	c := Changelog{
		ID:       NewID(1),
		Title:    r.cfg.Page.Title,
		Subtitle: r.cfg.Page.Subtitle,
	}

	if r.cfg.Page.Logo != nil {
		l := r.cfg.Page.Logo
		c.Logo = struct {
			Src    null.String
			Link   null.String
			Alt    null.String
			Height null.String
			Width  null.String
		}{
			Src:    null.NewString(l.Src, l.Src != ""),
			Link:   null.NewString(l.Link, l.Link != ""),
			Alt:    null.NewString(l.Alt, l.Alt != ""),
			Height: null.NewString(l.Height, l.Height != ""),
			Width:  null.NewString(l.Width, l.Width != ""),
		}
	}

	if r.cfg.Github != nil {
		gh := source.GHSource{
			ID:    source.NewID(1),
			Owner: r.cfg.Github.Owner,
			Repo:  r.cfg.Github.Repo,
			Path:  r.cfg.Github.Path,
		}
		if r.cfg.Github.Auth != nil {
			gh.InstallationID = r.cfg.Github.Auth.AppInstallationId
		}
		c.Source = gh
	}
	return c, nil
}

func (r configRepo) DeleteChangelog(ctx context.Context, workspaceID string, id ID) error {
	return domain.NewError(domain.ErrBadRequest, errors.New("changelog deletion not allowed in local config mode"))
}

func (r configRepo) ListChangelogs(ctx context.Context, workspaceID string) ([]Changelog, error) {
	c, err := r.GetChangelog(ctx, workspaceID, NewID(1))
	if err != nil {
		return []Changelog{}, err
	}
	return []Changelog{c}, nil
}
