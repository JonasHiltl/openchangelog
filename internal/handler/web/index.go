package web

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/jonashiltl/openchangelog/components"
	"github.com/jonashiltl/openchangelog/internal/analytics"
	"github.com/jonashiltl/openchangelog/internal/changelog"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/handler"
	"github.com/jonashiltl/openchangelog/internal/handler/web/static"
	"github.com/jonashiltl/openchangelog/internal/handler/web/views"
)

func index(e *env, w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query()
	page, pageSize := handler.ParsePagination(q)
	l, err := handler.LoadChangelog(e.loader, e.cfg.IsDBMode(), r, changelog.NewPagination(pageSize, page))
	if err != nil {
		return err
	}

	_, isWidget := q["widget"]
	parsed := l.Parse(r.Context())

	if parsed.CL.Protected {
		if isWidget {
			return errs.NewBadRequest(errors.New("can't display protected changelog in widget"))
		}
		err = ensurePasswordProvided(r, parsed.CL.PasswordHash)
		if err != nil {
			slog.InfoContext(r.Context(), "blocked access to changelog", slog.String("changelog", parsed.CL.ID.String()))

			go e.getAnalyticsEmitter(parsed.CL).Emit(analytics.NewAccessDeniedEvent(r, parsed.CL))
			return views.PasswordProtection(views.PasswordProtectionArgs{
				CSS: static.BaseCSS,
				ThemeArgs: components.ThemeArgs{
					ColorScheme: parsed.CL.ColorScheme.ToApiTypes(),
				},
				FooterArgs: components.FooterArgs{
					HidePoweredBy: parsed.CL.HidePoweredBy,
				},
			}).Render(r.Context(), w)
		}
	}

	if _, ok := q["articles"]; ok {
		return handleArticles(e, w, r, &parsed, page, pageSize)
	}

	setCacheControlHeader(w, parsed.CL.Protected)
	return renderChangelog(e, w, r, &parsed, isWidget)
}

func ensurePasswordProvided(r *http.Request, pwHash string) error {
	value, err := getProtectedCookieValue(r)
	if err == nil && value == pwHash {
		// user already entered the password before
		return nil
	}

	authorize := r.URL.Query().Get(handler.AUTHORIZE_QUERY)
	return handler.ValidatePassword(pwHash, authorize)
}

func handleArticles(e *env, w http.ResponseWriter, r *http.Request, parsed *changelog.ParsedChangelog, page, pageSize int) error {
	if len(parsed.Articles) > 0 {
		return e.render.RenderArticleList(r.Context(), w, RenderArticleListArgs{
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

func setCacheControlHeader(w http.ResponseWriter, isProtected bool) {
	if isProtected {
		w.Header().Set("Cache-Control", "private,max-age=300")
	} else {
		w.Header().Set("Cache-Control", "public,max-age=300")
	}
}

func renderChangelog(e *env, w http.ResponseWriter, r *http.Request, parsed *changelog.ParsedChangelog, isWidget bool) error {
	args := RenderChangelogArgs{
		FeedURL:    handler.GetFeedURL(r),
		CurrentURL: handler.GetFullURL(r),
		CL:         parsed.CL,
		Articles:   parsed.Articles,
		HasMore:    parsed.HasMore,
	}

	go e.getAnalyticsEmitter(parsed.CL).Emit(analytics.NewEvent(r, parsed.CL))
	if isWidget {
		return e.render.RenderWidget(r.Context(), w, args)
	}
	return e.render.RenderChangelog(r.Context(), w, args)
}
