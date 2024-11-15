package events

import "github.com/jonashiltl/openchangelog/internal/source"

// Fired if sources data changed
type SourceChanged struct {
	WID    string
	Source source.Source
}
