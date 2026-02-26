package search

import (
	"fmt"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

// IndexedDocument represents a note's search-ready data structure.
type IndexedDocument struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Tags    []string  `json:"tags"`
	Path    string    `json:"path"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func buildMapping() mapping.IndexMapping {
	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Analyzer = "standard"

	keywordFieldMapping := bleve.NewKeywordFieldMapping()

	storedOnlyFieldMapping := bleve.NewTextFieldMapping()
	storedOnlyFieldMapping.Index = false
	storedOnlyFieldMapping.Store = true

	dateFieldMapping := bleve.NewDateTimeFieldMapping()

	docMapping := bleve.NewDocumentMapping()
	docMapping.AddFieldMappingsAt("title", textFieldMapping)
	docMapping.AddFieldMappingsAt("content", textFieldMapping)
	docMapping.AddFieldMappingsAt("tags", keywordFieldMapping)
	docMapping.AddFieldMappingsAt("path", storedOnlyFieldMapping)
	docMapping.AddFieldMappingsAt("created", dateFieldMapping)
	docMapping.AddFieldMappingsAt("updated", dateFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultMapping = docMapping
	indexMapping.DefaultAnalyzer = "standard"

	return indexMapping
}

// BuildIndex initializes an in-memory Bleve index and populates it with documents from the specified notesPath.
func BuildIndex(notesPath string) (bleve.Index, error) {
	m := buildMapping()

	idx, err := bleve.NewMemOnly(m)
	if err != nil {
		return nil, fmt.Errorf("creating in-memory index: %w", err)
	}

	if err := ReindexAll(idx, notesPath); err != nil {
		idx.Close()
		return nil, fmt.Errorf("indexing notes: %w", err)
	}

	return idx, nil
}
