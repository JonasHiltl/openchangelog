package changelog

import (
	"context"
	"errors"

	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/naveensrinivasan/httpcache"
)

// Loader groups multiple ways of loading a changelog.
type Loader struct {
	cfg   config.Config
	store store.Store
	cache httpcache.Cache
}

// NewLoader creates a new Loader instance.
func NewLoader(cfg config.Config, store store.Store, cache httpcache.Cache) *Loader {
	return &Loader{
		cfg:   cfg,
		store: store,
		cache: cache,
	}
}

// LoadedChangelog represents a loaded changelog with its metadata.
type LoadedChangelog struct {
	cl   store.Changelog
	page Pagination
	res  LoadResult
}

// ParsedChangelog represents a parsed changelog with its articles.
type ParsedChangelog struct {
	CL       store.Changelog
	Articles []ParsedArticle
	HasMore  bool
}

func (l *Loader) FromConfig(ctx context.Context, page Pagination) (LoadedChangelog, error) {
	store := store.NewConfigStore(l.cfg)
	cl, err := store.GetChangelog(ctx, "", "")
	if err != nil {
		return LoadedChangelog{}, err
	}

	res, err := l.load(ctx, cl, page)
	if err != nil {
		return LoadedChangelog{}, err
	}

	return LoadedChangelog{
		cl:   cl,
		res:  res,
		page: page,
	}, nil
}

// Tries to load the corresponding changelog for the host, either by it's subdomain or domain.
func (l *Loader) FromHost(ctx context.Context, host string, page Pagination) (LoadedChangelog, error) {
	subdomain, serr := store.SubdomainFromHost(host)
	domain, derr := store.ParseDomain(host)
	if derr != nil && serr != nil {
		return LoadedChangelog{}, errs.NewBadRequest(errors.New("host is not a valid url"))
	}

	cl, err := l.store.GetChangelogByDomainOrSubdomain(ctx, domain, subdomain)
	if err != nil {
		return LoadedChangelog{}, err
	}

	res, err := l.load(ctx, cl, page)
	if err != nil {
		return LoadedChangelog{}, err
	}

	return LoadedChangelog{
		cl:   cl,
		res:  res,
		page: page,
	}, nil
}

func (l *Loader) FromWorkspace(ctx context.Context, wID, cID string, page Pagination) (LoadedChangelog, error) {
	parsedWID, err := store.ParseWID(wID)
	if err != nil {
		return LoadedChangelog{}, err
	}

	parsedCID, err := store.ParseCID(cID)
	if err != nil {
		return LoadedChangelog{}, err
	}
	cl, err := l.store.GetChangelog(ctx, parsedWID, parsedCID)
	if err != nil {
		return LoadedChangelog{}, err
	}

	res, err := l.load(ctx, cl, page)
	if err != nil {
		return LoadedChangelog{}, err
	}

	return LoadedChangelog{
		cl:   cl,
		res:  res,
		page: page,
	}, nil
}

func (l *Loader) load(ctx context.Context, cl store.Changelog, page Pagination) (LoadResult, error) {
	var source Source
	if cl.LocalSource.Valid {
		source = newLocalSourceFromStore(cl.LocalSource.ValueOrZero())
	} else if cl.GHSource.Valid {
		s, err := newGHSourceFromStore(l.cfg, cl.GHSource.ValueOrZero(), l.cache)
		if err != nil {
			return LoadResult{}, err
		}
		source = s
	}

	if source != nil {
		res, err := source.Load(ctx, page)
		if err != nil {
			return LoadResult{}, err
		}
		return res, nil
	}
	return LoadResult{}, nil
}

var parser = NewParser(createGoldmark())

func (c LoadedChangelog) Parse(ctx context.Context) ParsedChangelog {
	parsed := parser.Parse(ctx, c.res.Articles, c.page)

	return ParsedChangelog{
		CL:       c.cl,
		Articles: parsed.Articles,
		// parsed.HasMore might be true if keep-a-changelog parser finds more releases
		HasMore: c.res.HasMore || parsed.HasMore,
	}
}
