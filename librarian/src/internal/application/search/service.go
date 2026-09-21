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
	"sort"
	"strings"

	embeddings "github.com/sentzunhat/hawp/librarian/src/internal/domain/providers/embeddings"
	domainsearch "github.com/sentzunhat/hawp/librarian/src/internal/domain/search"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
	sqlite "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repositories/index"
)

type EmbedderFactory func(backend, model string) (embeddings.Embedder, error)

type projectResolver interface {
	ResolveSearchIndexPath(repoRoot string) (string, error)
}

type defaultProjectResolver struct{}

func (defaultProjectResolver) ResolveSearchIndexPath(repoRoot string) (string, error) {
	project := filesystem.ResolveHawpProject(repoRoot)
	return project.GetSearchIndexPath(), nil
}

type Service struct {
	resolver projectResolver
	open     func(path string) (domainsearch.Index, error)
	embedder EmbedderFactory
}

type QueryOptions struct {
	Query       string
	Limit       int
	Semantic    bool
	HybridRatio float32
}

type QueryExecution struct {
	Rows       []map[string]interface{}
	HasVectors bool
}

type IndexNotFoundError struct {
	Path string
	Err  error
}

func (e IndexNotFoundError) Error() string {
	return fmt.Sprintf("index not found at %s; run `hawp search index` first: %v", e.Path, e.Err)
}

func (e IndexNotFoundError) Unwrap() error {
	return e.Err
}

// isSchemaError reports whether err is a SQLite "no such table" error, which
// indicates the database file exists but the schema has not been applied (i.e.
// the index was created by os.Create but never populated by the build step).
func isSchemaError(err error) bool {
	return strings.Contains(err.Error(), "no such table")
}

func NewService(resolver projectResolver, open func(path string) (domainsearch.Index, error)) Service {
	return NewServiceWithEmbedder(resolver, open, nil)
}

func NewServiceWithEmbedder(resolver projectResolver, open func(path string) (domainsearch.Index, error), embedder EmbedderFactory) Service {
	if resolver == nil {
		resolver = defaultProjectResolver{}
	}
	if open == nil {
		open = func(path string) (domainsearch.Index, error) {
			return sqlite.Open(path)
		}
	}
	return Service{
		resolver: resolver,
		open:     open,
		embedder: embedder,
	}
}

func DefaultService() Service {
	return NewService(nil, nil)
}

func (s Service) Execute(repoRoot string, opts QueryOptions) (QueryExecution, error) {
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	dbPath, err := s.resolver.ResolveSearchIndexPath(repoRoot)
	if err != nil {
		return QueryExecution{}, fmt.Errorf("resolve search index path: %w", err)
	}
	db, err := s.open(dbPath)
	if err != nil {
		return QueryExecution{}, IndexNotFoundError{Path: dbPath, Err: err}
	}
	defer db.Close()

	hasVectors, err := db.HasVectors()
	if err != nil {
		if isSchemaError(err) {
			return QueryExecution{}, IndexNotFoundError{Path: dbPath, Err: err}
		}
		return QueryExecution{}, fmt.Errorf("check index vectors: %w", err)
	}
	if opts.Semantic {
		if !hasVectors {
			return QueryExecution{HasVectors: false}, nil
		}
		rows := semanticSearch(opts.Query, db, opts.Limit, s.embedder)
		return QueryExecution{
			Rows:       rows,
			HasVectors: hasVectors,
		}, nil
	}

	rows, err := db.QueryChunksLexical(opts.Query, opts.Limit*3)
	if err != nil {
		if isSchemaError(err) {
			return QueryExecution{}, IndexNotFoundError{Path: dbPath, Err: err}
		}
		return QueryExecution{}, fmt.Errorf("lexical search: %w", err)
	}
	if len(rows) == 0 {
		return QueryExecution{
			Rows:       nil,
			HasVectors: hasVectors,
		}, nil
	}
	if hasVectors {
		rows = hybridRank(rows, opts.Query, db, opts.Limit, opts.HybridRatio, s.embedder)
	} else if len(rows) > opts.Limit {
		rows = rows[:opts.Limit]
	}
	return QueryExecution{
		Rows:       rows,
		HasVectors: hasVectors,
	}, nil
}

// Query runs query against the index at repoRoot/.hawp/db/index.sqlite:
// lexical search first, then hybrid re-ranking (using whichever
// backend/model actually embedded the index — see
// sqlite.IndexDB.GetEmbeddingMetadata) when vectors exist. Returns at most
// limit results, already ranked.
func Query(repoRoot, query string, limit int) ([]domainsearch.Result, error) {
	return QueryWithEmbedder(repoRoot, query, limit, nil)
}

