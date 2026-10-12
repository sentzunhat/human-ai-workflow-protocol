# U1: recognize backlog section heading variants

**UUID:** `236b92e3-e4d4-4882-aae9-20c82ec4642e`
**Type:** bug
**Status:** in-progress
**Source:** Mochilada upstream report U1, 2026-10-11
**Base:** `origin/development` at `d7c46a4e` (live remote tip checked 2026-10-11)
**Risk:** high

## Investigation

**Directly verified in source:** librarian/src/internal/domain/work/normalization/scan.go; librarian/src/internal/domain/work/backlog/parse.go; distribution/sources/install/script-core.md; distribution/sources/update/script-core.md

**Root cause:** Exact-case heading checks miss lowercase and slash-spacing variants; A8 can add duplicate sections.

**Downstream evidence:** The Mochilada report describes the symptom. Its counts and scratch-copy result were not reproduced in this source checkout.

## Options

1. Keep current behavior and require downstream manual canonicalization.
2. Share tolerant heading recognition across validation and normalization; align install/update awk and add variant fixtures.

## Recommended Fix

Share tolerant heading recognition across validation and normalization; align install/update awk and add variant fixtures.

## Overlap Check

Can implement now: yes. Work is isolated on `codex/mochilada-upstream-triage` from `development`; the source worktree was clean before intake. Shared parser and normalization files are coordinated within this single branch.

## Verification Plan

Keep matching restricted to level-two section headings; verify no duplicate A8 operation and closed-plan reconciliation.

## Outcome

Local source fix: the backlog parser and normalization parser recognize case and slash-spacing variants; install/update reconciliation recognizes lowercase Recently Closed. Regression fixtures protect the behavior.

## Verification

Directly verified on the source branch: focused Go regression tests passed; `go test ./...` passed; `hawp kit validate`, `hawp work validate`, `hawp work normalize --dry-run --validate`, and `hawp check` passed. `hawp distribution sync` regenerated 40 guides and validated them. Downstream Mochilada runtime behavior after installation remains unverified.

## Close Checklist

- [ ] Outcome recorded
- [ ] Verification evidence recorded
- [ ] Backlog row moved to Recently Closed after verification
