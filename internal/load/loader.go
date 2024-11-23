package load

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	mint "github.com/btvoidx/mint/context"
	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/events"
	"github.com/jonashiltl/openchangelog/internal/parse"
	"github.com/jonashiltl/openchangelog/internal/source"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/internal/xcache"
	"github.com/jonashiltl/openchangelog/internal/xlog"
)

type LoadedChangelog struct {
	CL      store.Changelog
	Notes   []parse.ParsedReleaseNote
	HasMore bool
}

// Creates a new Loader.
func NewLoader(
	cfg config.Config,
	store store.Store,
	cache xcache.Cache,
	parser parse.Parser,
	e *mint.Emitter,
) *Loader {
	return &Loader{
		cfg:    cfg,
		store:  store,
		cache:  cache,
		parser: parser,
		e:      e,
	}
}

// The loader combines the source and parse package.
// It first loads the raw release notes using the source package and then parses it using the parse package.
type Loader struct {
	cfg    config.Config
	store  store.Store
	cache  xcache.Cache
	parser parse.Parser
	e      *mint.Emitter
}

// Returns the changelog of the request.
func (l *Loader) GetChangelog(r *http.Request) (store.Changelog, error) {
	host := r.Host
	if r.Header.Get("X-Forwarded-Host") != "" {
		host = r.Header.Get("X-Forwarded-Host")
	}

	if l.cfg.IsConfigMode() {
		return l.store.GetChangelog(r.Context(), "", "")
	}
	return l.fromHost(r.Context(), host)
}

// Loads the changelog and parses it's release notes for the specified http request.
func (l *Loader) LoadAndParse(r *http.Request, page internal.Pagination) (LoadedChangelog, error) {
	cl, err := l.GetChangelog(r)
	if err != nil {
		return LoadedChangelog{}, err
	}

	return l.LoadAndParseReleaseNotes(r.Context(), cl, page)
}

// Loads and parses the release notes for the specified changelog.
func (l *Loader) LoadAndParseReleaseNotes(ctx context.Context, cl store.Changelog, page internal.Pagination) (LoadedChangelog, error) {
	s, err := source.NewSourceFromStore(l.cfg, cl, l.cache)
	if err != nil {
		return LoadedChangelog{}, err
	}

	if s != nil {
		loaded, err := s.Load(ctx, page)
		if err != nil {
			return LoadedChangelog{}, err
		}
		// emit event if release notes have changed
		if loaded.HasChanged() {
			err = mint.Emit(l.e, ctx, events.SourceContentChanged{
				CL:     cl,
				Source: s,
			})
			if err != nil {
				slog.Debug("failed to emit source changed event", xlog.ErrAttr(err))
			}
		}
		parsed := l.parser.Parse(ctx, loaded.Raw, page)
		return LoadedChangelog{
			CL:      cl,
			Notes:   parsed.ReleaseNotes,
			HasMore: loaded.HasMore || parsed.HasMore,
		}, nil
	}

	return LoadedChangelog{CL: cl}, nil
}

func (l *Loader) fromHost(ctx context.Context, host string) (store.Changelog, error) {
	subdomain, err1 := store.SubdomainFromHost(host)
	domain, err2 := store.ParseDomain(host)
	if err1 != nil && err2 != nil {
		return store.Changelog{}, errs.NewBadRequest(errors.New("host & subdomain is not a valid url"))
	}

	return l.store.GetChangelogByDomainOrSubdomain(ctx, domain, subdomain)
}
