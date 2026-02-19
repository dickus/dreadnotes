package search

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/dickus/dreadnotes/internal/frontmatter"
)

func ReindexAll(idx bleve.Index, notesPath string) error {
	if strings.HasPrefix(notesPath, "$HOME") {
		homeDir, _ := os.UserHomeDir()
		notesPath = strings.Replace(notesPath, "$HOME", homeDir, 1)
	}

	entries, err := os.ReadDir(notesPath)
	if err != nil {
		return fmt.Errorf("reading notes dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		fullPath := filepath.Join(notesPath, entry.Name())

		path, fm, content := frontmatter.SplitFile(fullPath)
		doc := frontmatter.Parser(path, fm, content)

		indexed := DocToIndexed(doc)
		if err := IndexDocument(idx, indexed); err != nil {
			fmt.Fprintf(os.Stderr, "index error %s: %v\n", fullPath, err)
		}
	}

	return nil
}
