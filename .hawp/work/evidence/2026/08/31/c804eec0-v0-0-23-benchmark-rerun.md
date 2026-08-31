---
work-item: c804eec0
date: 2026-08-31
type: benchmark
title: "v0.0.23 token benchmark rerun after verification-clarity cleanup"
status: complete
---

# v0.0.23 Token Benchmark Rerun

## Scope

Rerun the token-savings benchmark for the current `v0.0.23` branch-local work
after the verification-clarity cleanup, and confirm that HAWP validation is
clean on the live source path.

## Directly Observed

### Source-backed validation

- `mise exec node@26.5.0 -- go run ./cmd/hawp work validate`
  - `BACKLOG CONSISTENCY` -> `PASS`
  - `CLOSED TASK COMPLETENESS` -> `PASS`
  - `EVIDENCE INTEGRITY` -> `PASS`
  - `VERIFICATION CLARITY` -> `PASS`
  - `DEAD LINKS` -> `PASS`
  - summary -> `VALIDATION PASS (0 issues, 0 warnings)`

- `mise exec node@26.5.0 -- go run ./cmd/hawp check`
  - `kit validate` -> `3 checks passed, 0 issues`
  - `work validate` -> `VALIDATION PASS (0 issues, 0 warnings)`
  - `links check` -> `Checked 134 Markdown file(s): local links are valid.`

### Token benchmark rerun

- Command: `mise exec node@26.5.0 -- go run ./cmd/hawp search benchmark --token`
- Shape budget: `2000`
- Aggregate result:
  - raw tokens: `22866`
  - shaped tokens: `18028`
  - saved: `4838` (`21%`)

### Per-query benchmark highlights

- Strong savings remained on dense queries:
  - `Plan file structure and fields` -> `1233` tokens saved (`38%`)
  - `Context transfer between sessions` -> `865` tokens saved (`30%`)
  - `Binary update and install flow` -> `854` tokens saved (`30%`)

- Sparse-query behavior improved versus the previous negative cases:
  - `Intake process and investigation ordering` -> `0` saved (`0%`)
  - `Core HAWP protocol shape fields` -> `0` saved (`0%`)
  - `Provider distribution and materialization` -> `-1` token (`-0%`)

## Inference

- The current shaping path still behaves as intended on dense contexts and now
  avoids the earlier materially negative sparse-query overhead.
- The remaining `-1` token case is negligible enough to treat as formatting
  noise rather than a benchmark blocker.

## Constraints

- The benchmark and validation were run from `librarian/src` via `go run`
  because the source-backed path was the trusted verification path for this
  cleanup pass.

## Suggested Next Step

- Keep this artifact as the current `v0.0.23` token-benchmark checkpoint for
  the benchmark/evidence lane, and only reopen shaping changes if a future run
  reintroduces materially negative sparse-query cases.
