# v0.0.23 legacy work-item UUID canonicalization follow-up

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `89cf7a85-0dc4-43cf-9b61-c31a4c770616`
**Type:** improvement
**Reported:** 2026-08-30

---

## Input (verbatim)

> Continue the v0.0.23 work migration by assigning or reconciling canonical UUID-native IDs for remaining non-UUID active/parked items, updating plan headers and backlog references safely, and defining the archive policy boundary for legacy closed records that should not be rewritten casually.

## Intake Summary

Track the next cleanup after the validated folder-migration checkpoint: decide
how far to canonicalize remaining non-UUID active/parked records, prove
older-repo safety, and avoid casual archive rewrites that would create churn
without release value.

## Current Context

- The active `v0.0.23` migration work already proved the folder-based work
  layout on this repository:
  `go run ./cmd/hawp work normalize --apply --migrate-folders --force-dirty --validate`
  made no further changes after the migration pass
- `go run ./cmd/hawp work validate` and `go run ./cmd/hawp check` now pass on
  the live source path with `0 issues, 0 warnings`
- The current repo contains a mix of:
  - UUID-native active items
  - legacy/non-UUID parked items that predate the folder migration
  - legacy flat closed archives that should not be rewritten casually without a
    concrete repair reason
- The user strategy for this lane is to preserve old-repo safety and keep the
  release patch simple

## Initial Analysis

**Directly verified:**

- Folder-per-item active and recent closed layouts are already working in this
  repository without further migration writes
- The current validation cleanup was able to support both legacy flat closed
  files and folder-based `closed/YYYY/MM/DD/<id>/plan.md` records
- The remaining non-UUID items are mostly parked backlog entries rather than a
  release-critical active-path breakage

**Inferred (not yet proven):**

- Older repositories may still benefit from a safer canonicalization pass for
  parked items and backlog references, but that proof belongs in its own lane
- Rewriting legacy closed archives for cosmetic consistency would add churn and
  risk without helping the `v0.0.23` release

**Likely scope:**

- Audit non-UUID active/parked items and choose the smallest safe
  canonicalization boundary
- Prove older-repo behavior on a pre-upgrade-style fixture or checkout
- Keep closed-archive rewrite policy conservative unless a broken reference
  demands repair

## Risk + Review Gate

**Risk:** medium
**Gate:** review first on medium/high

## Backlog + Plan Link

**Status now:** parked
**Plan file:** work/parked/89cf7a85/plan.md

## Next Step

- [x] Investigation recorded above (required before planning)
- [x] Write or update the plan file
- [x] Move backlog status accordingly

## Outcome

Parked 2026-08-31.

This lane is intentionally deferred out of the `v0.0.23` release. The current
repository already proves that folder-based work records validate cleanly, so
the remaining UUID canonicalization work is now about older-repo safety and
conservative archive policy rather than fixing a live release blocker.

## Verification

- [x] Source-backed migration is already stable in this repo. Evidence:
      [03299078 plan](../../closed/2026/08/31/03299078/plan.md)
      records the no-op `--migrate-folders` validation checkpoint
- [x] Current validation supports both legacy flat and folder-based closed
      layouts. Evidence: source updates in
      `librarian/src/internal/domain/work/evidence.go` and
      `librarian/src/internal/domain/work/validations_test.go` collect both
      closed-plan layouts and cover the folder-plan case
- [x] No release-critical breakage remains in the live repo. Evidence:
      [benchmark rerun](../../evidence/2026/08/31/c804eec0-v0-0-23-benchmark-rerun.md)
      records clean `hawp check` and `work validate` runs

## Park Checklist

- [x] Release scope kept simple
- [x] Older-repo proof isolated as a future lane
- [x] Archive rewrite policy kept conservative
