package changelog

import (
	"context"
	"io"
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

	expectedContent := `
### Added

- Arabic translation (#444).
- v1.1 French translation.
- v1.1 Dutch translation (#371).

### Fixed

- Improve French translation (#377).
- Improve id-ID translation (#416).
- Improve Persian translation (#457).

### Changed

- Upgrade dependencies: Ruby 3.2.1, Middleman, etc.

### Removed

- Unused normalize.css file.
- Identical links assigned in each translation file.
- Duplicate index file for the english version.
`
	all, err := io.ReadAll(article.Content)
	if err != nil {
		t.Fatal(err)
	}
	content := string(all)
	if expectedContent != content {
		t.Errorf("Expected \n%s to equal \n%s", content, expectedContent)
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

	expectedPublishedAt := time.Time{}
	if article.Meta.PublishedAt != expectedPublishedAt {
		t.Errorf("Expected %s to equal %s", article.Meta.PublishedAt, expectedPublishedAt)
	}

	expectedContent := `
### Added

- v1.1 Brazilian Portuguese translation.
- v1.1 German Translation
- v1.1 Spanish translation.
- v1.1 Italian translation.
- v1.1 Polish translation.
- v1.1 Ukrainian translation.

### Changed

- Use frontmatter title & description in each language version template
- Replace broken OpenGraph image with an appropriately-sized Keep a Changelog 
  image that will render properly (although in English for all languages)
- Fix OpenGraph title & description for all languages so the title and 
description when links are shared are language-appropriate

### Removed

- Trademark sign previously shown after the project description in version 
0.3.0
`
	all, err := io.ReadAll(article.Content)
	if err != nil {
		t.Fatal(err)
	}
	content := string(all)
	if expectedContent != content {
		t.Errorf("Expected \n%s to equal \n%s", content, expectedContent)
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
		t.Fatalf("Expected 15 parsed article but got %d", len(parsed))
	}
}
