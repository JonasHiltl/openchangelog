package web

import (
	"errors"
	"log"
	"net/http"

	"github.com/jonashiltl/openchangelog/components"
	"github.com/jonashiltl/openchangelog/internal/changelog"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/handler"
	"github.com/jonashiltl/openchangelog/internal/handler/web/views"
)

func index(e *env, w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query()
	page, pageSize := handler.ParsePagination(q)
	l, err := handler.LoadChangelog(e.loader, e.cfg.IsDBMode(), r, changelog.NewPagination(pageSize, page))
	if err != nil {
		return err
	}

	_, widget := q["widget"]

	parsed := l.Parse(r.Context())

	if parsed.CL.Protected {
		if widget {
			return errs.NewBadRequest(errors.New("can't display protected changelog in widget"))
		}
		err = ensurePasswordProvided(r, parsed.CL.PasswordHash)
		if err != nil {
			log.Printf("Blocked access to protected changelog: %s\n", parsed.CL.ID)
			return views.PasswordProtection(views.PasswordProtectionArgs{
				CSS: baseCSS,
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

	args := RenderChangelogArgs{
		FeedURL:    handler.GetFeedURL(r),
		CurrentURL: handler.GetFullURL(r),
		CL:         parsed.CL,
		Articles:   parsed.Articles,
		HasMore:    parsed.HasMore,
		PageSize:   pageSize,
		NextPage:   page + 1,
	}

	if widget {
		return e.render.RenderWidget(r.Context(), w, args)
	}
	return e.render.RenderChangelog(r.Context(), w, args)
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
