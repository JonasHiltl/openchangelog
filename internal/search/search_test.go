package search

import (
	"context"
	"fmt"
	"html"
	"strings"
	"testing"
	"time"

	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/source"
)

var sid = source.NewGitHubID("owner", "repo", "path")
var indexData = BatchIndexArgs{
	SID: sid.String(),
	ReleaseNotes: []parse.ParsedReleaseNote{
		{
			Meta: parse.Meta{
				ID:          "v0.5.0-analytics",
				Title:       "Analytics",
				Description: "Gain real-time insights into your changelog visitors",
				PublishedAt: time.Date(2024, 11, 8, 0, 0, 0, 0, time.UTC),
				Tags:        []string{"Cloud"},
			},
			Content: strings.NewReader(`
				<p><img src="https://i.ibb.co/bX101D5/Group-10.png" alt="An image to describe post "></p>
				<p>We're excited to introduce Analytics to Openchangelog!<br>You can now monitor key metrics, like daily visitor counts and country-based location data, directly from your changelog dashboard. Allowing you to understand your audience and optimize your changelog.</p>
				<p>For <strong>self-hosting</strong> we currently support <a href="https://www.tinybird.co">Tinybird</a> for storing analytics events.</p>
			`),
		},
		{
			Meta: parse.Meta{
				ID:          "v0.1.2-custom-domains",
				Title:       "Custom Domains",
				Description: "Bring your own domain to showcase your changelog",
				PublishedAt: time.Date(2024, 9, 14, 0, 0, 0, 0, time.UTC),
				Tags:        []string{"Feature", "Cloud"},
			},
			Content: strings.NewReader(`
				<p><img src="https://github.com/user-attachments/assets/ebc15809-bd1d-4a0e-abd8-2967627a1aec" alt="An image to describe post ">Want to host your changelog on a custom, branded domain like <strong>changelog.company.com</strong>?</p>
				<p>Now, with our new <strong>Custom Domain</strong> feature, you can easily point your changelog to any domain you own. SSL certificates are automatically managed by us, ensuring your changelog is secure without any extra effort on your end.</p>
			`),
		},
	},
}

func newMemorySearcher(t *testing.T) Searcher {
	s, err := NewSearcher(config.Config{
		Search: &config.SearchConfig{
			Type: config.SearchMem,
		},
	})
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	err = s.BatchIndex(ctx, indexData)
	if err != nil {
		t.Error(err)
	}
	return s
}

func TestSearch(t *testing.T) {
	searcher := newMemorySearcher(t)

	tests := []struct {
		name          string
		args          SearchArgs
		expectedTotal uint64
	}{
		{
			name: "only sid",
			args: SearchArgs{
				SID: sid.String(),
			},
			expectedTotal: 2,
		},
		{
			name: "single tag",
			args: SearchArgs{
				SID:  sid.String(),
				Tags: []string{"Cloud"},
			},
			expectedTotal: 2,
		},
		{
			name: "multiple tags",
			args: SearchArgs{
				SID:  sid.String(),
				Tags: []string{"Cloud", "Feature"},
			},
			expectedTotal: 1,
		},
		{
			name: "title",
			args: SearchArgs{
				SID:   sid.String(),
				Query: "Analytics",
			},
			expectedTotal: 1,
		},
		{
			name: "description",
			args: SearchArgs{
				SID:   sid.String(),
				Query: "showcase",
			},
			expectedTotal: 1,
		},
		{
			name: "content",
			args: SearchArgs{
				SID:   sid.String(),
				Query: "monitor key metrics",
			},
			expectedTotal: 1,
		},
		{
			name: "content case insensitive",
			args: SearchArgs{
				SID:   sid.String(),
				Query: "OWN",
			},
			expectedTotal: 1,
		},
		{
			name: "query and tags",
			args: SearchArgs{
				SID:   sid.String(),
				Tags:  []string{"Feature"},
				Query: "Analytics",
			},
			expectedTotal: 0,
		},
		{
			name: "no html",
			args: SearchArgs{
				SID:   sid.String(),
				Query: "<p>",
			},
			expectedTotal: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := searcher.Search(context.Background(), test.args)
			if err != nil {
				t.Error(err)
			}
			if res.Total != test.expectedTotal {
				t.Errorf("expected total %d to be %d", res.Total, test.expectedTotal)
			}
		})
	}
}

func TestBatchRemove(t *testing.T) {
	searcher := newMemorySearcher(t)
	ctx := context.Background()
	err := searcher.BatchRemove(ctx, BatchRemoveArgs{
		SID:          sid.String(),
		ReleaseNotes: indexData.ReleaseNotes,
	})
	if err != nil {
		t.Error(err)
	}

	res, err := searcher.Search(ctx, SearchArgs{
		SID: sid.String(),
	})
	if err != nil {
		t.Error(err)
	}

	if len(res.Hits) != 0 {
		t.Errorf("expected 0 hits, but got %d", len(res.Hits))
	}
}

