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

	titlePrefix := bleve.NewPrefixQuery(queryStr)
	titlePrefix.SetField("title")

	titleFuzzy := bleve.NewFuzzyQuery(queryStr)
	titleFuzzy.Fuzziness = 1
	titleFuzzy.SetField("title")

	contentPrefix := bleve.NewPrefixQuery(queryStr)
	contentPrefix.SetField("content")

	contentFuzzy := bleve.NewFuzzyQuery(queryStr)
	contentFuzzy.Fuzziness = 1
	contentFuzzy.SetField("content")

	combined := bleve.NewDisjunctionQuery(titlePrefix, titleFuzzy, contentPrefix, contentFuzzy)

	req := bleve.NewSearchRequest(combined)
	req.Size = limit
	req.Fields = []string{"title", "content", "path"}

	return idx.Search(req)
}

func SearchByTag(idx bleve.Index, tag string, limit int) (*bleve.SearchResult, error) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return &bleve.SearchResult{}, nil
	}

	q := query.NewTermQuery(strings.ToLower(tag))
	q.SetField("tags")

	req := bleve.NewSearchRequestOptions(q, limit, 0, false)
	req.Fields = []string{"title", "content", "path"}

	return idx.Search(req)
}

