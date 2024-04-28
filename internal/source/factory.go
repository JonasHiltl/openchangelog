package source

import (
	"errors"
	"fmt"

	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/source"
)

type SourceFactory struct {
	cfg config.Config
}

func NewSourceFactory(cfg config.Config) *SourceFactory {
	return &SourceFactory{
		cfg: cfg,
	}
}

func (l *SourceFactory) FromDB(changelog store.Changelog, src store.ChangelogSource) (source.Source, error) {
	if !changelog.SourceType.Valid {
		return nil, errors.New("changelog has no stored source")
	}

	switch changelog.SourceType.SourceType {
	case store.SourceTypeGitHub:
		{
			if !src.ID.Valid {
				return nil, errors.New("referenced github source has invalid id")
			}
			return source.Github(source.GithubSourceOptions{
				Owner:             src.Owner.String,
				Repository:        src.Repo.String,
				Path:              src.Path.String,
				AppInstallationId: src.InstallationID.Int64,
				AppPrivateKey:     l.cfg.Github.Auth.AppPrivateKey,
				AccessToken:       l.cfg.Github.Auth.AccessToken,
			})
		}
	default:
		return nil, fmt.Errorf("changelog has invalid source %s", string(store.SourceTypeGitHub))
	}
}

func (l *SourceFactory) FromConfig() (source.Source, error) {
	if l.cfg.Local != nil {
		return source.LocalFile(l.cfg.Local.FilesPath), nil
	}
	if l.cfg.Github != nil && l.cfg.Github.Auth != nil {
		return source.Github(source.GithubSourceOptions{
			Owner:             l.cfg.Github.Owner,
			Repository:        l.cfg.Github.Repo,
			Path:              l.cfg.Github.Path,
			AppPrivateKey:     l.cfg.Github.Auth.AppPrivateKey,
			AppInstallationId: l.cfg.Github.Auth.AppInstallationId,
			AccessToken:       l.cfg.Github.Auth.AccessToken,
		})
	}

	return nil, errors.New("markdown file source not specififed in config")
}
