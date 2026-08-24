---
evidence-for: v0.0.4 search benchmark
date: 2026-08-24
method: hawp search benchmark (15 queries × 3 patterns)
machine: Apple M1 Max, CPU-only
binary: dev build on feature/v0.0.5 (0.0.4-equivalent search stack)
---

# Search Benchmark — v0.0.4

Canonical results: `librarian/docs/benchmarks-v004.md`

## Key Numbers

| Pattern  | Avg Latency | Queries with hits |
|----------|-------------|-------------------|
| lexical  | 0.5ms       | 7 / 15            |
| semantic | 0.3ms       | 7 / 15            |
| hybrid   | 0.0ms       | 7 / 15            |

## FTS5 Fix Confirmed

All 15 queries completed without crash. Pre-0.0.4, queries containing
dots/hyphens would panic the FTS5 index. `sanitizeFTSQuery()` resolves this.

## Latency vs Baseline

Semantic 0.3ms (was ~100ms in v0.0.3 live-embed mode) — delta is cached vectors.
Hybrid sub-millisecond for same reason. Lexical unchanged at ~0.5ms.
