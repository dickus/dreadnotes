// Package search provides searching for notes using bleve library.
package search

import (
	"path/filepath"
	"strings"

	"github.com/dickus/dreadnotes/internal/frontmatter"
)

// DocToIndexed converts a parsed markdown document with frontmatter into an IndexedDocument suitable for the search engine.
func DocToIndexed(d frontmatter.Document) IndexedDocument {
	title := d.Meta.Title
	if title == "" {
		// Fallback: use filename without extension as the title
		baseName := filepath.Base(d.Path)
		title = strings.TrimSuffix(baseName, filepath.Ext(baseName))
	}

	return IndexedDocument{
		Title:   title,
		Content: string(d.Content),
		Tags:    d.Meta.Tags,
		Path:    d.Path,
		Created: d.Meta.Created.Time,
		Updated: d.Meta.Updated.Time,
	}
}
