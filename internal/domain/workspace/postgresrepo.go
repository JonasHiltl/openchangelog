package workspace

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jonashiltl/openchangelog/internal/domain"
	"github.com/jonashiltl/openchangelog/internal/domain/workspace/db"
)

type postgresRepo struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewPGRepo(pool *pgxpool.Pool) Repo {
	queries := db.New(pool)
	return &postgresRepo{
		queries: queries,
		pool:    pool,
	}
}

func (r *postgresRepo) SaveWorkspace(ctx context.Context, w Workspace) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	q := r.queries.WithTx(tx)
	err = q.SaveWorkspace(ctx, db.SaveWorkspaceParams{
		ID:   w.ID.String(),
		Name: w.Name,
	})
	if err != nil {
		return err
	}

	if w.Token.IsSet() {
		err = q.CreateToken(ctx, db.CreateTokenParams{
			WorkspaceID: w.ID.String(),
			Key:         w.Token.String(),
		})

		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *postgresRepo) DeleteWorkspace(ctx context.Context, id ID) error {
	return r.queries.DeleteWorkspace(ctx, id.String())
}

func (r *postgresRepo) GetWorkspace(ctx context.Context, id ID) (Workspace, error) {
	row, err := r.queries.GetWorkspace(ctx, id.String())
	if err != nil {
		if err == pgx.ErrNoRows {
			return Workspace{}, domain.NewError(domain.ErrNotFound, errors.New("workspace not found"))
		}
		return Workspace{}, err
	}
	t, err := ParseToken(row.Token.Key)
	if err != nil {
		return Workspace{}, err
	}

	return Workspace{
		ID:    id,
		Name:  row.Workspace.Name,
		Token: t,
	}, nil
}

func (r *postgresRepo) GetWorkspaceIDByToken(ctx context.Context, token Token) (ID, error) {
	t, err := r.queries.GetToken(ctx, token.String())
	if err != nil {
		return "", err
	}
	return ParseID(t.WorkspaceID)
}
