# U8: show the effective search index scope

**UUID:** `63700659-6b4f-4707-a4e7-17e3c0887c88`
**Type:** improvement
**Status:** in-progress
**Source:** Mochilada upstream report U8, 2026-10-11
**Base:** `origin/development` at `d7c46a4e` (live remote tip checked 2026-10-11)
**Risk:** low

## Investigation

**Directly verified in source:** librarian/src/internal/application/search/config.go; librarian/src/internal/platform/cli/index/ingest/command.go; .hawp/kit/usage/search.md

**Root cause:** Missing search.json selects built-in .hawp/kit and .hawp/work; the command already printed effective paths, but did not identify config-file presence and guides implied the file existed.

**Downstream evidence:** The Mochilada report describes the symptom. Its counts and scratch-copy result were not reproduced in this source checkout.

## Options

1. Keep current behavior and require downstream manual canonicalization.
2. Print effective index paths and default-source notice; document optional project configuration and adding docs/.

## Recommended Fix

Print effective index paths and default-source notice; document optional project configuration and adding docs/.

## Overlap Check

Can implement now: yes. Work is isolated on `codex/mochilada-upstream-triage` from `development`; the source worktree was clean before intake. Shared parser and normalization files are coordinated within this single branch.

## Verification Plan

Preserve current default scope and precedence.

## Outcome

Source correction: search index already printed effective paths. It now reports whether project/home search.json is present or built-in defaults apply; the search guide explains optional config and docs/ inclusion.

## Verification

Directly verified on the source branch: focused Go regression tests passed; `go test ./...` passed; `hawp kit validate`, `hawp work validate`, `hawp work normalize --dry-run --validate`, and `hawp check` passed. `hawp distribution sync` regenerated 40 guides and validated them. Downstream Mochilada runtime behavior after installation remains unverified.

## Close Checklist

- [ ] Outcome recorded
- [ ] Verification evidence recorded
- [ ] Backlog row moved to Recently Closed after verification
