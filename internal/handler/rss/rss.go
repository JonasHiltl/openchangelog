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
		"CL":       loaded.CL,
		"Articles": loaded.Notes,
		"HasMore":  loaded.HasMore,
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
