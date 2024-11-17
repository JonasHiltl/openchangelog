package events

import "github.com/jonashiltl/openchangelog/internal/source"

// Fired if sources data changed
type SourceContentChanged struct {
	WID    string
	Source source.Source
}
