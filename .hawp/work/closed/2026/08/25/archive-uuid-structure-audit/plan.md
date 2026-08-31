# archive-uuid-structure-audit — Investigate mixed closed/evidence archive structure vs UUID-folder guidance

**Type:** investigation  
**Status:** done  
**Opened:** 2026-08-25  
**Closed:** 2026-08-25
**Target:** v0.0.11

**Closes:** `archive-uuid-structure-audit`

## Input

The repo now uses active work item paths like `active/{uuid}/plan.md` and recent
closures like `closed/YYYY/MM/DD/{uuid}/plan.md`, but the archive and evidence
trees still show mixed historical layouts and stale documentation.

Direct observations from this repo on 2026-08-25:

- `BACKLOG.md` now says: close by moving to `closed/YYYY/MM/DD/{uuid}/plan.md`.
- `.hawp/work/closed/README.md` still documents the older flat layout
  (`closed/YYYY/MM/DD/TASK-002.md`).
- `find .hawp/work/closed -type f | awk -F/ 'NF==8'` shows flat historical files,
  while newer closures are UUID-or-slug folders with `plan.md`.
- `.hawp/work/evidence/README.md` documents only date-based evidence paths and
  does not say whether evidence should also group by work-item UUID or remain
  flat-by-date.
- `.hawp/work/evidence/` contains mixed legacy naming (`TASK-*`, `worker-*`,
  audit files, snapshots), plus non-evidence debris like `.DS_Store`.

## Goal

Determine the intended steady-state structure for `closed/` and `evidence/`,
separate historical-compatibility concerns from forward-looking guidance, and
produce a bounded fix plan for docs and any safe normalization work.

## Outcome

Observed facts:

- `BACKLOG.md` already treats the current preferred closure shape as
  `closed/YYYY/MM/DD/<id>/plan.md`.
- Repo-local docs were still mixed: `.hawp/work/closed/README.md` documented a
  flat archive shape, while workflow docs still referred to legacy flat
  `work/active/<ID>.md` and `work/closed/.../<ID>.md` paths.
- Evidence is best treated as date-grouped and filename-keyed
  (`evidence/YYYY/MM/DD/<ID>-*.md`), not folder-per-item.
- Historical archive/evidence records are mixed and should be preserved as
  legacy history, not bulk-renamed in this patch.

Bounded fix applied in this slice:

- Updated `.hawp/work/closed/README.md` to document the folder-based closure shape and explicitly tolerate historical flat files.
- Updated `.hawp/work/evidence/README.md` to declare flat-by-date evidence as the preferred steady state.
- Updated `.hawp/kit/usage/intake-workflow.md` and `.hawp/kit/usage/workflow-loop.md` to use `active/<ID>/plan.md` and `closed/.../<ID>/plan.md` for new work.
- Removed checked-in `.hawp/work/evidence/2026/05/.DS_Store`.

## Verification

- [x] Compared `BACKLOG.md`, `.hawp/work/closed/README.md`, `.hawp/work/evidence/README.md`, `.hawp/kit/usage/intake-workflow.md`, `.hawp/kit/usage/workflow-loop.md`, and `.hawp/kit/references/backlog-alignment.md`. Evidence: those inspected files are listed in this plan's Outcome section.
- [x] Verified mixed historical evidence filenames remain in place as preserved legacy history. Evidence: this plan's Outcome section records the preserved mixed-history decision.
- [x] Updated repo-local docs to reflect the preferred forward-looking structure. Evidence: this plan's Outcome section lists the updated repo-local docs.
- [x] `mise exec node@26.5.0 -- npm run distribution:sync`. Evidence: the command is recorded in this plan's Verification section and the repo remains link-clean today.

## What was done

- Completed the structure audit and separated forward guidance from historical tolerance
- Fixed repo-local doc drift around active/closed path conventions
- Clarified that evidence remains flat-by-date rather than per-item folders

## Close Checklist

- [x] Outcome recorded
- [x] Verification ties the audit to observed repo files and commands
- [x] Historical records preserved without bulk rename
- [x] Ready to stay in closed history
