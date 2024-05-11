package source

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jonashiltl/openchangelog/internal/domain/source/db"
)

type postgresRepo struct {
	queries *db.Queries
}

func NewPGRepo(pool *pgxpool.Pool) Repo {
	queries := db.New(pool)
	return &postgresRepo{
		queries: queries,
	}
}

func (r *postgresRepo) CreateGHSource(ctx context.Context, s GHSource) (GHSource, error) {
	dbS, err := r.queries.CreateGHSource(ctx, db.CreateGHSourceParams{
		WorkspaceID:    s.WorkspaceID,
		Owner:          s.Owner,
		Repo:           s.Repo,
		Path:           s.Path,
		InstallationID: s.InstallationID,
	})
	if err != nil {
		return GHSource{}, err
	}
	return GHSource{
		ID:             NewID(dbS.ID),
		WorkspaceID:    dbS.WorkspaceID,
		Owner:          dbS.Owner,
		Repo:           dbS.Repo,
		Path:           dbS.Path,
		InstallationID: dbS.InstallationID,
	}, nil
}

func (r *postgresRepo) DeleteGHSource(ctx context.Context, workspaceID string, id ID) error {
	return r.queries.DeleteGHSource(ctx, db.DeleteGHSourceParams{
		ID:          id.ToInt(),
		WorkspaceID: workspaceID,
	})
}

func (r *postgresRepo) GetGHSource(ctx context.Context, workspaceID string, id ID) (GHSource, error) {
	dbS, err := r.queries.GetGHSource(ctx, db.GetGHSourceParams{
		WorkspaceID: workspaceID,
		ID:          id.ToInt(),
	})
	if err != nil {
		return GHSource{}, err
	}
	return GHSource{
		ID:             NewID(dbS.ID),
		WorkspaceID:    dbS.WorkspaceID,
		Owner:          dbS.Owner,
		Repo:           dbS.Repo,
		Path:           dbS.Path,
		InstallationID: dbS.InstallationID,
	}, nil
}

func (r *postgresRepo) ListGHSource(ctx context.Context, workspaceID string) ([]GHSource, error) {
	dbS, err := r.queries.ListGHSources(ctx, workspaceID)
	if err != nil {
		return []GHSource{}, err
	}

	sources := make([]GHSource, len(dbS))
	for i, s := range dbS {
		sources[i] = GHSource{
			ID:             NewID(s.ID),
			WorkspaceID:    s.WorkspaceID,
			Owner:          s.Owner,
			Repo:           s.Repo,
			Path:           s.Path,
			InstallationID: s.InstallationID,
		}
	}
	return sources, nil
}
