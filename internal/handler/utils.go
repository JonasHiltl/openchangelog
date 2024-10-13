package handler

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jonashiltl/openchangelog/internal/changelog"
	"golang.org/x/crypto/bcrypt"
)

const (
	WS_ID_QUERY     = "wid"
	CL_ID_QUERY     = "cid"
	AUTHORIZE_QUERY = "authorize"
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

func GetQueryIDs(r *http.Request) (wID string, cID string) {
	query := r.URL.Query()
	wID = query.Get(WS_ID_QUERY)
	cID = query.Get(CL_ID_QUERY)

	if wID == "" && cID == "" {
		u, err := url.Parse(r.Header.Get("HX-Current-URL"))
		if err == nil {
			query = u.Query()
			return query.Get(WS_ID_QUERY), query.Get(CL_ID_QUERY)
		}
	}
	return wID, cID
}

// If in db-mode => load changelog by query ids or host.
//
// If in config mode => load changelog from config.
func LoadChangelog(loader *changelog.Loader, isDBMode bool, r *http.Request, page changelog.Pagination) (*changelog.LoadedChangelog, error) {
	if isDBMode {
		return loadChangelogDBMode(loader, r, page)
	} else {
		return loader.FromConfig(r.Context(), page)
	}
}

func loadChangelogDBMode(loader *changelog.Loader, r *http.Request, page changelog.Pagination) (*changelog.LoadedChangelog, error) {
	wID, cID := GetQueryIDs(r)
	if wID != "" && cID != "" {
		return loader.FromWorkspace(r.Context(), wID, cID, page)
	}

	host := r.Host
	if r.Header.Get("X-Forwarded-Host") != "" {
		host = r.Header.Get("X-Forwarded-Host")
	}

	return loader.FromHost(r.Context(), host, page)
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
