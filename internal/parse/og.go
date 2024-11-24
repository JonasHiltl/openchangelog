package parse

import (
	"bytes"
	"fmt"
	"io"

	enclave "github.com/quail-ink/goldmark-enclave"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	gmparser "github.com/yuin/goldmark/parser"
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
func CreateGoldmark() goldmark.Markdown {
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

func NewOGParser(gm goldmark.Markdown) *ogparser {
	return &ogparser{
		gm: gm,
	}
}

// Takes a raw article in our original markdown format and parses it.
func (g *ogparser) parseReleaseNote(article io.Reader) (ParsedReleaseNote, error) {
	source, err := io.ReadAll(article)
	if err != nil {
		return ParsedReleaseNote{}, err
	}

	return g.parseReleaseNoteBytes(source)
}

// Parses the raw article content, but expects a part of the content to be already read (to detect the file format).
func (g *ogparser) parseReleaseNoteRead(read string, rest io.Reader) (ParsedReleaseNote, error) {
	source, err := io.ReadAll(rest)
	if err != nil {
		return ParsedReleaseNote{}, err
	}

	full := append([]byte(read), source...)

	return g.parseReleaseNoteBytes(full)
}

// Don't use diretly, use parseArticle() and parseArticleRead() instead.
func (g *ogparser) parseReleaseNoteBytes(content []byte) (ParsedReleaseNote, error) {
	ctx := gmparser.NewContext()

	var target bytes.Buffer
	err := g.gm.Convert(content, &target, gmparser.WithContext(ctx))
	if err != nil {
		return ParsedReleaseNote{}, err
	}

	data := frontmatter.Get(ctx)
	if data == nil {
		return ParsedReleaseNote{
			Content: &target,
		}, nil
	}
	var meta Meta
	err = data.Decode(&meta)
	if err != nil {
		return ParsedReleaseNote{}, err
	}

	meta.ID = fmt.Sprint(meta.PublishedAt.Unix())

	return ParsedReleaseNote{
		Meta:    meta,
		Content: &target,
	}, nil
}
