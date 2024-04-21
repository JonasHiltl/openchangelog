package source

import "context"

type Article struct {
	Bytes []byte
}

// Represents a source of the Changelog Markdown files.
type Source interface {
	Load(ctx context.Context) ([]Article, error)
}
