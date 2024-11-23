package source

import (
	"context"
	"errors"
	"io"

	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/internal/xcache"
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

func NewSourceFromStore(cfg config.Config, cl store.Changelog, cache xcache.Cache) (Source, error) {
	if cl.LocalSource.Valid {
		return NewLocalSourceFromStore(cl.LocalSource.ValueOrZero(), cache), nil
	} else if cl.GHSource.Valid {
		return NewGHSourceFromStore(cfg, cl.GHSource.ValueOrZero(), cache)
	}
	return nil, errors.New("changelog has no active source")
}
