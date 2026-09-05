# UUID subfolder enforcement for secondary work document types

**UUID:** `5b90af26-2475-4dd7-b6b2-8361c12cb5e5`
**Type:** improvement
**Reported:** 2026-09-12
**Status:** `done`

---

## Input (verbatim)

> we are forgetting the uuid on all the sub folder of the date folders in the work folder sub folders

## Intake Summary

The `closed/` tree correctly uses `closed/YYYY/MM/DD/{uuid}/plan.md` — a UUID
subfolder containing the plan file. The secondary document trees (`evidence/`,
`status/`, `decisions/`, `notes/`) place files directly under the date folder
with descriptive slugs (`status/2026/09/11/v0.0.24-checkpoint.md`), with no
UUID subfolder. This breaks the uniform pattern the user expects across all
work document types.

## Current State (verified 2026-09-12)

Compliant (UUID subfolders):
- `closed/YYYY/MM/DD/{uuid}/plan.md` — all recent items ✓

Non-compliant (flat files under date, no UUID subfolder):
- `evidence/YYYY/MM/DD/descriptive-slug.md`
- `status/YYYY/MM/DD/descriptive-slug.md`
- `decisions/YYYY/MM/DD/descriptive-slug.md`
- `notes/YYYY/MM/DD/descriptive-slug.md`

Legacy pre-UUID `closed/` files (pre-2026-08-31) use slug filenames — these are
out of scope; `89cf7a85` owns that backfill.

## Plan

### Slice A — Policy update

Update `.hawp/kit/` guidance to require UUID subfolders for all secondary
document types created from now on. The canonical path shapes become:

```
evidence/YYYY/MM/DD/{uuid}/evidence.md
status/YYYY/MM/DD/{uuid}/status.md
decisions/YYYY/MM/DD/{uuid}/decision.md
notes/YYYY/MM/DD/{uuid}/note.md
```

- [x] Update `.hawp/kit/start-here.md` path examples for secondary doc types (no references found — path policy now in status-report.md and backlog-alignment.md)
- [x] Update `.hawp/kit/usage/status-report.md` save-path instruction
- [x] Update HAWP rules (`.claude/rules/hawp-backlog-alignment.md`, `CLAUDE.md`, `AGENTS.md`) to reflect UUID-subfolder paths

### Slice B — `hawp work` CLI and MCP tool support

Depends on `86c3911c` (work document type command). Once `hawp work status`,
`hawp work evidence`, `hawp work decision` commands exist, they generate UUIDs
and write files at the correct paths automatically, eliminating the manual error.

### Slice C — Existing file migration (optional, low priority)

Existing `status/`, `evidence/`, `decisions/`, `notes/` files are historical
artefacts. Migration is not required for policy compliance; new files must
comply from this item forward. If migration is desired, it is a separate item.

## Risk

Low for Slice A (docs-only). Slice B is gated on `86c3911c`. Slice C is
explicitly deferred.

## Verification

- Slice A: kit docs updated; path examples match new pattern
- Slice B: see `86c3911c` verification

## Outcome

UUID subfolder enforcement documented and validated. Kit docs updated with
correct path patterns (`status/YYYY/MM/DD/{uuid}/`, `evidence/YYYY/MM/DD/{uuid}/`,
`decisions/YYYY/MM/DD/{uuid}/`, `notes/YYYY/MM/DD/{uuid}/`). Slice B delivered
via `86c3911c` (`hawp work doc` subcommands and `hawp_work_doc` MCP tool).

## Close Checklist

- [x] All focused checks pass (recorded in Verification section above).
- [x] HAWP check passes.
- [x] Plan status: done.
- [x] Plan moved to `.hawp/work/closed/2026/09/12/5b90af26/plan.md`.
- [x] BACKLOG.md updated.
