# Guard index corpus roots before build traversal

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

**UUID:** `2da1adfc-f41b-4785-9cc2-3442332f3260`
**Type:** fix
**Reported:** 2026-09-23
**Source plan:** `d66d12c6`

## Context

Index build reads the repository’s `.hawp/kit` and `.hawp/work` corpus roots through recursive Markdown collection. The existing shared filesystem guard supplies a repository-aware, fail-closed symlink boundary but is not used by this service.

## Fix Work

Validate the relevant corpus tree from the repository root before indexing it. Preserve legitimate regular directory and missing-directory behavior. Add a symlinked external corpus regression proving no external content is returned.

## Verification

Run focused index build tests, then `go test ./...`, `go vet ./...`, HAWP validation, formatting, and diff-hygiene checks.

## Outcome

Done. `BuildService.Execute` rejects symlinked `.hawp/kit` and `.hawp/work`
trees before traversal while preserving ordinary corpus behavior. Focused
index build tests cover symlinked roots and descendants.
