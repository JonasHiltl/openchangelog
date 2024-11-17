package search

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/lgr"
	"github.com/jonashiltl/openchangelog/internal/parse"
)

type Searcher interface {
	Search(context.Context, SearchArgs) (SearchResults, error)
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

func (s storedReleaseNote) toDoc() mapping.Classifier {
	return s
}

func (s storedReleaseNote) Type() string {
	return "note"
}

func (s *bleveSearcher) Close() {
	err := s.idx.Close()
	if err != nil {
		slog.Error("failed to close search index", lgr.ErrAttr(err))
	}
}

type SearchResult struct {
	ID          string                 `json:"id"`
	Score       float64                `json:"score"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Content     string                 `json:"content"`
	Fragments   map[string][]string    `json:"fragments,omitempty"` // Highlighted fragments
	Fields      map[string]interface{} `json:"fields,omitempty"`
}

type SearchResults struct {
	Total    uint64         `json:"total"`
	Hits     []SearchResult `json:"hits"`
	MaxScore float64        `json:"max_score"`
}

type SearchArgs struct {
	WID   string
	SID   string
	Tags  []string
	Query string
}

func (s *bleveSearcher) Search(ctx context.Context, args SearchArgs) (SearchResults, error) {
	query := buildSearchQuery(args)
	req := bleve.NewSearchRequest(query)
	req.Fields = []string{"Title", "Description", "Content"}
	req.Highlight = bleve.NewHighlightWithStyle("html")

	res, err := s.idx.SearchInContext(ctx, req)
	if err != nil {
		return SearchResults{}, err
	}

	results := SearchResults{
		Total:    res.Total,
		MaxScore: res.MaxScore,
		Hits:     make([]SearchResult, len(res.Hits)),
	}

	// Convert each hit to our SearchResult struct
	for i, hit := range res.Hits {
		result := SearchResult{
			ID:        hit.ID,
			Score:     hit.Score,
			Fragments: hit.Fragments,
			Fields:    hit.Fields,
		}

		// Extract fields if they exist
		if title, exists := hit.Fields["Title"]; exists {
			result.Title = fmt.Sprint(title)
		}
		if desc, exists := hit.Fields["Description"]; exists {
			result.Description = fmt.Sprint(desc)
		}
		if content, exists := hit.Fields["Content"]; exists {
			result.Content = fmt.Sprint(content)
		}

		results.Hits[i] = result
	}
	return results, nil
}

func buildSearchQuery(args SearchArgs) query.Query {
	wIDQuery := bleve.NewMatchQuery(args.WID)
	wIDQuery.SetField("WID")
	sIDQuery := bleve.NewMatchQuery(args.SID)
	sIDQuery.SetField("SID")

	query := bleve.NewBooleanQuery()
	query.AddMust(wIDQuery, sIDQuery)

	if len(args.Tags) > 0 {
		for _, t := range args.Tags {
			tagQuery := bleve.NewMatchQuery(t)
			tagQuery.SetField("Tags")
			query.AddMust(tagQuery)
		}
	}

	if args.Query != "" {
		titleQuery := bleve.NewMatchQuery(strings.ToLower(args.Query))
		titleQuery.SetField("Title")
		titleQuery.SetBoost(4)

		descQuery := bleve.NewMatchQuery(strings.ToLower(args.Query))
		descQuery.SetField("Description")
		descQuery.SetBoost(2)

		contentQuery := bleve.NewMatchQuery(args.Query)
		contentQuery.SetField("Content")

		combinedQuery := bleve.NewDisjunctionQuery(titleQuery, descQuery, contentQuery)
		query.AddMust(combinedQuery) // ... AND title OR desc OR content
	}
	return query
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

	doc := storedReleaseNote{
		WID:         args.WID,
		SID:         args.SID,
		Title:       args.ReleaseNote.Meta.Title,
		Description: args.ReleaseNote.Meta.Description,
		PublishedAt: args.ReleaseNote.Meta.PublishedAt,
		Tags:        args.ReleaseNote.Meta.Tags,
		Content:     string(content),
	}

	return s.idx.Index(id, doc.toDoc())
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

		doc := storedReleaseNote{
			WID:         args.WID,
			SID:         args.SID,
			Title:       note.Meta.Title,
			Description: note.Meta.Description,
			PublishedAt: note.Meta.PublishedAt,
			Tags:        note.Meta.Tags,
			Content:     string(content),
		}

		err = b.Index(id, doc.toDoc())
		if err != nil {
			slog.DebugContext(ctx, fmt.Sprintf("failed to batch index %s, skipping it", id))
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
