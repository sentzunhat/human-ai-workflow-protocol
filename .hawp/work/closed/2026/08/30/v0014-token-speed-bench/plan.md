# v0014-token-speed-bench — v0.0.14 token-savings and speed benchmark gate

**Type:** feature
**Status:** done
**Branch:** `feature/v0.0.13`
**Opened:** 2026-08-26
**Closed:** 2026-08-30

## Problem

`v0.0.14` needed stronger proof than “feature works.” The release needed to
show that token counting was saving context budget and that the associated
search and shaping paths remained fast enough to justify the extra behavior.

## Goal

Make benchmark evidence a merge/release gate for `v0.0.14`, with two explicit
questions answered:

1. Are token-in / token-out measurements showing meaningful savings?
2. What is the speed cost or speed improvement of the associated feature work?

## Outcome

Closed 2026-08-30.

The benchmark gate is now satisfied by the checked-in evidence chain:

- `.hawp/work/evidence/2026/08/24/search-benchmark-v006.md` captures the speed
  side of the gate for the v0.0.6 search stack that later became the benchmark
  harness baseline used by the `v0.0.14` line:
  - lexical: `0.2ms`, `10/10` high-quality
  - hybrid: `60.5ms`, `10/10` high-quality
- `.hawp/work/evidence/2026/08/27/v0019-token-savings-benchmark.md` captures
  the token-savings side of the gate with the real `hawp search benchmark --token`
  command:
  - total raw tokens: `22089`
  - total shaped tokens: `17986`
  - total saved: `4103` (`19%`)
  - dense queries saved up to `38%`
  - sparse-result negatives were explicitly preserved rather than hidden

This item is closed as complete because both acceptance dimensions now have
direct evidence: one artifact for latency/quality and one artifact for
raw-vs-shaped token savings.

## Verification

- [x] `.hawp/work/evidence/2026/08/24/search-benchmark-v006.md` records speed
      evidence with query-count, latency, and quality notes
- [x] `.hawp/work/evidence/2026/08/27/v0019-token-savings-benchmark.md` records
      token-savings evidence with per-query and aggregate totals
- [x] `librarian/src/CHANGELOG.md` entry for `0.0.19` explicitly states
      `hawp search benchmark --token` closes this gate
- [x] Evidence includes neutral/negative sparse cases instead of cherry-picking

## Close Checklist

- [x] Outcome recorded
- [x] Verification cites the benchmark artifacts directly
- [x] Gate is closed without merging additional code in this session
