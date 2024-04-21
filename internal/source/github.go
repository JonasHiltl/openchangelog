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
	"github.com/google/go-github/github"
)

type GitHubSourceOptions struct {
	// The account name of the owner of the repository
	Owner string
	// The repository which holds the markdown files
	Repository string
	// The path to the root of the directory which holds all markdown files
	Path                string
	GHAppPrivateKey     string
	GHAppInstallationId int64
}

type githubSource struct {
	client *github.Client
	owner  string
	repo   string
	path   string
}

func GitHub(opts GitHubSourceOptions) (Source, error) {
	tr := http.DefaultTransport

	if opts.GHAppPrivateKey != "" && opts.GHAppInstallationId != 0 {
		// Wrap the shared transport for use with the app ID 1 authenticating with installation ID 99.
		itr, err := ghinstallation.NewKeyFromFile(tr, 881880, opts.GHAppInstallationId, opts.GHAppPrivateKey)
		if err != nil {
			return nil, err
		}
		tr = itr
	}

	client := github.NewClient(&http.Client{Transport: tr})
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

	dir = filter(dir, func(f *github.RepositoryContent) bool {
		return filepath.Ext(f.GetName()) == ".md"
	})

	if params.StartIdx() >= len(dir) {
		return LoadResult{
			Articles: []Article{},
			HasMore:  false,
		}, nil
	}

	// sort files in descending order by filename
	sort.Slice(dir, func(i, j int) bool {
		return dir[i].GetName() >= dir[j].GetName()
	})

	var wg sync.WaitGroup
	articles := make([]Article, 0, len(dir))
	mutex := &sync.Mutex{}

	for i := params.StartIdx(); i <= params.EndIdx() && i < len(dir); i++ {
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
		}(dir[i].GetName())
	}
	wg.Wait()

	return LoadResult{
		Articles: articles,
		HasMore:  params.EndIdx()+1 < len(dir),
	}, nil
}

func (s *githubSource) downloadFile(ctx context.Context, filename string) ([]byte, error) {
	read, err := s.client.Repositories.DownloadContents(ctx, s.owner, s.repo, fmt.Sprintf("%s/%s", s.path, filename), nil)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(read)
}
