package changelog

import (
	"context"
	"io"
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

type Parser interface {
	Parse(ctx context.Context, raw []RawArticle) ([]ParsedArticle, error)
}
