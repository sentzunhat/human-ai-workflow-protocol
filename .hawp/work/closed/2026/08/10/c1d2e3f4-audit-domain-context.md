---
work-item: c1d2e3f4
type: audit
title: "Recursive audit: domain context capability"
status: done
created: 2026-08-10
updated: 2026-08-10
parent: b6c4e8a2
follow-up: c1d2e3f5
---

# Audit: Domain Context Capability

## Mission

Audit every file under `librarian/go/internal/domain/context` for capability
cohesion, dependency direction, I/O ownership, and test seams.

## Constraints

- Inspect repository-local code only.
- Keep ports and adapters together by capability.
- Do not split files by arbitrary line count.
- Preserve `EnrichKit` and `EnrichWork` behavior while auditing.

## Required Output

- Confirmed findings with exact repository-relative paths.
- Before/after boundary proposal.
- Tests or evidence needed to prove behavior.
- Update follow-up item `c1d2e3f5` with implementation-ready scope.

## Known Starting Signal

`context/kit.go` and `context/work.go` directly use Markdown, filesystem, and
repository helpers. Determine which logic is domain enrichment and which is
corpus acquisition.

## Audit Evidence

### Confirmed Responsibilities

| Responsibility | Current location | Boundary classification |
| --- | --- | --- |
| `Document`, `Corpus`, and enrichment metadata shape | `domain/context/document.go` | domain-owned contract; keep here |
| Kit role classification and README summary extraction | `domain/context/kit.go` | pure enrichment rules; keep with context capability |
| Work role, ID, closed-date, backlog-row matching, and prefix construction | `domain/context/work.go` | pure enrichment rules; keep with context capability |
| Markdown file discovery | `domain/context/kit.go`, `domain/context/work.go` via `infrastructure/markdown.CollectFiles` | source acquisition; move behind a context source boundary |
| File reads | `domain/context/kit.go`, `domain/context/work.go` via `os.ReadFile` | source acquisition; move behind a context source boundary |
| Repo-relative path conversion | both files via `infrastructure/repo.ToRepoRelative` | source/path adaptation; move behind a context source boundary |
| Backlog parsing | `domain/context/work.go` via `domain/work.ParseBacklog` | repository acquisition plus work-domain coupling; move behind a context source boundary |
| Composition and corpus selection | `application/index/build-service.go` | application-owned wiring; keep outside domain rules |

### Confirmed Findings

1. **Domain-to-infrastructure coupling:** `domain/context/kit.go` imports
   concrete Markdown and repository infrastructure and calls `os.ReadFile`.
   `domain/context/work.go` does the same and also imports `domain/work` for
   backlog parsing.
2. **Acquisition and enrichment are mixed:** the same functions discover
   files, read content, derive paths, resolve metadata, and construct
   `Document` values. This prevents the enrichment rules from being tested
   independently of a filesystem.
3. **Silent content loss is part of the current behavior:** kit read errors
   skip a file, and work document read errors are ignored with `raw, _`. The
   follow-up must choose and test an explicit error policy rather than
   accidentally changing it.
4. **The work corpus has a fixed folder allow-list:** `EnrichWork` scans only
   `active`, `closed`, `parked`, `decisions`, `evidence`, `notes`, and `status`.
   This is a policy decision that must remain explicit when the source adapter
   is extracted; unknown `.hawp/work` folders are currently not indexed.
5. **The public call boundary is application-facing:**
   `application/index/build-service.go` calls `EnrichKit` and `EnrichWork`
   with absolute repository paths. Compatibility wrappers are therefore
   required while the concrete source is moved behind the new seam.

### Test Coverage Observed

- `go test ./internal/domain/context ./internal/application/index` passes.
- Existing tests prove role/prefix behavior, backlog metadata resolution,
  closed-date extraction, ID parsing, scope filtering, and JSON export.
- No focused test currently injects an in-memory corpus source.
- No focused test currently proves the read-error policy or unknown-folder
  policy.

## Boundary Proposal For `c1d2e3f5`

Use a context-specific source contract, not a generic repository interface:

- `domain/context/source/port.go`: typed source contract for kit/work input,
  kept beside the context capability that consumes it.
- `infrastructure/filesystem/context/adapter.go`: filesystem-backed adapter
  that owns Markdown discovery, reads, repo-relative paths, and backlog file
  acquisition for the context corpus.
- `application/index/build-service.go`: composition point that wires the
  filesystem adapter and calls context enrichment through the contract.
- `domain/context/kit.go` and `work.go`: retain pure enrichment functions and
  compatibility wrappers until callers migrate.

The port and adapter remain grouped by capability, with no global `ports/` or
`adapters/` bucket. The first implementation should preserve the current
skip-on-read-error and fixed-folder behavior unless a focused decision changes
those contracts.

## Handoff To `c1d2e3f5`

`c1d2e3f5` is now implementation-ready. Its first slice should introduce the
typed source seam, wire the default filesystem adapter in the build service,
retain `EnrichKit`/`EnrichWork` compatibility behavior, and add tests for
in-memory source input, read-error behavior, and ignored-folder behavior.

## Audit Result

The audit is complete for the inspected `domain/context` capability. The
follow-up implementation item is the next compoundable item; this audit does
not modify production code.

## Outcome

Completed full capability audit of `domain/context`. Identified five confirmed
findings: domain-to-infrastructure coupling, mixed acquisition and enrichment,
silent content loss, fixed folder allow-list policy, and application-facing
call boundary. Produced a concrete boundary proposal and handoff spec for
`c1d2e3f5`. No production code modified.

## Verification

- `go test ./internal/domain/context ./internal/application/index` passed.
- Existing tests confirm role/prefix behavior, backlog metadata resolution,
  closed-date extraction, ID parsing, scope filtering, and JSON export.
- Boundary proposal reviewed and accepted as implementation-ready handoff.

## Close Checklist

- [x] All confirmed findings documented with file paths
- [x] Boundary proposal produced for `c1d2e3f5`
- [x] No production code changed
- [x] Follow-up item `c1d2e3f5` updated with implementation-ready scope
- [x] Moved to `closed/2026/08/10/`
