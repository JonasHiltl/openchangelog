package changelog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sort"
	"sync"

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

func (g *gmark) Parse(ctx context.Context, raw []RawArticle) ([]ParsedArticle, error) {
	var wg sync.WaitGroup
	result := make([]ParsedArticle, 0, len(raw))
	mutex := &sync.Mutex{}

	for _, a := range raw {
		wg.Add(1)
		go func(a RawArticle) {
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

func (g *gmark) parseArticle(raw RawArticle) (ParsedArticle, error) {
	ctx := parser.NewContext()

	defer raw.Content.Close()
	source, err := io.ReadAll(raw.Content)
	if err != nil {
		return ParsedArticle{}, err
	}

	var target bytes.Buffer
	err = g.gm.Convert(source, &target, parser.WithContext(ctx))
	if err != nil {
		return ParsedArticle{}, err
	}

	data := frontmatter.Get(ctx)
	if data == nil {
		return ParsedArticle{
			Content: &target,
		}, nil
	}
	var meta Meta
	err = data.Decode(&meta)
	if err != nil {
		return ParsedArticle{}, err
	}

	meta.ID = fmt.Sprint(meta.PublishedAt.Unix())

	return ParsedArticle{
		Meta:    meta,
		Content: &target,
	}, nil
}
