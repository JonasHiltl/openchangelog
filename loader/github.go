package loader

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
	"github.com/google/go-github/v61/github"
)

const (
	ghAppID = 881880
)

type GithubLoaderOptions struct {
	// The account name of the owner of the repository
	Owner string
	// The repository which holds the markdown files
	Repository string
	// The path to the root of the directory which holds all markdown files
	Path              string
	AppPrivateKey     string
	AppInstallationId int64
	AccessToken       string
}

type GHLoader struct {
	client         *github.Client
	Owner          string
	Repo           string
	Path           string
	InstallationID int64
}

func NewGithub(opts GithubLoaderOptions) (Loader, error) {
	tr := http.DefaultTransport

	if opts.AppPrivateKey != "" && opts.AppInstallationId != 0 {
		// Wrap the shared transport for use with the app ID 1 authenticating with installation ID 99.
		itr, err := ghinstallation.NewKeyFromFile(tr, ghAppID, opts.AppInstallationId, opts.AppPrivateKey)
		if err != nil {
			return nil, err
		}
		tr = itr
	}

	client := github.NewClient(&http.Client{Transport: tr})
	if opts.AccessToken != "" {
		client = client.WithAuthToken(opts.AccessToken)
	}

	return GHLoader{
		client:         client,
		Owner:          opts.Owner,
		Repo:           opts.Repository,
		Path:           opts.Path,
		InstallationID: opts.AppInstallationId,
	}, nil
}

func (s GHLoader) Load(ctx context.Context, params Pagination) (LoadResult, error) {
	// sanitize params
	if params.PageSize() < 1 {
		return LoadResult{}, nil
	}

	file, dir, _, err := s.client.Repositories.GetContents(ctx, s.Owner, s.Repo, s.Path, nil)
	if err != nil {
		return LoadResult{}, err
	}

	if file != nil {
		c, err := file.GetContent()
		if err != nil {
			return LoadResult{}, err
		}
		return LoadResult{
			Articles: []RawArticle{
				{
					Content: io.NopCloser(strings.NewReader(c)),
				},
			},
		}, nil
	}
	return s.loadDir(ctx, dir, params)
}

func (s GHLoader) loadDir(ctx context.Context, files []*github.RepositoryContent, page Pagination) (LoadResult, error) {
	files = filter(files, func(f *github.RepositoryContent) bool {
		return filepath.Ext(f.GetName()) == ".md"
	})

	if page.StartIdx() >= len(files) {
		return LoadResult{
			Articles: []RawArticle{},
			HasMore:  false,
		}, nil
	}

	// sort files in descending order by filename
	sort.Slice(files, func(i, j int) bool {
		return files[i].GetName() >= files[j].GetName()
	})

	var wg sync.WaitGroup
	articles := make([]RawArticle, 0, page.PageSize())
	mutex := &sync.Mutex{}

	for i := page.StartIdx(); i <= page.EndIdx() && i < len(files); i++ {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			read, err := s.loadFile(ctx, name)
			if err != nil {
				return
			}
			mutex.Lock()
			articles = append(articles, RawArticle{
				Content: read,
			})
			mutex.Unlock()
		}(files[i].GetName())
	}
	wg.Wait()

	return LoadResult{
		Articles: articles,
		HasMore:  page.EndIdx()+1 < len(files),
	}, nil
}

func (s GHLoader) loadFile(ctx context.Context, filename string) (io.ReadCloser, error) {
	read, _, err := s.client.Repositories.DownloadContents(ctx, s.Owner, s.Repo, fmt.Sprintf("%s/%s", s.Path, filename), nil)
	if err != nil {
		return nil, err
	}
	return read, nil
}
