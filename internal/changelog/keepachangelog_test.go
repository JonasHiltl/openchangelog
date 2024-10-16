package changelog

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestKParseMinimal(t *testing.T) {
	ctx := context.Background()
	p := NewKeepAChangelogParser()

	file, err := os.Open("../../.testdata/keepachangelog/minimal.md")
	if err != nil {
		t.Fatal(err)
	}

	parsed, hasMore := p.Parse(ctx, RawArticle{Content: file}, NoPagination())
	if hasMore == true {
		t.Error("hasMore should be false")
	}

	if len(parsed) != 1 {
		t.Fatalf("Expected 1 parsed article but got %d", len(parsed))
	}

	article := parsed[0]

	expectedTags := []string{"Added", "Fixed", "Changed", "Removed"}
	if !reflect.DeepEqual(article.Meta.Tags, expectedTags) {
		t.Errorf("Expected %s to equal %s", article.Meta.Title, expectedTags)
	}

	expectedTitle := "1.1.1"
	if article.Meta.Title != expectedTitle {
		t.Errorf("Expected %s to equal %s", article.Meta.Title, expectedTitle)
	}

	expectedPublishedAt := time.Date(2023, 3, 5, 0, 0, 0, 0, time.UTC)
	if article.Meta.PublishedAt != expectedPublishedAt {
		t.Errorf("Expected %s to equal %s", article.Meta.PublishedAt, expectedPublishedAt)
	}
}

func TestKParseUnreleased(t *testing.T) {
	ctx := context.Background()
	p := NewKeepAChangelogParser()

	file, err := os.Open("../../.testdata/keepachangelog/unreleased.md")
	if err != nil {
		t.Fatal(err)
	}

	parsed, hasMore := p.Parse(ctx, RawArticle{Content: file}, NoPagination())
	if hasMore == true {
		t.Error("hasMore should be false")
	}

	if len(parsed) != 1 {
		t.Fatalf("Expected 1 parsed article but got %d", len(parsed))
	}

	article := parsed[0]

	expectedTags := []string{"Added", "Changed", "Removed"}
	if !reflect.DeepEqual(article.Meta.Tags, expectedTags) {
		t.Errorf("Expected %s to equal %s", article.Meta.Title, expectedTags)
	}

	expectedTitle := "Unreleased"
	if article.Meta.Title != expectedTitle {
		t.Errorf("Expected %s to equal %s", article.Meta.Title, expectedTitle)
	}

	expectedID := "unreleased"
	if article.Meta.ID != expectedID {
		t.Errorf("Expected %s to equal %s", article.Meta.ID, expectedID)
	}
}

func TestKParseFull(t *testing.T) {
	ctx := context.Background()
	p := NewKeepAChangelogParser()

	file, err := os.Open("../../.testdata/keepachangelog/full.md")
	if err != nil {
		t.Fatal(err)
	}

	parsed, hasMore := p.Parse(ctx, RawArticle{Content: file}, NoPagination())
	if hasMore == true {
		t.Error("hasMore should be false")
	}

	if len(parsed) != 15 {
		t.Errorf("Expected 15 parsed article but got %d", len(parsed))
	}
}

func TestKParsePagination(t *testing.T) {
	ctx := context.Background()
	p := NewKeepAChangelogParser()

	tables := []struct {
		size            int
		expectedSize    int
		page            int
		expectedHasMore bool
	}{
		{
			size:            3,
			page:            2,
			expectedSize:    3,
			expectedHasMore: true,
		},
		{
			size:            1,
			page:            1,
			expectedSize:    1,
			expectedHasMore: true,
		},
		{
			size:         15,
			page:         1,
			expectedSize: 15,
		},
		{
			size:            6,
			page:            2,
			expectedSize:    6,
			expectedHasMore: true,
		},
		{
			size:         8,
			page:         2,
			expectedSize: 7,
		},
		{
			size:         14,
			page:         2,
			expectedSize: 1,
		},
	}

	expectedTitle := []string{"Unreleased", "1.1.1", "1.1.0", "1.0.0", "0.3.0", "0.2.0", "0.1.0", "0.0.8", "0.0.7", "0.0.6", "0.0.5", "0.0.4", "0.0.3", "0.0.2", "0.0.1"}

	for _, table := range tables {
		file, err := os.Open("../../.testdata/keepachangelog/full.md")
		if err != nil {
			t.Fatal(err)
		}

		page := NewPagination(table.size, table.page)
		parsed, hasMore := p.Parse(ctx, RawArticle{Content: file}, NewPagination(table.size, table.page))

		if hasMore != table.expectedHasMore {
			t.Errorf("Expected hasMore %t but got %t", table.expectedHasMore, hasMore)
		}

		if len(parsed) != table.expectedSize {
			t.Errorf("Expected %d parsed article but got %d", table.expectedSize, len(parsed))
		}

		for i, a := range parsed {
			idx := page.StartIdx() + i
			if a.Meta.Title != expectedTitle[idx] {
				t.Errorf("Expected %s to equal %s", a.Meta.Title, expectedTitle[i])
			}
		}
	}
}
