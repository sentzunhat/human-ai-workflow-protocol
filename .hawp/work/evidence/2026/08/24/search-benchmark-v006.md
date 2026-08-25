---
evidence-for: v0.0.6 search benchmark
date: 2026-08-24
method: hawp search benchmark (10 queries × 2 patterns)
machine: Apple M1 Max, CPU-only
binary: dev build on development branch (v0.0.6 search stack)
---

# Search Benchmark — v0.0.6

Canonical results: `librarian/docs/benchmarks-v006.md`

## Key Numbers

| Pattern | Avg Latency | High-quality |
|---------|-------------|--------------|
| lexical | 0.2ms       | 10 / 10      |
| hybrid  | 60.5ms      | 10 / 10      |

## Corpus Expansion Confirmed

Index grew from 2873 → 3064 chunks (378 → 392 docs) after adding
`librarian/docs/` and `README.md` via `.hawp/config/search.json`.
Quality and latency unchanged — corpus expansion is additive, not disruptive.

## --semantic Removed

`--semantic` flag never had a separate implementation (fell through to
lexical path). Removed from registry, search.md, and help text in 0.0.6.
