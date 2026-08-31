# v0.0.23 benchmark evidence and release artifact refresh

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `c804eec0-e551-40d8-8941-67b63d2ec257`
**Type:** test
**Reported:** 2026-08-29

---

## Input (verbatim)

> Do the code audit refactoring needs please on a patch version branch and multiple agents working on the work and benchmark don't merge

## Intake Summary

Parallel benchmark/evidence lane for the v0.0.23 patch branch. This lane owns
benchmark artifacts, release-evidence notes, and benchmark-doc alignment only.

## Current Context

- Runs in a separate worktree/task from the main refactor lane
- Main implementation files are reserved by work item `387a37b7`
- Primary source files for this lane are benchmark docs, evidence artifacts, and
  this work item's own records

## Initial Analysis

**Directly verified:**

- A dedicated benchmark lane was requested by the user
- Existing benchmark gate plans already exist under
  `.hawp/work/active/v0014-token-speed-bench/plan.md` and
  `.hawp/work/active/release-benchmark-backfill/plan.md`

**Inferred (not yet proven):**

- The most useful near-term output is likely an evidence-gap audit plus any safe
  artifact refresh that does not collide with implementation changes

**Likely scope:**

- Audit benchmark coverage for the current patch/v0.1.0 direction
- Prepare or update evidence/docs artifacts only if non-conflicting
- Report remaining gaps back to the main lane

## Risk + Review Gate

**Risk:** low
**Gate:** auto-implement within docs/evidence scope only; do not merge

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/closed/2026/08/31/c804eec0/plan.md

## Next Step

- [x] Investigation recorded above
- [ ] Benchmark/evidence lane to continue in parallel task
- [ ] Return findings and artifact status to main lane

## 2026-08-31 Benchmark Checkpoint

**Directly verified:**

- `mise exec node@26.5.0 -- go run ./cmd/hawp work validate` now reports
  `VALIDATION PASS (0 issues, 0 warnings)` on the source-backed path
- `mise exec node@26.5.0 -- go run ./cmd/hawp check` passes with kit, work, and
  links all green
- `mise exec node@26.5.0 -- go run ./cmd/hawp search benchmark --token`
  reran successfully with aggregate savings of `4838` tokens (`21%`),
  improving on the prior `19%` aggregate result
- sparse-query negatives were reduced to one negligible `-1` token case on
  `Provider distribution and materialization`, while the earlier worst cases
  moved to `0`
- Evidence recorded at
  `.hawp/work/evidence/2026/08/31/c804eec0-v0-0-23-benchmark-rerun.md`

**Inference:**

- The benchmark lane now has a current branch-local token-savings checkpoint
  that is good enough to use as the live `v0.0.23` benchmark reference

## Outcome

Closed 2026-08-31.

The `v0.0.23` benchmark/evidence lane is complete. The repo now has a current
token-savings rerun tied to the final validation-cleanup state, and that rerun
is strong enough to serve as the release benchmark checkpoint.

## Verification

- [x] Source-backed validation is clean. Evidence:
      [.hawp/work/evidence/2026/08/31/c804eec0-v0-0-23-benchmark-rerun.md](/Users/beltrd/Desktop/projects/sentzunhat/human-ai-workflow-protocol/.hawp/work/evidence/2026/08/31/c804eec0-v0-0-23-benchmark-rerun.md)
      records `work validate` and `hawp check` passing
- [x] Token-savings rerun is recorded with aggregate and per-query results.
      Evidence:
      [.hawp/work/evidence/2026/08/31/c804eec0-v0-0-23-benchmark-rerun.md](/Users/beltrd/Desktop/projects/sentzunhat/human-ai-workflow-protocol/.hawp/work/evidence/2026/08/31/c804eec0-v0-0-23-benchmark-rerun.md)
      shows `22866` raw, `18028` shaped, `4838` saved (`21%`)
- [x] The remaining sparse negative is explicitly recorded instead of hidden.
      Evidence:
      [.hawp/work/evidence/2026/08/31/c804eec0-v0-0-23-benchmark-rerun.md](/Users/beltrd/Desktop/projects/sentzunhat/human-ai-workflow-protocol/.hawp/work/evidence/2026/08/31/c804eec0-v0-0-23-benchmark-rerun.md)
      notes the `-1` token outlier as negligible

## Close Checklist

- [x] Benchmark artifact refreshed
- [x] Validation state tied to the artifact
- [x] Release gate evidence is explicit
