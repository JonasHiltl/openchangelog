package web

import (
	"net/http"
	"net/url"
	"testing"
)

func TestCreateCookieKey(t *testing.T) {
	tables := []struct {
		requestURL string
		host       string
		expected   string
	}{
		{
			host:     "changelog.openchangelog.com",
			expected: "protected-changelog.openchangelog.com",
		},
		{
			requestURL: "/",
			host:       "openchangelog.com",
			expected:   "protected-openchangelog.com",
		},
	}

	for _, table := range tables {
		u, err := url.Parse(table.requestURL)
		if err != nil {
			t.Error(err)
		}
		r := &http.Request{
			URL:  u,
			Host: table.host,
		}
		key := createCookieKey(r)
		if key != table.expected {
			t.Errorf("Expected %s to equal %s", key, table.expected)
		}
	}
}

func TestGetHost(t *testing.T) {
	tables := []struct {
		r        *http.Request
		expected string
	}{
		{
			r: &http.Request{
				Host: "localhost:6001",
			},
			expected: "localhost",
		},
		{
			r: &http.Request{
				Host: "openchangelog.com",
			},
			expected: "openchangelog.com",
		},
		{
			r: &http.Request{
				Host: "subdomain.openchangelog.com",
			},
			expected: "subdomain.openchangelog.com",
		},
	}

	for _, table := range tables {
		host := getHost(table.r)
		if host != table.expected {
			t.Errorf("Expected %s to equal %s", host, table.expected)
		}
	}
}
