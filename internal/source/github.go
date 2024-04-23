package source

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"sort"
	"sync"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v61/github"
)

type GithubSourceOptions struct {
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

type githubSource struct {
	client *github.Client
	owner  string
	repo   string
	path   string
}

func Github(opts GithubSourceOptions) (Source, error) {
	tr := http.DefaultTransport

	if opts.AppPrivateKey != "" && opts.AppInstallationId != 0 {
		// Wrap the shared transport for use with the app ID 1 authenticating with installation ID 99.
		itr, err := ghinstallation.NewKeyFromFile(tr, 881880, opts.AppInstallationId, opts.AppPrivateKey)
		if err != nil {
			return nil, err
		}
		tr = itr
	}

	client := github.NewClient(&http.Client{Transport: tr})
	if opts.AccessToken != "" {
		client.WithAuthToken(opts.AccessToken)
	}

	return &githubSource{
		client: client,
		owner:  opts.Owner,
		repo:   opts.Repository,
		path:   opts.Path,
	}, nil
}

func (s *githubSource) Load(ctx context.Context, params LoadParams) (LoadResult, error) {
	// sanitize params
	if params.PageSize() < 1 {
		return LoadResult{}, nil
	}

	file, dir, _, err := s.client.Repositories.GetContents(ctx, s.owner, s.repo, s.path, nil)
	if err != nil {
		return LoadResult{}, err
	}

	if file != nil {
		c, err := file.GetContent()
		if err != nil {
			return LoadResult{}, err
		}
		return LoadResult{
			Articles: []Article{
				{
					Bytes: []byte(c),
				},
			},
		}, nil
	}
	return s.loadFiles(ctx, dir, params)
}

func (s *githubSource) loadFiles(ctx context.Context, files []*github.RepositoryContent, params LoadParams) (LoadResult, error) {
	files = filter(files, func(f *github.RepositoryContent) bool {
		return filepath.Ext(f.GetName()) == ".md"
	})

	if params.StartIdx() >= len(files) {
		return LoadResult{
			Articles: []Article{},
			HasMore:  false,
		}, nil
	}

	// sort files in descending order by filename
	sort.Slice(files, func(i, j int) bool {
		return files[i].GetName() >= files[j].GetName()
	})

	var wg sync.WaitGroup
	articles := make([]Article, 0, params.PageSize())
	mutex := &sync.Mutex{}

	for i := params.StartIdx(); i <= params.EndIdx() && i < len(files); i++ {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			b, err := s.downloadFile(ctx, name)
			if err != nil {
				return
			}
			mutex.Lock()
			articles = append(articles, Article{
				Bytes: b,
			})
			mutex.Unlock()
		}(files[i].GetName())
	}
	wg.Wait()

	return LoadResult{
		Articles: articles,
		HasMore:  params.EndIdx()+1 < len(files),
	}, nil
}

func (s *githubSource) downloadFile(ctx context.Context, filename string) ([]byte, error) {
	read, _, err := s.client.Repositories.DownloadContents(ctx, s.owner, s.repo, fmt.Sprintf("%s/%s", s.path, filename), nil)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(read)
}
