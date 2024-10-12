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
		l, err = loadFullChangelogDBMode(e, r)
	} else {
		l, err = loadFullChangelogConfigMode(e, r)
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
	link := handler.FeedToChangelogURL(r)
	args := map[string]any{
		"CL":       parsed.CL,
		"Articles": parsed.Articles,
		"HasMore":  parsed.HasMore,
		"Link":     strings.ReplaceAll(link, "&", "&amp;"), // & is reserved in xml
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

func loadFullChangelogDBMode(e *env, r *http.Request) (*changelog.LoadedChangelog, error) {
	query := r.URL.Query()
	wID := query.Get(handler.WS_ID_QUERY)
	cID := query.Get(handler.CL_ID_QUERY)
	if wID != "" && cID != "" {
		return e.loader.FromWorkspace(r.Context(), wID, cID, changelog.NoPagination())
	}

	host := r.Host
	if r.Header.Get("X-Forwarded-Host") != "" {
		host = r.Header.Get("X-Forwarded-Host")
	}

	return e.loader.FromHost(r.Context(), host, changelog.NoPagination())
}

func loadFullChangelogConfigMode(e *env, r *http.Request) (*changelog.LoadedChangelog, error) {
	return e.loader.FromConfig(r.Context(), changelog.NoPagination())
}
