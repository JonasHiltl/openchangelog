package source

import (
	"os"
	"path/filepath"
	"sync"
)

type localFileSource struct {
	path string
}

func LocalFileSource(path string) Source {
	return &localFileSource{
		path: path,
	}
}

func (s *localFileSource) Load() ([]Article, error) {
	info, err := os.Stat(s.path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return loadFromDir(s.path)
	}
	b, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	return []Article{{Bytes: b}}, nil
}

func loadFromDir(path string) ([]Article, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	results := make([]Article, 0)
	mutex := &sync.Mutex{}

	// load all files to memory in parallel
	for _, file := range files {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			if filepath.Ext(name) == ".md" {
				b, err := os.ReadFile(filepath.Join(path, name))
				if err != nil {
					return
				}
				mutex.Lock()
				results = append(results, Article{
					Bytes: b,
				})
				mutex.Unlock()
			}
		}(file.Name())
	}
	wg.Wait()

	return results, nil
}
