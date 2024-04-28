package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/labstack/echo/v4"
)

type createChangelogBody struct {
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	LogoSrc    string `json:"logoSrc"`
	LogoLink   string `json:"logoLink"`
	LogoAlt    string `json:"logoAlt"`
	LogoHeight string `json:"logoHeight"`
	LogoWidth  string `json:"logoWidth"`
}

func (a *api) CreateChangelog(c echo.Context) error {
	s, err := a.bearerAuth(c)
	if err != nil {
		return err
	}

	var req createChangelogBody
	err = c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}

	sCl, err := a.queries.CreateChangelog(c.Request().Context(), store.CreateChangelogParams{
		WorkspaceID: s.WorkspaceID,
		Title:       pgtype.Text{Valid: req.Title != "", String: req.Title},
		Subtitle:    pgtype.Text{Valid: req.Subtitle != "", String: req.Subtitle},
		LogoSrc:     pgtype.Text{Valid: req.LogoSrc != "", String: req.LogoSrc},
		LogoLink:    pgtype.Text{Valid: req.LogoLink != "", String: req.LogoLink},
		LogoAlt:     pgtype.Text{Valid: req.LogoAlt != "", String: req.LogoAlt},
		LogoHeight:  pgtype.Text{Valid: req.LogoHeight != "", String: req.LogoHeight},
		LogoWidth:   pgtype.Text{Valid: req.LogoWidth != "", String: req.LogoWidth},
	})
	if err != nil {
		return err
	}

	cl := changelog{}
	cl.FromStore(sCl)
	return c.JSON(201, cl)
}

var errNoChangelog = echo.NewHTTPError(echo.ErrNotFound.Code, "changelog not found")

func (a *api) GetChangelog(c echo.Context) error {
	s, err := a.bearerAuth(c)
	if err != nil {
		return err
	}

	cId, err := strconv.ParseInt(c.Param("changelog-id"), 10, 64)
	if err != nil {
		return err
	}

	stored, err := a.queries.GetChangelog(c.Request().Context(), store.GetChangelogParams{
		WorkspaceID: s.WorkspaceID,
		ID:          cId,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return errNoChangelog
		}
		return err
	}
	cl := changelog{}
	cl.FromStore(stored)
	return c.JSON(200, cl)
}

func (a *api) ListChangelogs(c echo.Context) error {
	s, err := a.bearerAuth(c)
	if err != nil {
		return err
	}

	changelogs, err := a.queries.ListChangelogs(c.Request().Context(), s.WorkspaceID)
	if err != nil {
		if err != pgx.ErrNoRows {
			empty := make([]changelog, 0)
			return c.JSON(200, empty)
		}
		return err
	}
	res := make([]changelog, len(changelogs))
	for i, stored := range changelogs {
		c := changelog{}
		c.FromStore(stored)
		res[i] = c
	}
	return c.JSON(200, res)
}

func (a *api) DeleteChangelog(c echo.Context) error {
	s, err := a.bearerAuth(c)
	if err != nil {
		return err
	}

	cId, err := strconv.ParseInt(c.Param("changelog-id"), 10, 64)
	if err != nil {
		return err
	}

	err = a.queries.DeleteChangelog(c.Request().Context(), store.DeleteChangelogParams{
		WorkspaceID: s.WorkspaceID,
		ID:          cId,
	})
	if err != nil {
		return err
	}
	return c.NoContent(204)
}

