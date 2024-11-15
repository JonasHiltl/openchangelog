package search

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/blevesearch/bleve"
	"github.com/jonashiltl/openchangelog/internal/config"
)

func createBleve(cfg config.Config) (bleve.Index, error) {
	if cfg.Search.Type == config.SearchDisk {
		if cfg.Search.Disk.Path == "" {
			return nil, errors.New("please define 'search.disk.path' as the directory path to store the search index")
		}

		_, err := os.Stat(cfg.Search.Disk.Path)
		if err != nil {
			if os.IsNotExist(err) {
				slog.Info("search index location doesn't exist, creating new index", slog.String("path", cfg.Search.Disk.Path))
				return bleve.New(cfg.Search.Disk.Path, bleve.NewIndexMapping())
			}
			return nil, err
		}

		idx, err := bleve.Open(cfg.Search.Disk.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to open search index at %s: %w", cfg.Search.Disk.Path, err)
		}
		slog.Info("successfully opened existing search index", slog.String("path", cfg.Search.Disk.Path))
		return idx, err
	}

	idx, err := bleve.NewMemOnly(bleve.NewIndexMapping())
	if err != nil {
		return nil, fmt.Errorf("failed to create in memory search index: %w", err)
	}
	slog.Info("successfully created in-memory search index")
	return idx, nil
}
