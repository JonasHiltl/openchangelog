package changelog

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/jonashiltl/openchangelog/internal/store"
)

type localSource struct {
	path string
}

func newLocalSourceFromStore(s store.LocalSource) Source {
	return localSource{
		path: s.Path,
	}
}

func (s localSource) Load(ctx context.Context, page Pagination) (LoadResult, error) {
	// sanitize params
	if page.IsDefined() && page.PageSize() < 1 {
		return LoadResult{}, nil
	}

	info, err := os.Stat(s.path)
	if err != nil {
		return LoadResult{}, err
	}

	if info.IsDir() {
		return loadDir(s.path, page)
	} else {
		return loadFile(s.path)
	}
}

func loadDir(path string, page Pagination) (LoadResult, error) {
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
			Articles: []RawArticle{},
			HasMore:  false,
		}, nil
	}

	// sort files in descending order by filename
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() >= files[j].Name()
	})

	var wg sync.WaitGroup
	results := make([]RawArticle, 0, page.PageSize())
	mutex := &sync.Mutex{}

	for i := startIdx; i <= endIdx && i < len(files); i++ {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			read, err := os.Open(filepath.Join(path, name))
			if err != nil {
				read.Close()
				return
			}
			mutex.Lock()
			results = append(results, RawArticle{
				Content: read,
			})
			mutex.Unlock()
		}(files[i].Name())
	}
	wg.Wait()

	return LoadResult{
		Articles: results,
		HasMore:  endIdx+1 < len(files),
	}, nil
}

func loadFile(path string) (LoadResult, error) {
	read, err := os.Open(path)
	if err != nil {
		read.Close()
		return LoadResult{}, err
	}
	return LoadResult{
		Articles: []RawArticle{{Content: read}},
		HasMore:  false,
	}, nil
}
