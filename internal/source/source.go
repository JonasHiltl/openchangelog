package source

import (
	"context"
	"errors"

	"github.com/jonashiltl/openchangelog/internal/config"
)

type Article struct {
	Bytes []byte
}

type LoadParams interface {
	PageSize() int
	Page() int
	StartIdx() int
	EndIdx() int
}

type loadParams struct {
	pageSize int
	page     int
}

func NewLoadParams(pageSize int, page int) LoadParams {
	return loadParams{
		pageSize: pageSize,
		page:     page,
	}
}

func (p loadParams) PageSize() int {
	return p.pageSize
}
func (p loadParams) Page() int {
	return p.page
}

func (p loadParams) StartIdx() int {
	return (p.page - 1) * p.pageSize
}

func (p loadParams) EndIdx() int {
	return p.page*p.pageSize - 1
}

type LoadResult struct {
	Articles []Article
	HasMore  bool
}

// Represents a source of the Changelog Markdown files.
type Source interface {
	Load(ctx context.Context, params LoadParams) (LoadResult, error)
}

func NewFromConfig(cfg config.Config) (Source, error) {
	if cfg.Local.FilesPath != "" {
		return LocalFile(cfg.Local.FilesPath), nil
	}
	if cfg.GitHub.AppInstallationId != 0 && cfg.GitHub.AppPrivateKey != "" {
		return GitHub(GitHubSourceOptions{
			Owner:             cfg.GitHub.Owner,
			Repository:        cfg.GitHub.Repo,
			Path:              cfg.GitHub.Path,
			AppPrivateKey:     cfg.GitHub.AppPrivateKey,
			AppInstallationId: cfg.GitHub.AppInstallationId,
		})
	}

	return nil, errors.New("markdown file source not specififed in config")
}
