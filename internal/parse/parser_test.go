package parse

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/source"
)

func TestParse(t *testing.T) {
	p := NewParser(CreateGoldmark())
	tables := []struct {
		name                  string
		files                 []string
		page                  internal.Pagination
		expectedHasMore       bool
		expectedArticleLength int
	}{
		{
			name:                  "OG one",
			files:                 []string{"v0.0.1-commonmark.md"},
			page:                  internal.NoPagination(),
			expectedHasMore:       false,
			expectedArticleLength: 1,
		},
		{
			name:                  "OG multiple",
			files:                 []string{"v0.0.1-commonmark.md", "v0.0.2-open-source.md", "v0.0.5-beta.md"},
			page:                  internal.NoPagination(),
			expectedHasMore:       false,
			expectedArticleLength: 3,
		},
		{
			name:                  "OG pagination not used",
			files:                 []string{"v0.0.1-commonmark.md", "v0.0.2-open-source.md", "v0.0.5-beta.md"},
			page:                  internal.NewPagination(2, 1),
			expectedHasMore:       false,
			expectedArticleLength: 3,
		},
		{
			name:                  "Keepachangelog no pagination",
			files:                 []string{"keepachangelog/minimal.md"},
			page:                  internal.NoPagination(),
			expectedHasMore:       false,
			expectedArticleLength: 1,
		},
		{
			name:                  "Keepachangelog full no pagination",
			files:                 []string{"keepachangelog/full.md"},
			page:                  internal.NoPagination(),
			expectedHasMore:       false,
			expectedArticleLength: 15,
		},
		{
			name:                  "Keepachangelog full paginated",
			files:                 []string{"keepachangelog/full.md"},
			page:                  internal.NewPagination(4, 2),
			expectedHasMore:       true,
			expectedArticleLength: 4,
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			notes := make([]source.RawReleaseNote, 0, len(table.files))
			for _, file := range table.files {
				content, err := os.Open(fmt.Sprintf("../../.testdata/%s", file))
				if err != nil {
					t.Fatal(err)
				}
				notes = append(notes, source.RawReleaseNote{Content: content})
			}

			parsed := p.Parse(context.Background(), notes, table.page)
			if len(parsed.Articles) != table.expectedArticleLength {
				t.Errorf("Expected article length %d but got %d", table.expectedArticleLength, len(parsed.Articles))
			}
			if parsed.HasMore != table.expectedHasMore {
				t.Errorf("Expected hasMore %t but got %t", table.expectedHasMore, parsed.HasMore)
			}
		})
	}
}
