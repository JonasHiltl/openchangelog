package workspace

import (
	"context"
)

type Service interface {
	CreateWorkspace(ctx context.Context, name string) (Workspace, error)
	GetWorkspace(ctx context.Context, id string) (Workspace, error)
	DeleteWorkspace(ctx context.Context, id string) error
	UpdateWorkspace(ctx context.Context, id string, name string) (Workspace, error)

	GetWorkspaceIDByToken(ctx context.Context, token string) (ID, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateWorkspace(ctx context.Context, name string) (Workspace, error) {
	ws := Workspace{
		ID:    NewID(),
		Name:  name,
		Token: NewToken(),
	}
	err := s.repo.SaveWorkspace(ctx, ws)
	if err != nil {
		return Workspace{}, err
	}

	return ws, nil
}

func (s *service) GetWorkspace(ctx context.Context, id string) (Workspace, error) {
	i, err := ParseID(id)
	if err != nil {
		return Workspace{}, err
	}
	return s.repo.GetWorkspace(ctx, i)
}

func (s *service) DeleteWorkspace(ctx context.Context, id string) error {
	i, err := ParseID(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteWorkspace(ctx, i)
}

func (s *service) UpdateWorkspace(ctx context.Context, id string, name string) (Workspace, error) {
	i, err := ParseID(id)
	if err != nil {
		return Workspace{}, err
	}
	ws := Workspace{
		ID:   i,
		Name: name,
	}
	err = s.repo.SaveWorkspace(ctx, ws)
	return ws, err
}

func (s *service) GetWorkspaceIDByToken(ctx context.Context, token string) (ID, error) {
	tkn, err := ParseToken(token)
	if err != nil {
		return "", err
	}
	return s.repo.GetWorkspaceIDByToken(ctx, tkn)
}
