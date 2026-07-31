package load

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

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
// Release notes scheduled to be published in the future are excluded from the result.
func (l *Loader) LoadAndParse(r *http.Request, page internal.Pagination) (LoadedChangelog, error) {
	cl, err := l.GetChangelog(r)
	if err != nil {
		return LoadedChangelog{}, err
	}

	loaded, err := l.LoadAndParseReleaseNotes(r.Context(), cl, page)
	if err != nil {
		return LoadedChangelog{}, err
	}

	loaded.Notes = removeUnpublished(loaded.Notes)
	return loaded, nil
}

// Filters out release notes with a publishedAt date in the future,
// so they stay hidden until that date/time arrives.
func removeUnpublished(notes []parse.ParsedReleaseNote) []parse.ParsedReleaseNote {
	now := time.Now()
	published := notes[:0]
	for _, n := range notes {
		if !n.Meta.PublishedAt.After(now) {
			published = append(published, n)
		}
	}
	return published
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
