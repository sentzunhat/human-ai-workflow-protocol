---
work-item: c1d2e3f7
type: fix
title: "Isolate kit content input from normalization and validation"
status: done
updated: 2026-08-13
created: 2026-08-10
updated: 2026-08-10
parent: b6c4e8a2
depends-on: c1d2e3f6
---

# Fix: Domain Kit Boundary

## Mission

Separate kit content acquisition and Markdown/repository concerns from pure
normalization and validation rules without changing existing CLI behavior.

## Done When

- Kit rules can be tested with in-memory content.
- Concrete repository/Markdown access is behind a capability-local boundary.
- Existing validation and normalization entry points remain compatible.
- Focused tests, build, and diff checks pass.

## Audit Handoff

The preceding audit confirmed that `domain/kit` currently mixes pure naming,
required-file, and link rules with Markdown discovery, file reads, existence
checks, renames, and writes. Add a kit-specific workspace contract under
`domain/kit/source` and a filesystem implementation under
`infrastructure/filesystem/kit`.

Stage the implementation internally: first move validation and planning to
typed file snapshots, then move rename/write application to the same adapter
boundary. Preserve command names, flags, output, dry-run behavior,
dirty-worktree protection, link rewriting, conflict handling, and errors.

Add in-memory tests for pure rules plus filesystem adapter tests for discovery,
missing files, rename conflicts, and link-update mutation.

## Implementation

- Added `domain/kit/source/port.go` with typed workspace entries, Markdown
  file snapshots, parsed links, and explicit mutation plans.
- Refactored domain kit validation and normalization to consume snapshots and
  produce plans without importing filesystem, Markdown, or repository helpers.
- Added `infrastructure/filesystem/kit/adapter.go` for discovery, fenced-code
  link extraction, read behavior, renames, and link writes.
- Updated application kit services to compose the default adapter and retain
  existing report rendering, dry-run/apply behavior, and dirty-worktree guard.
- Converted domain tests to in-memory snapshots and added filesystem adapter
  tests for discovery, fenced links, mutation, and conflict handling.

## Verification

- `go test ./internal/domain/kit ./internal/infrastructure/filesystem/kit ./internal/application/kit ./internal/application/check ./internal/platform/cli ./internal/application/search ./internal/application/index`
- `CGO_ENABLED=0 make build` from `librarian/go`
- `git diff --check`
- A package scan confirms `domain/kit` has no concrete infrastructure imports,
  filesystem calls, or Markdown parser calls.

## Outcome

The kit capability now follows the same capability-local boundary as context:
domain rules consume typed inputs, filesystem operations are isolated in the
filesystem kit adapter, and application services own use-case composition.
The next compoundable item is `c1d2e3f8`, the recursive audit of `domain/work`.

## Close Checklist

- [x] `domain/kit/source/port.go` added with typed workspace contract
- [x] `infrastructure/filesystem/kit/adapter.go` added for all I/O concerns
- [x] Domain kit validation and normalization consume snapshots, no concrete infrastructure imports
- [x] Application kit services compose the adapter; dry-run/apply/dirty-worktree preserved
- [x] In-memory domain tests and filesystem adapter tests added
- [x] Build and diff checks pass
