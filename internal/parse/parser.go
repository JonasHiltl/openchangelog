package parse

import (
	"context"
	"io"
	"slices"
	"sync"
	"time"

	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/source"
	"github.com/yuin/goldmark"
)

type Meta struct {
	// unique id, we use the published date as unix timestampe
	ID          string
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	PublishedAt time.Time `yaml:"publishedAt"`
	Tags        []string  `yaml:"tags"`
}

type ParsedReleaseNote struct {
	Meta    Meta
	Content io.Reader
}

func (a *ParsedReleaseNote) AddTag(t string) {
	if !slices.Contains(a.Meta.Tags, t) {
		a.Meta.Tags = append(a.Meta.Tags, t)
	}
}

type ParseResult struct {
	ReleaseNotes []ParsedReleaseNote
	HasMore      bool
}

// Creates a new Parser that can be used to parse RawReleaseNote files to ParsedReleaseNote.
// Converts all Markdown content to HTML.
func NewParser(gm goldmark.Markdown) Parser {
	return Parser{
		og: NewOGParser(gm),
		k:  NewKeepAChangelogParser(gm),
	}
}

type Parser struct {
	og *ogparser
	k  *kparser
}

// Parses a raw release note, using either the keep-a-changelog parser or our own format og parser.
// Uses our own og parser if the raw release notes starts with "---". Else uses keep-a-changelog parser.
// Pagination is only applied when using keep-a-changelog parser.
func (p *Parser) ParseRawRelease(ctx context.Context, raw source.RawReleaseNote, kPage internal.Pagination) ParseResult {
	// sanitize pagination
	if kPage.IsDefined() && kPage.PageSize() < 1 {
		return ParseResult{
			ReleaseNotes: []ParsedReleaseNote{},
			HasMore:      false,
		}
	}

	return p.parseOne(raw, kPage)
}

// Parses all the raw articles, uses either the keep-a-changelog parser or our og parser.
// Uses the keep-a-changelog parser if only a single article in the keep-a-changelog format is provided.
// Pagination is only applied when using the keep-a-changelog parser.
// Else parses using the original parser.
func (p *Parser) Parse(ctx context.Context, raw []source.RawReleaseNote, kPage internal.Pagination) ParseResult {
	// sanitize pagination
	if kPage.IsDefined() && kPage.PageSize() < 1 {
		return ParseResult{
			ReleaseNotes: []ParsedReleaseNote{},
			HasMore:      false,
		}
	}

	if len(raw) == 1 {
		return p.parseOne(raw[0], kPage)
	}

	result := make([]ParsedReleaseNote, len(raw))
	var wg sync.WaitGroup
	for i, a := range raw {
		wg.Add(1)
		go func(index int, a source.RawReleaseNote) {
			defer wg.Done()
			parsed, err := p.og.parseReleaseNote(a.Content)
			if err != nil {
				return
			}
			// Store at the correct index to maintain order
			result[index] = parsed
		}(i, a)
	}
	wg.Wait()

	// Create ParseResult with the correctly ordered results
	parseResult := ParseResult{
		ReleaseNotes: result,
		HasMore:      false, // hasMore is only true for keep-a-changelog parser
	}

	// Sort by descending order
	slices.SortFunc(parseResult.ReleaseNotes, sortArticleDesc)

	return parseResult
}

func (p *Parser) parseOne(raw source.RawReleaseNote, kPage internal.Pagination) ParseResult {
	format, read := detectFileFormat(raw.Content)
	if format == KeepAChangelog {
		// use keep a changelog parser
		return p.k.parse(read, raw.Content, kPage)
	}

	// use og parser
	parsed, err := p.og.parseReleaseNoteRead(read, raw.Content)
	if err != nil {
		return ParseResult{}
	}
	return ParseResult{
		ReleaseNotes: []ParsedReleaseNote{parsed},
		HasMore:      false, // hasMore can only be true with the keep-a-changelog parser
	}
}
