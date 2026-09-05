package search_test

import (
	"errors"
	"testing"

	appsearch "github.com/sentzunhat/hawp/librarian/src/internal/application/search"
	domainsearch "github.com/sentzunhat/hawp/librarian/src/internal/domain/search"
)

// Unused vector operations deliberately panic if lexical retrieval calls them.
type lexicalIndex struct {
	domainsearch.Index
	vectorErr error
	queryErr  error
	closed    bool
	queried   bool
}

func (i *lexicalIndex) Close() error              { i.closed = true; return nil }
func (i *lexicalIndex) HasVectors() (bool, error) { return false, i.vectorErr }
func (i *lexicalIndex) QueryChunksLexical(query string, limit int) ([]map[string]interface{}, error) {
	i.queried = true
	return []map[string]interface{}{{"text": query}}, i.queryErr
}

func TestServiceInjectedRepositoryLifecycle(t *testing.T) {
	failure := errors.New("repository failure")
	for _, stage := range []string{"success", "vectors", "query", "open"} {
		t.Run(stage, func(t *testing.T) {
			index := &lexicalIndex{}
			if stage == "vectors" {
				index.vectorErr = failure
			}
			if stage == "query" {
				index.queryErr = failure
			}
			service := appsearch.NewService(nil, func(path string) (domainsearch.Index, error) {
				if stage == "open" {
					return nil, failure
				}
				return index, nil
			})
			result, err := service.Execute(t.TempDir(), appsearch.QueryOptions{Query: "example", Limit: 1})
			if stage == "success" {
				if err != nil || len(result.Rows) != 1 || result.Rows[0]["text"] != "example" {
					t.Fatalf("unexpected retrieval: %+v, %v", result, err)
				}
			} else if !errors.Is(err, failure) {
				t.Fatalf("lost repository error: %v", err)
			}
			if index.closed != (stage != "open") {
				t.Fatalf("repository close = %v", index.closed)
			}
			if index.queried != (stage == "success" || stage == "query") {
				t.Fatalf("repository queried = %v", index.queried)
			}
		})
	}
}
