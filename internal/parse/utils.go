package parse

import (
	"bytes"
	"io"
)

type FileFormat int

const (
	OG FileFormat = iota
	KeepAChangelog
)

// Detects the file format of r and returns the string read to detect the file format.
// The read string can not be read again from r.
func detectFileFormat(r io.Reader) (FileFormat, string) {
	var buf bytes.Buffer
	_, err := io.CopyN(&buf, r, 3)
	if err != nil {
		return OG, ""
	}
	start := buf.String()
	if start == "---" {
		// if content has frontmatter => it's probably our own file format
		return OG, start
	}
	return KeepAChangelog, start
}

// Sorts ParsedArticles by their published date.
func sortArticleDesc(a ParsedReleaseNote, b ParsedReleaseNote) int {
	if a.Meta.PublishedAt.IsZero() && b.Meta.PublishedAt.IsZero() {
		return 0
	}
	if a.Meta.PublishedAt.IsZero() {
		return -1
	}
	if b.Meta.PublishedAt.IsZero() {
		return 1
	}

	if a.Meta.PublishedAt.After(b.Meta.PublishedAt) {
		return -1
	}

	return 1
}
