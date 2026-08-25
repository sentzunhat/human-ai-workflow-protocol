# Search Benchmark Results — v0.0.7

**Date:** 2026-08-24  
**Version:** 0.0.7  
**Corpus:** 392 docs / 3064 chunks (kit + work + librarian/docs + README.md)  
**Embedding backend:** Ollama (nomic-embed-text, 768d, warm)  
**Machine:** macOS, M-series  
**Test:** `hawp search benchmark` — 10 queries × 3 patterns

## Three-way results

| Pattern | Avg Latency | Min | Max | High Quality | Low Quality |
|---------|-------------|-----|-----|-------------|-------------|
| lexical | 0.1ms | 0.0ms | 1.0ms | 10/10 | 0/10 |
| semantic | 478.6ms | 466.0ms | 497.0ms | 9/10 | 1/10 |
| hybrid | 72.0ms | 66.0ms | 75.0ms | 10/10 | 0/10 |

## Per-query detail

| Query | Lexical | Semantic | Hybrid | Semantic results |
|-------|---------|----------|--------|-----------------|
| backlog alignment rules | 1ms / 10 / high | 497ms / 10 / high | 72ms / 10 / high | — |
| status report handoff | 0ms / 10 / high | 470ms / 10 / high | 73ms / 10 / high | — |
| evidence discipline patterns | 0ms / 10 / high | 474ms / 10 / high | 66ms / 10 / high | — |
| intake workflow investigation first | 0ms / 4 / high | 471ms / 10 / high | 75ms / 4 / high | +6 vs lexical |
| provider overlay sync | 0ms / 4 / high | 491ms / 10 / high | 71ms / 4 / high | +6 vs lexical |
| hawp mcp stdio server tools | 0ms / 3 / high | 487ms / 10 / high | 75ms / 3 / high | +7 vs lexical |
| HAWP shape template mission constraints output | 0ms / 1 / high | 476ms / 10 / high | 73ms / 1 / high | +9 vs lexical |
| work item plan file format | 0ms / 10 / high | 475ms / 10 / high | 72ms / 10 / high | — |
| kit normalize validate | 0ms / 10 / high | 479ms / 10 / **low** | 72ms / 10 / high | — |
| hawp update binary install | 0ms / 10 / high | 466ms / 10 / high | 71ms / 10 / high | — |

## Observations

- **Semantic widens recall significantly** for concept-heavy queries: "HAWP shape template" goes from 1 FTS5 hit to 10 semantic results. "MCP stdio server tools" goes from 3 to 10. This is the primary use-case for `--semantic`.
- **Hybrid remains the quality optimum**: 10/10 at 72ms vs semantic's 9/10 at 479ms. Hybrid's lexical pre-filter ensures keyword relevance is preserved.
- **Semantic miss**: "kit normalize validate" returns low quality — the top-1 keyword (`normalize`, `validate`, `kit`) wasn't in the highest-cosine chunk. Hybrid recovers it via FTS5 pre-filter.
- **Lexical wins on raw quality/latency ratio**: 0.1ms + 10/10 for this keyword-representative query set. Semantic's advantage shows when queries are conceptual rather than keyword-aligned.

## Benchmark winner (this query set)

**Lexical** — 10/10 quality, 0.1ms avg. For real use:

- Keyword queries → lexical or hybrid (same quality, hybrid adds re-ranking)
- Conceptual/semantic queries → `--semantic` (wider recall, ~480ms)
- Default (no flag, vectors present) → hybrid (best balance)

## Comparison vs v0.0.6

| Metric | v0.0.6 | v0.0.7 |
|--------|--------|--------|
| Patterns | lexical, hybrid | lexical, semantic, hybrid |
| Semantic | not measured | 478.6ms / 9/10 |
| Hybrid latency | 60.5ms | 72.0ms |
| Corpus | 392 docs / 3064 chunks | 392 docs / 3064 chunks (same) |
