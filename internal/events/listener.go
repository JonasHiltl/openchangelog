package events

import (
	"context"
	"log/slog"

	"github.com/btvoidx/mint"
	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/search"
	"github.com/jonashiltl/openchangelog/internal/source"
	"github.com/jonashiltl/openchangelog/internal/xcache"
	"github.com/jonashiltl/openchangelog/internal/xlog"
)

type offFunc func() <-chan struct{}

type EventListener struct {
	e        *mint.Emitter
	parser   parse.Parser
	searcher search.Searcher
	cache    xcache.Cache
	cfg      config.Config
	offs     []offFunc
}

func NewListener(
	cfg config.Config,
	e *mint.Emitter,
	parser parse.Parser,
	searcher search.Searcher,
	cache xcache.Cache,
) *EventListener {
	return &EventListener{
		e:        e,
		parser:   parser,
		searcher: searcher,
		cfg:      cfg,
		cache:    cache,
	}
}

// Starts listening to all events
func (l *EventListener) Start() {
	off1 := mint.On(l.e, l.OnSourceChanged)
	off2 := mint.On(l.e, l.OnChangelogUpdated)
	// save all off functions of mint to cleanup later
	l.offs = append(l.offs, off1, off2)
}

// Stops listening to all events
func (l *EventListener) Close() {
	for _, off := range l.offs {
		off()
	}
}

func (l *EventListener) OnSourceChanged(e SourceContentChanged) {
	slog.Debug("source content changed event", slog.String("sid", e.Source.ID().String()))
	if e.CL.Searchable {
		go l.reindexSource(e.Source)
	}
}

func (l *EventListener) OnChangelogUpdated(e ChangelogUpdated) {
	slog.Debug("changelog updated event", slog.String("cid", e.CL.ID.String()))
	if e.Args.Searchable != nil && *e.Args.Searchable {
		souce, err := source.NewSourceFromStore(l.cfg, e.CL, l.cache)
		if err == nil {
			go l.reindexSource(souce)
		} else {
			slog.Error("failed to create source", xlog.ErrAttr(err))
		}
	} else if e.Args.Searchable != nil && !*e.Args.Searchable {
		souce, err := source.NewSourceFromStore(l.cfg, e.CL, l.cache)
		if err == nil {
			go l.removeIndex(souce)
		} else {
			slog.Error("failed to create source", xlog.ErrAttr(err))
		}
	}
}

func (l *EventListener) reindexSource(source source.Source) {
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

func (l *EventListener) removeIndex(source source.Source) {
	if source == nil {
		return
	}

	slog.Debug("removing search index of source", slog.String("sid", source.ID().String()))
	ctx := context.Background()
	loaded, err := source.Load(ctx, internal.NoPagination())
	if err != nil {
		slog.Error("failed to load source content for search indexing", xlog.ErrAttr(err))
		return
	}
	parsed := l.parser.Parse(ctx, loaded.Raw, internal.NoPagination())
	err = l.searcher.BatchRemove(ctx, search.BatchRemoveArgs{
		SID:          source.ID().String(),
		ReleaseNotes: parsed.ReleaseNotes,
	})
	if err != nil {
		slog.Error("failed to index parsed release notes", xlog.ErrAttr(err))
		return
	}
}
