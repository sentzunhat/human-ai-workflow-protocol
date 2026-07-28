// Package search provides vector and lexical search capabilities.
package search

import "math"

// CosineSimilarity computes the cosine similarity between two vectors.
// Both vectors should have the same length. Returns a value between -1 and 1,
// where 1 is identical, 0 is orthogonal, and -1 is opposite.
func CosineSimilarity(vec1, vec2 []float32) float32 {
	if len(vec1) != len(vec2) {
		return 0
	}

	var dotProduct, norm1, norm2 float64

	for i := range vec1 {
		v1 := float64(vec1[i])
		v2 := float64(vec2[i])
		dotProduct += v1 * v2
		norm1 += v1 * v1
		norm2 += v2 * v2
	}

	if norm1 == 0 || norm2 == 0 {
		return 0
	}

	similarity := dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
	return float32(similarity)
}

// HybridScore combines lexical FTS5 rank and cosine similarity into a single score.
// lexicalRank: typically negative (FTS5 scoring); closer to 0 is better.
// cosineSim: [-1, 1] vector similarity; 1 is identical.
// weights: (lexicalWeight, semanticWeight); typical (0.3, 0.7) for semantic-heavy.
func HybridScore(lexicalRank float32, cosineSim float32, lexicalWeight, semanticWeight float32) float32 {
	// FTS5 rank is negative; normalize using exponential decay
	// rank -1 (best) → normalized ~1, rank -100 → normalized ~0
	// Use e^(rank/50) to map negative values to (0, 1]
	var normalizedLexical float32
	if lexicalRank == 0 {
		normalizedLexical = 0
	} else {
		// Exponential mapping: e^(rank/50) maps -1 to ~0.98, -100 to ~0.14
		normalizedLexical = float32(math.Exp(float64(lexicalRank) / 50))
		// Clamp to [0, 1]
		if normalizedLexical < 0 {
			normalizedLexical = 0
		} else if normalizedLexical > 1 {
			normalizedLexical = 1
		}
	}

	normalizedSemantic := (cosineSim + 1) / 2 // Shift from [-1, 1] to [0, 1]

	return (normalizedLexical * lexicalWeight) + (normalizedSemantic * semanticWeight)
}
