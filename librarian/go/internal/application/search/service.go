// Package search is the one shared implementation of "run a query against
// the local index" — lexical (FTS5) plus hybrid semantic re-ranking when
// vectors exist. Both `hawp search <query>` (internal/platform/cli) and the
// RAG pipeline's Retrieve() (internal/application/context) call this same
// code, so there is exactly one search implementation, not two that could
// drift apart.
package search

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/embeddings"
	domainsearch "github.com/sentzunhat/hawp/librarian/go/internal/domain/search"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/sqlite"
)

// Query runs query against the index at repoRoot/.hawp/db/index.sqlite:
// lexical search first, then hybrid re-ranking (using whichever
// backend/model actually embedded the index — see
// sqlite.IndexDB.GetEmbeddingMetadata) when vectors exist. Returns at most
// limit results, already ranked.
func Query(repoRoot, query string, limit int) ([]domainsearch.Result, error) {
	dbPath := filepath.Join(repoRoot, ".hawp", "db", "index.sqlite")
	db, err := sqlite.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("index not found at %s; run `hawp search index` first: %w", dbPath, err)
	}
	defer db.Close()

	// Over-fetch for re-ranking headroom, same as the CLI always has.
	rows, err := db.QueryChunksLexical(query, limit*3)
	if err != nil {
		return nil, fmt.Errorf("lexical search: %w", err)
	}
	if len(rows) == 0 {
		return nil, nil
	}

	hasVectors, _ := db.HasVectors()
	if hasVectors {
		rows = HybridRank(rows, query, db, limit)
	} else if len(rows) > limit {
		rows = rows[:limit]
	}

	results := make([]domainsearch.Result, len(rows))
	for i, r := range rows {
		relevance := float32(0.95) // no hybrid score available (lexical-only fallback)
		if hasVectors {
			relevance = float32(getFloat(r, "_hybrid_score"))
		}
		results[i] = domainsearch.Result{
			Source:    getStr(r, "path"),
			Title:     getStr(r, "folder_role"),
			Content:   getStr(r, "text"),
			Relevance: relevance,
			Embedding: []float32{},
		}
	}
	return results, nil
}

// HybridRank blends lexical rank with semantic similarity: the query is
// embedded with whichever backend/model actually embedded the index
// (sqlite.IndexDB.GetEmbeddingMetadata) — a query vector from a different
// model would be comparing an incompatible vector space, even when
// dimensions happen to match — then blended 30% lexical-rank / 70%
// semantic-cosine-similarity. Falls back to lexical order alone if no
// embedding metadata is recorded yet, the embedder can't be constructed, or
// query embedding fails, so hybrid ranking degrades gracefully rather than
// erroring the whole search.
func HybridRank(lexicalResults []map[string]interface{}, query string, db *sqlite.IndexDB, limit int) []map[string]interface{} {
	if len(lexicalResults) == 0 {
		return lexicalResults
	}
	fallback := func() []map[string]interface{} {
		if len(lexicalResults) > limit {
			return lexicalResults[:limit]
		}
		return lexicalResults
	}

	meta, ok, err := db.GetEmbeddingMetadata()
	if err != nil || !ok {
		return fallback() // no vectors embedded yet, or metadata unreadable
	}

	embedder, err := embeddings.NewEmbedder(meta.Backend, meta.Model)
	if err != nil {
		return fallback()
	}
	defer embedder.Close()

	queryVectors, err := embedder.EmbedBatch(context.Background(), []string{query})
	if err != nil || len(queryVectors) == 0 {
		return fallback()
	}
	queryVector := queryVectors[0]

	// Extract chunk IDs from lexical results
	chunkIDs := make([]int64, len(lexicalResults))
	for i, result := range lexicalResults {
		chunkIDs[i] = int64(getInt(result, "id"))
	}

	// Get vectors for all candidates
	vectors, err := db.GetChunkVectors(chunkIDs)
	if err != nil || len(vectors) == 0 {
		return fallback()
	}

	// Compute hybrid scores (lexical + semantic)
	for i, result := range lexicalResults {
		chunkID := int64(getInt(result, "id"))

		if chunkVec, ok := vectors[chunkID]; ok {
			// Cosine similarity (0 to 1)
			semanticScore := float64(domainsearch.CosineSimilarity(queryVector, chunkVec))

			// Hybrid blend: 30% lexical rank (normalized), 70% semantic score.
			// Lexical rank is implicit in FTS5 order; use position as proxy.
			lexicalRank := 1.0 / (1.0 + float64(i)/10.0) // Decay by position
			blendedScore := (lexicalRank * 0.3) + (semanticScore * 0.7)

			result["_semantic_score"] = semanticScore
			result["_lexical_rank"] = lexicalRank
			result["_hybrid_score"] = blendedScore
		} else {
			// No vector for this chunk, use lexical only
			lexicalRank := 1.0 / (1.0 + float64(i)/10.0)
			result["_semantic_score"] = 0.0
			result["_lexical_rank"] = lexicalRank
			result["_hybrid_score"] = lexicalRank * 0.3
		}
	}

	// Sort by hybrid score (descending)
	sort.Slice(lexicalResults, func(i, j int) bool {
		scoreI := getFloat(lexicalResults[i], "_hybrid_score")
		scoreJ := getFloat(lexicalResults[j], "_hybrid_score")
		return scoreI > scoreJ
	})

	if len(lexicalResults) > limit {
		return lexicalResults[:limit]
	}
	return lexicalResults
}

func getStr(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(*string); ok && s != nil {
			return *s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		if i, ok := v.(int64); ok {
			return i
		}
	}
	return 0
}

func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0
}
