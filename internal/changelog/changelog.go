package changelog

import (
	"context"
	"errors"

	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/naveensrinivasan/httpcache"
)

// Groups multiple ways of loading a changelog. Either from the config, by it's subdomain or workspace.
// After loading the changelog it can easily be parsed.
type Loader struct {
	cfg    config.Config
	store  store.Store
	cache  httpcache.Cache
	parser Parser
}

func NewLoader(cfg config.Config, store store.Store, cache httpcache.Cache) *Loader {
	return &Loader{
		cfg:    cfg,
		store:  store,
		cache:  cache,
		parser: NewParser(),
	}
}

func (l *Loader) FromConfig(ctx context.Context, page Pagination) (*LoadedChangelog, error) {
	store := store.NewConfigStore(l.cfg)
	cl, err := store.GetChangelog(ctx, "", "")
	if err != nil {
		return nil, err
	}

	res, err := l.load(ctx, cl, page)
	if err != nil {
		return nil, err
	}

	return &LoadedChangelog{
		cl:     cl,
		res:    res,
		parser: l.parser,
	}, nil
}

// Tries to load the corresponding changelog for the host, either by it's subdomain or domain.
func (l *Loader) FromHost(ctx context.Context, host string, page Pagination) (*LoadedChangelog, error) {
	subdomain, serr := store.SubdomainFromHost(host)
	domain, derr := store.ParseDomain(host)
	if derr != nil && serr != nil {
		return nil, errs.NewBadRequest(errors.New("host is not a valid url"))
	}

	cl, err := l.store.GetChangelogByDomainOrSubdomain(ctx, domain, subdomain)
	if err != nil {
		return nil, err
	}

	res, err := l.load(ctx, cl, page)
	if err != nil {
		return nil, err
	}

	return &LoadedChangelog{
		cl:     cl,
		res:    res,
		parser: l.parser,
	}, nil
}

func (l *Loader) FromWorkspace(ctx context.Context, wID, cID string, page Pagination) (*LoadedChangelog, error) {
	parsedWID, err := store.ParseWID(wID)
	if err != nil {
		return nil, err
	}

	parsedCID, err := store.ParseCID(cID)
	if err != nil {
		return nil, err
	}
	cl, err := l.store.GetChangelog(ctx, parsedWID, parsedCID)
	if err != nil {
		return nil, err
	}

	res, err := l.load(ctx, cl, page)
	if err != nil {
		return nil, err
	}

	return &LoadedChangelog{
		cl:     cl,
		res:    res,
		parser: l.parser,
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

type LoadedChangelog struct {
	cl     store.Changelog
	res    LoadResult
	parser Parser
}

type ParsedChangelog struct {
	CL       store.Changelog
	Articles []ParsedArticle
	HasMore  bool
}

func (c *LoadedChangelog) Parse(ctx context.Context) (ParsedChangelog, error) {
	parsed, err := c.parser.Parse(ctx, c.res.Articles)
	if err != nil {
		return ParsedChangelog{}, err
	}

	return ParsedChangelog{
		CL:       c.cl,
		Articles: parsed,
		HasMore:  c.res.HasMore,
	}, nil
}
