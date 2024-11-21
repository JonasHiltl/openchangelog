package source

import (
	"context"
	"io"

	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/store"
)

type RawReleaseNote struct {
	Content    io.Reader
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
	ID() ID
}

// A unique identifier for a source
type ID string

func (i ID) String() string {
	return string(i)
}

func NewIDFromChangelog(cl store.Changelog) ID {
	if cl.LocalSource.Valid {
		return NewLocalID(cl.LocalSource.V.Path)
	} else if cl.GHSource.Valid {
		return NewGitHubID(cl.GHSource.V.Owner, cl.GHSource.V.Repo, cl.GHSource.V.Path)
	}
	return ""
}
