package source

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v62/github"
	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/naveensrinivasan/httpcache"
)

type ghSource struct {
	client         *github.Client
	Owner          string
	Repo           string
	Path           string
	InstallationID int64
}

func NewGHSourceFromStore(cfg config.Config, gh store.GHSource, cache httpcache.Cache) (Source, error) {
	tr := http.DefaultTransport

	if cfg.HasGithubAuth() && cfg.Github.Auth.AppPrivateKey != "" && gh.InstallationID != 0 {
		// Wrap the shared transport for use with the app ID 1 authenticating with installation ID 99.
		itr, err := ghinstallation.NewKeyFromFile(tr, cfg.Github.Auth.AppID, gh.InstallationID, cfg.Github.Auth.AppPrivateKey)
		if err != nil {
			return nil, err
		}
		tr = itr
	}

	if cache != nil {
		cachedTransport := httpcache.NewTransport(cache)
		cachedTransport.Transport = tr
		tr = cachedTransport
	}

	client := github.NewClient(&http.Client{Transport: tr})
	if cfg.HasGithubAuth() && cfg.Github.Auth.AccessToken != "" {
		client = client.WithAuthToken(cfg.Github.Auth.AccessToken)
	}

	return &ghSource{
		client:         client,
		Owner:          gh.Owner,
		Repo:           gh.Repo,
		Path:           gh.Path,
		InstallationID: gh.InstallationID,
	}, nil
}

func (s *ghSource) Load(ctx context.Context, page internal.Pagination) (LoadResult, error) {
	// sanitize params
	if page.IsDefined() && page.PageSize() < 1 {
		return LoadResult{}, nil
	}

	file, dir, resp, err := s.client.Repositories.GetContents(ctx, s.Owner, s.Repo, s.Path, nil)
	if err != nil {
		return LoadResult{}, err
	}
	if file != nil {
		c, err := file.GetContent()
		if err != nil {
			return LoadResult{}, err
		}
		return LoadResult{
			Raw: []RawReleaseNote{
				{
					HasChanged: !fromCache(resp.Header),
					Content:    io.NopCloser(strings.NewReader(c)),
				},
			},
		}, nil
	}
	return s.loadDir(ctx, dir, page)
}

func (s *ghSource) loadDir(ctx context.Context, files []*github.RepositoryContent, page internal.Pagination) (LoadResult, error) {
	files = filter(files, func(f *github.RepositoryContent) bool {
		return filepath.Ext(f.GetName()) == ".md"
	})

	startIdx := page.StartIdx()
	endIdx := page.EndIdx()

	// If pagination is not applied, process all files
	if !page.IsDefined() {
		startIdx = 0
		endIdx = len(files) - 1
	}

	if startIdx >= len(files) {
		return LoadResult{
			Raw:     []RawReleaseNote{},
			HasMore: false,
		}, nil
	}

	// sort files in descending order by filename
	sort.Slice(files, func(i, j int) bool {
		return files[i].GetName() >= files[j].GetName()
	})

	var wg sync.WaitGroup
	notes := make([]RawReleaseNote, 0, page.PageSize())
	mutex := &sync.Mutex{}

	for i := startIdx; i <= endIdx && i < len(files); i++ {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			note, err := s.loadFile(ctx, name)
			if err != nil {
				return
			}
			mutex.Lock()
			notes = append(notes, note)
			mutex.Unlock()
		}(files[i].GetName())
	}
	wg.Wait()

	return LoadResult{
		Raw:     notes,
		HasMore: endIdx+1 < len(files),
	}, nil
}

func (s *ghSource) loadFile(ctx context.Context, filename string) (RawReleaseNote, error) {
	read, resp, err := s.client.Repositories.DownloadContents(ctx, s.Owner, s.Repo, fmt.Sprintf("%s/%s", s.Path, filename), nil)
	if err != nil {
		return RawReleaseNote{}, err
	}
	return RawReleaseNote{
		HasChanged: !fromCache(resp.Header),
		Content:    read,
	}, nil
}

// Returns true if the headers indicate that the response comes from the cache, else returns false.
func fromCache(h http.Header) bool {
	return h.Get(httpcache.XFromCache) != ""
}

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
