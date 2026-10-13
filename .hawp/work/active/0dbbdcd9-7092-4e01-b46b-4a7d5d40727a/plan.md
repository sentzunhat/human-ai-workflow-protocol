# U6: keep installed runtime artifacts out of Git

**UUID:** `0dbbdcd9-7092-4e01-b46b-4a7d5d40727a`
**Type:** bug
**Status:** in-progress
**Source:** Mochilada upstream report U6, 2026-10-11
**Base:** `origin/development` at `d7c46a4e` (live remote tip checked 2026-10-11)
**Risk:** medium

## Investigation

**Directly verified in source:** distribution/sources/install/script-core.md; distribution/sources/update/script-core.md; distribution/sources/install/script-footer.md

**Root cause:** Install writes a binary to .hawp/bin while neither install nor update seeds ignore guidance for .hawp/bin and .hawp/db.

**Downstream evidence:** The Mochilada report describes the symptom. Its counts and scratch-copy result were not reproduced in this source checkout.

## Options

1. Keep current behavior and require downstream manual canonicalization.
2. Seed a scoped .hawp/.gitignore without overwriting user content; document existing tracked-file cleanup and legacy shims.

## Recommended Fix

Seed a scoped .hawp/.gitignore without overwriting user content; document existing tracked-file cleanup and legacy shims.

## Overlap Check

Can implement now: yes. Work is isolated on `codex/mochilada-upstream-triage` from `development`; the source worktree was clean before intake. Shared parser and normalization files are coordinated within this single branch.

## Verification Plan

Verify generated guides and install/update behavior; do not remove downstream files.

## Outcome

Local source fix: the distributed core seeds a non-clobbering .hawp/.gitignore for bin/ and db/. Generated install/update guides mention already tracked files and preserved legacy shims. Existing downstream Git tracking is unchanged.

## Verification

Directly verified on the source branch: focused Go regression tests passed; `go test ./...` passed; `hawp kit validate`, `hawp work validate`, `hawp work normalize --dry-run --validate`, and `hawp check` passed. `hawp distribution sync` regenerated 40 guides and validated them. Downstream Mochilada runtime behavior after installation remains unverified.

## Close Checklist

- [ ] Outcome recorded
- [ ] Verification evidence recorded
- [ ] Backlog row moved to Recently Closed after verification