func (a *api) UpdateChangelog(c echo.Context) error {
	s, err := a.bearerAuth(c)
	if err != nil {
		return err
	}

	cId, err := strconv.ParseInt(c.Param("changelog-id"), 10, 64)
	if err != nil {
		return err
	}

	var req createChangelogBody
	err = c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}

	updated, err := a.queries.UpdateChangelog(c.Request().Context(), store.UpdateChangelogParams{
		WorkspaceID: s.WorkspaceID,
		ID:          cId,
		Title:       pgtype.Text{Valid: req.Title != "", String: req.Title},
		Subtitle:    pgtype.Text{Valid: req.Subtitle != "", String: req.Subtitle},
		LogoSrc:     pgtype.Text{Valid: req.LogoSrc != "", String: req.LogoSrc},
		LogoLink:    pgtype.Text{Valid: req.LogoLink != "", String: req.LogoLink},
		LogoAlt:     pgtype.Text{Valid: req.LogoAlt != "", String: req.LogoAlt},
		LogoHeight:  pgtype.Text{Valid: req.LogoHeight != "", String: req.LogoHeight},
		LogoWidth:   pgtype.Text{Valid: req.LogoWidth != "", String: req.LogoWidth},
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, "changelog does not exist")
		}
		return err
	}

	cl := changelog{}
	cl.FromStore(updated)
	return c.JSON(200, cl)
}

var errNoSource = echo.NewHTTPError(echo.ErrNotFound.Code, "Changelog has no source")

func (a *api) GetChangelogSource(c echo.Context) error {
	s, err := a.bearerAuth(c)
	if err != nil {
		return err
	}

	cId, err := strconv.ParseInt(c.Param("changelog-id"), 10, 64)
	if err != nil {
		return err
	}

	stored, err := a.queries.GetChangelogSource(c.Request().Context(), store.GetChangelogSourceParams{
		WorkspaceID: s.WorkspaceID,
		ID:          cId,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return errNoSource
		}
		return err
	}

	if stored.Changelog.SourceType.Valid {
		switch stored.Changelog.SourceType.SourceType {
		case store.SourceTypeGitHub:
			g := ghSource{}
			g.FromChangelogSource(stored.ChangelogSource)
			return c.JSON(200, g)
		}
	}
	return errNoSource
}

type setChangelogSourceRequest struct {
	ChangelogID int64      `param:"changelog-id" validate:"required"`
	SourceID    int64      `json:"id" validate:"required"`
	Type        SourceType `json:"type" validate:"oneof=GitHub"`
}

func (a *api) SetChangelogSource(c echo.Context) error {
	s, err := a.bearerAuth(c)
	if err != nil {
		return err
	}

	var req setChangelogSourceRequest
	err = c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}

	err = c.Validate(req)
	if err != nil {
		return err
	}

	tx, err := a.pool.Begin(c.Request().Context())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	qtx := a.queries.WithTx(tx)

	switch req.Type {
	case SourceTypeGitHub:
		_, err = qtx.GetGHSource(c.Request().Context(), store.GetGHSourceParams{
			WorkspaceID: s.WorkspaceID,
			ID:          req.SourceID,
		})
		if err != nil {
			return echo.NewHTTPError(echo.ErrBadRequest.Code, fmt.Sprintf("source with id %d does not exist", req.SourceID))
		}
	}

	err = qtx.UpdateChangelogSource(c.Request().Context(), store.UpdateChangelogSourceParams{
		SourceID:    pgtype.Int8{Valid: req.SourceID != 0, Int64: req.SourceID},
		SourceType:  store.NullSourceType{Valid: req.Type != "", SourceType: store.SourceType(req.Type)},
		WorkspaceID: s.WorkspaceID,
		ID:          req.ChangelogID,
	})
	if err != nil {
		return err
	}

	err = tx.Commit(c.Request().Context())
	if err != nil {
		return err
	}

	return c.NoContent(204)
}

func (a *api) DeleteChangelogSource(c echo.Context) error {
	s, err := a.bearerAuth(c)
	if err != nil {
		return err
	}

	cId, err := strconv.ParseInt(c.Param("changelog-id"), 10, 64)
	if err != nil {
		return err
	}

	err = a.queries.UpdateChangelogSource(c.Request().Context(), store.UpdateChangelogSourceParams{
		WorkspaceID: s.WorkspaceID,
		ID:          cId,
		SourceID:    pgtype.Int8{Valid: false},
		SourceType:  store.NullSourceType{Valid: false},
	})
	if err != nil {
		return err
	}
	return c.NoContent(204)
}
