package changelog

import (
	"bufio"
	"bytes"
	"context"
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
func (g *kparser) Parse(ctx context.Context, raw RawArticle, page Pagination) ([]ParsedArticle, error) {
	defer raw.Content.Close()

	// sanitize pagination
	if page.IsDefined() && page.PageSize() < 1 {
		return make([]ParsedArticle, 0), nil
	}

	parsed, err := g.parseChangelog(raw, page)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

// 1. Split by ##
//  1. section is ignored
//  1. section is an optional unreleased section
//     ... each section is it's own article
func (g *kparser) parseChangelog(raw RawArticle, page Pagination) ([]ParsedArticle, error) {
	sc := bufio.NewScanner(raw.Content)

	var articles []ParsedArticle
	// required since g.parseArticle() returns the first line of the next article
	var line string
	var currentArticleIdx int = 0

	// scans line per line
	for sc.Scan() || line != "" {
		if line == "" {
			line = sc.Text()
		}

		// start of new release
		if strings.HasPrefix(line, "## ") {
			if !page.IsDefined() || (currentArticleIdx >= page.StartIdx() && currentArticleIdx <= page.EndIdx()) {
				a, nextLine, err := g.parseArticle(line, sc)
				if err == nil {
					articles = append(articles, a)
				}
				line = nextLine
			} else {
				line = ""
			}

			currentArticleIdx++
		} else {
			line = ""
		}
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

// Called on each new ## section of the changelog file.
// Returns the currently parsed article and the new line of the next article if another article exists.
func (g *kparser) parseArticle(firstLine string, sc *bufio.Scanner) (ParsedArticle, string, error) {
	h2Parts := strings.SplitN(strings.TrimPrefix(firstLine, "## "), " - ", 2)
	title := cleanTitle(h2Parts[0])

	var content bytes.Buffer
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

	var line string
	for sc.Scan() {
		line = sc.Text()

		if strings.HasPrefix(line, "## ") {
			// begin of new article
			break
		} else if strings.HasPrefix(line, "### ") {
			// type of change
			parts := strings.Split(line, " ")
			if len(parts) > 1 {
				a.AddTag(parts[1])
			}

			content.WriteString(line + "\n")
		} else {
			// content line
			content.WriteString(line + "\n")
		}
	}

	if err := sc.Err(); err != nil {
		return ParsedArticle{}, line, err
	}

	var htmlContent bytes.Buffer
	err := g.gm.Convert(content.Bytes(), &htmlContent)
	if err != nil {
		return ParsedArticle{}, line, err
	}

	a.Content = &htmlContent

	return a, line, nil
}

func cleanTitle(title string) string {
	title = strings.TrimSpace(title)
	title = strings.Replace(title, "[", "", 1)
	return strings.Replace(title, "]", "", 1)
}
