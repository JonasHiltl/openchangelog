package source

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/domain"
	"github.com/jonashiltl/openchangelog/loader"
)

type Source interface {
	Type() SourceType
	ToLoader(cfg config.Config) (loader.Loader, error)
}

type SourceType string

const (
	GitHub SourceType = "github"
	Local  SourceType = "local"
	String SourceType = "string"
)

func ParseSourceType(t string) (SourceType, error) {
	switch strings.ToLower(t) {
	case "github", "local", "string":
		return SourceType(t), nil
	default:
		return "", domain.NewError(domain.ErrBadRequest, fmt.Errorf("invalid source type %s", t))
	}

}

type GHSource struct {
	ID             ID
	WorkspaceID    string
	Owner          string
	Repo           string
	Path           string
	InstallationID int64
}

func (g GHSource) Type() SourceType {
	return GitHub
}

func (g GHSource) ToLoader(cfg config.Config) (loader.Loader, error) {
	if cfg.Github == nil || cfg.Github.Auth == nil {
		return nil, errors.New("github authentication not setup")
	}

	return loader.NewGithub(loader.GithubLoaderOptions{
		AppPrivateKey:     cfg.Github.Auth.AppPrivateKey,
		AccessToken:       cfg.Github.Auth.AccessToken,
		Owner:             g.Owner,
		Repository:        g.Repo,
		Path:              g.Path,
		AppInstallationId: g.InstallationID,
	})
}
