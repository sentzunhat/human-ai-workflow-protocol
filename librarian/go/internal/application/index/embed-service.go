// Package index provides embedding services for indexed documents.
package index

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/embeddings"
	"github.com/sentzunhat/hawp/librarian/go/internal/domain/index/store"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/sqlite"
	sqliteindex "github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/sqlite/index"
)

// DefaultEmbeddingModel is the model used when none is specified. Matches
// embeddings.DefaultModel exactly (same model, same casing) — this package
// no longer has its own separate model-pull convention (see the "Unifying
// the embedding pipeline" note below).
const DefaultEmbeddingModel = embeddings.DefaultModel

// EmbedResult summarizes an embedding run.
type EmbedResult struct {
	ChunksEmbedded int
	VectorDims     int
	ElapsedSeconds float64
	ThroughputMs   float64 // milliseconds per chunk
	Model          string
}

func (r EmbedResult) String() string {
	return fmt.Sprintf(`
Embedding complete:
  Chunks embedded:   %d
  Vector dimensions: %d
  Model:             %s
  Time:              %.2f seconds
  Speed:             %.1f ms/chunk
`,
		r.ChunksEmbedded, r.VectorDims, r.Model, r.ElapsedSeconds, r.ThroughputMs,
	)
}

// EmbedService orchestrates embedding of indexed chunks.
//
// # Unifying the embedding pipeline (2026-07-27)
//
// This used to be a second, parallel embedding pipeline: its own hugot
// session management (recreated per batch — a real performance bug, see
// git history), its own Xenova-org model-pull convention
// (`hawp model pull Xenova/bge-base-en-v1.5`), entirely separate from
// `internal/domain/embeddings` (BAAI-org models, auto-downloading, used by
// the RAG/reshape pipeline). Two embedding pipelines doing the same job
// with different model sources was real duplication, not intentional
// design — this package now uses `embeddings.Embedder` directly, so
// `hawp search embed` and the RAG pipeline share one implementation, one
// model catalog, and one auto-download path. No `hawp model pull` step
// needed anymore for `search embed`.
type EmbedService struct {
	dbPath string
	store  store.EmbeddingStore
}

// NewEmbedService creates an embed service that opens a SQLite database at
// dbPath on each Execute call (backward compatible constructor).
func NewEmbedService(dbPath string) *EmbedService {
	return &EmbedService{dbPath: dbPath}
}

// NewEmbedServiceFromStore creates an embed service backed by the given
// EmbeddingStore. InitSchema is the caller's responsibility.
func NewEmbedServiceFromStore(s store.EmbeddingStore) *EmbedService {
	return &EmbedService{store: s}
}

// Execute embeds all chunks with NULL vectors in the index database.
// backend/modelID follow embeddings.NewEmbedder's conventions: "onnx"
// (modelID "bge-base-en-v1.5" or "all-MiniLM-L6-v2", auto-downloaded) or
// "ollama" (modelID is any model available on the local Ollama server,
// e.g. "nomic-embed-text").
//
// Refuses to mix models into one index: if chunks were already embedded
// with a different backend/model, comparing a query vector from a new
// model against those old document vectors would be comparing incompatible
// vector spaces (even when dimensions coincidentally match) — see
// store.EmbeddingMetadata's doc comment.
//
// Metadata-read failures are surfaced as errors rather than silently ignored,
// so the mixed-model guard always runs when the index already has embeddings.
func (s *EmbedService) Execute(ctx context.Context, backend, modelID string) (EmbedResult, error) {
	if backend == "" {
		backend = "onnx"
	}

	if s.store != nil {
		return embed(ctx, backend, modelID, s.store)
	}

	// Backward-compatible path: open SQLite, create adapter, embed, close.
	db, err := sqlite.Open(s.dbPath)
	if err != nil {
		return EmbedResult{Model: modelID}, fmt.Errorf("open index db: %w", err)
	}
	defer db.Close()

	// InitSchema is idempotent (CREATE TABLE IF NOT EXISTS) — call it here so
	// the index_metadata table exists even for a DB indexed before that
	// table existed, without needing a full `search index` re-run.
	if err := db.InitSchema(); err != nil {
		return EmbedResult{Model: modelID}, fmt.Errorf("ensure schema: %w", err)
	}

	return embed(ctx, backend, modelID, sqliteindex.NewAdapter(db))
}

