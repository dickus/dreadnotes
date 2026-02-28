package search

import (
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
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
func Search(idx bleve.Index, queryStr, tagInput string, start, end time.Time, dateField string, limit int) (*bleve.SearchResult, error) {
	var conjuncts []query.Query

	queryStr = strings.TrimSpace(queryStr)
	if queryStr != "" {
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
		conjuncts = append(conjuncts, textQuery)
	}

	tagInput = strings.TrimSpace(tagInput)
	if tagInput != "" {
		tags := strings.Split(tagInput, ",")

		for _, t := range tags {
			t = strings.TrimSpace(t)
			if t == "" {
				continue
			}

			tq := bleve.NewTermQuery(strings.ToLower(t))
			tq.SetField("tags")

			conjuncts = append(conjuncts, tq)
		}
	}

	if !start.IsZero() {
		dateQuery := bleve.NewDateRangeQuery(start, end)
		dateQuery.SetField(dateField)
		conjuncts = append(conjuncts, dateQuery)
	}

	var combined query.Query
	if len(conjuncts) == 0 {
		combined = bleve.NewMatchAllQuery()
	} else {
		combined = bleve.NewConjunctionQuery(conjuncts...)
	}

	req := bleve.NewSearchRequestOptions(combined, limit, 0, false)
	req.Fields = []string{"title", "content", "path", "created", "updated"}

	return idx.Search(req)
}
