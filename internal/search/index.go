package search

import (
	"fmt"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

type IndexedDocument struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	Path    string   `json:"path"`
}

func buildMapping() mapping.IndexMapping {
	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Analyzer = "standard"

	keywordFieldMapping := bleve.NewKeywordFieldMapping()

	storedOnlyFieldMapping := bleve.NewTextFieldMapping()
	storedOnlyFieldMapping.Index = false
	storedOnlyFieldMapping.Store = true

	docMapping := bleve.NewDocumentMapping()
	docMapping.AddFieldMappingsAt("title", textFieldMapping)
	docMapping.AddFieldMappingsAt("content", textFieldMapping)
	docMapping.AddFieldMappingsAt("tags", keywordFieldMapping)
	docMapping.AddFieldMappingsAt("path", storedOnlyFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultMapping = docMapping
	indexMapping.DefaultAnalyzer = "standard"

	return indexMapping
}

func BuildIndex(notesPath string) (bleve.Index, error) {
	mapping := buildMapping()

	idx, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return nil, fmt.Errorf("creating in-memory index: %w", err)
	}

	if err := ReindexAll(idx, notesPath); err != nil {
		idx.Close()
		return nil, fmt.Errorf("indexing notes: %w", err)
	}

	return idx, nil
}

