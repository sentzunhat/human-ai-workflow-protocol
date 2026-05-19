# TASK-019 Downstream Drift Follow-ups

## Decision

TASK-019 is accepted as complete in the HAWP repo. No Tekit or Mictlan files were changed from this repository.

## Tekit Follow-up (to execute in Tekit repo)

- BUG-156 is a modern closed file (post-2026-05-10) missing required close sections:
  - Outcome
  - Verification
  - Close Checklist
- Fix location: Tekit repository only.

## Mictlan Follow-up (to execute in Mictlan repo)

- TASK-030 appears in backlog but its active target file is missing.
- active/2026/05/06/TASK-011.md appears orphaned relative to backlog.
- TASK-033 and TASK-034 are modern closed files (post-2026-05-10) missing required close sections.
- Fix location: Mictlan repository only.

## HAWP Follow-up Item (only)

- Kind: task
- Title: Add validator CLI help output
- Scope:
  - Add --help output documenting default local behavior.
  - Add --help output documenting --hawp-root.
  - Add --help output documenting --work-root.
  - Add --help output documenting exit code behavior.
  - Add --help output documenting WARN vs FAIL meaning.
  - Do not change validation rules.
  - Do not add SQLite, indexing, search, or queueing.
  - Do not start UUID migration.

## Recommendation Order

1. Close TASK-019 in HAWP (done).
2. Add tiny CLI help task in HAWP (logged as TASK-020).
3. Run and fix Tekit drift inside Tekit.
4. Run and fix Mictlan drift inside Mictlan.
5. Revisit UUID migration or SQLite/indexing only after downstream drift is clean.
