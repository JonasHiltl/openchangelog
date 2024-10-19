package changelog

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func openKTestData(name string) (*os.File, error) {
	file, err := os.Open(fmt.Sprintf("../../.testdata/keepachangelog/%s.md", name))
	if err != nil {
		return nil, err
	}
	return file, nil
}

func openKTestDataAndDetect(name string) (*os.File, string, error) {
	file, err := openKTestData(name)
	if err != nil {
		return nil, "", err
	}
	_, read := detectFileFormat(file)
	return file, read, nil
}

func TestKParseMinimal(t *testing.T) {
	p := NewKeepAChangelogParser(createGoldmark())
	file, read, err := openKTestDataAndDetect("minimal")
	if err != nil {
		t.Fatal(err)
	}

	parsed := p.parse(read, file, NoPagination())
	if parsed.HasMore == true {
		t.Error("hasMore should be false")
	}

	if len(parsed.Articles) != 1 {
		t.Fatalf("Expected 1 parsed article but got %d", len(parsed.Articles))
	}

	article := parsed.Articles[0]

	expectedTags := []string{"Added", "Fixed", "Changed", "Removed"}
	if !reflect.DeepEqual(article.Meta.Tags, expectedTags) {
		t.Errorf("Expected %s to equal %s", article.Meta.Tags, expectedTags)
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
	p := NewKeepAChangelogParser(createGoldmark())
	file, read, err := openKTestDataAndDetect("unreleased")
	if err != nil {
		t.Fatal(err)
	}

	parsed := p.parse(read, file, NoPagination())
	if parsed.HasMore == true {
		t.Error("hasMore should be false")
	}

	if len(parsed.Articles) != 1 {
		t.Fatalf("Expected 1 parsed article but got %d", len(parsed.Articles))
	}

	article := parsed.Articles[0]

	expectedTags := []string{"Added", "Changed", "Removed"}
	if !reflect.DeepEqual(article.Meta.Tags, expectedTags) {
		t.Errorf("Expected %s to equal %s", article.Meta.Tags, expectedTags)
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
	p := NewKeepAChangelogParser(createGoldmark())
	file, read, err := openKTestDataAndDetect("full")
	if err != nil {
		t.Fatal(err)
	}

	parsed := p.parse(read, file, NoPagination())
	if parsed.HasMore == true {
		t.Error("hasMore should be false")
	}

	if len(parsed.Articles) != 15 {
		t.Errorf("Expected 15 parsed article but got %d", len(parsed.Articles))
	}
}

func TestKParsePagination(t *testing.T) {
	p := NewKeepAChangelogParser(createGoldmark())

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
		file, read, err := openKTestDataAndDetect("full")
		if err != nil {
			t.Fatal(err)
		}

		page := NewPagination(table.size, table.page)

		parsed := p.parse(read, file, NewPagination(table.size, table.page))

		if parsed.HasMore != table.expectedHasMore {
			t.Errorf("Expected hasMore %t but got %t", table.expectedHasMore, parsed.HasMore)
		}

		if len(parsed.Articles) != table.expectedSize {
			t.Errorf("Expected %d parsed article but got %d", table.expectedSize, len(parsed.Articles))
		}

		for i, a := range parsed.Articles {
			idx := page.StartIdx() + i
			if a.Meta.Title != expectedTitle[idx] {
				t.Errorf("Expected %s to equal %s", a.Meta.Title, expectedTitle[i])
			}
		}
	}
}
