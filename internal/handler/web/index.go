package web

import (
	"log"
	"net/http"

	"github.com/jonashiltl/openchangelog/components"
	"github.com/jonashiltl/openchangelog/internal/changelog"
	"github.com/jonashiltl/openchangelog/internal/handler"
	"github.com/jonashiltl/openchangelog/internal/handler/web/views"
	"github.com/jonashiltl/openchangelog/render"
)

func index(e *env, w http.ResponseWriter, r *http.Request) error {
	page, pageSize := handler.ParsePagination(r.URL.Query())

	var err error
	var l *changelog.LoadedChangelog
	if e.cfg.IsDBMode() {
		l, err = loadChangelogDBMode(e, r, changelog.NewPagination(pageSize, page))
	} else {
		l, err = loadChangelogConfigMode(e, r, changelog.NewPagination(pageSize, page))
	}
	if err != nil {
		return err
	}

	parsed, err := l.Parse(r.Context())
	if err != nil {
		return err
	}

	if parsed.CL.Protected {
		err = ensurePasswordProvided(r, parsed.CL.PasswordHash)
		if err != nil {
			log.Printf("Blocked access to protected changelog: %s\n", parsed.CL.ID)
			return views.PasswordProtection(views.PasswordProtectionArgs{
				ThemeArgs: components.ThemeArgs{
					ColorScheme: parsed.CL.ColorScheme.ToApiTypes(),
				},
				FooterArgs: components.FooterArgs{
					HidePoweredBy: parsed.CL.HidePoweredBy,
				},
				BaseCSSVersion: e.baseCSSVersion,
			}).Render(r.Context(), w)
		}
	}

	if htmxHeader := r.Header.Get("HX-Request"); len(htmxHeader) > 0 {
		if len(parsed.Articles) > 0 {
			return e.render.RenderArticleList(r.Context(), w, render.RenderArticleListArgs{
				WID:      parsed.CL.WorkspaceID,
				CID:      parsed.CL.ID,
				Articles: parsed.Articles,
				HasMore:  parsed.HasMore,
				NextPage: page + 1,
				PageSize: pageSize,
			})
		} else {
			w.WriteHeader(http.StatusNoContent)
			return nil
		}
	}

	return e.render.RenderChangelog(r.Context(), w, render.RenderChangelogArgs{
		FeedURL:        handler.ChangelogToFeedURL(r),
		CL:             parsed.CL,
		Articles:       parsed.Articles,
		HasMore:        parsed.HasMore,
		PageSize:       pageSize,
		NextPage:       page + 1,
		BaseCSSVersion: e.baseCSSVersion,
	})
}

func ensurePasswordProvided(r *http.Request, pwHash string) error {
	value, err := getProtectedCookieValue(r)
	if err == nil && value == pwHash {
		// user already entered the password before
		return nil
	}

	authorize := r.URL.Query().Get("authorize")
	return validatePassword(pwHash, authorize)
}

func loadChangelogDBMode(e *env, r *http.Request, page changelog.Pagination) (*changelog.LoadedChangelog, error) {
	query := r.URL.Query()
	wID := query.Get(handler.WS_ID_QUERY)
	cID := query.Get(handler.CL_ID_QUERY)
	if wID != "" && cID != "" {
		return e.loader.FromWorkspace(r.Context(), wID, cID, page)
	}

	host := r.Host
	if r.Header.Get("X-Forwarded-Host") != "" {
		host = r.Header.Get("X-Forwarded-Host")
	}

	return e.loader.FromHost(r.Context(), host, changelog.NoPagination())
}

func loadChangelogConfigMode(e *env, r *http.Request, page changelog.Pagination) (*changelog.LoadedChangelog, error) {
	return e.loader.FromConfig(r.Context(), page)
}
