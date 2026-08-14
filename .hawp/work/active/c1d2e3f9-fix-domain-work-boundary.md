---
work-item: c1d2e3f9
type: fix
title: "Extract work source and link-resolution boundaries"
status: in-progress
created: 2026-08-10
updated: 2026-08-13
parent: b6c4e8a2
depends-on: c1d2e3f8
---

# Fix: Domain Work Boundary

## Mission

Move repository traversal, Markdown access, and link resolution behind
capability-local contracts while leaving work validation and normalization
rules focused on work data.

## Done When

- Pure work rules accept testable inputs without concrete repository access.
- Existing validation, normalization, and backlog behavior is preserved.
- Focused tests cover repository and non-repository paths separately.

## Audit Handoff

Start with read-only validation only. Introduce a typed work validation
snapshot under `domain/work/source` and a filesystem work adapter under
`infrastructure/filesystem/work`; move backlog reads, closed-plan acquisition,
and Markdown link extraction there. Add content-based backlog parsers so
domain rules no longer read `BACKLOG.md` directly.

Preserve UUID matching, legacy cutoff behavior, evidence containment, and the
existing dead-link scope (`BACKLOG.md`, flat `active/`, flat `parked/`). Keep
normalization scans, record moves, and writes out of this item; those are
tracked separately in `c1d2e402`.

## Progress

- Added `ParseBacklogContent`, separating HAWP backlog semantics from the
  existing filesystem-backed `ParseBacklog` wrapper.
- Added a direct content-parser test; existing callers retain the same path
  API and behavior.
- Added `domain/work/source` and the capability-local filesystem work adapter.
  Application validation now acquires a typed snapshot and parses the backlog
  from snapshot content rather than opening `BACKLOG.md` in the domain path.
- Verified the domain, application, filesystem adapter, composite check, and
  CLI packages after the composition change.
- Verified domain work, application work, composite check, and CLI tests.

## Next Slice

Migrate backlog consistency, closed-plan checks, evidence, clarity, and
dead-link validation onto the existing snapshot. Normalization scan/mutation
remains explicitly deferred to `c1d2e402`.
