package internal

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

// A source is used to download the changelog articles from a targe
type Source interface {
	Load(ctx context.Context, page Pagination) (LoadResult, error)
}
