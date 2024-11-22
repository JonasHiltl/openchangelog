package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestFeedToChangelogURL(t *testing.T) {
	tables := []struct {
		requestURL string
		host       string
		expected   string
	}{
		{
			requestURL: "https://tenant.openchangelog.com/feed",
			host:       "tenant.openchangelog.com",
			expected:   "https://tenant.openchangelog.com",
		},
		{
			requestURL: "/feed",
			host:       "tenant.openchangelog.com",
			expected:   "https://tenant.openchangelog.com",
		},
		{
			requestURL: "/feed",
			host:       "localhost:6001",
			expected:   "http://localhost:6001",
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

func TestGetFeedURL(t *testing.T) {
	tables := []struct {
		requestURL string
		host       string
		expected   string
	}{
		{
			requestURL: "/",
			host:       "openchangelog.com",
			expected:   "https://openchangelog.com/feed",
		},
		{
			requestURL: "https://tenant.openchangelog.com",
			host:       "tenant.openchangelog.com",
			expected:   "https://tenant.openchangelog.com/feed",
		},
		{
			requestURL: "/",
			host:       "tenant.openchangelog.com",
			expected:   "https://tenant.openchangelog.com/feed",
		},
		{
			requestURL: "/?page-size=5&page=2",
			host:       "tenant.openchangelog.com",
			expected:   "https://tenant.openchangelog.com/feed",
		},
		{
			requestURL: "/?page-size=5&page=2",
			host:       "localhost:6001",
			expected:   "http://localhost:6001/feed",
		},
	}

	for _, table := range tables {
		u, _ := url.Parse(table.requestURL)
		r := &http.Request{
			URL:  u,
			Host: table.host,
		}
		feedURL := GetFeedURL(r)
		if feedURL != table.expected {
			t.Fatalf("expected %s to equal %s", feedURL, table.expected)
		}
	}
}

func TestGetFullURL(t *testing.T) {
	tables := []struct {
		requestURL string
		hxURL      string
		host       string
		expected   string
	}{
		{
			requestURL: "/",
			host:       "openchangelog.com",
			expected:   "https://openchangelog.com/",
		},
		{
			requestURL: "https://tenant.openchangelog.com",
			host:       "tenant.openchangelog.com",
			expected:   "https://tenant.openchangelog.com",
		},
		{
			requestURL: "/",
			host:       "tenant.openchangelog.com",
			expected:   "https://tenant.openchangelog.com/",
		},
		{
			requestURL: "/?page-size=5&page=2",
			host:       "tenant.openchangelog.com",
			expected:   "https://tenant.openchangelog.com/?page-size=5&page=2",
		},
		{
			requestURL: "/?page-size=5&page=2",
			host:       "localhost:6001",
			expected:   "http://localhost:6001/?page-size=5&page=2",
		},
		{
			hxURL:    "http://localhost:6001/?page-size=2",
			host:     "localhost:6002",
			expected: "http://localhost:6001/?page-size=2",
		},
	}

	for _, table := range tables {
		u, _ := url.Parse(table.requestURL)
		r := &http.Request{
			URL:    u,
			Host:   table.host,
			Header: http.Header{},
		}
		if table.hxURL != "" {
			r.Header.Set("HX-Current-URL", table.hxURL)
		}

		fullURL := GetFullURL(r)
		if fullURL != table.expected {
			t.Fatalf("expected %s to equal %s", fullURL, table.expected)
		}
	}
}

func TestParsePagination(t *testing.T) {
	tables := []struct {
		page     int
		pageSize int
	}{
		{
			page:     0,
			pageSize: 0,
		},
		{
			page:     1,
			pageSize: 1,
		},
		{
			page:     10,
			pageSize: 10,
		},
	}

	for _, table := range tables {
		q := url.Values{}
		q.Set("page", fmt.Sprint(table.page))
		q.Set("page-size", fmt.Sprint(table.pageSize))
		p, s := ParsePagination(q)
		if table.page != p {
			t.Errorf("expected %d to equal %d", p, table.page)
		}
		if table.pageSize != s {
			t.Errorf("expected %d to equal %d", s, table.pageSize)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tables := []struct {
		password string
		hash     string
		isValid  bool
	}{
		{
			password: "password",
			hash:     "$2a$12$tlMf1uwAtIYLKjrDeEpQlORscaoxSiMeQ0eHbigRVk/UlVkRMUe9G",
			isValid:  true,
		},
		{
			password: "password2",
			hash:     "$2a$12$tlMf1uwAtIYLKjrDeEpQlORscaoxSiMeQ0eHbigRVk/UlVkRMUe9G",
			isValid:  false,
		},
		{
			password: "",
			hash:     "",
			isValid:  false,
		},
	}

	for _, table := range tables {
		err := ValidatePassword(table.hash, table.password)
		if table.isValid != (err == nil) {
			t.Errorf("Expected valid password but got error: %s", err)
		}
	}
}
