package handler

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	WS_ID_QUERY     = "wid"
	CL_ID_QUERY     = "cid"
	AUTHORIZE_QUERY = "authorize"
)

// Turns the changelog request into the feed url of the changelog
func GetFeedURL(r *http.Request) string {
	rq := r.URL.Query()
	// only copy the query params we want, don't want page or page-size
	q := url.Values{}
	if len(rq.Get(WS_ID_QUERY)) > 0 {
		q.Add(WS_ID_QUERY, rq.Get(WS_ID_QUERY))
	}
	if len(rq.Get(CL_ID_QUERY)) > 0 {
		q.Add(CL_ID_QUERY, rq.Get(CL_ID_QUERY))
	}
	if len(rq.Get(AUTHORIZE_QUERY)) > 0 {
		q.Add(AUTHORIZE_QUERY, rq.Get(AUTHORIZE_QUERY))
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
	if strings.Contains(newURL.Host, "localhost") {
		newURL.Scheme = "http"
	} else {
		newURL.Scheme = "https"
	}
	return newURL.String()
}

// Turns the rss feed request into the changelog url.
// Done by stripping away the request path (/feed)
func FeedToChangelogURL(r *http.Request) string {
	newURL := &url.URL{
		Scheme:   r.URL.Scheme,
		Host:     r.URL.Host,
		RawQuery: r.URL.RawQuery,
	}

	if newURL.Host == "" {
		newURL.Host = r.Host
	}
	if strings.Contains(newURL.Host, "localhost") {
		newURL.Scheme = "http"
	} else {
		newURL.Scheme = "https"
	}

	return newURL.String()
}

// Returns the full url of the current request.
// If request is htmx request (password submit) will use HX-Current-URL
func GetFullURL(r *http.Request) string {
	var newURL *url.URL
	if r.Header.Get("HX-Current-URL") != "" {
		newURL, _ = url.Parse(r.Header.Get("HX-Current-URL"))
	} else {
		newURL, _ = url.Parse(r.URL.String()) // deep clone (dirty)
	}

	if newURL.Host == "" {
		newURL.Host = r.Host
	}
	if strings.Contains(newURL.Host, "localhost") {
		newURL.Scheme = "http"
	} else {
		newURL.Scheme = "https"
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

func ValidatePassword(hash, plaintext string) error {
	if hash == "" {
		return errors.New("protection is enabled, please configure the password")
	}
	if plaintext == "" {
		return errors.New("missing password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return errors.New("invalid password")
	}
	if err != nil {
		return err
	}
	return nil
}
