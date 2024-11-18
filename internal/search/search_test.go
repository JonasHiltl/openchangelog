package search

import (
	"context"
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

func TestHighlightTitle(t *testing.T) {
	searcher := newMemorySearcher(t)
	res, err := searcher.Search(context.Background(), SearchArgs{
		SID:   sid.String(),
		Query: "Domains",
	})
	if err != nil {
		t.Error(err)
	}
	if res.Total != 1 {
		t.Errorf("expected total %d to be 1", res.Total)
	}
	hit := res.Hits[0]
	if hit.Title != "Custom Domains" {
		t.Errorf("expected title \"%s\" to be \"Custom Domains\"", res.Hits[0].Title)
	}
	highlightTitle := hit.Fragments["Title"][0]
	if highlightTitle != "Custom <mark>Domains</mark>" {
		t.Errorf("expected highlight \"%s\" to be \"Custom <mark>Domains</mark>\"", highlightTitle)
	}
}
