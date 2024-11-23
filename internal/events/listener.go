package events

import (
	"context"
	"log/slog"

	"github.com/btvoidx/mint"
	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/search"
	"github.com/jonashiltl/openchangelog/internal/source"
	"github.com/jonashiltl/openchangelog/internal/xlog"
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

func (l EventListener) OnSourceChanged(e SourceContentChanged) {
	slog.Debug("source content changed", slog.String("sid", e.Source.ID().String()))
	if e.CL.Searchable {
		go l.reindexSource(e.Source)
	}
}

func (l EventListener) reindexSource(source source.Source) {
	if source == nil {
		return
	}

	slog.Debug("reindexing content of source", slog.String("sid", source.ID().String()))
	ctx := context.Background()
	loaded, err := source.Load(ctx, internal.NoPagination())
	if err != nil {
		slog.Error("failed to load source content for search indexing", xlog.ErrAttr(err))
		return
	}
	parsed := l.parser.Parse(ctx, loaded.Raw, internal.NoPagination())
	err = l.searcher.BatchIndex(ctx, search.BatchIndexArgs{
		SID:          source.ID().String(),
		ReleaseNotes: parsed.ReleaseNotes,
	})
	if err != nil {
		slog.Error("failed to index parsed release notes", xlog.ErrAttr(err))
		return
	}
}
