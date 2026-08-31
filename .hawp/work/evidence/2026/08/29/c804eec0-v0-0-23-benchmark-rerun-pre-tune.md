---
work-item: c804eec0
date: 2026-08-29
type: benchmark
title: "v0.0.23 branch-local benchmark rerun before shaping tune"
status: measured
---

# v0.0.23 Benchmark Rerun — Pre-Tune

## Scope

Branch-local rerun after rebuilding `.hawp/db/index.sqlite` on
`feature/v0.0.23-refactor-v010`, before any shaping-path tuning.

## Directly Verified

- `./.hawp/bin/hawp search index` rebuilt the local index successfully.
- `.hawp/db/index.sqlite` now contains live tables and data.
- `select count(*) from documents` -> `434`
- `select count(*) from chunks` -> `3382`
- `select count(*) from chunks where embedding is not null` -> `35`

## Search Benchmark

Command:

```text
./.hawp/bin/hawp search benchmark
```

Results summary:

| Pattern | Avg latency | Min/Max | Queries | Quality (high/medium/low) |
| --- | --- | --- | --- | --- |
| lexical | 1.1ms | 0.0/3.0ms | 10 | 10 / 0 / 0 |
| semantic | 61.5ms | 53.0/93.0ms | 10 | 5 / 0 / 5 |
| hybrid | 49.9ms | 46.0/60.0ms | 10 | 10 / 0 / 0 |

Observed result:

- Lexical was the latency winner.
- Hybrid matched lexical on the simple keyword-quality check and stayed well
  below semantic latency.

## Token-Savings Benchmark

Command:

```text
./.hawp/bin/hawp search benchmark --token
```

Results:

| Query intent | Results | Raw tokens | Shaped tokens | Saved | % saved |
| --- | --- | --- | --- | --- | --- |
| Work tracking policy | 10 | 3064 | 1987 | 1077 | 35% |
| Context transfer between sessions | 10 | 3031 | 1995 | 1036 | 34% |
| Evidence standards for findings | 10 | 2372 | 1993 | 379 | 16% |
| Intake process and investigation ordering | 5 | 1707 | 1896 | -189 | -11% |
| Provider distribution and materialization | 5 | 1796 | 1981 | -185 | -10% |
| MCP server configuration for AI agents | 8 | 2043 | 1994 | 49 | 2% |
| Core HAWP protocol shape fields | 2 | 591 | 707 | -116 | -20% |
| Plan file structure and fields | 10 | 3292 | 1992 | 1300 | 39% |
| Kit maintenance and validation commands | 10 | 2581 | 1996 | 585 | 23% |
| Binary update and install flow | 10 | 3044 | 1996 | 1048 | 34% |
| **TOTAL** | — | **23521** | **18537** | **4984** | **21%** |

## Interpretation

- Context shaping is net positive on this branch-local rerun: **21%** total
  estimated token reduction.
- Three sparse-result queries showed negative savings, consistent with
  formatting overhead outweighing truncation or deduplication benefit.

## Next Step

Tune the shaping path to reduce overhead on sparse/under-budget result sets, then
rerun the same benchmark commands and compare the delta.