func QueryWithEmbedder(repoRoot, query string, limit int, embedder EmbedderFactory) ([]domainsearch.Result, error) {
	execution, err := NewServiceWithEmbedder(nil, nil, embedder).Execute(repoRoot, QueryOptions{
		Query: query,
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}
	return RowsToResults(execution.Rows, execution.HasVectors), nil
}

func RowsToResults(rows []map[string]interface{}, hasVectors bool) []domainsearch.Result {
	results := make([]domainsearch.Result, len(rows))
	for i, r := range rows {
		relevance := float32(0.95) // lexical-only fallback
		if hasVectors {
			relevance = float32(getFloat(r, "_hybrid_score"))
		}
		results[i] = domainsearch.Result{
			ChunkID:       fmt.Sprintf("%d", getInt(r, "id")),
			Content:       getStr(r, "text"),
			Source:        getStr(r, "path"),
			Title:         getStr(r, "folder_role"),
			Relevance:     relevance,
			LexicalRank:   float32(getFloat(r, "_lexical_rank")),
			SemanticScore: float32(getFloat(r, "_semantic_score")),
			Embedding:     []float32{},
			LineStart:     int(getInt(r, "line_start")),
			LineEnd:       int(getInt(r, "line_end")),
			Priority:      i,
		}
	}
	return results
}

// SemanticSearch performs a pure-vector search: embeds the query using the same
// backend/model recorded in the index, then ranks all stored chunk vectors by
// cosine similarity and returns the top-limit rows.  No FTS5 is involved.
// Returns nil (not an error) when vectors are absent, the embedder can't be
// constructed, or the query fails — the caller decides how to handle that.
func SemanticSearch(query string, db domainsearch.Index, limit int) []map[string]interface{} {
	return semanticSearch(query, db, limit, nil)
}

// SemanticSearchWithEmbedder is like SemanticSearch but uses the provided
// EmbedderFactory instead of nil, so actual vector ranking is performed.
func SemanticSearchWithEmbedder(query string, db domainsearch.Index, limit int, newEmbedder EmbedderFactory) []map[string]interface{} {
	return semanticSearch(query, db, limit, newEmbedder)
}

func semanticSearch(query string, db domainsearch.Index, limit int, newEmbedder EmbedderFactory) []map[string]interface{} {
	meta, ok, err := db.GetEmbeddingMetadata()
	if err != nil || !ok {
		return nil
	}

	if newEmbedder == nil {
		return nil
	}
	embedder, err := newEmbedder(meta.Backend, meta.Model)
	if err != nil {
		return nil
	}
	defer embedder.Close()

	queryVectors, err := embedder.EmbedBatch(context.Background(), []string{query})
	if err != nil || len(queryVectors) == 0 {
		return nil
	}
	queryVector := queryVectors[0]

	allVectors, err := db.GetAllChunkVectors()
	if err != nil || len(allVectors) == 0 {
		return nil
	}

	type scored struct {
		id    int64
		score float32
	}
	ranked := make([]scored, 0, len(allVectors))
	for id, vec := range allVectors {
		ranked = append(ranked, scored{id: id, score: domainsearch.CosineSimilarity(queryVector, vec)})
	}
	sort.Slice(ranked, func(i, j int) bool { return ranked[i].score > ranked[j].score })

	if len(ranked) > limit {
		ranked = ranked[:limit]
	}

	chunkIDs := make([]int64, len(ranked))
	for i, s := range ranked {
		chunkIDs[i] = s.id
	}

	rows, err := db.QueryChunksByIDs(chunkIDs)
	if err != nil {
		return nil
	}

	// Attach semantic score from the ranking so callers can display it.
	scoreByID := make(map[int64]float32, len(ranked))
	for _, s := range ranked {
		scoreByID[s.id] = s.score
	}
	for _, row := range rows {
		if id, ok := row["id"]; ok {
			if idInt, ok := id.(int64); ok {
				row["_semantic_score"] = float64(scoreByID[idInt])
			}
		}
	}

	return rows
}

// HybridRank blends lexical rank with semantic similarity: the query is
// embedded with whichever backend/model actually embedded the index
// (sqlite.IndexDB.GetEmbeddingMetadata) — a query vector from a different
// model would be comparing an incompatible vector space, even when
// dimensions happen to match — then blended by lexicalWeight (lexical-rank
// fraction) and (1 - lexicalWeight) (semantic cosine-similarity fraction).
// When lexicalWeight <= 0 the default 30% lexical / 70% semantic blend is
// used. Falls back to lexical order alone if no embedding metadata is
// recorded yet, the embedder can't be constructed, or query embedding fails,
// so hybrid ranking degrades gracefully rather than erroring the whole search.
func HybridRank(lexicalResults []map[string]interface{}, query string, db domainsearch.Index, limit int, lexicalWeight float32) []map[string]interface{} {
	return hybridRank(lexicalResults, query, db, limit, lexicalWeight, nil)
}

// HybridRankWithEmbedder is like HybridRank but uses the provided
// EmbedderFactory so actual semantic re-ranking is performed.
func HybridRankWithEmbedder(lexicalResults []map[string]interface{}, query string, db domainsearch.Index, limit int, lexicalWeight float32, newEmbedder EmbedderFactory) []map[string]interface{} {
	return hybridRank(lexicalResults, query, db, limit, lexicalWeight, newEmbedder)
}

func hybridRank(lexicalResults []map[string]interface{}, query string, db domainsearch.Index, limit int, lexicalWeight float32, newEmbedder EmbedderFactory) []map[string]interface{} {
	if len(lexicalResults) == 0 {
		return lexicalResults
	}
	// Default blend: 30% lexical, 70% semantic.
	lw := float64(lexicalWeight)
	if lw <= 0 {
		lw = 0.3
	}
	sw := 1.0 - lw

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

	if newEmbedder == nil {
		return fallback()
	}
	embedder, err := newEmbedder(meta.Backend, meta.Model)
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

			// Hybrid blend: lw lexical rank (normalized), sw semantic score.
			// Lexical rank is implicit in FTS5 order; use position as proxy.
			lexicalRank := 1.0 / (1.0 + float64(i)/10.0) // Decay by position
			blendedScore := (lexicalRank * lw) + (semanticScore * sw)

			result["_semantic_score"] = semanticScore
			result["_lexical_rank"] = lexicalRank
			result["_hybrid_score"] = blendedScore
		} else {
			// No vector for this chunk, use lexical only
			lexicalRank := 1.0 / (1.0 + float64(i)/10.0)
			result["_semantic_score"] = 0.0
			result["_lexical_rank"] = lexicalRank
			result["_hybrid_score"] = lexicalRank * lw
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
