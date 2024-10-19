package changelog

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	p := NewParser(createGoldmark())
	tables := []struct {
		name                  string
		files                 []string
		page                  Pagination
		expectedHasMore       bool
		expectedArticleLength int
	}{
		{
			name:                  "OG one",
			files:                 []string{"v0.0.1-commonmark.md"},
			page:                  NoPagination(),
			expectedHasMore:       false,
			expectedArticleLength: 1,
		},
		{
			name:                  "OG multiple",
			files:                 []string{"v0.0.1-commonmark.md", "v0.0.2-open-source.md", "v0.0.5-beta.md"},
			page:                  NoPagination(),
			expectedHasMore:       false,
			expectedArticleLength: 3,
		},
		{
			name:                  "Keepachangelog no pagination",
			files:                 []string{"keepachangelog/minimal.md"},
			page:                  NoPagination(),
			expectedHasMore:       false,
			expectedArticleLength: 1,
		},
		{
			name:                  "Keepachangelog full no pagination",
			files:                 []string{"keepachangelog/full.md"},
			page:                  NoPagination(),
			expectedHasMore:       false,
			expectedArticleLength: 15,
		},
		{
			name:                  "Keepachangelog full paginated",
			files:                 []string{"keepachangelog/full.md"},
			page:                  NewPagination(4, 2),
			expectedHasMore:       true,
			expectedArticleLength: 4,
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			articles := make([]RawArticle, 0, len(table.files))
			for _, file := range table.files {
				content, err := os.Open(fmt.Sprintf("../../.testdata/%s", file))
				if err != nil {
					t.Fatal(err)
				}
				articles = append(articles, RawArticle{Content: content})
			}

			parsed := p.Parse(context.Background(), articles, table.page)
			if len(parsed.Articles) != table.expectedArticleLength {
				t.Errorf("Expected article length %d but got %d", table.expectedArticleLength, len(parsed.Articles))
			}
			if parsed.HasMore != table.expectedHasMore {
				t.Errorf("Expected hasMore %t but got %t", table.expectedHasMore, parsed.HasMore)
			}
		})
	}
}

func TestSortArticleDesc(t *testing.T) {
	tables := []struct {
		name   string
		a      time.Time
		b      time.Time
		aFirst bool
	}{
		{
			name:   "a zero",
			b:      time.Now(),
			aFirst: true,
		},
		{
			name:   "b zero",
			a:      time.Now(),
			aFirst: false,
		},
		{
			name:   "a earlier",
			a:      time.Date(2024, 10, 20, 0, 0, 0, 0, time.UTC),
			b:      time.Date(2024, 10, 19, 0, 0, 0, 0, time.UTC),
			aFirst: true,
		},
		{
			name: "b earlier",
			a:    time.Date(2024, 10, 19, 0, 0, 0, 0, time.UTC),
			b:    time.Date(2024, 10, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {

			slice := []ParsedArticle{
				{Meta: Meta{Title: "a", PublishedAt: table.a}},
				{Meta: Meta{Title: "b", PublishedAt: table.b}},
			}
			slices.SortFunc(slice, sortArticleDesc)

			aFirst := slice[0].Meta.Title == "a"
			if aFirst != table.aFirst {
				t.Error("Expected a to be first but got b")
			}
		})
	}
}

func TestDetectFileFormat(t *testing.T) {
	testCases := []struct {
		name     string
		file     string
		expected FileFormat
	}{
		{
			name:     "OG format with frontmatter",
			file:     "---\ntitle: Test\n---",
			expected: OG,
		},
		{
			name:     "KeepAChangelog format",
			file:     "# Changelog",
			expected: KeepAChangelog,
		},
		{
			name:     "Empty line",
			file:     "",
			expected: OG,
		},
		{
			name:     "Line without frontmatter",
			file:     "This is a regular line",
			expected: KeepAChangelog,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.file)
			result, read := detectFileFormat(r)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}

			rest, _ := io.ReadAll(r)
			expectedRest, _ := strings.CutPrefix(tc.file, read)
			if string(rest) != expectedRest {
				t.Errorf("Expected rest=%q got=%q", expectedRest, string(rest))
			}
		})
	}
}
