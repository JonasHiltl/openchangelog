package load

import (
	"testing"
	"time"

	"github.com/jonashiltl/openchangelog/internal/parse"
)

func TestRemoveUnpublished(t *testing.T) {
	now := time.Now()

	tables := []struct {
		name     string
		notes    []parse.ParsedReleaseNote
		expected []string
	}{
		{
			name: "keeps past and zero-value publishedAt, drops future",
			notes: []parse.ParsedReleaseNote{
				{Meta: parse.Meta{ID: "past", PublishedAt: now.Add(-time.Hour)}},
				{Meta: parse.Meta{ID: "future", PublishedAt: now.Add(time.Hour)}},
				{Meta: parse.Meta{ID: "no-date", PublishedAt: time.Time{}}},
			},
			expected: []string{"past", "no-date"},
		},
		{
			name:     "empty input",
			notes:    []parse.ParsedReleaseNote{},
			expected: []string{},
		},
		{
			name: "all future",
			notes: []parse.ParsedReleaseNote{
				{Meta: parse.Meta{ID: "future1", PublishedAt: now.Add(time.Hour)}},
				{Meta: parse.Meta{ID: "future2", PublishedAt: now.Add(24 * time.Hour)}},
			},
			expected: []string{},
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			got := removeUnpublished(table.notes)
			if len(got) != len(table.expected) {
				t.Fatalf("expected %d notes but got %d", len(table.expected), len(got))
			}
			for i, n := range got {
				if n.Meta.ID != table.expected[i] {
					t.Errorf("expected note at index %d to be %q but got %q", i, table.expected[i], n.Meta.ID)
				}
			}
		})
	}
}
