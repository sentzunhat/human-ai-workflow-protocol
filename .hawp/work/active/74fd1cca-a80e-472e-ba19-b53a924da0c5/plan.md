# U5: detect parked rows whose plans remain active

**UUID:** `74fd1cca-a80e-472e-ba19-b53a924da0c5`
**Type:** bug
**Status:** in-progress
**Source:** Mochilada upstream report U5, 2026-10-11
**Base:** `origin/development` at `d7c46a4e` (live remote tip checked 2026-10-11)
**Risk:** medium

## Investigation

**Directly verified in source:** librarian/src/internal/domain/work/validation/consistency.go; librarian/src/internal/application/work/normalize/normalize.go

**Root cause:** Parked consistency currently checks link existence but ignores active/ versus parked/ folder role.

**Downstream evidence:** The Mochilada report describes the symptom. Its counts and scratch-copy result were not reproduced in this source checkout.

## Options

1. Keep current behavior and require downstream manual canonicalization.
2. Report status-folder drift and propose a safe, reviewed migration only when source and destination are unambiguous.

## Recommended Fix

Report status-folder drift and propose a safe, reviewed migration only when source and destination are unambiguous.

## Overlap Check

Can implement now: yes. Work is isolated on `codex/mochilada-upstream-triage` from `development`; the source worktree was clean before intake. Shared parser and normalization files are coordinated within this single branch.

## Verification Plan

Do not move files or rewrite links on dry-run; avoid overwriting a destination.

## Outcome

Local source fix: validation warns when a parked row points into active/, normalization reports B8 with a proposed parked path, and --migrate-folders dry-run reports review-needed. Apply refuses that drift until an owner reviews the move and relink; no automatic move was implemented.

## Verification

Directly verified on the source branch: focused Go regression tests passed; `go test ./...` passed; `hawp kit validate`, `hawp work validate`, `hawp work normalize --dry-run --validate`, and `hawp check` passed. `hawp distribution sync` regenerated 40 guides and validated them. Downstream Mochilada runtime behavior after installation remains unverified.

## Close Checklist

- [ ] Outcome recorded
- [ ] Verification evidence recorded
- [ ] Backlog row moved to Recently Closed after verification
