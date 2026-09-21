# Reject symlinked roots before indexing repository content

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

**UUID:** `d66d12c6-5e63-48f3-bea9-a95245cac659`
**Type:** bug
**Reported:** 2026-09-23

## Intake Summary

Copilot reports that index build traverses `.hawp/kit` and `.hawp/work` without validating their repository containment. A symlinked corpus root can cause external Markdown to be included in index build and export output.

## Current Context

`librarian/src/internal/application/index/build-service.go` passes the repository root and corpus roots directly to `domain/context.EnrichKit` and `EnrichWork`; their file lister ultimately calls `os.ReadDir` recursively.

## Initial Analysis

**Directly verified:** `BuildService.Execute` has no `RejectSymlinkAncestors` or tree-root check before `EnrichKit` and `EnrichWork`; shared filesystem helpers already reject symlinked roots and ancestors for other repository-bound operations.

**Inferred:** `os.ReadDir` on a requested directory follows the directory symlink, so a symlinked `.hawp`, `kit`, or `work` root can expose external corpus documents.

**Likely scope:** application index build boundary and its direct service regression test.

## Root Cause

The build service treats corpus roots as trusted read paths despite the index/export contract requiring repository-owned content.

## Options

1. Validate the selected corpus root with the existing shared `RejectSymlinksInTree` boundary before traversal. This rejects root and nested symlinks and matches existing repository traversal hardening.
2. Change the generic Markdown collector. This would affect unrelated read-only callers and still lacks repository-root context.

## Recommended Fix

Use shared filesystem containment validation in `BuildService.Execute` for each selected corpus root, before enrichment, with a regression proving an external symlinked kit root is rejected.

## Risk + Review Gate

**Risk:** high — can disclose external Markdown through indexed/exported corpus output.
**Gate:** user explicitly authorized sequential remediation in this PR loop.

## Outcome

Done. The build service validates both selected corpus roots with the shared
repository-aware symlink guard before enrichment. Root and descendant
symlink regressions pass.

See also the completed implementation item `2da1adfc`.
