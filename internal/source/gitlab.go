package source

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"sync"

	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/internal/xcache"
	"github.com/naveensrinivasan/httpcache"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type glSource struct {
	client  *gitlab.Client
	baseURL string
	project string
	path    string
	ref     string
}

func NewGLSourceFromStore(cfg config.Config, gl store.GLSource, cache xcache.Cache) (Source, error) {
	tr := http.DefaultTransport
	if cache != nil {
		cachedTransport := httpcache.NewTransport(cache)
		cachedTransport.Transport = tr
		tr = cachedTransport
	}
	client, err := gitlab.NewClient(
		cfg.Gitlab.Token,
		gitlab.WithHTTPClient(&http.Client{
			Transport: tr,
		}))
	if err != nil {
		return nil, err
	}
	return &glSource{
		client:  client,
		baseURL: gl.BaseURL,
		project: gl.Project,
		path:    gl.Path,
		ref:     gl.Ref,
	}, nil

}

func (s *glSource) ID() ID {
	return ID(fmt.Sprintf("gl/%s/%s", s.project, s.path))
}

func (s *glSource) Load(ctx context.Context, page internal.Pagination) (LoadResult, error) {
	return s.load(ctx, page)
}

func (s *glSource) load(ctx context.Context, page internal.Pagination) (LoadResult, error) {
	if page.IsDefined() && page.PageSize() < 1 {
		return LoadResult{}, nil
	}

	file, _, err := s.client.RepositoryFiles.GetFile(
		s.project,
		s.path,
		&gitlab.GetFileOptions{Ref: gitlab.Ptr(s.ref)},
		gitlab.WithContext(ctx),
	)
	if err == nil {
		note, err := s.loadFile(ctx, file.FileName)
		if err != nil {
			return LoadResult{}, err
		}
		return LoadResult{
			Raw: []RawReleaseNote{note},
		}, nil
	}

	dir, _, err := s.client.Repositories.ListTree(
		s.project,
		&gitlab.ListTreeOptions{
			Path: gitlab.Ptr(s.path),
			Ref:  gitlab.Ptr(s.ref),
		},
		gitlab.WithContext(ctx),
	)
	if err == nil {
		return s.loadDir(ctx, dir, page)
	}
	return LoadResult{}, err
}

func (s *glSource) loadFile(ctx context.Context, filename string) (RawReleaseNote, error) {
	content, resp, err := s.client.RepositoryFiles.GetRawFile(
		s.project,
		filename,
		&gitlab.GetRawFileOptions{Ref: gitlab.Ptr(s.ref)},
		gitlab.WithContext(ctx),
	)
	if err != nil {
		return RawReleaseNote{}, err
	}
	if resp.StatusCode >= 400 {
		return RawReleaseNote{}, fmt.Errorf("failed to download file from gitlab, status %d", resp.StatusCode)
	}

	return RawReleaseNote{
		hasChanged: !fromCache(resp.Header),
		Content:    bytes.NewReader(content),
	}, nil
}

func (s *glSource) loadDir(ctx context.Context, files []*gitlab.TreeNode, page internal.Pagination) (LoadResult, error) {
	files = filter(files, gitlabFileIsMD)
	totalFiles := len(files)
	start, end := calculatePaginationIndices(page, totalFiles)
	if start >= totalFiles {
		return LoadResult{}, nil
	}

	// sort files in descending order by filename
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name >= files[j].Name
	})

	notes := make([]RawReleaseNote, end-start)
	var wg sync.WaitGroup

	for i, file := range files[start:end] {
		wg.Add(1)
		go func(index int, path string) {
			defer wg.Done()
			note, err := s.loadFile(ctx, path)
			if err != nil {
				return
			}
			notes[index] = note
		}(i, file.Path)
	}
	wg.Wait()

	return LoadResult{
		Raw:     notes,
		HasMore: end < len(files),
	}, nil
}

// using some functions from github.go
func gitlabFileIsMD(f *gitlab.TreeNode) bool {
	return filepath.Ext(f.Name) == ".md"
}
