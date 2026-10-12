# U2: resolve legacy closed-row plan links without treating prose as paths

**UUID:** `8655c5d6-df73-4155-ba6c-e1fd1a28739f`
**Type:** bug
**Status:** in-progress
**Source:** Mochilada upstream report U2, 2026-10-11
**Base:** `origin/development` at `d7c46a4e` (live remote tip checked 2026-10-11)
**Risk:** high

## Investigation

**Directly verified in source:** librarian/src/internal/domain/work/normalization/scan.go; librarian/src/internal/domain/work/normalization/rules.go; librarian/src/internal/domain/work/backlog/parse.go

**Root cause:** Closed-row path extraction reads Detail prose while overlooking a markdown link in ID; Title ignores Work alias.

**Downstream evidence:** The Mochilada report describes the symptom. Its counts and scratch-copy result were not reproduced in this source checkout.

## Options

1. Keep current behavior and require downstream manual canonicalization.
2. Prefer explicit links in supported cells, map Work to Title, and emit an actionable format issue for absent links rather than a false missing-path error.

## Recommended Fix

Prefer explicit links in supported cells, map Work to Title, and emit an actionable format issue for absent links rather than a false missing-path error.

## Overlap Check

Can implement now: yes. Work is isolated on `codex/mochilada-upstream-triage` from `development`; the source worktree was clean before intake. Shared parser and normalization files are coordinated within this single branch.

## Verification Plan

Preserve canonical table behavior; test a legacy five-column closed row and avoid automatic rewriting of prose.

## Outcome

Local source fix: linked IDs resolve to their labels for validation, normalization selects an explicit plan link instead of prose, and Work maps to Title. The source does not rewrite legacy table columns; that remains an optional migration.

## Verification

Directly verified on the source branch: focused Go regression tests passed; `go test ./...` passed; `hawp kit validate`, `hawp work validate`, `hawp work normalize --dry-run --validate`, and `hawp check` passed. `hawp distribution sync` regenerated 40 guides and validated them. Downstream Mochilada runtime behavior after installation remains unverified.

## Close Checklist

- [ ] Outcome recorded
- [ ] Verification evidence recorded
- [ ] Backlog row moved to Recently Closed after verification
