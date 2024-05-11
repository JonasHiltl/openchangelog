package changelog

import (
	"context"
	"errors"
	"fmt"

	"github.com/jonashiltl/openchangelog/internal/domain"
	"github.com/jonashiltl/openchangelog/internal/domain/source"
)

type CreateChangelogArgs struct {
	WorkspaceID string
	Title       string
	Subtitle    string
	Logo        struct {
		Src    string
		Link   string
		Alt    string
		Height string
		Width  string
	}
}

type UpdateChangelogArgs struct {
	Title    string
	Subtitle string
	Logo     struct {
		Src    string
		Link   string
		Alt    string
		Height string
		Width  string
	}
	Source source.Source
}

type SetChangelogSourceArgs struct {
	Type string
	ID   int64
}

type Service interface {
	SetChangelogSource(ctx context.Context, workspaceID string, id int64, args SetChangelogSourceArgs) error

	CreateChangelog(ctx context.Context, args CreateChangelogArgs) (Changelog, error)
	UpdateChangelog(ctx context.Context, workspaceID string, id int64, args UpdateChangelogArgs) (Changelog, error)
	GetChangelog(ctx context.Context, workspaceID string, id int64) (Changelog, error)
	DeleteChangelog(ctx context.Context, workspaceID string, id int64) error
	ListChangelogs(ctx context.Context, workspaceID string) ([]Changelog, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateChangelog(ctx context.Context, args CreateChangelogArgs) (Changelog, error) {
	if args.WorkspaceID == "" {
		return Changelog{}, domain.NewError(domain.ErrBadRequest, errors.New("workspace id is missing"))
	}
	return s.repo.CreateChangelog(ctx, Changelog{
		WorkspaceID: args.WorkspaceID,
		Title:       args.Title,
		Subtitle:    args.Subtitle,
		Logo:        args.Logo,
	})
}

func (s *service) UpdateChangelog(ctx context.Context, workspaceID string, id int64, args UpdateChangelogArgs) (Changelog, error) {
	return s.repo.UpdateChangelog(ctx, Changelog{
		ID:          NewID(id),
		WorkspaceID: workspaceID,
		Title:       args.Title,
		Subtitle:    args.Subtitle,
		Logo:        args.Logo,
		Source:      args.Source,
	})
}

func (s *service) GetChangelog(ctx context.Context, workspaceID string, id int64) (Changelog, error) {
	return s.repo.GetChangelog(ctx, workspaceID, NewID(id))
}

func (s *service) DeleteChangelog(ctx context.Context, workspaceID string, id int64) error {
	return s.repo.DeleteChangelog(ctx, workspaceID, NewID(id))
}

func (s *service) ListChangelogs(ctx context.Context, workspaceID string) ([]Changelog, error) {
	return s.repo.ListChangelogs(ctx, workspaceID)
}

func (s *service) SetChangelogSource(ctx context.Context, workspaceID string, id int64, args SetChangelogSourceArgs) error {
	srcType, err := source.ParseSourceType(args.Type)
	if err != nil {
		return err
	}
	switch srcType {
	case source.GitHub:
		src := source.GHSource{
			ID:          source.NewID(args.ID),
			WorkspaceID: workspaceID,
		}
		_, err = s.repo.UpdateChangelog(ctx, Changelog{
			ID:          NewID(id),
			WorkspaceID: workspaceID,
			Source:      src,
		})
		return err
	default:
		return domain.NewError(domain.ErrBadRequest, fmt.Errorf("we currently don't support %s source in cloud mode", srcType))
	}
}
