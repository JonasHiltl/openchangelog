package source

import (
	"bytes"
	"context"
	"encoding/base64"
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

	token := gl.Token
	if token == "" && cfg.Gitlab != nil {
		token = cfg.Gitlab.Token
	}

	opts := []gitlab.ClientOptionFunc{
		gitlab.WithHTTPClient(&http.Client{
			Transport: tr,
		}),
	}
	if gl.BaseURL != "" {
		opts = append(opts, gitlab.WithBaseURL(gl.BaseURL))
	}

	client, err := gitlab.NewClient(token, opts...)
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

func NewGitLabID(project, path string) ID {
	return ID(fmt.Sprintf("gl/%s/%s", project, path))
}

func (s *glSource) ID() ID {
	return NewGitLabID(s.project, s.path)
}

func (s *glSource) Load(ctx context.Context, page internal.Pagination) (LoadResult, error) {
	if page.IsDefined() && page.PageSize() < 1 {
		return LoadResult{}, nil
	}

	//-> Load as FIle if that works-> return that else load as Folder
	file, resp, err := s.client.RepositoryFiles.GetFile(
		s.project,
		s.path,
		&gitlab.GetFileOptions{Ref: gitlab.Ptr(s.ref)},
		gitlab.WithContext(ctx),
	)
	if err == nil && file != nil {
		return s.singleFileResult(resp, file)
	}

	return s.loadDir(ctx, page)
}

func (s *glSource) singleFileResult(resp *gitlab.Response, file *gitlab.File) (LoadResult, error) {
	decoded, err := base64.StdEncoding.DecodeString(file.Content)
	if err != nil {
		return LoadResult{}, fmt.Errorf("failed to decode file content: %w", err)
	}
	return LoadResult{
		Raw: []RawReleaseNote{
			{
				hasChanged: !fromCache(resp.Header),
				Content:    bytes.NewReader(decoded),
			},
		},
	}, nil
}

func (s *glSource) loadDir(ctx context.Context, page internal.Pagination) (LoadResult, error) {
	nodes, _, err := s.client.Repositories.ListTree(
		s.project,
		&gitlab.ListTreeOptions{
			Path: gitlab.Ptr(s.path),
			Ref:  gitlab.Ptr(s.ref),
		},
		gitlab.WithContext(ctx),
	)
	if err != nil {
		return LoadResult{}, err
	}

	// Only keep markdown files and sort newest first by filename
	mdFiles := filter(nodes, gitlabFileIsMD)
	sort.Slice(mdFiles, func(i, j int) bool {
		return mdFiles[i].Name >= mdFiles[j].Name
	})

	totalFiles := len(mdFiles)
	start, end := calculatePaginationIndices(page, totalFiles)
	if start >= totalFiles {
		return LoadResult{}, nil
	}
	notes := make([]RawReleaseNote, end-start)
	var wg sync.WaitGroup
	for i, file := range mdFiles[start:end] {
		wg.Add(1)
		go func(index int, path string) {
			defer wg.Done()
			note, err := s.fetchFile(ctx, path)
			if err != nil {
				return
			}
			notes[index] = note
		}(i, file.Path)
	}
	wg.Wait()

	return LoadResult{
		Raw:     notes,
		HasMore: end < totalFiles,
	}, nil
}

func (s *glSource) fetchFile(ctx context.Context, path string) (RawReleaseNote, error) {
	content, resp, err := s.client.RepositoryFiles.GetRawFile(
		s.project,
		path,
		&gitlab.GetRawFileOptions{Ref: gitlab.Ptr(s.ref)},
		gitlab.WithContext(ctx),
	)
	if err != nil {
		return RawReleaseNote{}, err
	}
	if resp.StatusCode >= 400 {
		return RawReleaseNote{}, fmt.Errorf("gitlab returned status %d for file %s", resp.StatusCode, path)
	}
	return RawReleaseNote{
		hasChanged: !fromCache(resp.Header),
		Content:    bytes.NewReader(content),
	}, nil
}

// some functions are used from github.go itself,
func gitlabFileIsMD(f *gitlab.TreeNode) bool {
	return filepath.Ext(f.Name) == ".md"
}
