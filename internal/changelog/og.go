package changelog

import (
	"bytes"
	"fmt"
	"io"

	enclave "github.com/quail-ink/goldmark-enclave"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
	"mvdan.cc/xurls/v2"
)

// This is the original parser that expects one markdown file per release note.
// All the meta information should be defined with Frontmatter.
type ogparser struct {
	gm goldmark.Markdown
}

// Creates a new goldmark instance, used to parse Markdown to HTML.
func createGoldmark() goldmark.Markdown {
	return goldmark.New(
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
	)
}

func newOGParser() *ogparser {
	return &ogparser{
		gm: createGoldmark(),
	}
}

// Takes a raw article in our original markdown format and parses it.
func (g *ogparser) parseArticle(article io.ReadCloser) (ParsedArticle, error) {
	defer article.Close()
	source, err := io.ReadAll(article)
	if err != nil {
		return ParsedArticle{}, err
	}

	return g.parseArticleBytes(source)
}

// Parses the raw article content, but expects a part of the content to be already read (to detect the file format).
func (g *ogparser) parseArticleRead(read string, rest io.ReadCloser) (ParsedArticle, error) {
	defer rest.Close()
	source, err := io.ReadAll(rest)
	if err != nil {
		return ParsedArticle{}, err
	}

	full := append([]byte(read), source...)

	return g.parseArticleBytes(full)
}

// Don't use diretly, use parseArticle() and parseArticleRead() instead.
func (g *ogparser) parseArticleBytes(content []byte) (ParsedArticle, error) {
	ctx := parser.NewContext()

	var target bytes.Buffer
	err := g.gm.Convert(content, &target, parser.WithContext(ctx))
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
