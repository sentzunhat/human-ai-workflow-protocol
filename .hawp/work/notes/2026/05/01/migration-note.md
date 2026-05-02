# Migration Note — 2026-05-01

## Date

2026-05-01

## Summary

Migrated source repo `.work/` from flat `adrs/` + `status/` layout to the date-based structure defined in ADR-004.

## Source Paths Found

- `.work/adrs/GUARDRAIL_ADR.md`
- `.work/adrs/INSTALL_BOUNDARY_ADR.md`
- `.work/adrs/PUBLICATION_SAFETY_ADR.md`
- `.work/status/2026-04-26-hawp-adr-template-review.md`
- `.work/status/2026-04-30-bug-001-update-flow-migrate-old-layouts.md`
- `.work/status/2026-04-30-task-001-parallel-work-guardrails-adr.md`
- `.work/status/2026-04-30-task-002-intake-template-adr.md`
- `.work/status/2026-04-30-task-003-kit-work-restructure-adr.md`
- `.work/status/2026-04-30-task-004-readme-ux-refactor.md`

## Target Paths Created

- `.work/decisions/2026/05/01/GUARDRAIL_ADR.md`
- `.work/decisions/2026/05/01/INSTALL_BOUNDARY_ADR.md`
- `.work/decisions/2026/05/01/PUBLICATION_SAFETY_ADR.md`
- `.work/closed/2026/04/26/2026-04-26-hawp-adr-template-review.md`
- `.work/closed/2026/04/30/2026-04-30-bug-001-update-flow-migrate-old-layouts.md`
- `.work/closed/2026/04/30/2026-04-30-task-001-parallel-work-guardrails-adr.md`
- `.work/closed/2026/04/30/2026-04-30-task-002-intake-template-adr.md`
- `.work/closed/2026/04/30/2026-04-30-task-003-kit-work-restructure-adr.md`
- `.work/closed/2026/04/30/2026-04-30-task-004-readme-ux-refactor.md`
- `.work/decisions/2026/05/01/ADR-004-date-based-work-folder-structure.md` (new ADR)

## Files Skipped

None — all targets were empty before migration.

## New Folders Created

- `.work/active/`
- `.work/closed/`
- `.work/decisions/`
- `.work/notes/`

## Other Changes

- `BACKLOG.md` plan links updated from `status/...` to `closed/YYYY/MM/DD/...`
- `.github/copilot-instructions.md` updated to reference `.work/active/` and `.work/decisions/YYYY/MM/DD/`
- Old `adrs/` and `status/` directories removed after successful migration
- `core/install.md` and `core/update.md` updated with new folder structure and migration steps
- `core/.hawp/work/` scaffold updated to match new structure

## Manual Review

None required.
