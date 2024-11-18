package search

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/char/html"
	"github.com/blevesearch/bleve/v2/analysis/token/ngram"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/web"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/jonashiltl/openchangelog/internal/config"
)

func createBleve(cfg config.Config) (bleve.Index, error) {
	idxMapping, err := buildIndexMapping()
	if err != nil {
		return nil, fmt.Errorf("failed to build index mapping: %w", err)
	}

	if cfg.Search.Type == config.SearchDisk {
		if cfg.Search.Disk.Path == "" {
			return nil, errors.New("please define 'search.disk.path' as the directory path to store the search index")
		}

		_, err := os.Stat(cfg.Search.Disk.Path)
		if err != nil {
			if os.IsNotExist(err) {
				slog.Info("search index location doesn't exist, creating new index", slog.String("path", cfg.Search.Disk.Path))
				return bleve.New(cfg.Search.Disk.Path, idxMapping)
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

	idx, err := bleve.NewMemOnly(idxMapping)
	if err != nil {
		return nil, fmt.Errorf("failed to create in memory search index: %w", err)
	}
	slog.Info("successfully created in-memory search index")
	return idx, nil
}

func buildIndexMapping() (mapping.IndexMapping, error) {
	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultMapping.DefaultAnalyzer = keyword.Name

	if err := indexMapping.AddCustomAnalyzer("custom-html", map[string]interface{}{
		"type":         custom.Name,
		"tokenizer":    web.Name,
		"char_filters": []interface{}{html.Name},
	}); err != nil {
		return nil, err
	}

	if err := indexMapping.AddCustomTokenFilter("ngram_min_3_max_3",
		map[string]interface{}{
			"min":  3,
			"max":  3,
			"type": ngram.Name,
		},
	); err != nil {
		return nil, err
	}

	if err := indexMapping.AddCustomAnalyzer("custom_ngram",
		map[string]interface{}{
			"type":         custom.Name,
			"char_filters": []interface{}{},
			"tokenizer":    unicode.Name,
			"token_filters": []interface{}{
				`to_lower`,
				`ngram_min_3_max_3`,
			},
		},
	); err != nil {
		return nil, err
	}

	releaseNoteMapping := bleve.NewDocumentMapping()

	releaseNoteMapping.AddFieldMappingsAt("SID", sidFieldMapping())
	releaseNoteMapping.AddFieldMappingsAt("Title", titleFieldMapping())
	releaseNoteMapping.AddFieldMappingsAt("Description", descriptionFieldMapping())
	releaseNoteMapping.AddFieldMappingsAt("PublishedAt", publishedAtFieldMapping())
	releaseNoteMapping.AddFieldMappingsAt("Tags", tagsFieldMapping())
	releaseNoteMapping.AddFieldMappingsAt("Content", contentFieldMapping())

	indexMapping.AddDocumentMapping("note", releaseNoteMapping)
	return indexMapping, nil
}

func sidFieldMapping() *mapping.FieldMapping {
	fm := bleve.NewTextFieldMapping()
	fm.Analyzer = keyword.Name
	fm.Store = false
	fm.IncludeInAll = false
	return fm
}

func titleFieldMapping() *mapping.FieldMapping {
	fm := bleve.NewTextFieldMapping()
	fm.Analyzer = "custom_ngram"
	return fm
}

func descriptionFieldMapping() *mapping.FieldMapping {
	return titleFieldMapping()
}

func publishedAtFieldMapping() *mapping.FieldMapping {
	fm := mapping.NewDateTimeFieldMapping()
	fm.IncludeInAll = false
	return fm
}

func tagsFieldMapping() *mapping.FieldMapping {
	fm := mapping.NewTextFieldMapping()
	fm.Analyzer = keyword.Name
	return fm
}

func contentFieldMapping() *mapping.FieldMapping {
	fm := mapping.NewTextFieldMapping()
	fm.Analyzer = "custom-html"
	return fm
}
