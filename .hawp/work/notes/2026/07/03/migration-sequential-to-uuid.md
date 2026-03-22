# Migration: sequential IDs → UUID-based work item IDs

**Work item:** TASK-013 / `361fb08e-6457-4ed5-80bd-76337b6f0e89`
**Phase 1 completed:** 2026-07-03

## Why

Sequential IDs (TASK-013, BUG-005) collide when agents create work items in parallel.
UUIDs remove the coordination bottleneck; the type label moves to the Type column.

## Phase plan (Option B — phased)

| Phase | Scope | Status |
|---|---|---|
| 1 | Dual format: validator accepts UUIDs, templates carry a UUID field, new items UUID-native | done (2026-07-03) |
| 2 | Existing active/parked items migrate as they close; short-UUID display + prefix matching in validator | done (2026-07-03, `0e1c4afa-9668-4d61-b5b6-1e27be42ca23`) |
| 3 (optional) | Retroactive migration of closed records to UUID paths | not planned |

## What Phase 1 changed

- `librarian/scripts/hawp/work-validate/validations/id-parser.ts` — recognizes full
  UUID v4 (case-insensitive input, lowercase output) alongside `TASK-NNN`/`BUG-NNN`
  and date-prefixed forms.
- `librarian/scripts/hawp/work-validate/orchestrate.ts` — backlog row IDs resolve
  from Legacy ID first, then the UUID cell; placeholder cells (`—`) are skipped.
- `.hawp/kit/templates/intake-plan.md`, `.hawp/kit/templates/work-intake.md` — header
  now carries `**Backlog ID (Legacy):**` plus `**UUID:**`.
- `.hawp/kit/usage/intake-workflow.md` — "Work item IDs" section documents the dual
  format and the rules for new vs. existing items.

## Rules going forward

- New items: UUID-native (`active/<uuid>.md`, full UUID in the UUID cell, Legacy ID `—`).
- Existing items: keep sequential IDs until closed; no retroactive renames.
- Short display forms (first 8 chars) are cosmetic only until Phase 2 lands prefix matching.

## Legacy ID → UUID mapping

| Legacy ID | UUID |
|---|---|
| TASK-013 | `361fb08e-6457-4ed5-80bd-76337b6f0e89` |

(Extend this table as existing items are assigned UUIDs.)

Generate new UUIDs with `./.hawp/bin/hawp uuid` (or `hawp uuid --short` for the display form).
