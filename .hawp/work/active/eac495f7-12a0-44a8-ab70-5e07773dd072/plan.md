# U4: classify inherited closures and supporting files accurately

**UUID:** `eac495f7-12a0-44a8-ab70-5e07773dd072`
**Type:** bug
**Status:** in-progress
**Source:** Mochilada upstream report U4, 2026-10-11
**Base:** `origin/development` at `d7c46a4e` (live remote tip checked 2026-10-11)
**Risk:** medium

## Investigation

**Directly verified in source:** librarian/src/internal/domain/work/validation/completeness.go; librarian/src/internal/domain/work/validation/types.go; .hawp/kit/references/backlog-alignment.md

**Root cause:** A fixed 2026-05-10 cutoff marks inherited September closures as failures; untyped supporting markdown in closed item folders is classified as a plan.

**Downstream evidence:** The Mochilada report describes the symptom. Its counts and scratch-copy result were not reproduced in this source checkout.

## Options

1. Keep current behavior and require downstream manual canonicalization.
2. Support an explicit inherited-closure marker and classify non-plan siblings in an item folder as supporting; document marking without inventing evidence.

## Recommended Fix

Support an explicit inherited-closure marker and classify non-plan siblings in an item folder as supporting; document marking without inventing evidence.

## Overlap Check

Can implement now: yes. Work is isolated on `codex/mochilada-upstream-triage` from `development`; the source worktree was clean before intake. Shared parser and normalization files are coordinated within this single branch.

## Verification Plan

Marker must be explicit and preserve warning visibility; current unmarked plans must continue to fail.

## Outcome

Local source fix: an explicit inherited-closure marker downgrades missing sections to a visible warning. Non-plan files within item folders are supporting, and flat supporting files can use an explicit marker. Downstream records are unchanged.

## Verification

Directly verified on the source branch: focused Go regression tests passed; `go test ./...` passed; `hawp kit validate`, `hawp work validate`, `hawp work normalize --dry-run --validate`, and `hawp check` passed. `hawp distribution sync` regenerated 40 guides and validated them. Downstream Mochilada runtime behavior after installation remains unverified.

## Close Checklist

- [ ] Outcome recorded
- [ ] Verification evidence recorded
- [ ] Backlog row moved to Recently Closed after verification
