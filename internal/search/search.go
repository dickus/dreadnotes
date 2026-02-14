package search

import (
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

func IndexDocument(idx bleve.Index, doc IndexedDocument) error { return idx.Index(doc.Path, doc) }

func DeleteDocument(idx bleve.Index, path string) error { return idx.Delete(path) }

func Search(idx bleve.Index, queryStr string, limit int) (*bleve.SearchResult, error) {
	queryStr = strings.TrimSpace(queryStr)
	if queryStr == "" {
		return &bleve.SearchResult{}, nil
	}

	prefix := bleve.NewPrefixQuery(queryStr)

	fuzzy := bleve.NewFuzzyQuery(queryStr)
	fuzzy.Fuzziness = 1

	combined := bleve.NewDisjunctionQuery(prefix, fuzzy)

	req := bleve.NewSearchRequest(combined)
	req.Size = limit
	req.Fields = []string{"title", "content", "path"}

	return idx.Search(req)
}

func SearchByTag(idx bleve.Index, tag string, limit int) (*bleve.SearchResult, error) {
	q := query.NewTermQuery(tag)
	q.SetField("tags")

	req := bleve.NewSearchRequestOptions(q, limit, 0, false)
	req.Fields = []string{"title", "path"}

	return idx.Search(req)
}

