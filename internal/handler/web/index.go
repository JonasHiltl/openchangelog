package web

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/parse"
	"github.com/jonashiltl/openchangelog/render"
)

const (
	default_page      = 1
	default_page_size = 10
)

func index(e *env, w http.ResponseWriter, r *http.Request) error {
	subdomain := parseSubdomain(r.Host)
	if subdomain == "" && e.cfg.IsDBMode() {
		return errs.NewError(errs.ErrServiceUnavailable, errors.New("openchangelog is backed by sqlite (multi tenant mode), use the /:subdomain or /:workspace-id/:changelog-id route"))
	}

	var cl store.Changelog
	var err error
	if subdomain != "" {
		cl, err = e.store.GetChangelogBySubdomain(r.Context(), subdomain)
	} else {
		cl, err = e.store.GetChangelog(r.Context(), store.WS_DEFAULT_ID, store.CL_DEFAULT_ID)
	}
	if err != nil {
		return err
	}

	return renderIndex(e, w, r, cl)
}

func tenantIndex(e *env, w http.ResponseWriter, r *http.Request) error {
	if !e.cfg.IsDBMode() {
		return errs.NewError(errs.ErrServiceUnavailable, errors.New("openchangelog is in config mode, use the / route"))
	}
	wID := r.PathValue("workspace")
	parsedWID, err := store.ParseWID(wID)
	if err != nil {
		return err
	}

	cID := r.PathValue("changelog")
	parsedCID, err := store.ParseCID(cID)
	if err != nil {
		return err
	}

	cl, err := e.store.GetChangelog(r.Context(), parsedWID, parsedCID)
	if err != nil {
		return err
	}
	return renderIndex(e, w, r, cl)
}

func parseSubdomain(host string) string {
	// Remove port if present
	host = strings.Split(host, ":")[0]
	log.Println(host)
	parts := strings.Split(host, ".")
	if parts[0] == "www" {
		parts = parts[1:]
	}

	// subdomain exists, e.g. tenant.openchangelog.com
	if len(parts) > 2 {
		return parts[0]
	}
	return ""
}

func renderIndex(
	e *env,
	w http.ResponseWriter,
	r *http.Request,
	cl store.Changelog,
) error {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page-size")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = default_page
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = default_page_size
	}

	var source internal.Source
	if cl.LocalSource.Valid {
		source = internal.NewLocalSourceFromStore(cl.LocalSource.ValueOrZero())
	} else if cl.GHSource.Valid {
		s, err := internal.NewGHSourceFromStore(e.cfg, cl.GHSource.ValueOrZero(), e.cache)
		if err != nil {
			return err
		}
		source = s
	}

	var loadResult internal.LoadResult
	var parseResult parse.ParseResult
	if source != nil {
		loaded, err := source.Load(r.Context(), internal.NewPagination(pageSize, page))
		if err != nil {
			return err
		}

		parsed, err := e.parse.Parse(r.Context(), loaded.Articles)
		if err != nil {
			return err
		}

		loadResult = loaded
		parseResult = parsed
	}

	if htmxHeader := r.Header.Get("HX-Request"); len(htmxHeader) > 0 {
		if len(parseResult.Articles) > 0 {
			return e.render.RenderArticleList(r.Context(), w, render.RenderArticleListArgs{
				Articles: parseResult.Articles,
				HasMore:  loadResult.HasMore,
				NextPage: page + 1,
				PageSize: pageSize,
			})
		} else {
			w.WriteHeader(http.StatusNoContent)
			return nil
		}
	}

	return e.render.RenderIndex(r.Context(), w, render.RenderIndexArgs{
		CL:       cl,
		Articles: parseResult.Articles,
		HasMore:  loadResult.HasMore,
		PageSize: pageSize,
		NextPage: page + 1,
	})
}
