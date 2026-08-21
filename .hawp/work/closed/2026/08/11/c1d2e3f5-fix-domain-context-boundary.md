---
work-item: c1d2e3f5
type: fix
title: "Extract domain context corpus/source boundary"
status: done
created: 2026-08-10
updated: 2026-08-11
parent: b6c4e8a2
depends-on: c1d2e3f4
---

# Fix: Domain Context Boundary

## Mission

Move repository file walking and reading behind a capability-local context
corpus/source boundary while keeping context enrichment rules in the domain
context group.

## Constraints

- Keep the port beside the context capability and the concrete adapter beside
  its filesystem capability.
- Preserve repository-facing build behavior while moving internal enrichment
  to typed corpus input.
- Do not change search output or repository layout.

## Done When

- Domain context no longer imports concrete Markdown/repository infrastructure
  for corpus acquisition.
- Focused context tests cover kit/work enrichment and missing-file behavior.
- Targeted Go tests, build, and diff checks pass.

## Audit Handoff

The preceding audit confirmed that context enrichment currently combines pure
role/metadata rules with Markdown discovery, file reads, repo-relative path
conversion, and backlog loading. Preserve the current skip-on-read-error and
fixed work-folder behavior in the first slice unless a separate decision is
recorded.

Recommended capability grouping:

- `domain/context/source/port.go` for the typed context source contract.
- `infrastructure/filesystem/context/adapter.go` for filesystem acquisition.
- `application/index/build-service.go` for default adapter composition.

The internal enrichment functions are migrated directly to typed corpus input
because they are repository-internal APIs. The application build service keeps
the existing repository-root and scope behavior. Add tests that prove
in-memory source behavior and the current missing-file and ignored-folder
policies.

## Implementation

- Added `domain/context/source/port.go` with typed kit/work corpus inputs and a
  `CorpusSource` acquisition contract.
- Moved kit/work enrichment to typed corpus inputs so the domain context
  package owns role, metadata, and prefix rules without filesystem imports.
- Added `infrastructure/filesystem/context/adapter.go` for Markdown discovery,
  reads, repo-relative paths, and backlog acquisition.
- Wired the default adapter in `application/index/build-service.go` and added
  `NewBuildServiceWithSource` for focused source-level tests.
- Preserved README exclusion, fixed work-folder allow-list, and skip-on-read-
  error behavior in the filesystem adapter.
- Converted context tests to in-memory typed corpus fixtures and added adapter
  policy coverage.

The internal enrichment API now accepts typed corpus input rather than paths;
the application build service preserves the previous repository-facing build
behavior while owning the adapter composition.

## Verification

- `go test ./internal/domain/context ./internal/infrastructure/filesystem/context ./internal/application/index`
- `go test ./internal/platform/cli ./internal/application/search`
- `make build` from `librarian/go` with `CGO_ENABLED=0`
- `git diff --check`
- `rg` confirms no infrastructure imports or filesystem acquisition calls in
  `domain/context`.

## Outcome

Introduced `domain/context/source/port.go` with typed kit/work corpus contracts
and moved all filesystem acquisition behind `infrastructure/filesystem/context/adapter.go`.
Domain context now has no direct infrastructure imports. Application build service
wires the default adapter; `NewBuildServiceWithSource` enables focused in-memory
testing. Existing repository-facing behavior, README exclusion, fixed work-folder
allow-list, and skip-on-read-error policy all preserved.

## Result

The context corpus/source boundary is implemented and verified. The next
compoundable work item is the recursive audit of `domain/kit` (`c1d2e3f6`).

## Close Checklist

- [x] Domain context imports no concrete Markdown/repository infrastructure
- [x] Typed corpus source contract added in `domain/context/source/`
- [x] Filesystem adapter added in `infrastructure/filesystem/context/`
- [x] Default adapter wired in `application/index/build-service.go`
- [x] In-memory source tests and policy coverage added
- [x] All verification commands passed
- [x] Moved to `closed/2026/08/11/`
