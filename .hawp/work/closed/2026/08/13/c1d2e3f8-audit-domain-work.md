---
work-item: c1d2e3f8
type: audit
title: "Recursive audit: domain work capability"
status: done
created: 2026-08-10
updated: 2026-08-13
parent: b6c4e8a2
follow-up: c1d2e3f9
---

# Audit: Domain Work Capability

## Mission

Audit `librarian/go/internal/domain/work` for repository-layout assumptions,
link resolution, backlog parsing, normalization, and pure work invariants.

## Required Output

- Confirmed findings separated from inference.
- Capability-local port/adapter proposal.
- Compatibility and verification requirements for `c1d2e3f9`.

## Constraints

Do not rewrite work rules during the audit. Preserve existing HAWP backlog and
UUID semantics.

## Audit Evidence

### Responsibility Map

| Concern | Current implementation | Classification |
| --- | --- | --- |
| IDs, row matching, headings, evidence/verification interpretation, record classification | `domain/work` rule files | domain rules; retain in the work capability |
| backlog file read | `backlog.go`, `normalize_scan.go` | source acquisition; separate from parse rules |
| closed-plan collection, content reads, folder traversal | `completeness.go`, `evidence.go`, `clarity.go`, `consistency.go`, `deadlinks.go` | read-only workspace acquisition; move behind a work source boundary |
| Markdown fence blanking and link extraction | `deadlinks.go` via concrete Markdown infrastructure | adapter parsing concern; pass typed links to domain rules |
| plan existence and repo-relative path projection | `consistency.go`, `deadlinks.go`, `normalize_rules.go` | workspace/path concern; adapter input or application composition |
| normalization backlog/plan scans | `normalize_scan.go` | distinct scan input boundary; do after validation extraction |
| closed-record moves and content rewrites | `normalize_apply.go` | mutation boundary; do after scan extraction and preserve dry-run safeguards |
| CLI options, output, dirty-worktree gate, export behavior | `application/work/normalize.go` | application orchestration; retain there |

### Confirmed Findings

1. **Domain-to-infrastructure coupling is broader here than in context or
   kit.** `domain/work` reads files, traverses dated archive trees, resolves
   paths, checks existence, uses concrete Markdown helpers, emits stderr
   warnings, moves records, and writes content.
2. **Read-only validation has a coherent first seam.** Backlog consistency,
   completeness, evidence integrity, verification clarity, and dead-link
   reporting can consume a typed work snapshot containing rows, plan content,
   paths, and pre-parsed local links. Their rules do not need filesystem calls.
3. **Normalization is a separate lifecycle.** `normalize_scan.go`,
   `normalize_rules.go`, and `normalize_apply.go` mix scan data, canonical-path
   checks, detection, and destructive mutations. Combining this with the
   validation extraction would make one change too broad to prove safely.
4. **Backlog parsing has a pure core obscured by file I/O.** `ParseBacklog`
   and `ParseNormalizeBacklog` read a path before applying line/table parsing;
   content-based variants should be introduced before applications compose a
   source adapter.
5. **Current tests are useful but mostly filesystem fixtures.** They cover
   pass/fail validation, evidence, clarity, dead links, normalization rules,
   and closed-record application, but they do not isolate snapshot acquisition
   or mutation adapter behavior.

### Verification

- `go test ./internal/domain/work ./internal/application/work ./internal/application/check ./internal/platform/cli` passes.
- The concrete I/O scan confirms the identified read, traversal, Markdown,
  path, move, and write dependencies.

## Boundary Proposal

Use capability-local work boundaries, not generic repository interfaces:

- `domain/work/source/port.go`: typed read-only validation snapshot: backlog
  text/rows, work files, closed-plan metadata/content, and pre-parsed links.
- `infrastructure/filesystem/work/adapter.go`: backlog and plan acquisition,
  archive traversal, Markdown fence/link parsing, safe relative-path facts,
  and read-error policy.
- `application/work/validate.go`: compose the filesystem adapter and pass its
  typed snapshot to pure validation rules.
- `domain/work/normalize/source` and `infrastructure/filesystem/work` mutation
  operations: a later, distinct normalization scan/mutation extraction.

The validation adapter must preserve the current scope: dead links scan only
`BACKLOG.md`, flat `active/`, and flat `parked/`; closed records use their
existing date/legacy semantics; unreadable files warn/skip where they do now.

## Handoff

`c1d2e3f9` is implementation-ready for the read-only validation boundary:

- add content-based backlog parsers and typed validation snapshots;
- migrate validation rules first, keeping report output unchanged;
- add in-memory rule tests and filesystem adapter policy tests;
- preserve UUID matching, legacy cutoff, evidence containment, and dead-link
  scan scope.

`c1d2e402` is added as the follow-up for normalization scan/mutation
extraction. Do not fold it into `c1d2e3f9`.

## Audit Result

The recursive work audit is complete. No production behavior changed in this
audit; the two linked improvements are ordered to keep the refactor safe.
