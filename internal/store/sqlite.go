package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jonashiltl/openchangelog/internal/errs"
	_ "github.com/mattn/go-sqlite3"

	"github.com/guregu/null/v5"
)

func (cl changelog) ToExported(source changelogSource) Changelog {
	c := Changelog{
		WorkspaceID: WorkspaceID(cl.WorkspaceID),
		ID:          ChangelogID(cl.ID),
		Title:       null.NewString(cl.Title.String, cl.Title.Valid),
		Subtitle:    null.NewString(cl.Subtitle.String, cl.Subtitle.Valid),
		LogoSrc:     null.NewString(cl.LogoSrc.String, cl.LogoSrc.Valid),
		LogoLink:    null.NewString(cl.LogoLink.String, cl.LogoLink.Valid),
		LogoAlt:     null.NewString(cl.LogoAlt.String, cl.LogoAlt.Valid),
		LogoHeight:  null.NewString(cl.LogoHeight.String, cl.LogoHeight.Valid),
		LogoWidth:   null.NewString(cl.LogoWidth.String, cl.LogoWidth.Valid),
		CreatedAt:   time.Unix(cl.CreatedAt, 0),
		GHSource:    null.NewValue(GHSource{}, false),
	}

	if source.ID.Valid && source.WorkspaceID.Valid {
		c.GHSource = null.NewValue(GHSource{
			ID:             GHSourceID(source.ID.String),
			WorkspaceID:    WorkspaceID(source.WorkspaceID.String),
			Owner:          source.Owner.String,
			Repo:           source.Repo.String,
			Path:           source.Path.String,
			InstallationID: source.InstallationID.Int64,
		}, true)
	}
	return c
}

func (gh ghSource) ToExported() GHSource {
	return GHSource{
		ID:             GHSourceID(gh.ID),
		WorkspaceID:    WorkspaceID(gh.WorkspaceID),
		Owner:          gh.Owner,
		Repo:           gh.Repo,
		Path:           gh.Path,
		InstallationID: gh.InstallationID,
	}
}

func NewSQLiteStore(conn string) (Store, error) {
	db, err := sql.Open("sqlite3", conn)
	if err != nil {
		return nil, err
	}

	q := New(db)

	return &sqlite{
		q:  q,
		db: db,
	}, nil
}

type sqlite struct {
	q  *Queries
	db *sql.DB
}

func (s *sqlite) CreateChangelog(ctx context.Context, cl Changelog) (Changelog, error) {
	c, err := s.q.createChangelog(ctx, createChangelogParams{
		ID:          cl.ID.String(),
		WorkspaceID: cl.WorkspaceID.String(),
		Title:       cl.Title.NullString,
		Subtitle:    cl.Subtitle.NullString,
		LogoSrc:     cl.LogoSrc.NullString,
		LogoLink:    cl.LogoLink.NullString,
		LogoAlt:     cl.LogoAlt.NullString,
		LogoHeight:  cl.LogoHeight.NullString,
		LogoWidth:   cl.LogoWidth.NullString,
	})
	if err != nil {
		return Changelog{}, err
	}

	// TODO get source
	return c.ToExported(changelogSource{}), nil
}

func (s *sqlite) GetChangelog(ctx context.Context, wID WorkspaceID, cID ChangelogID) (Changelog, error) {
	cl, err := s.q.getChangelog(ctx, getChangelogParams{
		WorkspaceID: wID.String(),
		ID:          cID.String(),
	})
	if err != nil {
		return Changelog{}, err
	}

	return cl.changelog.ToExported(cl.ChangelogSource), nil
}

func (s *sqlite) ListChangelogs(ctx context.Context, wID WorkspaceID) ([]Changelog, error) {
	cls, err := s.q.listChangelogs(ctx, wID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]Changelog, 0), nil
		}
		return nil, err
	}

	res := make([]Changelog, len(cls))
	for i, cl := range cls {
		res[i] = cl.changelog.ToExported(cl.ChangelogSource)
	}
	return res, nil
}

func (s *sqlite) UpdateChangelog(ctx context.Context, wID WorkspaceID, cID ChangelogID, args UpdateChangelogArgs) (Changelog, error) {
	c, err := s.q.updateChangelog(ctx, updateChangelogParams{
		ID:          cID.String(),
		WorkspaceID: wID.String(),
		Title:       args.Title.NullString,
		Subtitle:    args.Subtitle.NullString,
		LogoSrc:     args.LogoSrc.NullString,
		LogoLink:    args.LogoLink.NullString,
		LogoAlt:     args.LogoAlt.NullString,
		LogoHeight:  args.LogoHeight.NullString,
		LogoWidth:   args.LogoWidth.NullString,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Changelog{}, errs.NewError(errs.ErrNotFound, errors.New("changelog not found"))
		}
		return Changelog{}, err
	}
	return c.ToExported(changelogSource{}), nil
}

