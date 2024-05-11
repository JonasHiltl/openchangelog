package workspace

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

func (c configRepo) SaveWorkspace(ctx context.Context, w Workspace) error {
	return domain.NewError(domain.ErrBadRequest, errors.New("workspace creation not allowed in local config mode"))
}

func (c configRepo) DeleteWorkspace(ctx context.Context, id ID) error {
	return domain.NewError(domain.ErrBadRequest, errors.New("workspace deletion not allowed in local config mode"))
}

func (c configRepo) GetWorkspace(ctx context.Context, id ID) (Workspace, error) {
	return Workspace{}, domain.NewError(domain.ErrBadRequest, errors.New("get workspace not allowed in local config mode"))
}

func (c configRepo) GetWorkspaceIDByToken(ctx context.Context, tkn Token) (ID, error) {
	return "", domain.NewError(domain.ErrBadRequest, errors.New("get workspace id by token not allowed in local config mode"))
}
