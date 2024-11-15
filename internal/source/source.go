package source

import (
	"context"
	"io"

	"github.com/jonashiltl/openchangelog/internal"
)

type RawReleaseNote struct {
	Content    io.ReadCloser
	hasChanged bool // only available if caching is enabled
}

type LoadResult struct {
	Raw     []RawReleaseNote
	HasMore bool
}

// Returns if any of the loaded release notes have changed since last access.
func (r LoadResult) HasChanged() bool {
	for _, note := range r.Raw {
		if note.hasChanged {
			return true
		}
	}
	return false
}

// A source can be used to load raw release notes from a (remote) source like GitHub.
type Source interface {
	Load(ctx context.Context, page internal.Pagination) (LoadResult, error)
	ID() string
}
