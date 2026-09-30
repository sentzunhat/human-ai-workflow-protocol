package context

import (
	"math"
	"strings"

	"github.com/sentzunhat/hawp/librarian/src/internal/domain/search"
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

// ContentJaccardDedup removes near-duplicate search results using word-set
// Jaccard similarity on chunk content. Results are assumed to be in rank order
// (index 0 = highest rank). Any chunk whose Jaccard overlap with a
// higher-ranked kept chunk exceeds threshold is dropped.
//
// Jaccard similarity: |A ∩ B| / |A ∪ B|, where A and B are the word sets
// of two chunks. A threshold of 0.70 catches near-duplicate paragraphs from
// the same document section without requiring embeddings.
//
// Returns the kept results and the count of dropped chunks.
func ContentJaccardDedup(results []search.Result, threshold float64) ([]search.Result, int) {
	if len(results) == 0 {
		return results, 0
	}

	kept := make([]search.Result, 0, len(results))
	keptSets := make([]map[string]struct{}, 0, len(results))
	dropped := 0

	for _, r := range results {
		words := wordSet(r.Content)
		isDuplicate := false
		for _, ks := range keptSets {
			if jaccardSim(words, ks) > threshold {
				isDuplicate = true
				break
			}
		}
		if isDuplicate {
			dropped++
		} else {
			kept = append(kept, r)
			keptSets = append(keptSets, words)
		}
	}

	return kept, dropped
}

// wordSet lowercases text and returns the set of alphanumeric tokens longer
// than one character. Single-character tokens (articles, punctuation remnants)
// are excluded to reduce noise.
func wordSet(text string) map[string]struct{} {
	lower := strings.ToLower(text)
	set := make(map[string]struct{})
	for _, word := range strings.FieldsFunc(lower, func(r rune) bool {
		return !('a' <= r && r <= 'z') && !('0' <= r && r <= '9')
	}) {
		if len(word) > 1 {
			set[word] = struct{}{}
		}
	}
	return set
}

// jaccardSim computes |A ∩ B| / |A ∪ B| for two word sets.
func jaccardSim(a, b map[string]struct{}) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 1.0
	}
	intersection := 0
	for w := range a {
		if _, ok := b[w]; ok {
			intersection++
		}
	}
	union := len(a) + len(b) - intersection
	if union == 0 {
		return 0
	}
	return float64(intersection) / float64(union)
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
