package events

import (
	"github.com/jonashiltl/openchangelog/internal/source"
	"github.com/jonashiltl/openchangelog/internal/store"
)

// Fired if sources data changed
type SourceContentChanged struct {
	CL     store.Changelog
	Source source.Source
}
