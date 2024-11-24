package parse

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func openOGTestData(name string) (*os.File, error) {
	file, err := os.Open(fmt.Sprintf("../../.testdata/%s.md", name))
	if err != nil {
		return nil, err
	}
	return file, nil
}

func TestOGParseArticle(t *testing.T) {
	p := NewOGParser(CreateGoldmark())
	file, err := openOGTestData("v0.0.1-commonmark")
	if err != nil {
		t.Fatal(err)
	}

	parsed, err := p.parseReleaseNote(file)
	if err != nil {
		t.Fatal(err)
	}

	expectedTags := []string{"Improvement"}
	if !reflect.DeepEqual(parsed.Meta.Tags, expectedTags) {
		t.Errorf("Expected %s to equal %s", parsed.Meta.Tags, expectedTags)
	}

	expectedTitle := "CommonMark 0.31.2 compliance"
	if parsed.Meta.Title != expectedTitle {
		t.Errorf("Expected %s to equal %s", parsed.Meta.Title, expectedTitle)
	}

	expectedPublishedAt := time.Date(2024, 4, 3, 0, 0, 0, 0, time.UTC)
	if parsed.Meta.PublishedAt != expectedPublishedAt {
		t.Errorf("Expected %s to equal %s", parsed.Meta.PublishedAt, expectedPublishedAt)
	}
}

func TestOGParseArticleRead(t *testing.T) {
	p := NewOGParser(CreateGoldmark())
	file, err := openOGTestData("v0.0.5-beta")
	if err != nil {
		t.Fatal(err)
	}
	_, read := detectFileFormat(file)
	parsed, err := p.parseReleaseNoteRead(read, file)
	if err != nil {
		t.Fatal(err)
	}

	expectedTags := []string{"Community", "Cloud"}
	if !reflect.DeepEqual(parsed.Meta.Tags, expectedTags) {
		t.Errorf("Expected %s to equal %s", parsed.Meta.Tags, expectedTags)
	}

	expectedTitle := "Open Beta"
	if parsed.Meta.Title != expectedTitle {
		t.Errorf("Expected %s to equal %s", parsed.Meta.Title, expectedTitle)
	}

	expectedPublishedAt := time.Date(2024, 8, 26, 0, 0, 0, 0, time.UTC)
	if parsed.Meta.PublishedAt != expectedPublishedAt {
		t.Errorf("Expected %s to equal %s", parsed.Meta.PublishedAt, expectedPublishedAt)
	}
}