// embed runs the core embedding logic against any EmbeddingStore implementation.
func embed(ctx context.Context, backend, modelID string, es store.EmbeddingStore) (EmbedResult, error) {
	result := EmbedResult{Model: modelID}
	start := time.Now()

	chunks, err := es.GetChunksNeedingEmbedding()
	if err != nil {
		return result, fmt.Errorf("fetch chunks: %w", err)
	}

	if len(chunks) == 0 {
		return result, fmt.Errorf("no chunks to embed")
	}

	// GetEmbeddingMetadata failures are now propagated: a real read error
	// (e.g. schema mismatch, disk error) must not silently skip the
	// mixed-model guard. Previously the guard was wrapped in
	// `if metaErr == nil && ok`, which swallowed real errors. The correct
	// control flow is: propagate real errors, skip the guard only when ok==false
	// (no prior embedding run), allow embedding when ok==true and models match,
	// and reject when ok==true but models differ.
	existing, ok, metaErr := es.GetEmbeddingMetadata()
	if metaErr != nil {
		return result, fmt.Errorf("read embedding metadata: %w", metaErr)
	}
	if ok {
		if existing.Backend != backend || existing.Model != modelID {
			return result, fmt.Errorf(
				"this index was already embedded with %s/%s; embedding the rest with %s/%s would mix incompatible vector spaces in one index. "+
					"Re-run with --backend %s --model %s, or clear the index (`hawp search index`) to start over with a different model",
				existing.Backend, existing.Model, backend, modelID, existing.Backend, existing.Model)
		}
	}

	embedder, err := embeddings.NewEmbedder(backend, modelID)
	if err != nil {
		return result, fmt.Errorf("init %s embedder for %s: %w", backend, modelID, err)
	}
	defer embedder.Close()

	// Embed and store with proper transaction management.
	// maxCharsPerChunk is a coarse pre-truncation safety net, not a token
	// budget — it exists only to bound obviously-oversized text before it
	// reaches the tokenizer; the real length enforcement is the per-chunk
	// fallback below, since dense markdown (code blocks, tables, JSON) can
	// tokenize far more densely than a chars-per-token average assumes.
	const maxCharsPerChunk = 2000
	const batchSize = 8
	const commitFreq = 64

	totalEmbedded := 0
	totalSkipped := 0
	vectorDims := 0

	// Start transaction for entire embedding process
	if err := es.BeginTx(); err != nil {
		return result, fmt.Errorf("start transaction: %w", err)
	}

	for i := 0; i < len(chunks); i += batchSize {
		end := i + batchSize
		if end > len(chunks) {
			end = len(chunks)
		}

		// Prepare text batch
		batchTexts := make([]string, end-i)
		batchIndices := make([]int, end-i)
		for j := 0; j < end-i; j++ {
			text := chunks[i+j].Text
			if len(text) > maxCharsPerChunk {
				text = text[:maxCharsPerChunk]
			}
			batchTexts[j] = text
			batchIndices[j] = i + j
		}

		// Embed this batch. A batch-level failure (e.g. one chunk's tokenized
		// length exceeds the model's max sequence length — confirmed
		// 2026-07-27: a chunk under maxCharsPerChunk still tokenized to 619
		// tokens against BGE-base's 512 limit) falls back to embedding this
		// batch one chunk at a time, so one oversized document doesn't abort
		// embedding for the rest of the corpus. A chunk that still fails on
		// its own is skipped (left unembedded, logged) rather than crashing.
		batchVectors, batchIdx, err := embedBatchResilient(ctx, embedder, batchTexts)
		if err != nil {
			// Every chunk in this batch failed even individually (e.g. a
			// batch made entirely of oversized documents) — skip the whole
			// batch rather than aborting the run; a rare, unlucky batch
			// composition shouldn't cost every chunk after it.
			fmt.Printf("\nWarning: skipping batch %d entirely (%v)\n", i/batchSize+1, err)
			totalSkipped += len(batchTexts)
			continue
		}
		totalSkipped += len(batchTexts) - len(batchVectors)

		// Capture dimensions from first successful vector
		if vectorDims == 0 && len(batchVectors) > 0 {
			vectorDims = len(batchVectors[0])
		}

		// Write vectors and commit periodically
		for k, vec := range batchVectors {
			chunkIdx := batchIndices[batchIdx[k]]
			vectorJSON, _ := json.Marshal(vec)
			if err := es.UpdateChunkEmbedding(chunks[chunkIdx].ID, vectorJSON); err != nil {
				es.Rollback()
				return result, fmt.Errorf("store embedding for chunk %d: %w", chunks[chunkIdx].ID, err)
			}
			totalEmbedded++

			// Commit every N vectors to avoid OOM
			if totalEmbedded%commitFreq == 0 {
				if err := es.Commit(); err != nil {
					es.Rollback()
					return result, fmt.Errorf("commit transaction at %d vectors: %w", totalEmbedded, err)
				}
				if err := es.BeginTx(); err != nil {
					return result, fmt.Errorf("restart transaction: %w", err)
				}
				fmt.Printf("\rStored %d/%d vectors...", totalEmbedded, len(chunks))
			}
		}
	}

	// Final commit
	if err := es.Commit(); err != nil {
		return result, fmt.Errorf("final commit: %w", err)
	}
	fmt.Println() // Newline after progress
	if totalSkipped > 0 {
		fmt.Printf("Skipped %d chunk(s) that failed to embed individually (see warnings above).\n", totalSkipped)
	}

	if totalEmbedded > 0 {
		if metaErr := es.SetEmbeddingMetadata(store.EmbeddingMetadata{
			Backend: backend,
			Model:   modelID,
			Dim:     vectorDims,
		}); metaErr != nil {
			return result, fmt.Errorf("record embedding metadata: %w", metaErr)
		}
	}

	result.ChunksEmbedded = totalEmbedded
	result.VectorDims = vectorDims
	result.ElapsedSeconds = time.Since(start).Seconds()
	result.ThroughputMs = (result.ElapsedSeconds * 1000) / float64(len(chunks))

	return result, nil
}

