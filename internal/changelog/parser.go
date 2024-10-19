package changelog

import (
	"bytes"
	"context"
	"io"
	"slices"
	"sync"
	"time"

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

type ParsedArticle struct {
	Meta    Meta
	Content io.Reader
}

func (a *ParsedArticle) AddTag(t string) {
	if !slices.Contains(a.Meta.Tags, t) {
		a.Meta.Tags = append(a.Meta.Tags, t)
	}
}

type ParseResult struct {
	Articles []ParsedArticle
	HasMore  bool
}

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
// Pagination is only applied when using the keep-a-changelog parser
// Else parses using the original parser.
func (p *Parser) Parse(ctx context.Context, raw []RawArticle, kPage Pagination) ParseResult {

	// sanitize pagination
	if kPage.IsDefined() && kPage.PageSize() < 1 {
		return ParseResult{
			Articles: []ParsedArticle{},
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
			Articles: []ParsedArticle{parsed},
			HasMore:  false, // hasMore can only be true with the keep-a-changelog parser
		}
	}

	result := ParseResult{}
	var wg sync.WaitGroup
	mutex := &sync.Mutex{}
	for _, a := range raw {
		wg.Add(1)
		go func(a RawArticle) {
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

type FileFormat int

const (
	OG FileFormat = iota
	KeepAChangelog
)

// Detects the file format of r and returns the string read to detect the file format.
// The read string can not be read again from r.
func detectFileFormat(r io.Reader) (FileFormat, string) {
	var buf bytes.Buffer
	_, err := io.CopyN(&buf, r, 3)
	if err != nil {
		return OG, ""
	}
	start := buf.String()
	if start == "---" {
		// if content has frontmatter => it's probably our own file format
		return OG, start
	}
	return KeepAChangelog, start
}

// Sorts ParsedArticles by their published date.
func sortArticleDesc(a ParsedArticle, b ParsedArticle) int {
	if a.Meta.PublishedAt.IsZero() && b.Meta.PublishedAt.IsZero() {
		return 0
	}
	if a.Meta.PublishedAt.IsZero() {
		return -1
	}
	if b.Meta.PublishedAt.IsZero() {
		return 1
	}

	if a.Meta.PublishedAt.After(b.Meta.PublishedAt) {
		return -1
	}

	return 1
}
