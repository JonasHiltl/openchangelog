package api

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/internal/utils"
	"github.com/labstack/echo/v4"
)

type createWorkspaceBody struct {
	Name string `json:"name" validate:"required"`
}

func (a *api) CreateWorkspace(c echo.Context) error {
	var req createWorkspaceBody
	err := c.Bind(&req)
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
	stored, err := qtx.CreateWorkspace(c.Request().Context(), store.CreateWorkspaceParams{
		ID:   utils.NewWorkspaceID(),
		Name: req.Name,
	})
	if err != nil {
		return err
	}

	token, err := qtx.CreateToken(c.Request().Context(), store.CreateTokenParams{
		Key:         utils.NewSessionID(),
		WorkspaceID: stored.ID,
	})
	if err != nil {
		return err
	}

	err = tx.Commit(c.Request().Context())
	if err != nil {
		return err
	}

	wt := workspaceWithToken{}
	wt.FromStore(stored, token)
	return c.JSON(201, wt)
}

func (a *api) UpdateWorkspace(c echo.Context) error {
	s, err := a.bearerAuth(c)
	if err != nil {
		return err
	}

	wId := c.Param("workspace-id")
	if wId != s.WorkspaceID {
		return errNoWorkspace
	}

	var req createWorkspaceBody
	err = c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	stored, err := a.queries.UpdateWorkspace(c.Request().Context(), store.UpdateWorkspaceParams{
		ID:   s.WorkspaceID,
		Name: pgtype.Text{Valid: req.Name != "", String: req.Name},
	})
	if err != nil {
		return err
	}

	w := workspace{}
	w.FromStore(stored)
	return c.JSON(200, w)
}

var errNoWorkspace = echo.NewHTTPError(echo.ErrNotFound.Code, "workspace not found")

func (a *api) GetWorkspace(c echo.Context) error {
	s, err := a.bearerAuth(c)
	if err != nil {
		return err
	}

	wId := c.Param("workspace-id")
	if err != nil {
		return err
	}
	if wId != s.WorkspaceID {
		return errNoWorkspace
	}

	stored, err := a.queries.GetWorkspace(c.Request().Context(), s.WorkspaceID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errNoWorkspace
		}
		return err
	}

	wt := workspaceWithToken{}
	wt.FromStore(stored.Workspace, stored.Token)
	return c.JSON(200, wt)
}

func (a *api) DeleteWorkspace(c echo.Context) error {
	s, err := a.bearerAuth(c)
	if err != nil {
		return err
	}

	wId := c.Param("workspace-id")
	if wId != s.WorkspaceID {
		return echo.NewHTTPError(echo.ErrForbidden.Code, "Not allowed to delete this workspace")
	}

	err = a.queries.DeleteWorkspace(c.Request().Context(), s.WorkspaceID)
	if err != nil {
		return err
	}
	return c.NoContent(204)
}
