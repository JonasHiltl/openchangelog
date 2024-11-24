package parse

import (
	"io"
	"slices"
	"strings"
	"testing"
	"time"
)

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

			slice := []ParsedReleaseNote{
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
