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

	parsed, err := p.Parse(ctx, RawArticle{Content: file}, NoPagination())
	if err != nil {
		t.Fatal(err)
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

	parsed, err := p.Parse(ctx, RawArticle{Content: file}, NoPagination())
	if err != nil {
		t.Fatal(err)
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

	parsed, err := p.Parse(ctx, RawArticle{Content: file}, NoPagination())
	if err != nil {
		t.Fatal(err)
	}

	if len(parsed) != 15 {
		t.Errorf("Expected 15 parsed article but got %d", len(parsed))
	}
}

func TestKParsePagination(t *testing.T) {
	ctx := context.Background()
	p := NewKeepAChangelogParser()

	tables := []struct {
		size int
		page int
	}{
		{
			size: 3,
			page: 2,
		},
		{
			size: 1,
			page: 1,
		},
		{
			size: 15,
			page: 1,
		},
		{
			size: 6,
			page: 2,
		},
	}

	expectedTitle := []string{"Unreleased", "1.1.1", "1.1.0", "1.0.0", "0.3.0", "0.2.0", "0.1.0", "0.0.8", "0.0.7", "0.0.6", "0.0.5", "0.0.4", "0.0.3", "0.0.2", "0.0.1"}

	for _, table := range tables {
		file, err := os.Open("../../.testdata/keepachangelog/full.md")
		if err != nil {
			t.Fatal(err)
		}

		page := NewPagination(table.size, table.page)
		parsed, err := p.Parse(ctx, RawArticle{Content: file}, NewPagination(table.size, table.page))
		if err != nil {
			t.Fatal(err)
		}

		if len(parsed) != table.size {
			t.Errorf("Expected %d parsed article but got %d", table.size, len(parsed))
		}

		for i, a := range parsed {
			idx := page.StartIdx() + i
			if a.Meta.Title != expectedTitle[idx] {
				t.Errorf("Expected %s to equal %s", a.Meta.Title, expectedTitle[i])
			}
		}
	}
}
