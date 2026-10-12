# U7: correct the search index next-step hint

**UUID:** `7f224fa6-dba2-4caf-8627-e7ab11e3238d`
**Type:** bug
**Status:** in-progress
**Source:** Mochilada upstream report U7, 2026-10-11
**Base:** `origin/development` at `d7c46a4e` (live remote tip checked 2026-10-11)
**Risk:** low

## Investigation

**Directly verified in source:** librarian/src/internal/platform/cli/index/ingest/command.go; librarian/src/internal/platform/cli/registry.go

**Root cause:** Successful ingest prints hawp search query even though the command is hawp search <query>.

**Downstream evidence:** The Mochilada report describes the symptom. Its counts and scratch-copy result were not reproduced in this source checkout.

## Options

1. Keep current behavior and require downstream manual canonicalization.
2. Correct the one-line hint and verify it matches CLI usage.

## Recommended Fix

Correct the one-line hint and verify it matches CLI usage.

## Overlap Check

Can implement now: yes. Work is isolated on `codex/mochilada-upstream-triage` from `development`; the source worktree was clean before intake. Shared parser and normalization files are coordinated within this single branch.

## Verification Plan

No search behavior change.

## Outcome

Local source fix: the successful index hint now shows hawp search "<query>".

## Verification

Directly verified on the source branch: focused Go regression tests passed; `go test ./...` passed; `hawp kit validate`, `hawp work validate`, `hawp work normalize --dry-run --validate`, and `hawp check` passed. `hawp distribution sync` regenerated 40 guides and validated them. Downstream Mochilada runtime behavior after installation remains unverified.

## Close Checklist

- [ ] Outcome recorded
- [ ] Verification evidence recorded
- [ ] Backlog row moved to Recently Closed after verification
