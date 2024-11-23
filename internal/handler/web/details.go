package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jonashiltl/openchangelog/components"
	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/handler"
	"github.com/jonashiltl/openchangelog/internal/handler/web/static"
	"github.com/jonashiltl/openchangelog/internal/handler/web/views"
	"github.com/jonashiltl/openchangelog/internal/parse"
)

func details(e *env, w http.ResponseWriter, r *http.Request) error {
	noteID := r.PathValue("nid")
	loaded, err := e.loader.LoadAndParse(r, internal.NoPagination())
	if err != nil {
		return err
	}

	if loaded.CL.Protected {
		err = ensurePasswordProvided(r, loaded.CL.PasswordHash)
		if err != nil {
			slog.InfoContext(
				r.Context(),
				"blocked access to changelog details",
				slog.String("changelog", loaded.CL.ID.String()),
				slog.String("release", noteID),
			)
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

	setCacheControlHeader(r, w, loaded.CL.Protected)

	for i, note := range loaded.Notes {
		if note.Meta.ID == noteID {
			var prev, next parse.ParsedReleaseNote
			if i > 0 {
				prev = loaded.Notes[i-1]
			}
			if i < len(loaded.Notes)-1 {
				next = loaded.Notes[i+1]
			}

			return e.render.RenderDetails(r.Context(), w, RenderDetailsArgs{
				CL:          loaded.CL,
				ReleaseNote: note,
				FeedURL:     handler.GetFeedURL(r),
				Prev:        prev,
				Next:        next,
				HasMetaKey:  requestFromMac(r.Header),
			})
		}
	}

	return errs.NewNotFound(fmt.Errorf("release note %s not found", noteID))
}
