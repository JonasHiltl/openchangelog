package source

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/internal/xcache"
	"github.com/naveensrinivasan/httpcache"

	forgejo "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
)

type fjSource struct {
	client  *forgejo.Client
	baseUrl string
	project string
	path    string
	ref     string
}

func NewFJSourceFromStore(cfg config.Config, fj store.FJSource, cache xcache.Cache) (Source, error) {

	tr := http.DefaultTransport

	if cache != nil {
		cachedTransport := httpcache.NewTransport(cache)
		cachedTransport.Transport = tr
		tr = cachedTransport
	}

	token := fj.Token

	if token == "" && cfg.Forgejo != nil {
		token = cfg.Forgejo.Token
	}

	url := fj.BaseURL

	if url != "" {
		url = cfg.Forgejo.BaseURL
	}

	// Default to https. but overridden in the cfg
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = fmt.Sprintf("https://%s", url)
	}

	client, err := forgejo.NewClient(
		url,
		forgejo.SetToken(token),
	)

	if err != nil {
		return nil, err
	}

	return &fjSource{
		client:  client,
		baseUrl: fj.BaseURL,
		project: fj.Project,
		path:    fj.Path,
		ref:     fj.Ref,
	}, nil

}

func NewForgejoID(project, path string) ID {
	return ID(fmt.Sprintf("fj/%s/%s", project, path))
}

func (f *fjSource) ID() ID {
	return NewForgejoID(f.project, f.path)
}

func (f *fjSource) Load(ctx context.Context, page internal.Pagination) (LoadResult, error) {
	if page.IsDefined() && page.PageSize() < 1 {
		return LoadResult{}, nil
	}

	owner := strings.Split(f.project, "/")

	file, resp, err := f.client.GetFile(
		owner[0],
		owner[1],
		f.ref,
		f.path,
		false,
	)

	if err == nil && file != nil {
		return f.singleFileResult(resp, file)
	}

	return f.loadDir(ctx, page)
}

func (f *fjSource) singleFileResult(resp *forgejo.Response, file []byte) (LoadResult, error) {
	return LoadResult{
		Raw: []RawReleaseNote{
			{
				hasChanged: !fromCache(resp.Header),
				Content:    bytes.NewReader(file),
			},
		},
	}, nil
}

func (f *fjSource) loadDir(ctx context.Context, page internal.Pagination) (LoadResult, error) {

	owner := strings.Split(f.project, "/")

	treeResp, _, err := f.client.GetTrees(owner[0], owner[1], f.ref, forgejo.GetTreesOptions{})
	if err != nil {
		return LoadResult{}, fmt.Errorf("failed to get trees: %w", err)
	}

	nodes := treeResp.Entries

	var mdFiles []forgejo.GitEntry
	for _, node := range nodes {
		if node.Type == "blob" && forgejoFileIsMD(node.URL) {
			mdFiles = append(mdFiles, node)
		}
	}

	sort.Slice(mdFiles, func(i, j int) bool {
		return mdFiles[i].Path >= mdFiles[j].Path
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

			note, err := f.fetchFile(ctx, path)
			if err != nil {
				log.Printf("Warning: failed to fetch file %s: %v", path, err)
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

func (f *fjSource) fetchFile(ctx context.Context, path string) (RawReleaseNote, error) {

	owner := strings.Split(f.project, "/")

	content, resp, err := f.client.GetFile(
		owner[0],
		owner[1],
		f.ref,
		f.path,
		false,
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

type ForgejoFile struct {
	URL string
}

func forgejoFileIsMD(path string) bool {
	return filepath.Ext(path) == ".md"
}
