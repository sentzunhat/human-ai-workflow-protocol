# vector embedding via ONNX (fbf12a93 Slice 2)

**UUID:** `77a6879a-72f9-4779-840d-53d6eea58b51`
**Type:** feature
**Reported:** 2026-07-22
**Risk Level:** medium

---

### Input (what was reported)

> Wire real ONNX model embedding via hugot into `hawp search embed`. Load model from `~/.hawp/models/`, embed all chunks, store vectors in SQLite, enable hybrid lexical+semantic search.

---

### Context

Builds on **4b114c0c** (ingest + FTS5 lexical search, now complete).
Depends on:
- hugot already integrated (`hawp embed` command works)
- SQLite chunks table has embedding column (BLOB, nullable)
- 2,445 chunks ready to embed
- EmbedService scaffolding in place

---

### Analysis

**Root cause (or most likely cause):**
_Without real vectors, hybrid search is incomplete. Lexical-only misses semantic meaning; vectors enable retrieval of "closely related" items even if keywords don't match._

**Directly verified:**
_hugot embedding (`hawp embed`) works. SQLite schema ready. Chunks table ready. EmbedService framework in place._

**Inferred (not yet proven):**
_hugot can load ONNX models fast enough (<500ms startup). Embedding 2,445 chunks takes ~20s total (8ms/chunk). Cosine similarity on 384-768 dims is instant (<1ms per query). Hybrid ranking (FTS5 score + cosine score) is straightforward._

**Scope — what else is affected:**
_Model selection (all-MiniLM-L6-v2 vs BGE-base-en-v1.5), vector storage format (JSON vs binary), search ranking algorithm, query time benchmarks._

---

### Recommended Fix

**Implementation steps:**

1. **Embed pipeline** (`internal/application/index/embed-service.go`):
   - Load ONNX model from `~/.hawp/models/{modelID}/` via hugot
   - Fetch all chunks with NULL embeddings
   - For each chunk.text, call model.Embed(text) → []float32
   - Serialize vector as JSON array, store in chunks.embedding

2. **Vector storage** (`internal/infrastructure/sqlite/`):
   - UpdateChunkEmbedding already exists; implement it properly
   - Method: store vector as JSON (human-readable, debuggable)
   - Add: GetChunkVector(id) to retrieve vectors for search

3. **Hybrid search** (`internal/domain/index/`):
   - Implement cosine similarity: `CosineSimilarity(vec1, vec2) float32`
   - Implement search ranking: combine FTS5 rank + cosine score
   - Formula: `hybrid_score = (fts5_rank * 0.3) + (cosine * 0.7)`

4. **CLI wiring** (`internal/platform/cli/run.go`):
   - `runSearchEmbed`: wire hugot model loading
   - Progress bar: show embedding progress (X/2445)
   - `runSearch`: update to use vectors if available, fall back to lexical

5. **Tests**:
   - Unit: cosine similarity edge cases (orthogonal vectors, parallel, identical)
   - E2E: embed real chunks, verify vector storage, test hybrid search
   - Performance: measure embedding speed + search latency

**Model options (test both):**
- `all-MiniLM-L6-v2`: 384 dims, 43 MB, fast, decent quality
- `BGE-base-en-v1.5`: 768 dims, 230 MB, excellent quality, retrieval-optimized

**What to verify after:**

- [ ] Model loads without errors (<500ms)
- [ ] All 2,445 chunks get vectors stored
- [ ] Vectors are JSON-serialized and correct shape
- [ ] Cosine search finds semantically similar items
- [ ] Hybrid query combines lexical + semantic rankings
- [ ] Embedding speed is ~8ms/chunk (20s total for repo)
- [ ] Search latency <1ms with vectors cached
- [ ] Fallback to lexical works if model missing
- [ ] Both models (MiniLM + BGE) work; quality difference visible

---

## Outcome (filled at close)

_Pending._

## Verification (filled at close)

_Pending._
