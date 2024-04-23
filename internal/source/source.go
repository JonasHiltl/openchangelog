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
	if cfg.Local != nil {
		return LocalFile(cfg.Local.FilesPath), nil
	}
	if cfg.Github != nil && cfg.Github.Auth != nil {
		return Github(GithubSourceOptions{
			Owner:             cfg.Github.Owner,
			Repository:        cfg.Github.Repo,
			Path:              cfg.Github.Path,
			AppPrivateKey:     cfg.Github.Auth.AppPrivateKey,
			AppInstallationId: cfg.Github.Auth.AppInstallationId,
			AccessToken:       cfg.Github.Auth.AccessToken,
		})
	}

	return nil, errors.New("markdown file source not specififed in config")
}
