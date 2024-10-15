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

// A parser that parses changelogs using the https://keepachangelog.com/en/1.1.0/ format.
type keepachangelog struct {
	gm goldmark.Markdown
}

func NewKeepAChangelogParser() Parser {
	return &keepachangelog{
		gm: createGoldmark(),
	}
}

func (g *keepachangelog) Parse(ctx context.Context, raw []RawArticle) ([]ParsedArticle, error) {
	// optimize for this, since it's the most common case. Just asingle CHANGELOG.md file.
	if len(raw) == 1 {
		parsed, err := g.parseChangelog(raw[0])
		if err != nil {
			return nil, err
		}
		return parsed, nil
	}

	return nil, errors.New("keep a changelog format expects one raw article")
}

// 1. Split by ##
//  1. section is ignored
//  1. section is an optional unreleased section
//     ... each section is it's own article
func (g *keepachangelog) parseChangelog(raw RawArticle) ([]ParsedArticle, error) {
	sc := bufio.NewScanner(raw.Content)

	var articles []ParsedArticle
	// required since g.parseArticle() returns the first line of the next article
	var line string

	// scans line per line
	for sc.Scan() || line != "" {
		if line == "" {
			line = sc.Text()
		}

		if strings.HasPrefix(line, "## ") {
			a, nextLine, err := g.parseArticle(sc)
			if err == nil {
				articles = append(articles, a)
			}
			line = nextLine
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
func (g *keepachangelog) parseArticle(sc *bufio.Scanner) (ParsedArticle, string, error) {
	h2Parts := strings.SplitN(strings.TrimPrefix(sc.Text(), "## "), " - ", 2)

	b := new(bytes.Buffer)
	a := ParsedArticle{
		Meta: Meta{
			Title: cleanTitle(h2Parts[0]),
		},
		Content: b,
	}

	if len(h2Parts) > 1 {
		publishedAt, err := time.Parse("2006-01-02", strings.TrimSpace(h2Parts[1]))
		if err == nil {
			a.Meta.PublishedAt = publishedAt
		}
	}

	for sc.Scan() {
		line := sc.Text()

		if strings.HasPrefix(line, "## ") {
			// begin of new article
			return a, line, nil
		} else if strings.HasPrefix(line, "### ") {
			// type of change
			parts := strings.Split(line, " ")
			if len(parts) > 1 {
				a.AddTag(parts[1])
			}

			b.WriteString(line + "\n")
		} else {
			// content line
			b.WriteString(line + "\n")
		}
	}

	if err := sc.Err(); err != nil {
		return ParsedArticle{}, "", err
	}

	return a, "", nil
}

func cleanTitle(title string) string {
	title = strings.TrimSpace(title)
	title = strings.Replace(title, "[", "", 1)
	return strings.Replace(title, "]", "", 1)
}
