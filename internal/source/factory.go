package source

import (
	"errors"
	"fmt"

	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/store"
)

type SourceFactory struct {
	cfg config.Config
}

func NewSourceFactory(cfg config.Config) *SourceFactory {
	return &SourceFactory{
		cfg: cfg,
	}
}

func (l *SourceFactory) FromDB(changelog store.Changelog, source store.ChangelogSource) (Source, error) {
	if !changelog.SourceType.Valid {
		return nil, errors.New("changelog has no stored source")
	}

	switch changelog.SourceType.SourceType {
	case store.SourceTypeGitHub:
		{
			if !source.ID.Valid {
				return nil, errors.New("referenced github source has invalid id")
			}
			return Github(GithubSourceOptions{
				Owner:             source.Owner.String,
				Repository:        source.Repo.String,
				Path:              source.Path.String,
				AppInstallationId: source.InstallationID.Int64,
				AppPrivateKey:     l.cfg.Github.Auth.AppPrivateKey,
				AccessToken:       l.cfg.Github.Auth.AccessToken,
			})
		}
	default:
		return nil, fmt.Errorf("changelog has invalid source %s", string(store.SourceTypeGitHub))
	}
}

func (l *SourceFactory) FromConfig() (Source, error) {
	if l.cfg.Local != nil {
		return LocalFile(l.cfg.Local.FilesPath), nil
	}
	if l.cfg.Github != nil && l.cfg.Github.Auth != nil {
		return Github(GithubSourceOptions{
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
