package events

import (
	"context"
	"log/slog"

	"github.com/btvoidx/mint"
	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/lgr"
	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/search"
)

type offFunc func() <-chan struct{}

type EventListener struct {
	e        *mint.Emitter
	parser   parse.Parser
	searcher search.Searcher
	offs     []offFunc
}

func NewListener(e *mint.Emitter, parser parse.Parser, searcher search.Searcher) EventListener {
	return EventListener{
		e:        e,
		parser:   parser,
		searcher: searcher,
	}
}

// Starts listening to all events
func (l *EventListener) Start() {
	off := mint.On(l.e, l.OnSourceChanged)
	// save all off functions of mint to cleanup later
	l.offs = append(l.offs, off)
}

// Stops listening to all events
func (l EventListener) Close() {
	for _, off := range l.offs {
		off()
	}
}

func (l EventListener) OnSourceChanged(e SourceChanged) {
	go func() {
		ctx := context.Background()
		loaded, err := e.Source.Load(ctx, internal.NoPagination())
		if err != nil {
			slog.Error("failed to load source content for search indexing", lgr.ErrAttr(err))
			return
		}
		parsed := l.parser.Parse(ctx, loaded.Raw, internal.NoPagination())
		err = l.searcher.BatchIndex(ctx, search.BatchIndexArgs{
			WID:          e.WID,
			SID:          e.Source.ID(),
			ReleaseNotes: parsed.ReleaseNotes,
		})
		if err != nil {
			slog.Error("failed to index parsed release notes", lgr.ErrAttr(err))
			return
		}
	}()
}
