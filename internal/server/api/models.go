package api

import (
	"time"

	"github.com/jonashiltl/openchangelog/internal/store"
)

type changelog struct {
	ID          int64     `json:"id"`
	WorkspaceID string    `json:"workspaceID"`
	Title       string    `json:"title,omitempty"`
	Subtitle    string    `json:"subtitle,omitempty"`
	LogoSrc     string    `json:"logoSrc,omitempty"`
	LogoLink    string    `json:"logoLink,omitempty"`
	LogoAlt     string    `json:"logoAlt,omitempty"`
	LogoHeight  string    `json:"logoHeight,omitempty"`
	LogoWidth   string    `json:"logoWidth,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (c *changelog) FromStore(cl store.Changelog) {
	c.ID = cl.ID
	c.WorkspaceID = cl.WorkspaceID
	c.Title = cl.Title.String
	c.Subtitle = cl.Subtitle.String
	c.LogoSrc = cl.LogoSrc.String
	c.LogoLink = cl.LogoLink.String
	c.LogoAlt = cl.LogoAlt.String
	c.LogoHeight = cl.LogoHeight.String
	c.LogoWidth = cl.LogoWidth.String
	c.CreatedAt = cl.CreatedAt.Time.UTC()
}

type workspace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (w *workspace) FromStore(stored store.Workspace) {
	w.ID = stored.ID
	w.Name = stored.Name
}

type workspaceWithToken struct {
	Workspace workspace `json:"workspace"`
	Token     string    `json:"token"`
}

func (wt *workspaceWithToken) FromStore(stored store.Workspace, token store.Token) {
	w := workspace{}
	w.FromStore(stored)
	wt.Workspace = w
	wt.Token = token.Key
}

type SourceType string

const (
	SourceTypeGitHub SourceType = "GitHub"
)

type ghSource struct {
	Id             int64      `json:"id"`
	Type           SourceType `json:"type"`
	Owner          string     `json:"owner"`
	Repo           string     `json:"repo"`
	Path           string     `json:"path"`
	InstallationID int64      `json:"installationID"`
}

func (g *ghSource) FromStore(stored store.GhSource) {
	g.Id = stored.ID
	g.Owner = stored.Owner
	g.Repo = stored.Repo
	g.Path = stored.Path
	g.InstallationID = stored.InstallationID
	g.Type = SourceTypeGitHub
}

func (g *ghSource) FromChangelogSource(stored store.ChangelogSource) {
	g.Id = stored.ID.Int64
	g.InstallationID = stored.InstallationID.Int64
	g.Owner = stored.Owner.String
	g.Path = stored.Path.String
	g.Repo = stored.Repo.String
	g.Type = SourceTypeGitHub
}
