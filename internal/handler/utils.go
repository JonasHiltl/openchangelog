package handler

import (
	"net/http"
	"net/url"
	"strconv"
)

const (
	WS_ID_QUERY = "wid"
	CL_ID_QUERY = "cid"
)

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

func ParsePagination(q url.Values) (page int, size int) {
	const default_page, default_page_size = 1, 10
	page, err := strconv.Atoi(q.Get("page"))
	if err != nil {
		page = default_page
	}
	pageSize, err := strconv.Atoi(q.Get("page-size"))
	if err != nil {
		pageSize = default_page_size
	}

	return page, pageSize
}
