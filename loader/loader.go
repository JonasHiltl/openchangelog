package loader

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

// A loader downloads the changelog markdown files.
type Loader interface {
	Load(ctx context.Context, page Pagination) (LoadResult, error)
}
