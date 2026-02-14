package search

import "github.com/dickus/dreadnotes/internal/frontmatter"

func DocToIndexed(d frontmatter.Document) IndexedDocument {
	return IndexedDocument{
		Title:   d.Meta.Title,
		Content: string(d.Content),
		Tags:    d.Meta.Tags,
		Path:    d.Path,
	}
}

