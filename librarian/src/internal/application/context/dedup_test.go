package context

import (
	"testing"

	"github.com/sentzunhat/hawp/librarian/src/internal/domain/search"
)

func TestDeduplicateResults(t *testing.T) {
	tests := []struct {
		name      string
		results   []search.Result
		threshold float64
		want      int
	}{
		{
			name: "no duplicates",
			results: []search.Result{
				{Embedding: []float32{1, 0, 0}, Source: "a.md"},
				{Embedding: []float32{0, 1, 0}, Source: "b.md"},
				{Embedding: []float32{0, 0, 1}, Source: "c.md"},
			},
			threshold: 0.95,
			want:      3,
		},
		{
			name: "identical embeddings (duplicates)",
			results: []search.Result{
				{Embedding: []float32{1, 0, 0}, Source: "a.md"},
				{Embedding: []float32{1, 0, 0}, Source: "b.md"}, // duplicate
				{Embedding: []float32{0, 1, 0}, Source: "c.md"},
			},
			threshold: 0.95,
			want:      2,
		},
		{
			name: "high similarity threshold removes near-duplicates",
			results: []search.Result{
				{Embedding: []float32{1, 0, 0}, Source: "a.md"},
				{Embedding: []float32{0.99, 0.01, 0}, Source: "b.md"}, // very similar
				{Embedding: []float32{0, 1, 0}, Source: "c.md"},
			},
			threshold: 0.95,
			want:      2,
		},
		{
			name:      "empty results",
			results:   []search.Result{},
			threshold: 0.95,
			want:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeduplicateResults(tt.results, tt.threshold)
			if len(got) != tt.want {
				t.Errorf("DeduplicateResults() = %d results, want %d", len(got), tt.want)
			}
		})
	}
}

func TestGroupBySource(t *testing.T) {
	results := []search.Result{
		{Source: "README.md", Relevance: 0.9},
		{Source: "guide.md", Relevance: 0.8},
		{Source: "README.md", Relevance: 0.7},
		{Source: "example.md", Relevance: 0.6},
		{Source: "guide.md", Relevance: 0.5},
	}

	groups := GroupBySource(results)

	if len(groups) != 3 {
		t.Errorf("GroupBySource() = %d groups, want 3", len(groups))
	}

	if len(groups["README.md"]) != 2 {
		t.Errorf("README.md group = %d results, want 2", len(groups["README.md"]))
	}

	if len(groups["guide.md"]) != 2 {
		t.Errorf("guide.md group = %d results, want 2", len(groups["guide.md"]))
	}

	if len(groups["example.md"]) != 1 {
		t.Errorf("example.md group = %d results, want 1", len(groups["example.md"]))
	}
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name string
		a    []float32
		b    []float32
		want float64
		tol  float64 // tolerance for floating point comparison
	}{
		{
			name: "identical vectors",
			a:    []float32{1, 0, 0},
			b:    []float32{1, 0, 0},
			want: 1.0,
			tol:  0.0001,
		},
		{
			name: "orthogonal vectors",
			a:    []float32{1, 0, 0},
			b:    []float32{0, 1, 0},
			want: 0.0,
			tol:  0.0001,
		},
		{
			name: "opposite vectors",
			a:    []float32{1, 0, 0},
			b:    []float32{-1, 0, 0},
			want: -1.0,
			tol:  0.0001,
		},
		{
			name: "normalized unit vectors",
			a:    []float32{1, 1},
			b:    []float32{1, 1},
			want: 1.0,
			tol:  0.0001,
		},
		{
			name: "different lengths",
			a:    []float32{1, 0},
			b:    []float32{1, 0, 0},
			want: 0.0,
			tol:  0.0001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cosineSimilarity(tt.a, tt.b)
			diff := got - tt.want
			if diff < 0 {
				diff = -diff
			}
			if diff > tt.tol {
				t.Errorf("cosineSimilarity() = %f, want %f (±%f)", got, tt.want, tt.tol)
			}
		})
	}
}
