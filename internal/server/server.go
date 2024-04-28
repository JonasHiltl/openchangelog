package server

import (
	"fmt"
	"html/template"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/server/api"
	"github.com/jonashiltl/openchangelog/internal/source"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server interface {
	Start()
}

type server struct {
	e             *echo.Echo
	cfg           config.Config
	sourceFactory *source.SourceFactory
	parser        parse.Parser
	queries       *store.Queries
}

type ServerArgs struct {
	Parser  parse.Parser
	Cfg     config.Config
	Pool    *pgxpool.Pool
	Queries *store.Queries
}

func New(args ServerArgs) Server {
	e := echo.New()
	e.Renderer = newTemplate()
	e.Validator = newValidator()
	e.HTTPErrorHandler = customHTTPErrorHandler
	e.Use(middleware.Logger())
	e.Static("/static", "./web/css")

	srv := server{
		e:             e,
		cfg:           args.Cfg,
		sourceFactory: source.NewSourceFactory(args.Cfg),
		parser:        args.Parser,
		queries:       args.Queries,
	}

	api := api.New(args.Queries, args.Pool)

	e.GET("/", srv.renderChangeLog)
	e.GET("/:workspace-id/:changelog-id", srv.renderTenantChangeLog)

	e.POST("/api/changelogs", api.CreateChangelog)
	e.GET("/api/changelogs", api.ListChangelogs)
	e.GET("/api/changelogs/:changelog-id", api.GetChangelog)
	e.PATCH("/api/changelogs/:changelog-id", api.UpdateChangelog)
	e.DELETE("/api/changelogs/:changelog-id", api.DeleteChangelog)
	e.GET("/api/changelogs/:changelog-id/source", api.GetChangelogSource)
	e.PUT("/api/changelogs/:changelog-id/source", api.SetChangelogSource)
	e.DELETE("/api/changelogs/:changelog-id/source", api.DeleteChangelogSource)

	e.POST("/api/workspaces", api.CreateWorkspace)
	e.GET("/api/workspaces/:workspace-id", api.GetWorkspace)
	e.PATCH("/api/workspaces/:workspace-id", api.UpdateWorkspace)
	e.DELETE("/api/workspaces/:workspace-id", api.DeleteWorkspace)

	e.POST("/api/sources/gh", api.CreateGHSource)
	e.GET("/api/sources/gh", api.ListGHSources)
	e.GET("/api/sources/gh/:source-id", api.GetGHSource)
	e.DELETE("/api/sources/gh/:source-id", api.DeleteGHSource)

	return &srv
}

type changelogVars struct {
	PageSize int
	NextPage int
	HasMore  bool
	Title    string
	Subtitle string
	Logo     logo
	Articles []articleVars
}

func (v changelogVars) toMap() map[string]any {
	m := make(map[string]any)
	if v.PageSize != 0 {
		m["PageSize"] = v.PageSize
	}
	if v.NextPage != 0 {
		m["NextPage"] = v.NextPage
	}
	if v.Title != "" {
		m["Title"] = v.Title
	}
	if v.Subtitle != "" {
		m["Subtitle"] = v.Subtitle
	}
	if v.Logo.Src != "" {
		m["Logo"] = v.Logo
	}
	if len(v.Articles) > 0 {
		m["Articles"] = v.Articles
	}
	m["HasMore"] = v.HasMore
	return m
}

type articleVars struct {
	Id          string
	Title       string
	Description string
	PublishedAt string
	Content     template.HTML
}

type logo struct {
	Src    string
	Width  string
	Height string
	Alt    string
	Link   string
}

func (s *server) Start() {
	port := 8080
	if s.cfg.Port != 0 {
		port = 8080
	}
	addr := fmt.Sprintf(":%d", port)
	s.e.Logger.Fatal(s.e.Start(addr))
}
