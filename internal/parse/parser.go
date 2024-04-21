package parse

import (
	"bytes"
	"context"
	"time"

	"github.com/jonashiltl/openchangelog/internal/source"
)

type Meta struct {
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	PublishedAt time.Time `yaml:"publishedAt"`
}

type ParsedArticle struct {
	Meta    Meta
	Content *bytes.Buffer
}

type Parser interface {
	Parse(ctx context.Context, s source.Source) ([]ParsedArticle, error)
}
