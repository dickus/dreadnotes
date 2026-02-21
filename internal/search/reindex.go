package search

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/dickus/dreadnotes/internal/frontmatter"
	"github.com/dickus/dreadnotes/internal/utils"
)

// ReindexAll reads the notes directory, parses markdown files, and adds them to the provided Bleve search index.
func ReindexAll(idx bleve.Index, notesPath string) error {
	resolvedPath := utils.PathParse(notesPath)

	entries, err := os.ReadDir(resolvedPath)
	if err != nil {
		return fmt.Errorf("reading notes dir: %w", err)
	}

	for _, entry := range entries {
		// Skip directories and non-markdown files
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		fullPath := filepath.Join(resolvedPath, entry.Name())

		doc, err := frontmatter.ParseFile(fullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping invalid note %s: %v\n", fullPath, err)

			continue
		}

		indexed := DocToIndexed(doc)
		if err := IndexDocument(idx, indexed); err != nil {
			fmt.Fprintf(os.Stderr, "Index error %s: %v\n", fullPath, err)
		}
	}

	return nil
}
