package changelog

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

type kparser struct {
	gm goldmark.Markdown
}

func NewKeepAChangelogParser() *kparser {
	return &kparser{
		gm: createGoldmark(),
	}
}

// Parses a a markdown file in the https://keepachangelog.com/en/1.1.0/ format to multiple articles to be displayed by Openchangelog.
func (g *kparser) Parse(ctx context.Context, raw RawArticle, page Pagination) ([]ParsedArticle, bool) {
	defer raw.Content.Close()

	// sanitize pagination
	if page.IsDefined() && page.PageSize() < 1 {
		return make([]ParsedArticle, 0), false
	}

	return g.parseChangelog(raw, page)
}

// Returns the parsed articles and true if there are more articles to parse, else false.
func (g *kparser) parseChangelog(raw RawArticle, page Pagination) ([]ParsedArticle, bool) {
	sc := bufio.NewScanner(raw.Content)
	sc.Split(splitNewRelease)

	var articles []ParsedArticle
	var currentIdx = 0
	var hasMore = false

	for sc.Scan() {
		section := sc.Text()

		if strings.HasPrefix(section, "# ") {
			// ignore the section with the changelog title
			continue
		}

		if !page.IsDefined() || (currentIdx >= page.StartIdx() && currentIdx <= page.EndIdx()) {
			a, err := g.parseArticle(section)
			if err == nil {
				articles = append(articles, a)
			}
		}

		// check if we have one more article
		if page.IsDefined() && currentIdx == page.EndIdx()+1 {
			hasMore = true
			break
		}

		currentIdx++
	}

	return articles, hasMore
}

// Called on each new ## section of the changelog file.
// Returns the currently parsed article and the new line of the next article if another article exists.
func (g *kparser) parseArticle(article string) (ParsedArticle, error) {
	firstLineIdx := strings.Index(article, "\n")
	if firstLineIdx == -1 {
		return ParsedArticle{}, errors.New("no new line character found")
	}
	firstLine := article[:firstLineIdx]
	content := article[firstLineIdx+1:]

	h2Parts := strings.SplitN(strings.TrimPrefix(firstLine, "## "), " - ", 2)
	title := cleanTitle(h2Parts[0])

	a := ParsedArticle{
		Meta: Meta{
			Title: title,
			ID:    strings.ToLower(title),
		},
	}

	if len(h2Parts) > 1 {
		publishedAt, err := time.Parse("2006-01-02", strings.TrimSpace(h2Parts[1]))
		if err == nil {
			a.Meta.PublishedAt = publishedAt
		}
	}

	sc := bufio.NewScanner(strings.NewReader(content))
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "### ") {
			changeType := strings.TrimPrefix(line, "### ")
			a.AddTag(changeType)
		}
	}

	var htmlContent bytes.Buffer
	err := g.gm.Convert([]byte(content), &htmlContent)
	if err != nil {
		return ParsedArticle{}, err
	}

	a.Content = &htmlContent

	return a, nil
}

func cleanTitle(title string) string {
	title = strings.TrimSpace(title)
	title = strings.Replace(title, "[", "", 1)
	return strings.Replace(title, "]", "", 1)
}

func splitNewRelease(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// try to find index of the new release section
	if i := bytes.Index(data, []byte("\n## ")); i >= 0 {
		return i + 1, data[:i], nil
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}

	// Request more data.
	return 0, nil, nil
}
