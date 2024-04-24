package server

import (
	"fmt"
	"html/template"
	"strconv"

	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/source"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server interface {
	Start()
}

type server struct {
	e      *echo.Echo
	cfg    config.Config
	source source.Source
	parser parse.Parser
}

func New(s source.Source, p parse.Parser, cfg config.Config) Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/static", "./web/css")
	e.Renderer = newTemplate()

	srv := server{
		e:      e,
		cfg:    cfg,
		source: s,
		parser: p,
	}

	e.GET("/", srv.renderChangeLog)

	return &srv
}

type ArticleVars struct {
	Id          string
	Title       string
	Description string
	PublishedAt string
	Content     template.HTML
}

type Logo struct {
	Src    string
	Width  string
	Height string
	Alt    string
	Link   string
}

func (s *server) renderChangeLog(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.QueryParam("page-size"))
	if err != nil {
		pageSize = 10
	}
	res, err := s.parser.Parse(c.Request().Context(), s.source, source.NewLoadParams(pageSize, page))
	if err != nil {
		return err
	}

	articles := make([]ArticleVars, 0, len(res.Articles))
	for _, a := range res.Articles {
		articles = append(articles, ArticleVars{
			Id:          fmt.Sprint(a.Meta.PublishedAt.Unix()),
			Title:       a.Meta.Title,
			Description: a.Meta.Description,
			PublishedAt: a.Meta.PublishedAt.Format("02 Jan 2006"),
			Content:     template.HTML(a.Content.String()),
		})
	}

	vars := map[string]interface{}{
		"Articles": articles,
		"PageSize": pageSize,
		"NextPage": page + 1,
		"HasMore":  res.HasMore,
	}

	if s.cfg.Page != nil {
		vars["Title"] = s.cfg.Page.Title
		vars["Subtitle"] = s.cfg.Page.Subtitle
	}

	if s.cfg.Page != nil && s.cfg.Page.Logo != nil {
		vars["Logo"] = Logo{
			Src:    s.cfg.Page.Logo.Src,
			Width:  s.cfg.Page.Logo.Width,
			Height: s.cfg.Page.Logo.Height,
			Link:   s.cfg.Page.Logo.Link,
			Alt:    s.cfg.Page.Logo.Alt,
		}
	}

	// this is used by htmx to allow infinite scrolling
	if htmxHeader := c.Request().Header.Get("HX-Request"); len(htmxHeader) > 0 {
		if len(articles) > 0 {
			return c.Render(200, "changelog/article_list", vars)
		}
	} else {
		return c.Render(200, "changelog/index", vars)
	}

	return c.NoContent(200)
}

func (s *server) Start() {
	port := 8080
	if s.cfg.Port != 0 {
		port = 8080
	}
	addr := fmt.Sprintf(":%d", port)
	s.e.Logger.Fatal(s.e.Start(addr))
}