// embedBatchResilient embeds texts as one batch; on failure it retries one
// text at a time so a single oversized/malformed chunk doesn't abort the
// whole batch. Returns the successfully embedded vectors alongside the
// original index (into texts) each vector corresponds to — chunks that fail
// even alone are skipped and logged, not returned.
func embedBatchResilient(ctx context.Context, embedder embeddings.Embedder, texts []string) ([][]float32, []int, error) {
	vectors, err := embedder.EmbedBatch(ctx, texts)
	if err == nil {
		idx := make([]int, len(vectors))
		for i := range idx {
			idx[i] = i
		}
		return vectors, idx, nil
	}

	var okVectors [][]float32
	var okIdx []int
	for i, text := range texts {
		vecs, singleErr := embedder.EmbedBatch(ctx, []string{text})
		if singleErr != nil || len(vecs) == 0 {
			fmt.Printf("\nWarning: skipping one chunk that failed to embed (%v)\n", singleErr)
			continue
		}
		okVectors = append(okVectors, vecs[0])
		okIdx = append(okIdx, i)
	}
	if len(okVectors) == 0 {
		return nil, nil, fmt.Errorf("every chunk in this batch failed individually, first batch error: %w", err)
	}
	return okVectors, okIdx, nil
}

// CreateDummyVector creates a placeholder vector for testing.
// In real implementation, this will be actual ONNX model output.
func CreateDummyVector(dims int) []byte {
	// Create a JSON array of floats
	vec := make([]float32, dims)
	// Populate with dummy values
	for i := range vec {
		vec[i] = 0.5
	}
	data, _ := json.Marshal(vec)
	return data
}
