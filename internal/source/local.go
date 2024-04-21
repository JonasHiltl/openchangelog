package source

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type localFileSource struct {
	path string
}

func LocalFile(path string) Source {
	return &localFileSource{
		path: path,
	}
}

func (s *localFileSource) Load(ctx context.Context, params LoadParams) (LoadResult, error) {
	// sanitize params
	if params.PageSize() < 1 {
		return LoadResult{}, nil
	}

	info, err := os.Stat(s.path)
	if err != nil {
		return LoadResult{}, err
	}

	if info.IsDir() {
		return loadFromDir(s.path, params)
	} else {
		return loadFromFile(s.path)
	}
}

func loadFromDir(path string, params LoadParams) (LoadResult, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return LoadResult{}, err
	}

	if params.StartIdx() >= len(files) {
		return LoadResult{
			Articles: []Article{},
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
	results := make([]Article, 0, params.PageSize())
	mutex := &sync.Mutex{}

	for i := params.StartIdx(); i <= params.EndIdx() && i < len(files); i++ {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			b, err := os.ReadFile(filepath.Join(path, name))
			if err != nil {
				return
			}
			mutex.Lock()
			results = append(results, Article{
				Bytes: b,
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

func loadFromFile(path string) (LoadResult, error) {
	b, err := os.ReadFile(path)
	return LoadResult{
		Articles: []Article{{Bytes: b}},
		HasMore:  false,
	}, err
}
