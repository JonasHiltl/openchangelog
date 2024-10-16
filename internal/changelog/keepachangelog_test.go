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

	parsed, err := p.Parse(ctx, []RawArticle{
		{
			Content: file,
		},
	})
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

	parsed, err := p.Parse(ctx, []RawArticle{
		{
			Content: file,
		},
	})
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

	parsed, err := p.Parse(ctx, []RawArticle{
		{
			Content: file,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(parsed) != 15 {
		t.Errorf("Expected 15 parsed article but got %d", len(parsed))
	}
}
