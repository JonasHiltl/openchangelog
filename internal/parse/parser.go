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
	Articles []ParsedReleaseNote
	HasMore  bool
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

// Parses all the raw articles, uses either the keep-a-changelog parser or our og parser.
// Uses the keep-a-changelog parser if only a single article in the keep-a-changelog format is provided.
// Pagination is only applied when using the keep-a-changelog parser.
// Else parses using the original parser.
func (p *Parser) Parse(ctx context.Context, raw []source.RawReleaseNote, kPage internal.Pagination) ParseResult {

	// sanitize pagination
	if kPage.IsDefined() && kPage.PageSize() < 1 {
		return ParseResult{
			Articles: []ParsedReleaseNote{},
			HasMore:  false,
		}
	}

	if len(raw) == 1 {
		format, read := detectFileFormat(raw[0].Content)
		if format == KeepAChangelog {
			return p.k.parse(read, raw[0].Content, kPage)
		}
		parsed, err := p.og.parseArticleRead(read, raw[0].Content)
		if err != nil {
			return ParseResult{}
		}
		return ParseResult{
			Articles: []ParsedReleaseNote{parsed},
			HasMore:  false, // hasMore can only be true with the keep-a-changelog parser
		}
	}

	result := ParseResult{}
	var wg sync.WaitGroup
	mutex := &sync.Mutex{}
	for _, a := range raw {
		wg.Add(1)
		go func(a source.RawReleaseNote) {
			defer wg.Done()
			parsed, err := p.og.parseArticle(a.Content)
			if err != nil {
				return
			}
			mutex.Lock()
			result.Articles = append(result.Articles, parsed)
			mutex.Unlock()
		}(a)
	}
	wg.Wait()

	slices.SortFunc(result.Articles, sortArticleDesc)

	return result
}
