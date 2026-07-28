package search

import (
	"math"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name string
		vec1 []float32
		vec2 []float32
		want float32
	}{
		{
			name: "identical vectors",
			vec1: []float32{1, 0, 0},
			vec2: []float32{1, 0, 0},
			want: 1.0,
		},
		{
			name: "orthogonal vectors",
			vec1: []float32{1, 0, 0},
			vec2: []float32{0, 1, 0},
			want: 0.0,
		},
		{
			name: "opposite vectors",
			vec1: []float32{1, 0, 0},
			vec2: []float32{-1, 0, 0},
			want: -1.0,
		},
		{
			name: "parallel vectors",
			vec1: []float32{1, 2, 3},
			vec2: []float32{2, 4, 6},
			want: 1.0,
		},
		{
			name: "normalized vectors with small angle",
			vec1: []float32{1, 0},
			vec2: []float32{0.866, 0.5}, // ~30 degree angle
			want: 0.866,                  // cos(30°) ≈ 0.8660
		},
		{
			name: "zero vectors",
			vec1: []float32{0, 0, 0},
			vec2: []float32{0, 0, 0},
			want: 0.0,
		},
		{
			name: "one zero vector",
			vec1: []float32{1, 2, 3},
			vec2: []float32{0, 0, 0},
			want: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CosineSimilarity(tt.vec1, tt.vec2)
			// Allow small tolerance for floating point errors (use 1e-3 for better precision)
			if math.Abs(float64(got-tt.want)) > 1e-3 {
				t.Errorf("CosineSimilarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCosineSimilarityMismatchedLength(t *testing.T) {
	vec1 := []float32{1, 2, 3}
	vec2 := []float32{1, 2}
	got := CosineSimilarity(vec1, vec2)
	if got != 0 {
		t.Errorf("CosineSimilarity with mismatched lengths should return 0, got %v", got)
	}
}

func TestHybridScore(t *testing.T) {
	tests := []struct {
		name            string
		lexicalRank     float32
		cosineSim       float32
		lexicalWeight   float32
		semanticWeight  float32
		expectMin       float32
		expectMax       float32
	}{
		{
			name:           "perfect lexical and semantic match",
			lexicalRank:    -1,
			cosineSim:      1.0,
			lexicalWeight:  0.3,
			semanticWeight: 0.7,
			expectMin:      0.98,
			expectMax:      1.0,
		},
		{
			name:           "good semantic, poor lexical",
			lexicalRank:    -100,
			cosineSim:      0.9,
			lexicalWeight:  0.3,
			semanticWeight: 0.7,
			expectMin:      0.68,
			expectMax:      0.72,
		},
		{
			name:           "no semantic match",
			lexicalRank:    -1,
			cosineSim:      -1.0,
			lexicalWeight:  0.3,
			semanticWeight: 0.7,
			expectMin:      0.28,
			expectMax:      0.31,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HybridScore(tt.lexicalRank, tt.cosineSim, tt.lexicalWeight, tt.semanticWeight)
			if got < tt.expectMin || got > tt.expectMax {
				t.Errorf("HybridScore() = %v, expected [%v, %v]", got, tt.expectMin, tt.expectMax)
			}
		})
	}
}
