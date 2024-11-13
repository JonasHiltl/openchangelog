package load

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/gregjones/httpcache"
	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/handler"
	"github.com/jonashiltl/openchangelog/internal/source"
	"github.com/jonashiltl/openchangelog/internal/store"
)

type LoadedChangelog struct {
	CL    store.Changelog
	Notes source.LoadResult
}

func NewLoader(cfg config.Config, store store.Store, cache httpcache.Cache) *Loader {
	return &Loader{
		cfg:   cfg,
		store: store,
		cache: cache,
	}
}

type Loader struct {
	cfg   config.Config
	store store.Store
	cache httpcache.Cache
}

// Loads the changelog and it's release notes for this http request.
func (l *Loader) LoadChangelog(r *http.Request, page internal.Pagination) (LoadedChangelog, error) {
	wID, cID := GetQueryIDs(r)
	host := r.Host
	if r.Header.Get("X-Forwarded-Host") != "" {
		host = r.Header.Get("X-Forwarded-Host")
	}

	var cl store.Changelog
	var err error

	if l.cfg.IsConfigMode() {
		cl, err = l.store.GetChangelog(r.Context(), "", "")
	} else if wID != "" && cID != "" {
		cl, err = l.fromWorkspace(r.Context(), wID, cID)
	} else {
		cl, err = l.fromHost(r.Context(), host)
	}
	if err != nil {
		return LoadedChangelog{}, err
	}

	return l.LoadReleaseNotes(r.Context(), cl, page)
}

// Loads the release notes for the specified changelog
func (l *Loader) LoadReleaseNotes(ctx context.Context, cl store.Changelog, page internal.Pagination) (LoadedChangelog, error) {
	var err error
	var s source.Source
	if cl.LocalSource.Valid {
		s = source.NewLocalSourceFromStore(cl.LocalSource.ValueOrZero())
	} else if cl.GHSource.Valid {
		s, err = source.NewGHSourceFromStore(l.cfg, cl.GHSource.ValueOrZero(), l.cache)
	}
	if err != nil {
		return LoadedChangelog{}, err
	}

	res := LoadedChangelog{
		CL: cl,
	}
	if s != nil {
		res.Notes, err = s.Load(ctx, page)
		if err != nil {
			return LoadedChangelog{}, err
		}
	}

	return res, nil
}

func (l *Loader) fromWorkspace(ctx context.Context, wID, cID string) (store.Changelog, error) {
	parsedWID, err := store.ParseWID(wID)
	if err != nil {
		return store.Changelog{}, err
	}
	parsedCID, err := store.ParseCID(cID)
	if err != nil {
		return store.Changelog{}, err
	}

	return l.store.GetChangelog(ctx, parsedWID, parsedCID)
}

func (l *Loader) fromHost(ctx context.Context, host string) (store.Changelog, error) {
	subdomain, err1 := store.SubdomainFromHost(host)
	domain, err2 := store.ParseDomain(host)
	if err1 != nil && err2 != nil {
		return store.Changelog{}, errs.NewBadRequest(errors.New("host & subdomain is not a valid url"))
	}

	return l.store.GetChangelogByDomainOrSubdomain(ctx, domain, subdomain)
}

// Parses the workspace and changelog id from the request query params
func GetQueryIDs(r *http.Request) (wID string, cID string) {
	query := r.URL.Query()
	wID = query.Get(handler.WS_ID_QUERY)
	cID = query.Get(handler.CL_ID_QUERY)

	if wID == "" && cID == "" {
		u, err := url.Parse(r.Header.Get("HX-Current-URL"))
		if err == nil {
			query = u.Query()
			return query.Get(handler.WS_ID_QUERY), query.Get(handler.CL_ID_QUERY)
		}
	}
	return wID, cID
}
