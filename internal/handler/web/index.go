package web

import (
	"net/http"

	"github.com/jonashiltl/openchangelog/internal/changelog"
	"github.com/jonashiltl/openchangelog/internal/handler"
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

	return e.render.RenderIndex(r.Context(), w, render.RenderIndexArgs{
		FeedURL:        handler.ChangelogToFeedURL(r),
		CL:             parsed.CL,
		Articles:       parsed.Articles,
		HasMore:        parsed.HasMore,
		PageSize:       pageSize,
		NextPage:       page + 1,
		BaseCSSVersion: e.baseCSSVersion,
	})
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
