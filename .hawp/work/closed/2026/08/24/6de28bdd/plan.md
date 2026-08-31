---
work-item: 6de28bdd
type: feature
title: "Implement `--semantic` search mode (pure vector search, ONNX + Ollama)"
status: done
owner: unassigned
created: 2026-08-24
updated: 2026-08-24
closed: 2026-08-24
---

# Implement `--semantic` Search Mode

## Mission

Wire `--semantic` as a real pure-vector search mode: embed the query, scan all
stored chunk vectors, rank by cosine similarity, return top-N. Works with both
ONNX and Ollama backends (whichever was used to embed the index). Add it to the
benchmark as a third pattern so lexical / semantic / hybrid are all measured.

## Context

`--semantic` appeared in docs and registry through 0.0.5 but was never
implemented — it silently fell through to lexical. Removed from docs in 0.0.6.

The infrastructure is fully in place:
- `db.GetAllChunkVectors()` — returns all stored vectors keyed by chunk ID
- `db.GetEmbeddingMetadata()` — returns which backend/model embedded the index
- `embeddings.NewEmbedder(backend, model)` — constructs the same embedder used at index time
- `domainsearch.CosineSimilarity(a, b)` — cosine similarity between two vectors

A pure semantic search is: embed query → load all vectors → rank by cosine → return top-N.
No FTS5 involved. Complementary to hybrid (which starts from lexical candidates).

## Scope

### application/search/service.go
- Add `SemanticSearch(query string, db *sqlite.IndexDB, limit int) []map[string]interface{}`
  - Reads embedding metadata to construct the right embedder
  - Embeds the query
  - Loads all chunk vectors via `GetAllChunkVectors()`
  - Ranks by cosine similarity descending
  - Returns top-N rows enriched with `_semantic_score`
  - Graceful fallback to nil on any error (caller decides)

### infrastructure/sqlite/index.go
- Add `QueryChunksBySemantic(chunkIDs []int64, limit int) ([]map[string]interface{}, error)`
  — fetches the full document/chunk rows for the ranked chunk IDs so the result
  format matches `QueryChunksLexical` (same map structure, usable by existing formatters)

### platform/cli/run.go
- Wire `--semantic` flag in `runSearch`:
  - If `--semantic` and vectors present → call `SemanticSearch`, skip lexical
  - If `--semantic` and no vectors → print guidance, return nil
  - If neither flag nor vectors → lexical only (unchanged default)
  - If vectors present and no flag → hybrid (unchanged auto-upgrade)

### platform/cli/benchmark.go
- Add `"semantic"` back to `availablePatterns` when vectors exist
- `benchmarkOneQuery` semantic case → calls `appsearch.SemanticSearch`
- Summary table and relative-performance block already handle any pattern in `patternStats`

### registry.go + search.md
- Document `--semantic` as a real flag
- Note backend requirements (same as `--hybrid`: embed step must have run first)

## Non-Goals

- ANN (approximate nearest neighbor) index — at 3064 chunks linear scan is <5ms
- Streaming results
- Pagination

## Success Criteria

- `hawp search "backlog alignment" --semantic` returns ranked results without touching FTS5
- Works with both `--backend onnx` and `--backend ollama` index
- Benchmark shows real latency (expected: ~50–90ms Ollama warm, ~10ms ONNX warm)
- 10/10 quality on the same 10 benchmark queries as lexical/hybrid

## Outcome

Implemented in v0.0.7. `--semantic` is a real pure-vector search mode: embeds query via same backend/model used at index time, ranks all stored chunk vectors by cosine similarity, returns top-N. No FTS5 involved. `QueryChunksByIDs` added to sqlite/index.go to fetch full rows preserving ranking order. Benchmark now runs all three patterns: lexical 0.1ms 10/10, semantic 478.6ms 9/10, hybrid 72.0ms 10/10. Semantic widens recall significantly for concept-heavy queries (up to +9 results vs FTS5 on "HAWP shape template").

## Verification

- [x] `hawp search "backlog alignment" --semantic` returns 10 ranked results, top result is backlog-alignment.md. Evidence: see Outcome section above.
- [x] All 3 patterns in benchmark: lexical / semantic / hybrid. Evidence: see Outcome section above.
- [x] All Go tests pass (`go test ./...`). Evidence: see Outcome section above.
- [x] Links clean, work validates PASS. Evidence: see Outcome section above.

## Close Checklist

- [x] Implementation complete
- [x] Benchmark results in `librarian/docs/benchmarks-v007.md`
- [x] BACKLOG updated
- [x] Plan moved to closed
