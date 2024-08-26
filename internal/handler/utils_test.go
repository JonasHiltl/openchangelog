package handler

import (
	"net/http"
	"net/url"
	"testing"
)

func TestParseSubdomain(t *testing.T) {
	tables := []struct {
		host      string
		subdomain string
	}{
		{
			host:      "tenant.openchangelog.com",
			subdomain: "tenant",
		},
		{
			host:      "tenant-2.openchangelog.com",
			subdomain: "tenant-2",
		},
		{
			host:      "openchangelog.com",
			subdomain: "",
		},
		{
			host:      "www.openchangelog.com",
			subdomain: "",
		},
		{
			host:      "",
			subdomain: "",
		},
		{
			host:      ".",
			subdomain: "",
		},
		{
			host:      ".com",
			subdomain: "",
		},
	}

	for _, table := range tables {
		s := ParseSubdomain(table.host)
		if table.subdomain != s {
			t.Fatalf("expected %s to equal %s", s, table.subdomain)
		}
	}
}

func TestFeedToChangelogURL(t *testing.T) {
	tables := []struct {
		requestURL string
		host       string
		expected   string
	}{
		{
			requestURL: "/feed?wid=ws_cqj9svnd5lbga0eemd00&cid=cl_cqj9t0fd5lbga0eemd10",
			host:       "openchangelog.com",
			expected:   "http://openchangelog.com?wid=ws_cqj9svnd5lbga0eemd00&cid=cl_cqj9t0fd5lbga0eemd10",
		},
		{
			requestURL: "https://tenant.openchangelog.com/feed",
			host:       "tenant.openchangelog.com",
			expected:   "https://tenant.openchangelog.com",
		},
		{
			requestURL: "/feed",
			host:       "tenant.openchangelog.com",
			expected:   "http://tenant.openchangelog.com",
		},
	}

	for _, table := range tables {
		u, _ := url.Parse(table.requestURL)
		r := &http.Request{
			URL:  u,
			Host: table.host,
		}
		changelogURL := FeedToChangelogURL(r)
		if changelogURL != table.expected {
			t.Fatalf("expected %s to equal %s", changelogURL, table.expected)
		}
	}
}

func TestChangelogToFeedURL(t *testing.T) {
	tables := []struct {
		requestURL string
		host       string
		expected   string
	}{
		{
			requestURL: "/?wid=ws_cqj9svnd5lbga0eemd00&cid=cl_cqj9t0fd5lbga0eemd10",
			host:       "openchangelog.com",
			expected:   "http://openchangelog.com/feed?cid=cl_cqj9t0fd5lbga0eemd10&wid=ws_cqj9svnd5lbga0eemd00",
		},
		{
			requestURL: "https://tenant.openchangelog.com",
			host:       "tenant.openchangelog.com",
			expected:   "https://tenant.openchangelog.com/feed",
		},
		{
			requestURL: "/",
			host:       "tenant.openchangelog.com",
			expected:   "http://tenant.openchangelog.com/feed",
		},
		{
			requestURL: "/?page-size=5&page=2",
			host:       "tenant.openchangelog.com",
			expected:   "http://tenant.openchangelog.com/feed",
		},
	}

	for _, table := range tables {
		u, _ := url.Parse(table.requestURL)
		r := &http.Request{
			URL:  u,
			Host: table.host,
		}
		changelogURL := ChangelogToFeedURL(r)
		if changelogURL != table.expected {
			t.Fatalf("expected %s to equal %s", changelogURL, table.expected)
		}
	}
}
