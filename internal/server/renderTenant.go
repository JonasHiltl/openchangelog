package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jonashiltl/openchangelog/internal/source"
	"github.com/jonashiltl/openchangelog/internal/store"
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

	vars := changelogVars{
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

		articles := make([]articleVars, 0, len(res.Articles))
		for _, a := range res.Articles {
			articles = append(articles, articleVars{
				Id:          fmt.Sprint(a.Meta.PublishedAt.Unix()),
				Title:       a.Meta.Title,
				Description: a.Meta.Description,
				PublishedAt: a.Meta.PublishedAt.Format("02 Jan 2006"),
				Content:     template.HTML(a.Content.String()),
			})
		}
		vars.Articles = articles
		vars.HasMore = res.HasMore
	}

	if row.Changelog.Title.Valid {
		vars.Title = row.Changelog.Title.String
	}
	if row.Changelog.Subtitle.Valid {
		vars.Subtitle = row.Changelog.Subtitle.String
	}

	if row.Changelog.LogoSrc.Valid {
		vars.Logo = logo{
			Src:    row.Changelog.LogoSrc.String,
			Width:  row.Changelog.LogoWidth.String,
			Height: row.Changelog.LogoHeight.String,
			Link:   row.Changelog.LogoLink.String,
			Alt:    row.Changelog.LogoAlt.String,
		}
	}

	// this is used by htmx to allow infinite scrolling
	if htmxHeader := c.Request().Header.Get("HX-Request"); len(htmxHeader) > 0 {
		if len(vars.Articles) > 0 {
			return c.Render(200, "changelog/article_list", vars.toMap())
		}
	} else {
		return c.Render(200, "changelog/index", vars.toMap())
	}

	return c.NoContent(200)
}
