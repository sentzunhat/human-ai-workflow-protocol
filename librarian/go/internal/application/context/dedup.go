package context

import (
	"math"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/search"
)

// DeduplicateResults removes near-duplicate search results based on embedding similarity.
// Results with cosine similarity > threshold are considered duplicates.
func DeduplicateResults(results []search.Result, threshold float64) []search.Result {
	if len(results) == 0 {
		return results
	}

	// Track indices of results to keep (not duplicates)
	kept := make([]bool, len(results))
	kept[0] = true // Always keep the first result

	// Compare each result against all previous kept results
	for i := 1; i < len(results); i++ {
		isDuplicate := false

		for j := 0; j < i; j++ {
			if !kept[j] {
				continue // Skip already-discarded results
			}

			// Calculate cosine similarity between embeddings
			sim := cosineSimilarity(results[i].Embedding, results[j].Embedding)
			if sim > threshold {
				isDuplicate = true
				break
			}
		}

		kept[i] = !isDuplicate
	}

	// Collect kept results
	deduplicated := make([]search.Result, 0)
	for i, k := range kept {
		if k {
			deduplicated = append(deduplicated, results[i])
		}
	}

	return deduplicated
}

// GroupBySource organizes results by their source document/file.
func GroupBySource(results []search.Result) map[string][]search.Result {
	groups := make(map[string][]search.Result)

	for _, result := range results {
		// Use the document source as the group key
		source := result.Source
		groups[source] = append(groups[source], result)
	}

	return groups
}

// cosineSimilarity calculates the cosine similarity between two vectors.
func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float64

	for i := range a {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
