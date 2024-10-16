package changelog

import (
	"context"
	"io"
	"slices"
	"time"
)

type Meta struct {
	// unique id, we use the published date as unix timestampe
	ID          string
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	PublishedAt time.Time `yaml:"publishedAt"`
	Tags        []string  `yaml:"tags"`
}

type ParsedArticle struct {
	Meta    Meta
	Content io.Reader
}

func (a *ParsedArticle) AddTag(t string) {
	if !slices.Contains(a.Meta.Tags, t) {
		a.Meta.Tags = append(a.Meta.Tags, t)
	}
}

type Parser interface {
	Parse(ctx context.Context, raw []RawArticle) ([]ParsedArticle, error)
}
