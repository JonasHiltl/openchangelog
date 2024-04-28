package parse

import (
	"bytes"
	"context"
	"time"

	"github.com/jonashiltl/openchangelog/source"
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

type ParseResult struct {
	Articles []ParsedArticle
	HasMore  bool
}

type Parser interface {
	Parse(ctx context.Context, s source.Source, params source.LoadParams) (ParseResult, error)
}
