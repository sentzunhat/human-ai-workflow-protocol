---
work-item: c1d2e3f6
type: audit
title: "Recursive audit: domain kit capability"
status: done
created: 2026-08-10
updated: 2026-08-11
parent: b6c4e8a2
follow-up: c1d2e3f7
---

# Audit: Domain Kit Capability

## Mission

Audit `librarian/go/internal/domain/kit` and its nested tests for separation
between pure kit rules, Markdown parsing, repository access, and application
orchestration.

## Required Output

- Confirmed dependency and cohesion findings with paths.
- Smallest capability-local port/adapter boundary.
- Test and compatibility requirements for `c1d2e3f7`.

## Constraints

Keep ports and adapters beside their owning capability. Do not perform the
implementation in this audit item or broaden into unrelated kit tooling.

## Audit Evidence

### Responsibility Map

| Responsibility | Current location | Classification |
| --- | --- | --- |
| filename normalization and allowed-name policy | `domain/kit/normalize.go` | pure kit rule; keep in domain kit |
| required-file policy | `domain/kit/validate.go` | pure kit rule; keep in domain kit |
| link target and fenced-code rules | `domain/kit/validate.go`, `normalize.go` | pure rule when given typed file input; keep in domain kit |
| directory traversal and file reads | `domain/kit/normalize.go`, `validate.go` via `os.ReadDir`/`os.ReadFile` | source acquisition; move behind kit workspace boundary |
| Markdown scanning and link extraction | both files via `infrastructure/markdown` | external parsing utility; expose typed content to domain rules |
| path existence and repo-relative display | `validate.go` and application normalize via `infrastructure/repo` | filesystem/path adapter concern |
| rename and link-update writes | `domain/kit/normalize.go` via `os.Rename`/`os.WriteFile` | mutation adapter concern; remove from domain |
| dry-run/apply reporting and dirty-worktree guard | `application/kit/normalize.go` | application orchestration; keep here |

### Confirmed Findings

1. **Domain-to-infrastructure coupling:** both domain kit files import
   concrete Markdown/repository infrastructure. Validation and normalization
   also call filesystem APIs directly.
2. **Read-only and mutation responsibilities are mixed:** the same domain
   capability validates content, discovers files, plans changes, renames files,
   and writes link updates.
3. **Application orchestration is correctly placed:** `application/kit` owns
   output, dry-run versus apply mode, and the dirty-worktree safety gate. It
   currently delegates path-based operations to a domain package that owns too
   much I/O.
4. **The current domain tests are filesystem fixtures:** they prove useful
   end-to-end behavior, but do not prove validation or planning independently
   of a filesystem source.
5. **Current compatibility surface is internal:** CLI and composite check call
   application kit services, while application normalize calls the domain
   functions. The follow-up can migrate these internal calls without changing
   command names, flags, output, or safety behavior.

### Verification

- `go test ./internal/domain/kit ./internal/application/kit ./internal/application/check ./internal/platform/cli` passes.
- The inspected domain package still contains concrete infrastructure imports
  and filesystem calls, confirming the finding.
- Existing tests cover clean validation, naming, required files, internal
  links, planned renames, link rewrites, and rename conflicts.
- No focused in-memory snapshot tests or adapter mutation tests exist yet.

## Boundary Proposal For `c1d2e3f7`

Use a kit-specific workspace contract, not a generic repository interface:

- `domain/kit/source/port.go`: typed kit file/snapshot input plus explicit
  workspace operations needed by normalization.
- `infrastructure/filesystem/kit/adapter.go`: Markdown discovery, reads,
  existence checks, repo-relative paths, renames, and writes.
- `application/kit/validate.go` and `normalize.go`: compose the adapter,
  retain rendering, dry-run behavior, dirty-worktree protection, and command
  compatibility.
- `domain/kit`: accept typed content/path facts and retain pure naming,
  required-file, link, and rename-planning rules.

Keep the port beside the kit capability and the adapter beside the filesystem
kit capability. Do not introduce global `ports/` or `adapters/` folders.

The implementation should be staged internally: first move read-only
validation/planning to typed snapshots, then move rename/write application to
the same adapter boundary without changing the CLI contract.

## Handoff To `c1d2e3f7`

`c1d2e3f7` is implementation-ready. It must add in-memory tests for pure
rules, filesystem adapter tests for discovery and mutation/conflict behavior,
and preserve the existing dry-run/apply, dirty-worktree, link-rewrite, and
error behavior.

## Outcome

The recursive audit is complete for `domain/kit`. This audit did not modify
production code; its linked implementation item (`c1d2e3f7`) is the next
compoundable slice.

## Close Checklist

- [x] Findings confirmed with direct evidence (responsibility map + test run)
- [x] Boundary proposal written and handed off to `c1d2e3f7`
- [x] No production code changed in this audit
- [x] Linked implementation item exists and is implementation-ready