func TestGetTags(t *testing.T) {
	searcher := newMemorySearcher(t)

	tags := searcher.GetAllTags(context.Background(), sid.String())
	if len(tags) != 2 {
		t.Errorf("expected 2 tags but got %d", len(tags))
	}
}

func TestHighlightTitle(t *testing.T) {
	searcher := newMemorySearcher(t)
	res, err := searcher.Search(context.Background(), SearchArgs{
		SID:   sid.String(),
		Query: "Domains",
	})
	if err != nil {
		t.Error(err)
	}
	if res.Total != 2 {
		t.Errorf("expected total %d to be 2", res.Total)
	}
	hit := res.Hits[0]
	if hit.Title != "Custom <mark>Domains</mark>" {
		t.Errorf("expected title \"%s\" to be \"<mark>Custom Domains</mark>\"", res.Hits[0].Title)
	}
	highlightTitle := hit.Fragments["Title"][0]
	if highlightTitle != "Custom <mark>Domains</mark>" {
		t.Errorf("expected highlight \"%s\" to be \"Custom <mark>Domains</mark>\"", highlightTitle)
	}
}

func TestStripPartialHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "no partial tag",
			html:     "test <ul><li>a</li></ul> b",
			expected: "test a b",
		},
		{
			name:     "partial closing tag at start",
			html:     "...agaadaw\">test</a> <li>a</li>",
			expected: "test a",
		},
		{
			name:     "partial tag at end",
			html:     "test <li>a<a>b<",
			expected: "test ab",
		},
		{
			name:     "partial closing tag at end",
			html:     "test <li>a<a>b</",
			expected: "test ab",
		},
		{
			name:     "longer partial closing tag at end",
			html:     "test <li>a<a>b</veryLongTag",
			expected: "test ab",
		},
		{
			name:     "missing closing tags",
			html:     "test <li>a<a>b",
			expected: "test ab",
		},
		{
			name:     "no html tags",
			html:     "this has no html tags",
			expected: "this has no html tags",
		},
		{
			name:     "empty string",
			html:     "",
			expected: "",
		},
		{
			name:     "release notes demo",
			html:     "...ca2\">77731ec</a>)</li> </ul> <h3>🐛 Bug Fixes</h3> <ul> <li><em>(fixtures)</em> Fix test failures - (<a href=\"https://gi...",
			expected: "77731ec)  🐛 Bug Fixes  (fixtures) Fix test failures - (",
		},
		{
			name:     "release notes demo 2",
			html:     "...error message to failing test - (<a href=\"...\">7d7470b</a>)</li> <li>Fix keep a changelog test case - (<a href=\"...\"",
			expected: "...error message to failing test - (7d7470b) Fix keep a changelog test case - (",
		},
		{
			name:     "preserve mark tags",
			html:     "bug <mark>fixes</mark>",
			expected: "bug <mark>fixes</mark>",
		},
		{
			name:     "preserve mark tags with partial",
			html:     "...ca2\">bu<a>g</a> <mark>fixes</mark>b</",
			expected: "bug <mark>fixes</mark>b",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stripped := stripPartialHTML(test.html)
			if stripped != test.expected {
				t.Errorf("expected \"%s\" to equal \"%s\"", stripped, test.expected)
			}
		})
		t.Run(fmt.Sprintf("%s escaped", test.name), func(t *testing.T) {
			stripped := stripPartialHTML(html.EscapeString(test.html))
			if stripped != test.expected {
				t.Errorf("expected \"%s\" to equal \"%s\"", stripped, test.expected)
			}
		})
	}
}

func TestSurroundWithEllipsis(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no ellipsis",
			input:    "test",
			expected: "...test...",
		},
		{
			name:     "prefix ellipsis",
			input:    "...test",
			expected: "...test...",
		},
		{
			name:     "suffix ellipsis",
			input:    "test...",
			expected: "...test...",
		},
		{
			name:     "suffix and prefix ellipsis",
			input:    "...test...",
			expected: "...test...",
		},
		{
			name:     "with ellipsis character",
			input:    "…test…",
			expected: "...test...",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := surroundWithEllipsis(test.input)
			if res != test.expected {
				t.Errorf("expected \"%s\" to equal \"%s\"", res, test.expected)
			}
		})
	}
}

func TestFirstNWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		n        int
		expected string
	}{
		{
			name:     "3 words",
			input:    "one two three four five six",
			n:        3,
			expected: "one two three",
		},
		{
			name:     "0 words",
			input:    "one two three four five six",
			n:        0,
			expected: "",
		},
		{
			name:     "n larger than words of input",
			input:    "one two three four five six",
			n:        10,
			expected: "one two three four five six",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nwords := firstNWords(test.input, test.n)
			if nwords != test.expected {
				t.Errorf("expected %s to equal %s", nwords, test.expected)
			}
		})
	}
}
