package rss

import (
	_ "embed"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/handler"
)

//go:embed feed.tmpl
var feedTemplate string

func feedHandler(e *env, w http.ResponseWriter, r *http.Request) error {
	loaded, err := e.loader.LoadAndParse(r, internal.NoPagination())
	if err != nil {
		return err
	}

	if loaded.CL.Protected {
		authorize := r.URL.Query().Get(handler.AUTHORIZE_QUERY)
		if authorize == "" {
			return errs.NewBadRequest(errors.New("can't load rss feed of protected changelog, specify \"authorize\" query param to subscribe"))
		}

		err = handler.ValidatePassword(loaded.CL.PasswordHash, authorize)
		if err != nil {
			return errs.NewBadRequest(err)
		}
	}

	createdAt := loaded.CL.CreatedAt
	if createdAt.IsZero() && len(loaded.Notes) > 0 {
		createdAt = loaded.Notes[0].Meta.PublishedAt
	}

	tmpl, err := template.
		New("feed").
		Funcs(template.FuncMap{
			"addFragment": addFragment,
			"toPubDate":   toPubDate,
		}).
		Parse(feedTemplate)
	if err != nil {
		return errs.NewBadRequest(errors.New("failed to parse feed template"))
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	link := handler.FeedToChangelogURL(r)
	args := map[string]any{
		"CL":        loaded.CL,
		"Articles":  loaded.Notes,
		"HasMore":   loaded.HasMore,
		"CreatedAt": createdAt,
		"Link":      strings.ReplaceAll(link, "&", "&amp;"), // & is reserved in xml
	}
	return tmpl.Execute(w, args)
}

func toPubDate(t time.Time) string {
	// time.RFC822 produces an different format than the expected format for RSS. Day of the week is missing.
	// time.RFC822 and time.RFC1123 may produce "UTC" as timezone but the spec only allows "GMT", "UT", "Z", or "0000".
	return strings.ReplaceAll(t.Format(time.RFC1123), "UTC", "GMT")
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
