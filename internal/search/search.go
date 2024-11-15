package search

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/lgr"
	"github.com/jonashiltl/openchangelog/internal/parse"
)

type Searcher interface {
	Index(context.Context, IndexArgs) error
	BatchIndex(context.Context, BatchIndexArgs) error
	Close()
}

type bleveSearcher struct {
	idx bleve.Index
}

func NewSearcher(cfg config.Config) (Searcher, error) {
	idx, err := createBleve(cfg)
	if err != nil {
		return nil, err
	}
	return &bleveSearcher{
		idx: idx,
	}, nil
}

type storedReleaseNote struct {
	WID         string
	SID         string
	Title       string
	Description string
	PublishedAt time.Time
	Tags        []string
	Content     string
}

func (s *bleveSearcher) Close() {
	s.idx.Close()
}

type IndexArgs struct {
	WID         string // the workspace id of the release notes
	SID         string // the id of the source that was used to load the release notes
	ReleaseNote parse.ParsedReleaseNote
}

func (s *bleveSearcher) Index(ctx context.Context, args IndexArgs) error {
	content, err := io.ReadAll(args.ReleaseNote.Content)
	if err != nil {
		slog.DebugContext(ctx, "failed to read release note content for search indexing", lgr.ErrAttr(err))
		return err
	}

	id := createID(args.WID, args.SID, args.ReleaseNote.Meta.ID)
	slog.Debug("indexing document", slog.String("id", id))

	return s.idx.Index(
		id,
		storedReleaseNote{
			WID:         args.WID,
			SID:         args.SID,
			Title:       args.ReleaseNote.Meta.Title,
			Description: args.ReleaseNote.Meta.Description,
			PublishedAt: args.ReleaseNote.Meta.PublishedAt,
			Tags:        args.ReleaseNote.Meta.Tags,
			Content:     string(content),
		},
	)
}

type BatchIndexArgs struct {
	WID          string
	SID          string
	ReleaseNotes []parse.ParsedReleaseNote
}

func (s *bleveSearcher) BatchIndex(ctx context.Context, args BatchIndexArgs) error {
	b := s.idx.NewBatch()
	for _, note := range args.ReleaseNotes {
		content, err := io.ReadAll(note.Content)
		if err != nil {
			slog.DebugContext(ctx, "failed to read release note content for search indexing, skipping it", lgr.ErrAttr(err))
			return err
		}

		id := createID(args.WID, args.SID, note.Meta.ID)
		slog.Debug("indexing document", slog.String("id", id))

		err = b.Index(
			id,
			storedReleaseNote{
				WID:         args.WID,
				SID:         args.SID,
				Title:       note.Meta.Title,
				Description: note.Meta.Description,
				PublishedAt: note.Meta.PublishedAt,
				Tags:        note.Meta.Tags,
				Content:     string(content),
			},
		)
		if err != nil {
			slog.DebugContext(ctx, fmt.Sprintf("failed to index %s, skipping it", id))
		}
	}

	if b.Size() > 0 {
		if err := s.idx.Batch(b); err != nil {
			return err
		}
	}
	return nil
}

func createID(wID, sID, releaseNoteID string) string {
	return fmt.Sprintf("%s/%s/%s", wID, sID, releaseNoteID)
}
