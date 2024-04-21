package server

import (
	"fmt"
	"html/template"

	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/source"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server interface {
	Start()
}

type ServerCfg struct {
	Port int
}

func WithPort(p int) func(c *ServerCfg) {
	return func(c *ServerCfg) {
		c.Port = p
	}
}

func newConfig() ServerCfg {
	return ServerCfg{
		Port: 80,
	}
}

type server struct {
	e      *echo.Echo
	cfg    ServerCfg
	source source.Source
	parser parse.Parser
}

func New(s source.Source, p parse.Parser, config ...func(c *ServerCfg)) Server {
	c := newConfig()
	for _, apply := range config {
		apply(&c)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/static", "./web/css")
	e.Renderer = newTemplate()

	srv := server{
		e:      e,
		cfg:    c,
		source: s,
		parser: p,
	}

	e.GET("/", srv.renderChangeLog)

	return &srv
}

type ArticleVars struct {
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
	res, err := s.parser.Parse(s.source)
	if err != nil {
		return err
	}

	articles := make([]ArticleVars, 0, len(res))
	for _, a := range res {
		articles = append(articles, ArticleVars{
			Title:       a.Meta.Title,
			Description: a.Meta.Description,
			PublishedAt: a.Meta.PublishedAt.Format("02 Jan 2006"),
			Content:     template.HTML(a.Content.String()),
		})
	}

	vars := map[string]interface{}{
		"Articles": articles,
		"Logo": Logo{
			Src:    "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_92x30dp.png",
			Width:  "70px",
			Height: "25px",
			Link:   "https://www.google.com",
		},
	}

	return c.Render(200, "index", vars)
}

func (s *server) Start() {
	addr := fmt.Sprintf(":%d", s.cfg.Port)
	s.e.Logger.Fatal(s.e.Start(addr))
}
