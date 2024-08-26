package rss

import (
	_ "embed"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/jonashiltl/openchangelog/internal/changelog"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/handler"
)

//go:embed feed.tmpl
var feedTemplate string

func feedHandler(e *env, w http.ResponseWriter, r *http.Request) error {
	var l *changelog.LoadedChangelog
	var err error
	if e.cfg.IsDBMode() {
		l, err = loadChangelogDBMode(e, r)
	} else {
		l, err = loadChangelogConfigMode(e, r)
	}
	if err != nil {
		return err
	}

	parsed, err := l.Parse(r.Context())
	if err != nil {
		return err
	}

	tmpl, err := template.
		New("feed").
		Funcs(template.FuncMap{
			"addFragment": addFragment,
			"toRFC822":    toRFC822,
		}).
		Parse(feedTemplate)
	if err != nil {
		return errs.NewBadRequest(errors.New("failed to parse feed template"))
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	args := map[string]any{
		"CL":       parsed.CL,
		"Articles": parsed.Articles,
		"HasMore":  parsed.HasMore,
		"Link":     getChangelogURL(r),
	}
	return tmpl.Execute(w, args)
}

func toRFC822(t time.Time) string {
	return t.Format(time.RFC822)
}

// Adds a fragment to the specified url
func addFragment(u string, fragment string) string {
	parsed, err := url.Parse(u)
	if err != nil {
		return ""
	}
	parsed.Fragment = fragment
	return parsed.String()
}

// Returns the url at which the changelog is hosted.
func getChangelogURL(r *http.Request) string {
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

	return strings.ReplaceAll(newURL.String(), "&", "&amp;")
}

func loadChangelogDBMode(e *env, r *http.Request) (*changelog.LoadedChangelog, error) {
	query := r.URL.Query()
	wID := query.Get(handler.WS_ID_QUERY)
	cID := query.Get(handler.CL_ID_QUERY)
	if wID != "" && cID != "" {
		return e.loader.FromWorkspace(r.Context(), wID, cID, changelog.NoPagination())
	}

	subdomain := handler.ParseSubdomain(r.Host)
	if subdomain != "" {
		return e.loader.FromSubdomain(r.Context(), subdomain, changelog.NoPagination())
	}

	return nil, errs.NewServiceUnavailable(errors.New("you need to specify the subdomain or workspace & changelog id when running openchangelog in db mode"))
}

func loadChangelogConfigMode(e *env, r *http.Request) (*changelog.LoadedChangelog, error) {
	return e.loader.FromConfig(r.Context(), changelog.NoPagination())
}
