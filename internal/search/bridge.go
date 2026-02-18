package search

import (
	"path/filepath"
	"strings"

	"github.com/dickus/dreadnotes/internal/frontmatter"
)

func DocToIndexed(d frontmatter.Document) IndexedDocument {
	title := d.Meta.Title
	if title == "" {
		title = strings.TrimSuffix(filepath.Base(d.Path), ".md")
	}

	return IndexedDocument{
		Title:   title,
		Content: string(d.Content),
		Tags:    d.Meta.Tags,
		Path:    d.Path,
	}
}

