package search

import (
	"context"
	"fmt"
	"html"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
	mapset "github.com/deckarep/golang-set/v2"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/lgr"
	"github.com/jonashiltl/openchangelog/internal/parse"
)

type Searcher interface {
	Search(context.Context, SearchArgs) (SearchResults, error)
	GetAllTags(ctx context.Context, sid string) []string
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
	ID               string
	Title            string
	Description      string
	ContentHighlight string
	Score            float64
	Fragments        map[string][]string
	Fields           map[string]interface{}
}

type SearchResults struct {
	Total    uint64
	Hits     []SearchResult
	MaxScore float64
}

type SearchArgs struct {
	SID   string
	Tags  []string
	Query string
}

func (s *bleveSearcher) Search(ctx context.Context, args SearchArgs) (SearchResults, error) {
	query := buildSearchQuery(ctx, args)
	req := bleve.NewSearchRequestOptions(query, 10, 0, false)
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
			result.Title = title.(string)
		}
		if desc, exists := hit.Fields["Description"]; exists {
			result.Description = desc.(string)
		}

		if title, exists := hit.Fragments["Title"]; exists && len(title) > 0 {
			result.Title = title[0]
		}
		if description, exists := hit.Fragments["Description"]; exists && len(description) > 0 {
			result.Description = description[0]
		}

		if content, exists := hit.Fragments["Content"]; exists && len(content) > 0 {
			result.ContentHighlight = surroundWithEllipsis(stripPartialHTML(content[0]))
		} else if content, exists := hit.Fields["Content"]; exists {
			// use start of content if no content highlights exist
			nwords := firstNWords(fmt.Sprint(content), 10)
			content := stripPartialHTML(nwords)
			result.ContentHighlight = fmt.Sprintf("%s...", content)
		}

		results.Hits[i] = result
	}
	return results, nil
}

// stripPartialHTML removes all HTML tags and partial tags from the input string,
// except for <mark> tags, which are preserved. It also handles cases where
// the input starts or ends with incomplete HTML tags to ensure clean output.
func stripPartialHTML(input string) string {
	input = html.UnescapeString(input)

	contentStartIdx := 0
	for i, r := range input {
		// if partial tag exists at start of input, content start after it
		if r == '>' {
			contentStartIdx = i + 1
			break
		}
		// text start with a partial opening tag, will be removed later
		if r == '<' {
			break
		}
	}

	if contentStartIdx >= len(input) {
		return ""
	}

	input = input[contentStartIdx:]

	// replace <mark> tags so that they can later be recovered
	input = replaceMarks(input)
	// clean all html tags
	cleanedContent := strip.StripTags(input)
	// recover <mark> tags
	cleanedContent = replaceMarkPlaceholders(cleanedContent)

	cleanedRunes := []rune(cleanedContent)

	// Find where valid content ends
	// We still might have some partial tags at the end of content
	contentEnd := len(cleanedRunes)
	for i := len(cleanedRunes) - 1; i > 0; i-- {
		r := cleanedRunes[i]
		if r == '<' {
			contentEnd = i
		}
		// if </mark> is at end we can stop
		if r == '>' {
			break
		}
	}

	return string(cleanedRunes[:contentEnd])
}

// Gets the first n words of the input.
// If input has less than n words just returns input.
func firstNWords(input string, n int) string {
	words := strings.Fields(input) // Split the string into words

	if n > len(words) {
		return input
	}

	return strings.Join(words[:n], " ")
}

const mark_placeholder_start = "__MARK_START__"
const mark_placeholder_end = "__MARK_END__"

func replaceMarks(input string) string {
	input = strings.ReplaceAll(input, "<mark>", mark_placeholder_start)
	input = strings.ReplaceAll(input, "</mark>", mark_placeholder_end)
	return input
}

func replaceMarkPlaceholders(input string) string {
	input = strings.ReplaceAll(input, mark_placeholder_start, "<mark>")
	input = strings.ReplaceAll(input, mark_placeholder_end, "</mark>")
	return input
}

func surroundWithEllipsis(input string) string {
	input, _ = strings.CutPrefix(input, "...")
	input, _ = strings.CutPrefix(input, "…")
	input, _ = strings.CutSuffix(input, "...")
	input, _ = strings.CutSuffix(input, "…")
	return fmt.Sprintf("...%s...", input)
}

func (s *bleveSearcher) GetAllTags(ctx context.Context, sid string) []string {
	query := bleve.NewMatchQuery(sid)
	query.SetField("SID")
	req := bleve.NewSearchRequest(query)
	req.Fields = []string{"Tags"}

	res, err := s.idx.SearchInContext(ctx, req)
	if err != nil {
		return []string{}
	}

	set := mapset.NewThreadUnsafeSet("")
	for _, hit := range res.Hits {
		if tags, exists := hit.Fields["Tags"]; exists {
			if tagsSlice, ok := tags.([]any); ok {
				for _, t := range tagsSlice {
					set.Add(fmt.Sprint(t))
				}
			} else if tag, ok := tags.(string); ok {
				set.Add(tag)
			}
		}
	}

	return set.ToSlice()
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
