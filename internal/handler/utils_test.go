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

func TestGetQueryIDs(t *testing.T) {
	tables := []struct {
		url   string
		hxURL string
		wID   string
		cID   string
	}{
		{
			url: "/?wid=ws_1&cid=cl_1",
			wID: "ws_1",
			cID: "cl_1",
		},
		{
			url: "/",
			wID: "",
			cID: "",
		},
		{
			hxURL: "http://localhost:6001/?wid=ws_2&cid=cl_2",
			wID:   "ws_2",
			cID:   "cl_2",
		},
	}

	for _, table := range tables {
		u, err := url.Parse(table.url)
		if err != nil {
			t.Error(err)
		}

		r := &http.Request{URL: u, Header: http.Header{}}
		if table.hxURL != "" {
			r.Header.Set("HX-Current-URL", table.hxURL)
		}

		wID, cID := GetQueryIDs(r)
		if wID != table.wID {
			t.Errorf("Expected %s to equals %s", wID, table.wID)
		}
		if cID != table.cID {
			t.Errorf("Expected %s to equals %s", cID, table.cID)
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
