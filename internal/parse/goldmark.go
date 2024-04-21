package parse

import (
	"bytes"
	"sort"
	"sync"

	"github.com/jonashiltl/openchangelog/internal/source"
	enclave "github.com/quail-ink/goldmark-enclave"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
	"mvdan.cc/xurls/v2"
)

type gmark struct {
	gm goldmark.Markdown
}

func NewParser() Parser {
	return &gmark{
		gm: goldmark.New(
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
			goldmark.WithExtensions(
				extension.GFM,
				enclave.New(&enclave.Config{}),
				&frontmatter.Extender{},
				extension.NewLinkify(
					extension.WithLinkifyAllowedProtocols([]string{
						"http:",
						"https:",
					}),
					extension.WithLinkifyURLRegexp(
						xurls.Strict(),
					),
				),
			),
		),
	}
}

func (g *gmark) Parse(s source.Source) ([]ParsedArticle, error) {
	articles, err := s.Load()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	result := make([]ParsedArticle, 0, len(articles))
	mutex := &sync.Mutex{}

	for _, a := range articles {
		wg.Add(1)
		go func(a source.Article) {
			defer wg.Done()
			parsed, err := g.parseArticle(a)
			if err != nil {
				return
			}
			mutex.Lock()
			result = append(result, parsed)
			mutex.Unlock()
		}(a)
	}
	wg.Wait()

	sort.Slice(result, func(i, j int) bool {
		return result[i].Meta.PublishedAt.After(result[j].Meta.PublishedAt)
	})

	return result, nil
}

func (g *gmark) parseArticle(a source.Article) (ParsedArticle, error) {
	ctx := parser.NewContext()
	var buf bytes.Buffer
	err := g.gm.Convert(a.Bytes, &buf, parser.WithContext(ctx))
	if err != nil {
		return ParsedArticle{}, err
	}

	data := frontmatter.Get(ctx)
	var meta Meta
	err = data.Decode(&meta)
	if err != nil {
		return ParsedArticle{}, err
	}
	return ParsedArticle{
		Meta:    meta,
		Content: &buf,
	}, nil
}
