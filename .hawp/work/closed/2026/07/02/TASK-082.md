# TASK-082 — Rename script folders and npm commands to work: and hawp: domains

**Backlog ID:** TASK-082

**Type:** improvement  
**Status:** done  
**Created:** 2026-07-02  
**Closed:** 2026-07-02

## Goal

Rename the three librarian script folders and their npm commands so the domain names match what each script actually acts on (`work:` for workflow records, `hawp:` for the combined check). No new logic — pure renames and reference updates.

## Mapping

| Old folder | New folder | Old command | New command |
|---|---|---|---|
| `validate-hawp-workflow/` | `work-validate/` | `workflow:validate` | `work:validate` |
| `backlog-upgrade/` | `work-normalize/` | `workflow:normalize` | `work:normalize` |
| `backlog-validate/` | `hawp-check/` | `hawp:check` | `hawp:check` (no change) |

## Steps

- [x] Open TASK-082 in BACKLOG and plan file
- [x] `git mv` the three folders
- [x] Update `librarian/package.json` (script names and tsx paths)
- [x] Update `.hawp/bin/hawp` (bash paths)
- [x] Update internal script references (script.ts section headers, etc.)
- [x] Update `librarian/scripts/README.md` domain map
- [x] Update `librarian/README.md`
- [x] Update `CLAUDE.md`
- [x] Run `npm --prefix librarian run validate` — must pass
- [x] Close TASK-082

## Constraints

- Do NOT update closed work items, status archives, or evidence files — those are historical records.
- Do NOT change logic inside the scripts — renames only.
- Live references to update: package.json, .hawp/bin/hawp, CLAUDE.md, librarian/README.md, librarian/scripts/README.md, any internal cross-script references.

## Outcome

All three script folders renamed. Two internal cross-references fixed (`hawp-check/script.ts` path to `work-validate/index.ts`; `work-normalize/script.ts` import from `work-validate/orchestrate`). All npm scripts updated. `validate` suite passes: 3 checks, 0 issues.

## Verification


### Evidence Follow-Up

- [ ] Research evidence for: `npm --prefix librarian run validate` → VALIDATION PASS, 0 issues, 1 legacy warning (pre-existing)
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: `ls librarian/scripts/` shows `work-validate/`, `work-normalize/`, `hawp-check/`
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: `npm --prefix librarian run work:validate` resolves correctly (typecheck passes)
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.

- [x] `npm --prefix librarian run validate` → VALIDATION PASS, 0 issues, 1 legacy warning (pre-existing) — unproven: evidence not recorded at close (annotated 2026-07-20)
- [x] `ls librarian/scripts/` shows `work-validate/`, `work-normalize/`, `hawp-check/` — unproven: evidence not recorded at close (annotated 2026-07-20)
- [x] `npm --prefix librarian run work:validate` resolves correctly (typecheck passes) — unproven: evidence not recorded at close (annotated 2026-07-20)

## Close Checklist

- [x] All steps completed
- [x] Validation passes
- [x] Backlog updated (TASK-082 moved to Recently Closed)
- [x] Plan file moved to `closed/2026/07/02/`
