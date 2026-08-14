// Package search contains the application search workflow. Storage and model
// details enter through domain ports so CLI and RAG callers share one path.
package search

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/embeddings"
	domainsearch "github.com/sentzunhat/hawp/librarian/go/internal/domain/search"
	searchport "github.com/sentzunhat/hawp/librarian/go/internal/domain/search/index"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/sqlite"
	sqlitesearch "github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/sqlite/search"
)

// Query opens the repository-local index and executes the shared search flow.
// It remains as a compatibility entry point for existing callers.
func Query(repoRoot, query string, limit int) ([]domainsearch.Result, error) {
	dbPath := filepath.Join(repoRoot, ".hawp", "db", "index.sqlite")
	db, err := sqlite.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("index not found at %s; run `hawp search index` first: %w", dbPath, err)
	}
	defer db.Close()
	return QueryWithIndex(context.Background(), sqlitesearch.NewAdapter(db), query, limit)
}

// QueryWithIndex executes search against a domain port. Tests and future index
// implementations can provide a different adapter without changing search.
func QueryWithIndex(ctx context.Context, index searchport.IndexPort, query string, limit int) ([]domainsearch.Result, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("search limit must be positive")
	}
	candidates, err := index.LexicalSearch(ctx, query, limit*3)
	if err != nil {
		return nil, fmt.Errorf("lexical search: %w", err)
	}
	if len(candidates) == 0 {
		return nil, nil
	}

	hasVectors, err := index.HasVectors(ctx)
	if err == nil && hasVectors {
		candidates = HybridRank(ctx, candidates, query, index, limit)
	} else if len(candidates) > limit {
		candidates = candidates[:limit]
	}

	results := make([]domainsearch.Result, len(candidates))
	for i, candidate := range candidates {
		relevance := candidate.Relevance
		if relevance == 0 {
			relevance = 0.95
		}
		results[i] = domainsearch.Result{
			ChunkID:       fmt.Sprint(candidate.ID),
			Content:       candidate.Content,
			Source:        candidate.Source,
			Title:         candidate.Title,
			FolderContext: candidate.FolderContext,
			ChunkIndex:    candidate.ChunkIndex,
			Type:          candidate.Type,
			Category:      candidate.Category,
			WorkUUID:      candidate.WorkUUID,
			Status:        candidate.Status,
			Relevance:     relevance,
			LexicalRank:   candidate.LexicalRank,
			SemanticScore: candidate.SemanticScore,
		}
	}
	return results, nil
}

// HybridRank blends lexical rank with semantic similarity while keeping the
// ranking policy independent from the concrete index implementation.
func HybridRank(ctx context.Context, candidates []searchport.Candidate, query string, index searchport.IndexPort, limit int) []searchport.Candidate {
	if len(candidates) == 0 {
		return candidates
	}
	fallback := func() []searchport.Candidate {
		if len(candidates) > limit {
			return candidates[:limit]
		}
		return candidates
	}

	meta, ok, err := index.EmbeddingMetadata(ctx)
	if err != nil || !ok {
		return fallback()
	}
	embedder, err := embeddings.NewEmbedder(meta.Backend, meta.Model)
	if err != nil {
		return fallback()
	}
	defer embedder.Close()
	queryVectors, err := embedder.EmbedBatch(ctx, []string{query})
	if err != nil || len(queryVectors) == 0 {
		return fallback()
	}

	ids := make([]int64, len(candidates))
	for i := range candidates {
		ids[i] = candidates[i].ID
	}
	vectors, err := index.ChunkVectors(ctx, ids)
	if err != nil || len(vectors) == 0 {
		return fallback()
	}
	for i := range candidates {
		lexicalRank := float32(1.0 / (1.0 + float64(i)/10.0))
		semanticScore := float32(0)
		if vector, exists := vectors[candidates[i].ID]; exists {
			semanticScore = domainsearch.CosineSimilarity(queryVectors[0], vector)
		}
		candidates[i].LexicalRank = lexicalRank
		candidates[i].SemanticScore = semanticScore
		candidates[i].Relevance = (lexicalRank * 0.3) + (semanticScore * 0.7)
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Relevance > candidates[j].Relevance
	})
	return fallback()
}
