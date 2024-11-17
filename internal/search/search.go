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
	parse.ParsedReleaseNote
	ID        string
	Score     float64
	Fragments map[string][]string
	Fields    map[string]interface{}
}

type SearchResults struct {
	Total    uint64
	Hits     []SearchResult
	MaxScore float64
}

func (r SearchResults) GetParsedReleaseNotes() []parse.ParsedReleaseNote {
	notes := make([]parse.ParsedReleaseNote, len(r.Hits))
	for i, h := range r.Hits {
		notes[i] = h.ParsedReleaseNote
	}
	return notes
}

type SearchArgs struct {
	SID   string
	Tags  []string
	Query string
}

func (s *bleveSearcher) Search(ctx context.Context, args SearchArgs) (SearchResults, error) {
	query := buildSearchQuery(ctx, args)
	req := bleve.NewSearchRequest(query)
	req.SortBy([]string{"-PublishedAt"})
	req.Fields = []string{"Title", "Description", "Content", "Tags", "PublishedAt"}
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
			result.Meta.Title = fmt.Sprint(title)
		}
		if desc, exists := hit.Fields["Description"]; exists {
			result.Meta.Description = fmt.Sprint(desc)
		}
		if content, exists := hit.Fields["Content"]; exists {
			result.Content = strings.NewReader(fmt.Sprint(content))
		}
		if tags, exists := hit.Fields["Tags"]; exists {
			if tagsSlice, ok := tags.([]any); ok {
				for _, t := range tagsSlice {
					result.Meta.Tags = append(result.Meta.Tags, fmt.Sprint(t))
				}
			} else if tag, ok := tags.(string); ok {
				result.Meta.Tags = []string{tag}
			}
		}
		if publishedAt, exists := hit.Fields["PublishedAt"]; exists {
			t, err := time.Parse(time.RFC3339, fmt.Sprint(publishedAt))
			if err == nil {
				result.Meta.PublishedAt = t
			}
		}

		results.Hits[i] = result
	}
	return results, nil
}

func buildSearchQuery(ctx context.Context, args SearchArgs) query.Query {
	slog.DebugContext(
		ctx,
		"building search query",
		slog.String("sid", args.SID),
		slog.String("query", args.Query),
	)
	sIDQuery := bleve.NewMatchQuery(args.SID)
	sIDQuery.SetField("SID")

	query := bleve.NewBooleanQuery()
	query.AddMust(sIDQuery)

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
	SID         string // the id of the source that was used to load the release notes
	ReleaseNote parse.ParsedReleaseNote
}

func (s *bleveSearcher) Index(ctx context.Context, args IndexArgs) error {
	content, err := io.ReadAll(args.ReleaseNote.Content)
	if err != nil {
		slog.DebugContext(ctx, "failed to read release note content for search indexing", lgr.ErrAttr(err))
		return err
	}

	id := createID(args.SID, args.ReleaseNote.Meta.ID)
	slog.Debug("indexing document", slog.String("id", id))

	doc := storedReleaseNote{
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

		id := createID(args.SID, note.Meta.ID)
		slog.Debug("indexing document", slog.String("id", id))

		doc := storedReleaseNote{
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

func createID(sID, releaseNoteID string) string {
	return fmt.Sprintf("%s/%s", sID, releaseNoteID)
}
