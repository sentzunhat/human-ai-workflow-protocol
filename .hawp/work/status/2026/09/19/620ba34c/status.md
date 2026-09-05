# Status Report

## Intent

Complete the remaining `domain/work` migration boundary and decide whether
validation needs a separate capability package.

## Current State

Migration orchestration now lives in `internal/domain/work/normalization`.
The parent `domain/work` package retains the existing `WorkSource` methods as
compatibility adapters and supplies the concrete filesystem seam.

## What Was Inspected

- `librarian/src/internal/domain/work/normalize_migrate.go`
- `librarian/src/internal/domain/work/work_source.go`
- `librarian/src/internal/domain/work/normalization/*`
- Existing migration compatibility tests
- `internal/domain/work` callers and package topology

## What Changed

- Added `normalization.MigrationIO` as the explicit filesystem boundary.
- Moved folder migration orchestration, sidecar movement, moved-plan tracking,
  and backlog-link reconciliation into `normalization/migration_apply.go`.
- Reduced the parent migration file to preview/copy behavior and compatibility
  wiring.
- Preserved `WorkSource.ApplyWorkItemFolderMigration` and
  `WorkSource.PreviewWorkItemFolderMigration` unchanged for callers.
- Extracted reusable `model`, `validation`, `backlog`, and `intake` packages;
  the parent package retains compatibility aliases and adapters.
- Co-located backlog and intake tests with their new packages.

## What Was Directly Verified

- `go test ./...` passed.
- `go vet ./...` passed.
- `git diff --check` passed.
- Full repository checks still pass after the package batch and test moves.
- A runtime-import audit found 15 production `domain/work` files still using
  `os` or concrete filesystem traversal; the package split alone did not yet
  remove that dependency.
- Existing parent migration tests passed, including flat-plan movement,
  canonical folder renaming, idempotence, and preview/apply equivalence.
- The parent package no longer owns the migration orchestration body.

## What Remains Unproven

- Historical closed-record evidence quality remains a separate review concern.
- Normalization remains intentionally cohesive; its scan, rules, reports, and
  mutation files share one operation model and explicit seams.
- The strict filesystem boundary is still open for validation scans,
  normalization scans/cleanup, closed-record reconciliation, and migration
  preview copying.
- Validation policy now consumes an injected filesystem observation source;
  the application validation service supplies the concrete `os` functions.
- Normalization plan/backlog scanning now consumes an injected read-only
  `ScanSource`; the parent wrapper remains only as a compatibility adapter.
- Active-row cleanup now consumes an injected `ActiveSource` covering path
  resolution, symlink checks, stat, read, write, and repo-relative reporting.
- Normalization migration orchestration now uses `io/fs` contracts rather than
  importing `os`; the concrete adapter remains intentionally visible in the
  application work-normalize adapter while parent `WorkSource` methods remain
  compatible.
- Migration preview temp-directory creation, copying, and cleanup now flow
  through injected `WorkSource` capabilities; `domain/work/normalize_migrate.go`
  no longer imports `os`.
- Active-row cleanup now receives all filesystem operations, including
  `Lstat` and symlink evaluation, from `WorkSource`; the parent adapter no
  longer imports `os` for that operation.
- Closed-record normalization and date-folder reconciliation now live in the
  normalization capability package behind `ClosedSource`; the parent keeps a
  compatible `WorkSource` method while the application supplies filesystem
  mutation functions.

## Constraints

Preserve parent APIs, keep filesystem effects behind explicit seams, avoid
mechanical one-file-per-folder movement, and do not split validation solely to
reduce the parent file count.

## Help Wanted

Review whether the next architectural slice should extract shared validation
models/helpers or move on to `domain/context`.

## Suggested Next Step

Move the remaining scan compatibility helpers and legacy validation adapter's
concrete filesystem source out of `domain/work`. The remaining production
`os` imports are now limited to `normalize_scan.go` and
`validation_adapters.go`; after those are moved, verify zero production `os`
imports across all of `domain/work` before auditing `domain/kit` and
`domain/kitsync`.
