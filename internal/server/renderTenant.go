package server

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jonashiltl/openchangelog/components"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/source"
	"github.com/jonashiltl/openchangelog/web/views"
	"github.com/labstack/echo/v4"
)

type renderTenantRequest struct {
	WorkspaceID string `param:"workspace-id"`
	ChangelogID int64  `param:"changelog-id"`
	Page        int    `query:"page"`
	PageSize    int    `query:"page-size"`
}

func (s *server) renderTenantChangeLog(c echo.Context) error {
	var req renderTenantRequest
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	row, err := s.queries.GetChangelogSource(c.Request().Context(), store.GetChangelogSourceParams{
		WorkspaceID: req.WorkspaceID,
		ID:          req.ChangelogID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(echo.ErrNotFound.Code, "Changelog not found")
		}
		return err
	}

	articleListArgs := components.ArticleListArgs{
		Articles: []components.ArticleArgs{},
		PageSize: req.PageSize,
		NextPage: req.Page + 1,
		HasMore:  false,
	}

	src, err := s.sourceFactory.FromDB(row.Changelog, row.ChangelogSource)
	if err == nil {
		res, err := s.parser.Parse(c.Request().Context(), src, source.NewLoadParams(req.PageSize, req.Page))
		if err != nil {
			return echo.NewHTTPError(echo.ErrInternalServerError.Code, "failed to parse markdown fiels: %s", err)
		}

		articles := make([]components.ArticleArgs, 0, len(res.Articles))
		for _, a := range res.Articles {
			articles = append(articles, components.ArticleArgs{
				ID:          fmt.Sprint(a.Meta.PublishedAt.Unix()),
				Title:       a.Meta.Title,
				Description: a.Meta.Description,
				PublishedAt: a.Meta.PublishedAt.Format("02 Jan 2006"),
				Content:     a.Content.String(),
			})
		}
		articleListArgs.Articles = articles
		articleListArgs.HasMore = res.HasMore
	}

	if htmxHeader := c.Request().Header.Get("HX-Request"); len(htmxHeader) > 0 {
		if len(articleListArgs.Articles) > 0 {
			return components.ArticleList(articleListArgs).Render(c.Request().Context(), c.Response().Writer)
		} else {
			return c.NoContent(204)
		}
	}

	indexArgs := views.IndexArgs{
		ChangelogArgs: components.ChangelogArgs{
			Title:           row.Changelog.Title.String,
			Subtitle:        row.Changelog.Subtitle.String,
			ArticleListArgs: articleListArgs,
		},
		NavbarArgs: components.NavbarArgs{
			Logo: components.Logo{
				Src:    row.Changelog.LogoSrc.String,
				Width:  row.Changelog.LogoWidth.String,
				Height: row.Changelog.LogoHeight.String,
				Alt:    row.Changelog.LogoLink.String,
				Link:   row.Changelog.LogoAlt.String,
			},
		},
	}

	return views.Index(indexArgs).Render(c.Request().Context(), c.Response().Writer)
}
