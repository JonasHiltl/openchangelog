package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/jonashiltl/openchangelog/source"
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

	vars := changelogVars{
		Articles: articles,
		PageSize: req.PageSize,
		NextPage: req.Page + 1,
		HasMore:  res.HasMore,
	}

	if s.cfg.Page != nil {
		vars.Title = s.cfg.Page.Title
		vars.Subtitle = s.cfg.Page.Subtitle
	}

	if s.cfg.Page != nil && s.cfg.Page.Logo != nil {
		vars.Logo = logo{
			Src:  s.cfg.Page.Logo.Src,
			Link: s.cfg.Page.Logo.Link,
			Alt:  s.cfg.Page.Logo.Alt,
		}
	}

	// this is used by htmx to allow infinite scrolling
	if htmxHeader := c.Request().Header.Get("HX-Request"); len(htmxHeader) > 0 {
		if len(articles) > 0 {
			return c.Render(200, "changelog/article_list", vars.toMap())
		}
	} else {
		return c.Render(200, "changelog/index", vars.toMap())
	}

	return c.NoContent(200)
}
