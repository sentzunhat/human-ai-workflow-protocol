# U3: handle missing legacy work types without guessing

**UUID:** `ec70c29a-2fc4-48ee-965d-6c301df88510`
**Type:** improvement
**Status:** in-progress
**Source:** Mochilada upstream report U3, 2026-10-11
**Base:** `origin/development` at `d7c46a4e` (live remote tip checked 2026-10-11)
**Risk:** medium

## Investigation

**Directly verified in source:** librarian/src/internal/domain/work/normalization/rules.go; librarian/src/internal/domain/work/normalization/scan.go

**Root cause:** B1 blocks rows whose IDs provide no type and whose legacy table has no Type column.

**Downstream evidence:** The Mochilada report describes the symptom. Its counts and scratch-copy result were not reproduced in this source checkout.

## Options

1. Keep current behavior and require downstream manual canonicalization.
2. Treat absent legacy Type column distinctly from a missing value in a canonical Type column; keep ambiguous records explicitly unknown and avoid keyword guesses.

## Recommended Fix

Treat absent legacy Type column distinctly from a missing value in a canonical Type column; keep ambiguous records explicitly unknown and avoid keyword guesses.

## Overlap Check

Can implement now: yes. Work is isolated on `codex/mochilada-upstream-triage` from `development`; the source worktree was clean before intake. Shared parser and normalization files are coordinated within this single branch.

## Verification Plan

Check that current-schema rows still require an explicit type and no historical type is fabricated.

## Outcome

Local source policy: a closed table without a Type column no longer produces speculative B1 blocked operations. Rows with a Type column still require a value when the ID cannot infer one; no type is guessed from a title.

## Verification

Directly verified on the source branch: focused Go regression tests passed; `go test ./...` passed; `hawp kit validate`, `hawp work validate`, `hawp work normalize --dry-run --validate`, and `hawp check` passed. `hawp distribution sync` regenerated 40 guides and validated them. Downstream Mochilada runtime behavior after installation remains unverified.

## Close Checklist

- [ ] Outcome recorded
- [ ] Verification evidence recorded
- [ ] Backlog row moved to Recently Closed after verification
