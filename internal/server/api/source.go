package api

import (
	"strconv"

	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/labstack/echo/v4"
)

type createGHSourceBody struct {
	Owner          string `json:"owner" validate:"required"`
	Repo           string `json:"repo" validate:"required"`
	Path           string `json:"path"`
	InstallationID int64  `json:"installationID" validate:"required"`
}

func (a *api) CreateGHSource(c echo.Context) error {
	s, err := a.bearerAuth(c)
	if err != nil {
		return err
	}

	var req createGHSourceBody
	err = c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}

	err = c.Validate(req)
	if err != nil {
		return err
	}

	stored, err := a.queries.CreateGHSource(c.Request().Context(), store.CreateGHSourceParams{
		WorkspaceID:    s.WorkspaceID,
		Owner:          req.Owner,
		Repo:           req.Repo,
		Path:           req.Path,
		InstallationID: req.InstallationID,
	})
	if err != nil {
		return err
	}

	g := ghSource{}
	g.FromStore(stored)
	return c.JSON(201, g)
}

func (a *api) GetGHSource(c echo.Context) error {
	s, err := a.bearerAuth(c)
	if err != nil {
		return err
	}

	sId, err := strconv.ParseInt(c.Param("source-id"), 10, 64)
	if err != nil {
		return err
	}

	stored, err := a.queries.GetGHSource(c.Request().Context(), store.GetGHSourceParams{
		WorkspaceID: s.WorkspaceID,
		ID:          sId,
	})
	if err != nil {
		return err
	}

	g := ghSource{}
	g.FromStore(stored)
	return c.JSON(200, g)
}

func (a *api) ListGHSources(c echo.Context) error {
	s, err := a.bearerAuth(c)
	if err != nil {
		return err
	}

	stored, err := a.queries.ListGHSources(c.Request().Context(), s.WorkspaceID)
	if err != nil {
		return err
	}

	res := make([]ghSource, len(stored))
	for i, source := range stored {
		g := ghSource{}
		g.FromStore(source)
		res[i] = g
	}
	return c.JSON(200, res)
}

func (a *api) DeleteGHSource(c echo.Context) error {
	s, err := a.bearerAuth(c)
	if err != nil {
		return err
	}

	sId, err := strconv.ParseInt(c.Param("source-id"), 10, 64)
	if err != nil {
		return err
	}

	err = a.queries.DeleteGHSource(c.Request().Context(), store.DeleteGHSourceParams{
		WorkspaceID: s.WorkspaceID,
		ID:          sId,
	})
	if err != nil {
		return err
	}
	return c.NoContent(204)
}
