package parse

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/jonashiltl/openchangelog/internal"
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
	p := NewKeepAChangelogParser(CreateGoldmark())
	file, read, err := openKTestDataAndDetect("minimal")
	if err != nil {
		t.Fatal(err)
	}

	parsed := p.parse(read, file, internal.NoPagination())
	if parsed.HasMore == true {
		t.Error("hasMore should be false")
	}

	if len(parsed.ReleaseNotes) != 1 {
		t.Fatalf("Expected 1 parsed article but got %d", len(parsed.ReleaseNotes))
	}

	article := parsed.ReleaseNotes[0]

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
	p := NewKeepAChangelogParser(CreateGoldmark())
	file, read, err := openKTestDataAndDetect("unreleased")
	if err != nil {
		t.Fatal(err)
	}

	parsed := p.parse(read, file, internal.NoPagination())
	if parsed.HasMore == true {
		t.Error("hasMore should be false")
	}

	if len(parsed.ReleaseNotes) != 1 {
		t.Fatalf("Expected 1 parsed article but got %d", len(parsed.ReleaseNotes))
	}

	article := parsed.ReleaseNotes[0]

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
	p := NewKeepAChangelogParser(CreateGoldmark())
	file, read, err := openKTestDataAndDetect("full")
	if err != nil {
		t.Fatal(err)
	}

	parsed := p.parse(read, file, internal.NoPagination())
	if parsed.HasMore == true {
		t.Error("hasMore should be false")
	}

	if len(parsed.ReleaseNotes) != 15 {
		t.Errorf("Expected 15 parsed article but got %d", len(parsed.ReleaseNotes))
	}
}

func TestKParseGitCliff(t *testing.T) {
	p := NewKeepAChangelogParser(CreateGoldmark())
	file, read, err := openKTestDataAndDetect("git-cliff")
	if err != nil {
		t.Fatal(err)
	}

	parsed := p.parse(read, file, internal.NoPagination())
	if len(parsed.ReleaseNotes) != 1 {
		t.Fatalf("Expected 1 parsed article but got %d", len(parsed.ReleaseNotes))
	}

	article := parsed.ReleaseNotes[0]

	expectedTitle := "2.6.1"
	if article.Meta.Title != expectedTitle {
		t.Errorf("Expected %s to equal %s", article.Meta.Title, expectedTitle)
	}

}

func TestKParsePagination(t *testing.T) {
	p := NewKeepAChangelogParser(CreateGoldmark())

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

		page := internal.NewPagination(table.size, table.page)

		parsed := p.parse(read, file, internal.NewPagination(table.size, table.page))

		if parsed.HasMore != table.expectedHasMore {
			t.Errorf("Expected hasMore %t but got %t", table.expectedHasMore, parsed.HasMore)
		}

		if len(parsed.ReleaseNotes) != table.expectedSize {
			t.Errorf("Expected %d parsed article but got %d", table.expectedSize, len(parsed.ReleaseNotes))
		}

		for i, a := range parsed.ReleaseNotes {
			idx := page.StartIdx() + i
			if a.Meta.Title != expectedTitle[idx] {
				t.Errorf("Expected %s to equal %s", a.Meta.Title, expectedTitle[i])
			}
		}
	}
}

func TestParseTitle(t *testing.T) {
	tables := []struct {
		name          string
		firstLine     string
		expectedTitle string
	}{
		{
			name:          "basic keep a changelog title",
			firstLine:     "## [2.6.1] - 2024-09-27",
			expectedTitle: "2.6.1",
		},
		{
			name:          "no []",
			firstLine:     "## 2.6.1 - 2024-09-27",
			expectedTitle: "2.6.1",
		},
		{
			name:          "with link",
			firstLine:     "## [2.6.1](https://github.com/orhun/git-cliff/compare/v2.6.0..v2.6.1) - 2024-09-27",
			expectedTitle: "2.6.1",
		},
		{
			name:          "no release date",
			firstLine:     "## [2.6.1]",
			expectedTitle: "2.6.1",
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			title := parseTitle(table.firstLine)
			if title != table.expectedTitle {
				t.Errorf("Expected %s to be %s", title, table.expectedTitle)
			}
		})
	}
}
