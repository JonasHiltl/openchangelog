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
	return localSource{
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

func (s localSource) Load(ctx context.Context, page internal.Pagination) (LoadResult, error) {
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

	files = filter(files, func(f fs.DirEntry) bool {
		return filepath.Ext(f.Name()) == ".md"
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
		return files[i].Name() >= files[j].Name()
	})

	var wg sync.WaitGroup
	notes := make([]RawReleaseNote, 0, page.PageSize())
	mutex := &sync.Mutex{}

	for i := startIdx; i <= endIdx && i < len(files); i++ {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			raw, err := s.openAndCacheFile(filepath.Join(path, name))
			if err != nil {
				return
			}
			mutex.Lock()
			notes = append(notes, raw)
			mutex.Unlock()
		}(files[i].Name())
	}
	wg.Wait()

	return LoadResult{
		Raw:     notes,
		HasMore: endIdx+1 < len(files),
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
