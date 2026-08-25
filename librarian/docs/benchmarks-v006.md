# v0.0.6 Search Benchmarks

**Date:** 2026-08-24  
**Binary:** dev build (development branch), v0.0.6 search stack  
**Machine:** Apple M1 Max (CPU-only, no GPU)  
**Index:** 3064 chunks from 392 documents (kit + work + librarian/docs + README.md)  
**Vectors:** nomic-embed-text (768d) via Ollama — 32.9ms/chunk, 100.9s total  
**Method:** `hawp search benchmark` — 10 queries × 2 patterns (lexical, hybrid)

---

## Summary Table

| Pattern | Avg Latency | Min / Max     | High-quality results |
|---------|-------------|---------------|----------------------|
| lexical | 0.2ms       | 0.0 / 1.0ms   | 10 / 10              |
| hybrid  | 60.5ms      | 49.0 / 90.0ms | 10 / 10              |

**Winner:** lexical (equal quality, ~300× faster at this corpus size).

---

## Per-Query Results

| Query (intent)                             | Lexical            | Hybrid             |
|--------------------------------------------|--------------------|--------------------|
| Work tracking policy                       | 0ms, 10 hits [high]| 87ms, 10 hits [high]|
| Context transfer between sessions          | 0ms, 10 hits [high]| 55ms, 10 hits [high]|
| Evidence standards for findings            | 0ms, 10 hits [high]| 56ms, 10 hits [high]|
| Intake process and investigation ordering  | 0ms,  4 hits [high]| 53ms,  4 hits [high]|
| Provider distribution and materialization  | 0ms,  4 hits [high]| 48ms,  4 hits [high]|
| MCP server configuration for AI agents    | 0ms,  3 hits [high]| 55ms,  3 hits [high]|
| Core HAWP protocol shape fields            | 0ms,  1 hits [high]| 58ms,  1 hits [high]|
| Plan file structure and fields             | 0ms, 10 hits [high]| 51ms, 10 hits [high]|
| Kit maintenance and validation commands    | 0ms, 10 hits [high]| 53ms, 10 hits [high]|
| Binary update and install flow             | 0ms, 10 hits [high]| 49ms, 10 hits [high]|

---

## Comparison vs v0.0.5 Baseline

| Metric              | v0.0.5            | v0.0.6            | Change                          |
|---------------------|-------------------|-------------------|---------------------------------|
| Index size (docs)   | 378               | 392               | +14 (librarian/docs + README.md)|
| Index size (chunks) | 2873              | 3064              | +191                            |
| Lexical latency     | 1.4ms             | 0.2ms             | Consistent (run variation)      |
| Hybrid latency      | 56.7ms            | 60.5ms            | Consistent (run variation)      |
| High-quality hits   | 10 / 10           | 10 / 10           | Unchanged                       |

Quality and latency are stable. The +191 chunks from `librarian/docs/` and `README.md` are now searchable — benchmark queries don't cover them directly, but they are indexed and retrievable.

---

## Corpus Expansion Notes

`librarian/docs/` and `README.md` are now indexed via `.hawp/config/search.json`. Content now reachable via search that wasn't before:

- `backends.md` — backend architecture, performance characteristics
- `benchmarks-v003.md`, `benchmarks-v004.md` — historical benchmark results
- `benchmark-plan-v004plus.md` — benchmarking strategy
- `troubleshooting.md` — common issues and solutions
- `context-reshaping.md` — context reshaping usage guide
- `README.md` — project overview and install instructions

---

## Notes

Lexical latency shows as 0.0ms for most queries because `time.Since(start).Milliseconds()` truncates sub-millisecond durations. True latency is <1ms (consistent with v0.0.3 baseline of <1ms).

Hybrid latency (49–90ms) reflects a real Ollama embed call per query (nomic-embed-text, warm). At this corpus size, lexical FTS5 already surfaces keyword-matched results for all 10 queries — hybrid adds semantic re-ranking without improving quality for these queries, which are all strong keyword matches.
