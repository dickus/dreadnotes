package search

import (
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
)

// IndexDocument adds or updates a document in the search index.
func IndexDocument(idx bleve.Index, doc IndexedDocument) error {
	return idx.Index(doc.Path, doc)
}

// DeleteDocument removes a document from the search index by its path.
func DeleteDocument(idx bleve.Index, path string) error {
	return idx.Delete(path)
}

// Search queries the index for a given string across titles and contents, using a combination of prefix and fuzzy matching.
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
	// Requesting these fields to be returned in the search results
	req.Fields = []string{"title", "content", "path"}

	return idx.Search(req)
}

// SearchByTag finds documents that contain the specified tag exactly.
func SearchByTag(idx bleve.Index, tag string, limit int) (*bleve.SearchResult, error) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return &bleve.SearchResult{}, nil
	}

	q := bleve.NewTermQuery(strings.ToLower(tag))
	q.SetField("tags")

	req := bleve.NewSearchRequestOptions(q, limit, 0, false)
	req.Fields = []string{"title", "content", "path"}

	return idx.Search(req)
}

// SearchByDate filters documents by their creation or modification date.
func SearchByDate(idx bleve.Index, start, end time.Time, field string, limit int) (*bleve.SearchResult, error) {
	q := bleve.NewDateRangeQuery(start, end)
	q.SetField(field)

	req := bleve.NewSearchRequestOptions(q, limit, 0, false)
	req.Fields = []string{"title", "content", "path", "created", "updated"}

	return idx.Search(req)
}

// SearchWithDateFilter combines a text search with a date range filter.
// If queryStr is empty, it returns all notes within the specified date range.
func SearchWithDateFilter(idx bleve.Index, queryStr string, start, end time.Time, dateField string, limit int) (*bleve.SearchResult, error) {
	dateQuery := bleve.NewDateRangeQuery(start, end)
	dateQuery.SetField(dateField)

	queryStr = strings.TrimSpace(queryStr)

	if queryStr == "" {
		req := bleve.NewSearchRequestOptions(dateQuery, limit, 0, false)
		req.Fields = []string{"title", "content", "path", "created", "updated"}

		return idx.Search(req)
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

	textQuery := bleve.NewDisjunctionQuery(titlePrefix, titleFuzzy, contentPrefix, contentFuzzy)

	combined := bleve.NewConjunctionQuery(textQuery, dateQuery)

	req := bleve.NewSearchRequestOptions(combined, limit, 0, false)
	req.Fields = []string{"title", "content", "path", "created", "updated"}

	return idx.Search(req)
}
