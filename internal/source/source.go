package source

type Article struct {
	Bytes []byte
}

// Represents a source of the Changelog Markdown files.
type Source interface {
	Load() ([]Article, error)
}