func (s *sqlite) DeleteChangelog(ctx context.Context, wID WorkspaceID, cID ChangelogID) error {
	return s.q.deleteChangelog(ctx, deleteChangelogParams{
		WorkspaceID: wID.String(),
		ID:          cID.String(),
	})
}

func (s *sqlite) SetChangelogGHSource(ctx context.Context, wID WorkspaceID, cID ChangelogID, ghID GHSourceID) error {
	return s.q.setChangelogSource(ctx, setChangelogSourceParams{
		SourceID:    sql.NullString{String: ghID.String(), Valid: true},
		WorkspaceID: wID.String(),
		ID:          cID.String(),
	})
}

func (s *sqlite) DeleteChangelogSource(ctx context.Context, wID WorkspaceID, cID ChangelogID) error {
	return s.q.deleteChangelogSource(ctx, deleteChangelogSourceParams{
		WorkspaceID: wID.String(),
		ID:          cID.String(),
	})
}

func (s *sqlite) SaveWorkspace(ctx context.Context, ws Workspace) (Workspace, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return Workspace{}, err
	}
	defer tx.Rollback()
	q := s.q.WithTx(tx)

	c, err := q.saveWorkspace(ctx, saveWorkspaceParams{
		ID:   ws.ID.String(),
		Name: ws.Name,
	})
	if err != nil {
		return Workspace{}, err
	}

	if ws.Token != "" {
		err := q.createToken(ctx, createTokenParams{
			Key:         ws.Token.String(),
			WorkspaceID: ws.ID.String(),
		})
		if err != nil {
			return Workspace{}, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return Workspace{}, err
	}

	return Workspace{
		ID:    WorkspaceID(c.ID),
		Name:  c.Name,
		Token: ws.Token,
	}, nil
}

func (s *sqlite) GetWorkspace(ctx context.Context, wID WorkspaceID) (Workspace, error) {
	row, err := s.q.getWorkspace(ctx, wID.String())
	if err != nil {
		return Workspace{}, err
	}
	return Workspace{
		ID:    WorkspaceID(row.workspace.ID),
		Name:  row.workspace.Name,
		Token: Token(row.token.Key),
	}, nil
}

func (s *sqlite) GetWorkspaceIDByToken(ctx context.Context, token string) (WorkspaceID, error) {
	row, err := s.q.getToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.NewError(errs.ErrUnauthorized, errors.New("invalid bearer token"))
		}
		log.Println("failed to get token")
		return "", err
	}
	return WorkspaceID(row.WorkspaceID), nil
}

func (s *sqlite) DeleteWorkspace(ctx context.Context, wID WorkspaceID) error {
	return s.q.deleteWorkspace(ctx, wID.String())
}

func (s *sqlite) CreateGHSource(ctx context.Context, gh GHSource) (GHSource, error) {
	row, err := s.q.createGHSource(ctx, createGHSourceParams{
		WorkspaceID:    gh.WorkspaceID.String(),
		ID:             gh.ID.String(),
		Owner:          gh.Owner,
		Repo:           gh.Repo,
		Path:           gh.Path,
		InstallationID: gh.InstallationID,
	})
	if err != nil {
		return GHSource{}, err
	}
	return row.ToExported(), nil
}

func (s *sqlite) DeleteGHSource(ctx context.Context, wID WorkspaceID, ghID GHSourceID) error {
	return s.q.deleteGHSource(ctx, deleteGHSourceParams{
		WorkspaceID: wID.String(),
		ID:          ghID.String(),
	})
}

func (s *sqlite) GetGHSource(ctx context.Context, wID WorkspaceID, ghID GHSourceID) (GHSource, error) {
	row, err := s.q.getGHSource(ctx, getGHSourceParams{
		WorkspaceID: wID.String(),
		ID:          ghID.String(),
	})
	if err != nil {
		return GHSource{}, err
	}
	return row.ToExported(), nil
}

func (s *sqlite) ListGHSources(ctx context.Context, wID WorkspaceID) ([]GHSource, error) {
	rows, err := s.q.listGHSources(ctx, wID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]GHSource, 0), nil
		}
		return nil, err
	}

	sources := make([]GHSource, len(rows))
	for i, row := range rows {
		sources[i] = row.ToExported()
	}
	return sources, nil
}
