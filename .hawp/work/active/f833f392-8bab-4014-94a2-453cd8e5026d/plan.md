# U9: resolve metadata for folder-style work plans

**UUID:** `f833f392-8bab-4014-94a2-453cd8e5026d`
**Type:** bug
**Status:** in-progress
**Source:** Mochilada upstream report U9, 2026-10-11
**Base:** `origin/development` at `d7c46a4e` (live remote tip checked 2026-10-11)
**Risk:** low

## Investigation

**Directly verified in source:** librarian/src/internal/domain/context/work.go; librarian/src/internal/application/index/ingest-service.go

**Root cause:** The search ingester only accepted UUID-shaped folder names for metadata, so numeric folder IDs were omitted. Separately, context enrichment resolved plan.md from its basename instead of the parent item folder.

**Downstream evidence:** The Mochilada report describes the symptom. Its counts and scratch-copy result were not reproduced in this source checkout.

## Options

1. Keep current behavior and require downstream manual canonicalization.
2. Accept supported legacy numeric IDs in the search ingester and resolve plan.md from its parent item folder in context enrichment; test active, parked, and closed layouts.

## Recommended Fix

Accept supported legacy numeric IDs in the search ingester and resolve plan.md from its parent item folder in context enrichment; test active, parked, and closed layouts.

## Overlap Check

Can implement now: yes. Work is isolated on `codex/mochilada-upstream-triage` from `development`; the source worktree was clean before intake. Shared parser and normalization files are coordinated within this single branch.

## Verification Plan

The search ingester's metadata count reflects recognized work IDs; the separate context enrichment path uses backlog rows for type and status.

## Outcome

Source correction: search index metadata was limited by canonicalWorkUUID accepting only UUID-shaped folder names; the separate context enrichment path also used the basename plan. Both paths now resolve numbered folder plans, and regression fixtures cover active, parked, and closed records.

## Verification

Directly verified on the source branch: focused Go regression tests passed; `go test ./...` passed; `hawp kit validate`, `hawp work validate`, `hawp work normalize --dry-run --validate`, and `hawp check` passed. `hawp distribution sync` regenerated 40 guides and validated them. Downstream Mochilada runtime behavior after installation remains unverified.

## Close Checklist

- [ ] Outcome recorded
- [ ] Verification evidence recorded
- [ ] Backlog row moved to Recently Closed after verification
