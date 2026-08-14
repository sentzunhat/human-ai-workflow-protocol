package cli

import (
	"reflect"
	"testing"

	appcontext "github.com/sentzunhat/hawp/librarian/go/internal/application/context"
	"github.com/sentzunhat/hawp/librarian/go/internal/domain/search"
)

func TestPrepareSearchContextMatchesApplicationPreparation(t *testing.T) {
	results := []search.Result{
		{Source: "guide.md", Title: "Guide", Content: "First result.", Relevance: 0.9, Embedding: []float32{1, 0}},
		{Source: "guide.md", Title: "Guide", Content: "Duplicate result.", Relevance: 0.8, Embedding: []float32{1, 0}},
		{Source: "plan.md", Title: "Plan", Content: "Independent result.", Relevance: 0.7, Embedding: []float32{0, 1}},
	}

	want := appcontext.PrepareContext(results, "context query", 250)
	got := prepareSearchContext(results, "context query", 250)

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("CLI context preparation diverged from application preparation (-want +got):\nwant: %#v\ngot:  %#v", want, got)
	}
}
