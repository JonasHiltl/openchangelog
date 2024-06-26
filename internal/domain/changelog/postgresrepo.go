package changelog

import (
	"context"
	"errors"
	"fmt"

	"github.com/guregu/null/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jonashiltl/openchangelog/internal/domain"
	"github.com/jonashiltl/openchangelog/internal/domain/changelog/db"
	"github.com/jonashiltl/openchangelog/internal/domain/source"
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

func (r *postgresRepo) CreateChangelog(ctx context.Context, c Changelog) (Changelog, error) {
	dbC, err := r.queries.CreateChangelog(ctx, db.CreateChangelogParams{
		WorkspaceID: c.WorkspaceID,
		Title:       pgtype.Text{Valid: c.Title != "", String: c.Title},
		Subtitle:    pgtype.Text{Valid: c.Subtitle != "", String: c.Subtitle},
		LogoSrc:     pgtype.Text(c.Logo.Src.NullString),
		LogoLink:    pgtype.Text(c.Logo.Link.NullString),
		LogoAlt:     pgtype.Text(c.Logo.Alt.NullString),
		LogoHeight:  pgtype.Text(c.Logo.Height.NullString),
		LogoWidth:   pgtype.Text(c.Logo.Width.NullString),
	})
	if err != nil {
		return Changelog{}, err
	}

	if c.Source != nil {
		switch c.Source.Type() {
		case source.GitHub:
			s := c.Source.(source.GHSource)
			err := r.queries.UpdateChangelogSource(ctx, db.UpdateChangelogSourceParams{
				WorkspaceID: c.WorkspaceID,
				ID:          c.ID.ToInt(),
				SourceID:    pgtype.Int8{Int64: s.ID.ToInt(), Valid: true},
				SourceType:  db.NullSourceType{SourceType: db.SourceTypeGitHub, Valid: true},
			})
			if err != nil {
				return Changelog{}, err
			}
		case source.Local:
			fmt.Println("pg repo does not support updating a local data source")
		case source.String:
			fmt.Println("pg repo does not support updating a string data source")
		}
	}

	return Changelog{
		ID:          NewID(dbC.ID),
		WorkspaceID: dbC.WorkspaceID,
		Title:       dbC.Title.String,
		Subtitle:    dbC.Subtitle.String,
		Logo:        c.Logo,
		Source:      c.Source,
		CreatedAt:   dbC.CreatedAt.Time,
	}, nil
}

func (r *postgresRepo) UpdateChangelog(ctx context.Context, c Changelog) (Changelog, error) {
	dbC, err := r.queries.UpdateChangelog(ctx, db.UpdateChangelogParams{
		ID:          c.ID.ToInt(),
		WorkspaceID: c.WorkspaceID,
		Title:       pgtype.Text{Valid: c.Title != "", String: c.Title},
		Subtitle:    pgtype.Text{Valid: c.Subtitle != "", String: c.Subtitle},
		LogoSrc:     pgtype.Text(c.Logo.Src.NullString),
		LogoLink:    pgtype.Text(c.Logo.Link.NullString),
		LogoAlt:     pgtype.Text(c.Logo.Alt.NullString),
		LogoHeight:  pgtype.Text(c.Logo.Height.NullString),
		LogoWidth:   pgtype.Text(c.Logo.Width.NullString),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return Changelog{}, domain.NewError(domain.ErrNotFound, errors.New("changelog not found"))
		}
		return Changelog{}, err
	}

	if c.Source != nil {
		switch c.Source.Type() {
		case source.GitHub:
			s := c.Source.(source.GHSource)
			err := r.queries.UpdateChangelogSource(ctx, db.UpdateChangelogSourceParams{
				WorkspaceID: c.WorkspaceID,
				ID:          c.ID.ToInt(),
				SourceID:    pgtype.Int8{Int64: s.ID.ToInt(), Valid: true},
				SourceType:  db.NullSourceType{SourceType: db.SourceTypeGitHub, Valid: true},
			})
			if err != nil {
				return Changelog{}, err
			}
		case source.Local:
			fmt.Println("pg repo does not support updating a local data source")
		case source.String:
			fmt.Println("pg repo does not support updating a string data source")
		}
	}

	return Changelog{
		ID:          NewID(dbC.ID),
		WorkspaceID: dbC.WorkspaceID,
		Title:       dbC.Title.String,
		Subtitle:    dbC.Subtitle.String,
		Logo:        c.Logo,
		Source:      c.Source,
		CreatedAt:   dbC.CreatedAt.Time,
	}, nil
}

func (r *postgresRepo) GetChangelog(ctx context.Context, workspaceID string, id ID) (Changelog, error) {
	s, err := r.queries.GetChangelog(ctx, db.GetChangelogParams{
		WorkspaceID: workspaceID,
		ID:          id.ToInt(),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return Changelog{}, domain.NewError(domain.ErrNotFound, errors.New("changelog not found"))
		}
		return Changelog{}, err
	}
	return pgRowToChangelog(s.Changelog, s.ChangelogSource), nil
}

func (r *postgresRepo) DeleteChangelog(ctx context.Context, workspaceID string, id ID) error {
	return r.queries.DeleteChangelog(ctx, db.DeleteChangelogParams{
		WorkspaceID: workspaceID,
		ID:          id.ToInt(),
	})
}

func (r *postgresRepo) ListChangelogs(ctx context.Context, workspaceID string) ([]Changelog, error) {
	rows, err := r.queries.ListChangelogs(ctx, workspaceID)
	if err != nil {
		return []Changelog{}, err
	}

	res := make([]Changelog, len(rows))
	for i, row := range rows {
		res[i] = pgRowToChangelog(row.Changelog, row.ChangelogSource)
	}
	return res, nil
}

func pgRowToChangelog(cl db.Changelog, src db.ChangelogSource) Changelog {
	c := Changelog{
		ID:          NewID(cl.ID),
		WorkspaceID: cl.WorkspaceID,
		Title:       cl.Title.String,
		Subtitle:    cl.Subtitle.String,
		Logo: struct {
			Src    null.String
			Link   null.String
			Alt    null.String
			Height null.String
			Width  null.String
		}{
			Src:    null.NewString(cl.LogoSrc.String, cl.LogoSrc.Valid),
			Link:   null.NewString(cl.LogoLink.String, cl.LogoLink.Valid),
			Alt:    null.NewString(cl.LogoAlt.String, cl.LogoAlt.Valid),
			Height: null.NewString(cl.LogoHeight.String, cl.LogoHeight.Valid),
			Width:  null.NewString(cl.LogoWidth.String, cl.LogoWidth.Valid),
		},
		CreatedAt: cl.CreatedAt.Time,
	}

	if src.Owner.Valid && src.Repo.Valid {
		c.Source = source.GHSource{
			ID:             source.NewID(src.ID.Int64),
			WorkspaceID:    src.WorkspaceID.String,
			Owner:          src.Owner.String,
			Repo:           src.Repo.String,
			Path:           src.Path.String,
			InstallationID: src.InstallationID.Int64,
		}
	}
	return c
}
