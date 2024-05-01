package server

import (
	"fmt"
	"net/http"

	"github.com/jonashiltl/openchangelog/components"
	"github.com/jonashiltl/openchangelog/source"
	"github.com/jonashiltl/openchangelog/web/views"
	"github.com/labstack/echo/v4"
)

type renderRequest struct {
	Page     int `query:"page"`
	PageSize int `query:"page-size"`
}

func (s *server) renderChangeLog(c echo.Context) error {
	var req renderRequest
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

	src, err := s.sourceFactory.FromConfig()
	if err != nil {
		return err
	}

	res, err := s.parser.Parse(c.Request().Context(), src, source.NewLoadParams(req.PageSize, req.Page))
	if err != nil {
		return err
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

	articleListArgs := components.ArticleListArgs{
		Articles: articles,
		PageSize: req.PageSize,
		NextPage: req.Page + 1,
		HasMore:  res.HasMore,
	}

	if htmxHeader := c.Request().Header.Get("HX-Request"); len(htmxHeader) > 0 {
		if len(articles) > 0 {
			return components.ArticleList(articleListArgs).Render(c.Request().Context(), c.Response().Writer)
		} else {
			return c.NoContent(204)
		}
	}

	indexArgs := views.IndexArgs{
		ChangelogArgs: components.ChangelogArgs{
			Title:           s.cfg.Page.Title,
			Subtitle:        s.cfg.Page.Subtitle,
			ArticleListArgs: articleListArgs,
		},
		NavbarArgs: components.NavbarArgs{
			Logo: components.Logo{
				Src:    s.cfg.Page.Logo.Src,
				Width:  s.cfg.Page.Logo.Width,
				Height: s.cfg.Page.Logo.Height,
				Alt:    s.cfg.Page.Logo.Alt,
				Link:   s.cfg.Page.Logo.Link,
			},
		},
	}

	return views.Index(indexArgs).Render(c.Request().Context(), c.Response().Writer)
}
