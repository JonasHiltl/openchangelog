package source

import (
	"context"
	"io"

	"github.com/jonashiltl/openchangelog/internal"
)

type RawReleaseNote struct {
	HasChanged bool // only available if caching is enabled
	Content    io.ReadCloser
}

type LoadResult struct {
	Raw     []RawReleaseNote
	HasMore bool
}

// A source can be used to load raw release notes from a (remote) source like GitHub.
type Source interface {
	Load(ctx context.Context, page internal.Pagination) (LoadResult, error)
}
