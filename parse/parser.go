package parse

import (
	"context"
	"io"
	"time"

	"github.com/jonashiltl/openchangelog/internal"
)

type Meta struct {
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	PublishedAt time.Time `yaml:"publishedAt"`
}

type ParsedArticle struct {
	Meta    Meta
	Content io.Reader
}

type ParseResult struct {
	Articles []ParsedArticle
}

type Parser interface {
	Parse(ctx context.Context, raw []internal.RawArticle) (ParseResult, error)
}
