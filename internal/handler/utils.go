package handler

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	WS_ID_QUERY = "wid"
	CL_ID_QUERY = "cid"
)

func ParseSubdomain(host string) string {
	// Remove port if present
	host = strings.Split(host, ":")[0]
	parts := strings.Split(host, ".")
	if parts[0] == "www" {
		parts = parts[1:]
	}

	// subdomain exists, e.g. tenant.openchangelog.com
	if len(parts) > 2 {
		return parts[0]
	}
	return ""
}

func ChangelogToFeedURL(r *http.Request) string {
	rq := r.URL.Query()
	// only copy the query params we want
	q := url.Values{}
	if len(rq.Get(WS_ID_QUERY)) > 0 {
		q.Add(WS_ID_QUERY, rq.Get(WS_ID_QUERY))
	}
	if len(rq.Get(CL_ID_QUERY)) > 0 {
		q.Add(CL_ID_QUERY, rq.Get(CL_ID_QUERY))
	}

	newURL := &url.URL{
		Scheme:   r.URL.Scheme,
		Host:     r.URL.Host,
		RawQuery: q.Encode(),
		Path:     "feed",
	}

	if newURL.Host == "" {
		newURL.Host = r.Host
	}
	if newURL.Scheme == "" {
		if r.TLS != nil {
			newURL.Scheme = "https"
		} else {
			newURL.Scheme = "http"
		}
	}
	return newURL.String()
}

// Parses the changelog url from a request to a changelogs feed.
func FeedToChangelogURL(r *http.Request) string {
	newURL := &url.URL{
		Scheme:   r.URL.Scheme,
		Host:     r.URL.Host,
		RawQuery: r.URL.RawQuery,
	}

	if newURL.Host == "" {
		newURL.Host = r.Host
	}
	if newURL.Scheme == "" {
		if r.TLS != nil {
			newURL.Scheme = "https"
		} else {
			newURL.Scheme = "http"
		}
	}

	return newURL.String()
}
