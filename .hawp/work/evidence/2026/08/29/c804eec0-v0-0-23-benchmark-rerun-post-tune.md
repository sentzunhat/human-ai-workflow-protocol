---
work-item: c804eec0
date: 2026-08-29
type: benchmark
title: "v0.0.23 branch-local benchmark rerun after shaping tune"
status: measured
---

# v0.0.23 Benchmark Rerun — Post-Tune

## Scope

Benchmark rerun against a fresh source-built branch binary after reducing
context-format overhead in `librarian/src/internal/application/context/format.go`.

## Directly Verified

- Focused formatter tests passed:

```text
go test ./internal/application/context -run 'TestFormatAsMarkdown|TestFormatAsMarkdownRespectsBudget|TestFormatAsMarkdownEmpty|TestContextBlockString|TestContextBlockReferences|TestContextBlockStringInterleavesReferences'
```

- Source build used for the rerun:

```text
go build -o /tmp/hawp-v0023 ./cmd/hawp
```

- Benchmark commands run against the source-built branch binary:

```text
/tmp/hawp-v0023 search benchmark
/tmp/hawp-v0023 search benchmark --token
```

## Embedded Chunk Inspection

- `select count(*) from chunks where embedding is not null` -> `35`
- `select count(*) from chunks` -> `3382`
- `pragma table_info(index_metadata); select * from index_metadata;` shows:
  - `embedding_backend`: `ollama`
  - `embedding_model`: `nomic-embed-text`
  - `embedding_dim`: `768`
  - `updated_at`: `2026-08-25 16:11:23`
- Embedded chunks are concentrated in only three documents:
  - `librarian/docs/troubleshooting.md` -> `13`
  - `.hawp/work/active/t5l7m9n1/plan.md` -> `11`
  - `librarian/docs/context-reshaping.md` -> `11`

Interpretation:

- The branch-local index is rebuilt and usable, but vector coverage is still only
  partial. This looks like prior embedding state preserved from an earlier run,
  not a fresh full-index embed pass for the current branch.

## Search Benchmark

Results summary:

| Pattern | Avg latency | Min/Max | Queries | Quality (high/medium/low) |
| --- | --- | --- | --- | --- |
| lexical | 1.9ms | 0.0/4.0ms | 10 | 10 / 0 / 0 |
| semantic | 64.0ms | 54.0/93.0ms | 10 | 5 / 0 / 5 |
| hybrid | 50.1ms | 45.0/61.0ms | 10 | 10 / 0 / 0 |

Observed result:

- Lexical remains the latency winner.
- Hybrid remains the best balanced mode on this corpus: `10/10` high-quality
  hits with materially better latency than semantic.

## Token-Savings Benchmark

| Query intent | Results | Raw tokens | Shaped tokens | Saved | % saved |
| --- | --- | --- | --- | --- | --- |
| Work tracking policy | 10 | 3064 | 1987 | 1077 | 35% |
| Context transfer between sessions | 10 | 3031 | 1992 | 1039 | 34% |
| Evidence standards for findings | 10 | 2372 | 1991 | 381 | 16% |
| Intake process and investigation ordering | 5 | 1707 | 1795 | -88 | -5% |
| Provider distribution and materialization | 5 | 1796 | 1911 | -115 | -6% |
| MCP server configuration for AI agents | 8 | 2043 | 1992 | 51 | 2% |
| Core HAWP protocol shape fields | 2 | 591 | 642 | -51 | -9% |
| Plan file structure and fields | 10 | 3292 | 1991 | 1301 | 40% |
| Kit maintenance and validation commands | 10 | 2581 | 1992 | 589 | 23% |
| Binary update and install flow | 10 | 3044 | 1993 | 1051 | 35% |
| **TOTAL** | — | **23521** | **18286** | **5235** | **22%** |

## Delta vs Pre-Tune

| Query intent | Pre-tune shaped | Post-tune shaped | Improvement |
| --- | --- | --- | --- |
| Intake process and investigation ordering | 1896 | 1795 | -101 |
| Provider distribution and materialization | 1981 | 1911 | -70 |
| Core HAWP protocol shape fields | 707 | 642 | -65 |
| Total shaped tokens | 18537 | 18286 | -251 |
| Total savings | 4984 | 5235 | +251 |
| Total savings percent | 21% | 22% | +1 point |

## Interpretation

- The shaping tune improved all three previously negative sparse-result cases.
- Negative savings still remain for those sparse queries, but the overhead is now
  smaller and closer to neutral.
- The next meaningful benchmark improvement likely requires either:
  - a more aggressive sparse-result passthrough mode, or
  - a full fresh embedding pass so hybrid/semantic behavior is measured against
    broader vector coverage.
