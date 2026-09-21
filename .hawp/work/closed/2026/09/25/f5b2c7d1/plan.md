# Reject symlinked work trees before normalization traversal

## Outcome

The PR-review finding was implemented and reconciled with the current branch.

## Verification

The focused regression coverage and repository-wide Go, HAWP, distribution,
formatting, and diff-hygiene checks passed before close.

## Close Checklist

- [x] Outcome recorded.
- [x] Verification evidence recorded or referenced.
- [x] Backlog row removed from active coordination.
- [x] Plan archived under the close date.

**UUID:** `f5b2c7d1-9a42-4b2e-8a1c-6c3d7e9f1042`
**Type:** fix
**Reported:** 2026-09-25
**Status:** done

## Source Finding

The work-normalization scanner accepts any non-directory Markdown entry returned
by `ReadDir`, including symlinks. The normalization and duplicate-link paths
then read those paths and, in apply mode, can rewrite them. A symlinked
`.hawp/work` tree or nested plan can therefore redirect reads and writes outside
the repository.

## Problem

`normalization.ScanPlanFiles` and `WalkPlanMarkdown` do not establish that the
work corpus is repository-owned before traversal. The application entry points
must enforce that boundary before reading `BACKLOG.md`, scanning plan files, or
applying related-record edits.

## Evidence

- `librarian/src/internal/domain/work/normalization/scan.go` walks Markdown
  entries and reads them through the injected source without a symlink policy.
- `librarian/src/internal/application/work/normalize/duplicate_links.go` scans
  the work tree and writes changed plan files.
- `librarian/src/internal/application/work/normalize/normalize.go` performs
  multiple work-tree reads and mutations through the same normalization flow.
- The shared `RejectSymlinksInTree` guard already provides repository-root
  containment and nested symlink rejection for other corpus boundaries.

## Root Cause

The scanner treated `entry.IsDir() == false` as sufficient evidence of a safe
regular Markdown file, and the callers did not preflight the repository-owned
work corpus before invoking it.

## Scope

Guard the application normalization and duplicate-link entry points with the
existing repository-aware tree guard. Add focused regressions proving that a
symlinked work root and a nested symlinked plan are rejected before external
content is read or changed. Keep the domain scanner generic and capability
driven; do not add infrastructure imports to the domain package.

## Required Changes

1. Preflight `<repo>/.hawp/work` with `filesystem.RejectSymlinksInTree` before
   any normalization or duplicate-link traversal.
2. Fail closed with a clear `unsafe work tree` error and a non-zero command
   result.
3. Cover both the duplicate-link preview path and normalization dry-run/apply
   setup with symlink fixtures and external sentinels.
4. Preserve ordinary regular-file behavior and existing worktree semantics.

## Detailed Outcome

When the work root or any existing descendant is a symlink, normalization and
duplicate-link preview/apply stop before reading or writing the affected tree.
Regular repository-owned work trees continue to behave as before.

## Acceptance Criteria

- A symlinked `.hawp/work` root is rejected before `BACKLOG.md` or plan content
  is read.
- A nested symlinked plan is rejected before duplicate-link preview/apply can
  inspect or rewrite the target.
- No external sentinel file changes and no external content appears in command
  output.
- Existing normal and apply-mode tests remain green.
- The fix is represented by focused regression tests and does not broaden the
  domain package's filesystem dependencies.

## Validation Plan

- Focused normalization and duplicate-link tests.
- `go test ./...` from `librarian/src`.
- `go vet ./...` and `go build ./cmd/hawp` from `librarian/src`.
- `go run ./cmd/hawp check` and `go run ./cmd/hawp work validate`.
- Tracked Go formatting and `git diff --check`.

## Likely Files

- `librarian/src/internal/application/work/normalize/duplicate_links.go`
- `librarian/src/internal/application/work/normalize/duplicate_links_symlink_test.go`
- `librarian/src/internal/application/work/normalize/normalize.go`
- `librarian/src/internal/application/work/normalize/normalize_test.go`

## Dependencies

Depends on the existing shared `RejectSymlinksInTree` implementation and must
remain compatible with the already planned work-tree validation and indexing
hardening items.

## Risks / Notes

This intentionally rejects symlinks anywhere in the selected work corpus,
including symlinks that happen to point inside the repository. That is the
existing fail-closed policy for repository-controlled corpus traversal.

## Out of Scope

- Refactoring the generic domain scanner.
- Changing unrelated kit/index/provider symlink boundaries.
- Committing, pushing, merging, or claiming external PR review resolution.

## Implementation Evidence

The implementation was present as an uncommitted four-file patch when this work
item was created and was retained in place. The application normalization and
duplicate-link entry points now call `RejectSymlinksInTree` before traversal,
and focused regressions cover the symlinked-root and nested-plan cases.

## Direct Verification

- Focused uncached symlink regressions: pass.
- `go test ./...`: pass.
- `go vet ./...`: pass.
- `go build ./cmd/hawp`: pass.
- `go run ./cmd/hawp check`: pass; all three validations passed.
- `go run ./cmd/hawp work validate`: pass; one pre-existing verification-clarity warning remains.
- `git diff --check`: pass.
- Scoped tracked Go `gofmt -l` check: pass with no output.
