package web

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jonashiltl/openchangelog/components"
	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/analytics"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/handler"
	"github.com/jonashiltl/openchangelog/internal/handler/web/static"
	"github.com/jonashiltl/openchangelog/internal/handler/web/views"
	"github.com/jonashiltl/openchangelog/internal/load"
)

func index(e *env, w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query()
	page, pageSize := handler.ParsePagination(q)
	pagination := internal.NewPagination(pageSize, page)

	loaded, err := e.loader.LoadAndParse(r, pagination)
	if err != nil {
		return err
	}

	_, isWidget := q["widget"]

	if loaded.CL.Protected {
		if isWidget {
			return errs.NewBadRequest(errors.New("can't display protected changelog in widget"))
		}
		err = ensurePasswordProvided(r, loaded.CL.PasswordHash)
		if err != nil {
			slog.InfoContext(r.Context(), "blocked access to changelog", slog.String("changelog", loaded.CL.ID.String()))
			go e.getAnalyticsEmitter(loaded.CL).Emit(analytics.NewAccessDeniedEvent(r, loaded.CL))
			return views.PasswordProtection(views.PasswordProtectionArgs{
				CSS: static.BaseCSS,
				ThemeArgs: components.ThemeArgs{
					ColorScheme: loaded.CL.ColorScheme.ToApiTypes(),
				},
				FooterArgs: components.FooterArgs{
					HidePoweredBy: loaded.CL.HidePoweredBy,
				},
			}).Render(r.Context(), w)
		}
	}

	if _, ok := q["articles"]; ok {
		return handleArticles(e, w, r.Context(), loaded, page, pageSize)
	}

	setCacheControlHeader(w, loaded.CL.Protected)
	return renderChangelog(e, w, r, loaded, isWidget)
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

func handleArticles(
	e *env,
	w http.ResponseWriter,
	ctx context.Context,
	loaded load.LoadedChangelog,
	page, pageSize int,
) error {
	if len(loaded.Notes) > 0 {
		return e.render.RenderArticleList(ctx, w, RenderArticleListArgs{
			WID:      loaded.CL.WorkspaceID,
			CID:      loaded.CL.ID,
			Articles: loaded.Notes,
			HasMore:  loaded.HasMore,
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

func renderChangelog(
	e *env,
	w http.ResponseWriter,
	r *http.Request,
	loaded load.LoadedChangelog,
	isWidget bool,
) error {
	args := RenderChangelogArgs{
		FeedURL:      handler.GetFeedURL(r),
		CurrentURL:   handler.GetFullURL(r),
		CL:           loaded.CL,
		ReleaseNotes: loaded.Notes,
		HasMore:      loaded.HasMore,
	}

	go e.getAnalyticsEmitter(loaded.CL).Emit(analytics.NewEvent(r, loaded.CL))
	if isWidget {
		return e.render.RenderWidget(r.Context(), w, args)
	}
	return e.render.RenderChangelog(r.Context(), w, args)
}
