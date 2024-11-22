package source

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/internal/xcache"
)

type localSource struct {
	path  string
	cache xcache.Cache
}

func NewLocalSourceFromStore(s store.LocalSource, cache xcache.Cache) Source {
	return &localSource{
		path:  s.Path,
		cache: cache,
	}
}

func NewLocalID(path string) ID {
	return ID(fmt.Sprintf("lc/%s", path))
}

func (s localSource) ID() ID {
	return NewLocalID(s.path)
}

func (s *localSource) Load(ctx context.Context, page internal.Pagination) (LoadResult, error) {
	// sanitize params
	if page.IsDefined() && page.PageSize() < 1 {
		return LoadResult{}, nil
	}

	info, err := os.Stat(s.path)
	if err != nil {
		return LoadResult{}, err
	}

	if info.IsDir() {
		return s.loadDir(s.path, page)
	} else {
		return s.loadFile(s.path)
	}
}

func (s *localSource) loadDir(path string, page internal.Pagination) (LoadResult, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return LoadResult{}, err
	}

	files = filter(files, fileIsMD)
	totalFiles := len(files)
	start, end := calculatePaginationIndices(page, totalFiles)
	if start >= totalFiles {
		return LoadResult{}, nil
	}

	// sort files in descending order by filename
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() >= files[j].Name()
	})

	var wg sync.WaitGroup
	var mutex sync.Mutex
	notes := make([]RawReleaseNote, 0, page.PageSize())

	for _, file := range files[start:end] {
		wg.Add(1)
		go func(file fs.DirEntry) {
			defer wg.Done()
			raw, err := s.openAndCacheFile(filepath.Join(path, file.Name()))
			if err != nil {
				return
			}
			mutex.Lock()
			notes = append(notes, raw)
			mutex.Unlock()
		}(file)
	}
	wg.Wait()

	return LoadResult{
		Raw:     notes,
		HasMore: end < len(files),
	}, nil
}

func (s *localSource) loadFile(path string) (LoadResult, error) {
	raw, err := s.openAndCacheFile(path)
	if err != nil {
		return LoadResult{}, err
	}

	return LoadResult{
		Raw:     []RawReleaseNote{raw},
		HasMore: false,
	}, nil
}

// Loads the file from the disk and checks wether is has been modified by comparing the cached value.
// Updates the cached value if the content has changed.
func (s *localSource) openAndCacheFile(path string) (RawReleaseNote, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return RawReleaseNote{}, err
	}

	cached, found := s.cache.Get(path)
	equal := bytes.Equal(cached, content)

	if !found || !equal {
		s.cache.Set(path, content)
	}

	return RawReleaseNote{
		Content:    bytes.NewReader(content),
		hasChanged: !equal,
	}, nil
}

// Caculates the start and end index for the total files.
// Start is inclusive, end ist exclusive.
func calculatePaginationIndices(page internal.Pagination, totalFiles int) (start, end int) {
	if !page.IsDefined() {
		return 0, totalFiles
	}
	start = page.StartIdx()
	end = min(start+page.PageSize(), totalFiles)
	return
}
