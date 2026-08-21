---
work-item: c1d2e3f9
type: fix
title: "Extract work source and link-resolution boundaries"
status: in-progress
created: 2026-08-10
updated: 2026-08-21
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
- Migrated backlog consistency, evidence integrity, and verification clarity
  validators onto the typed snapshot. Added `ClosedFiles []File` to
  `source.Snapshot`; the filesystem adapter classifies and pre-loads
  `closed/YYYY/MM/DD/*.md` files into it.
- `CheckBacklogConsistency` now scans `snapshot.Files` RelPaths instead of
  calling `os.ReadDir`/`os.Stat`; no longer accepts a `workDir` string.
- `CheckEvidenceIntegrity` now accepts the full snapshot; checks evidence
  path existence via `snapshot.Files` RelPath scan rather than `os.Stat`.
- `CheckVerificationClarity` now accepts `[]source.File` and reads content
  from `file.Content` rather than `os.ReadFile`.
- `CollectClosedPlanFiles` removed from the domain; the adapter populates
  `snapshot.ClosedFiles` directly.
- Domain tests rewritten as pure in-memory (non-repository) snapshot tests;
  repository path covered by the new `TestAdapterReadPopulatesClosedFiles`
  infrastructure test.
- `BuildResearchQueue` in `normalize_apply.go` updated via a local
  `closedFilesFromPaths` adapter to compile against the new API; normalization
  mutation logic is otherwise untouched (deferred to `c1d2e402`).
- Build clean; all domain/work, application/work, filesystem, and CLI tests
  pass.

## Next Slice

No remaining slices for this item. Normalization scan/mutation is tracked in
`c1d2e402`.
