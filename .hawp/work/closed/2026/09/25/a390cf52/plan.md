# Prevent provider materializer writes through symlinked paths

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

**UUID:** `a390cf52-9dfb-45b5-a7db-8d619eee5a01`
**Type:** bug
**Reported:** 2026-09-23

## Intake Summary

Copilot reports that provider materialization computes repository-relative outputs but `providersync.Materialize` calls `os.MkdirAll` and `os.WriteFile` without validating the repository root, output ancestors, or destination. A symlink under `core/providers` can redirect generated provider writes outside the checkout.

## Current Context

The repository already has shared filesystem guards in `librarian/src/internal/infrastructure/filesystem/symlink_guard.go`, including root/ancestor validation and an atomic writer. Other mutation paths use those guards. Provider materialization currently performs direct filesystem mutation in the application adapter.

## Initial Analysis

**Directly verified:** `librarian/src/internal/application/providersync/providersync.go` calls `os.MkdirAll(filepath.Dir(output.OutputPath))` and `os.WriteFile(output.OutputPath, ...)` after `ComputeOutputs`.

**Inferred:** A symlinked output directory or destination can redirect writes outside `repoRoot`; `MkdirAll` also creates missing path components before any safety check.

**Likely scope:** provider materialization adapter, filesystem guard reuse, and symlink regression tests.

## Root Cause

The path computation layer establishes lexical paths but the materializer does not enforce the filesystem containment invariant before mutation.

## Options

1. Reuse `filesystem.RejectSymlinkAncestors` plus `filesystem.AtomicWriteFile` for every output. This aligns with existing repository-wide mutation contracts.
2. Add a provider-specific path walker and temporary-write implementation. This duplicates security-sensitive logic.

## Recommended Fix

Validate `repoRoot` and each output path before directory creation, use the shared atomic writer for writes, and revalidate before the write/rename. Add tests for symlinked output ancestors and destinations.

## Risk + Review Gate

**Risk:** high — generated repository writes.
**Gate:** user explicitly requested implementation and sequential remediation.

## Outcome

Done. Provider materialization now validates every output path before
directory creation and uses the shared atomic writer for generated files.
Symlinked output ancestors are covered by focused regression tests.
