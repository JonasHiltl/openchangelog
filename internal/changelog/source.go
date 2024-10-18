package changelog

import (
	"context"
	"io"
)

type RawArticle struct {
	Content io.ReadCloser
}

type LoadResult struct {
	Articles []RawArticle
	HasMore  bool
}

// A source is used to download the changelog articles from a target
type Source interface {
	Load(ctx context.Context, page Pagination) (LoadResult, error)
}
