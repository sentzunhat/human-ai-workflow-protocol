# TASK-013 Phase 1 Verification Evidence — 2026-07-03

**UUID:** `361fb08e-6457-4ed5-80bd-76337b6f0e89`

## Direct evidence

- `npm --prefix librarian run typecheck` — clean.
- `npm --prefix librarian test` — 49/49 pass (3 new UUID tests):
  - `parseBacklog falls back to the UUID cell when Legacy ID is a placeholder`
  - `extractIdFromFilename recognizes legacy, date-prefixed, and UUID formats`
  - `checkBacklogConsistency resolves UUID-named active plan files without orphans`
- `npm --prefix librarian run validate` — full suite PASS (typecheck, tests,
  markdown links, kit:validate, distribution:sync, work:validate 4 checks / 0 issues)
  against the real backlog with the UUID row present in Active Work.

## Code changed

- `librarian/scripts/hawp/work-validate/validations/id-parser.ts` — UUID v4 pattern,
  lowercase canonical output; short 8-char forms deliberately rejected (Phase 2).
- `librarian/scripts/hawp/work-validate/orchestrate.ts` — row ID resolution tries
  Legacy ID → ID → UUID cells, skipping `—` placeholders, preferring parseable values.

## Docs changed

- `.hawp/kit/templates/intake-plan.md`, `.hawp/kit/templates/work-intake.md` — dual header.
- `.hawp/kit/usage/intake-workflow.md` — "Work item IDs" section.
- `.hawp/work/notes/2026/07/03/migration-sequential-to-uuid.md` — phase table + mapping.

## Dogfood

This task closed as `closed/2026/07/03/361fb08e-6457-4ed5-80bd-76337b6f0e89.md` with the
full UUID in the Recently Closed ID cell — the first UUID-named record validated in place.
