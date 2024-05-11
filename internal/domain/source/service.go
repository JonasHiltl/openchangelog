package source

import (
	"context"
	"errors"

	"github.com/jonashiltl/openchangelog/internal/domain"
)

type CreateGHSourceArgs struct {
	WorkspaceID    string
	Owner          string
	Repo           string
	Path           string
	InstallationID int64
}

type Service interface {
	CreateGHSource(ctx context.Context, args CreateGHSourceArgs) (GHSource, error)
	GetGHSource(ctx context.Context, workspaceID string, id int64) (GHSource, error)
	ListGHSources(ctx context.Context, workspaceID string) ([]GHSource, error)
	DeleteGHSource(ctx context.Context, workspaceID string, id int64) error
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateGHSource(ctx context.Context, args CreateGHSourceArgs) (GHSource, error) {
	if args.Owner == "" {
		return GHSource{}, domain.NewError(domain.ErrBadRequest, errors.New("owner must be specified"))
	}
	if args.Repo == "" {
		return GHSource{}, domain.NewError(domain.ErrBadRequest, errors.New("repo must be specified"))
	}
	return s.repo.CreateGHSource(ctx, GHSource{
		WorkspaceID:    args.WorkspaceID,
		Owner:          args.Owner,
		Repo:           args.Repo,
		Path:           args.Path,
		InstallationID: args.InstallationID,
	})
}

func (s *service) GetGHSource(ctx context.Context, workspaceID string, id int64) (GHSource, error) {
	return s.repo.GetGHSource(ctx, workspaceID, NewID(id))
}

func (s *service) ListGHSources(ctx context.Context, workspaceID string) ([]GHSource, error) {
	return s.repo.ListGHSource(ctx, workspaceID)
}

func (s *service) DeleteGHSource(ctx context.Context, workspaceID string, id int64) error {
	return s.repo.DeleteGHSource(ctx, workspaceID, NewID(id))
}
