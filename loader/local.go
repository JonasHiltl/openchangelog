package loader

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type LocalFileLoader struct {
	path string
}

func NewLocalFile(path string) Loader {
	return LocalFileLoader{
		path: path,
	}
}

func (s LocalFileLoader) Load(ctx context.Context, page Pagination) (LoadResult, error) {
	// sanitize params
	if page.PageSize() < 1 {
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

func loadDir(path string, params Pagination) (LoadResult, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return LoadResult{}, err
	}

	if params.StartIdx() >= len(files) {
		return LoadResult{
			Articles: []RawArticle{},
			HasMore:  false,
		}, nil
	}

	files = filter(files, func(f fs.DirEntry) bool {
		return filepath.Ext(f.Name()) == ".md"
	})

	// sort files in descending order by filename
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() >= files[j].Name()
	})

	var wg sync.WaitGroup
	results := make([]RawArticle, 0, params.PageSize())
	mutex := &sync.Mutex{}

	for i := params.StartIdx(); i <= params.EndIdx() && i < len(files); i++ {
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
		HasMore:  params.EndIdx()+1 < len(files),
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
