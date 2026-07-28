# TASK-083 — Implement kit:validate script

**Backlog ID:** TASK-083

**Type:** feature  
**Status:** done  
**Created:** 2026-07-02  
**Closed:** 2026-07-02

## Goal

New TypeScript script at `librarian/scripts/kit-validate/` that validates the `.hawp/kit/` folder structure. Follows the three-file boundary pattern (index.ts / cli.ts / script.ts). Wired as `kit:validate` npm script.

## What it checks

1. **File naming** — all files under `.hawp/kit/` must be lowercase with hyphens (no spaces, no underscores, no uppercase except `README.md`)
2. **Required files** — key files must exist: `start-here.md`, `usage/status-report.md`, `usage/intake-workflow.md`, `usage/init.md`
3. **Broken internal links** — markdown links within kit files that point to non-existent kit files (relative links only)

## Steps

- [x] Open TASK-083 in BACKLOG and plan file
- [x] Create `librarian/scripts/kit-validate/index.ts`
- [x] Create `librarian/scripts/kit-validate/cli.ts`
- [x] Create `librarian/scripts/kit-validate/script.ts`
- [x] Create `librarian/scripts/kit-validate/validations/file-naming.ts`
- [x] Create `librarian/scripts/kit-validate/validations/required-files.ts`
- [x] Create `librarian/scripts/kit-validate/validations/internal-links.ts`
- [x] Add `__tests__/kit-validate.test.ts`
- [x] Wire `kit:validate` in `librarian/package.json`
- [x] Update `librarian/scripts/README.md` domain map
- [x] Update `CLAUDE.md`
- [x] Run `npm --prefix librarian run validate`
- [x] Close TASK-083

## Outcome

Script implemented across 7 files (index, cli, script + 3 validations + 1 test). `npm --prefix librarian run kit:validate` passes typecheck and all 7 unit tests. Running against the actual `.hawp/kit/` surfaces 7 real issues (1 file naming violation, 6 broken links) — tracked as TASK-085.

## Verification


### Evidence Follow-Up

- [ ] Research evidence for: `npm --prefix librarian run typecheck` — clean
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: `npm --prefix librarian run validate` — PASS, 0 issues
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: `npm --prefix librarian run kit:validate` — runs and reports 7 real kit issues (expected; tracked se
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.

- [x] `npm --prefix librarian run typecheck` — clean — unproven: evidence not recorded at close (annotated 2026-07-20)
- [x] `npm --prefix librarian run validate` — PASS, 0 issues — unproven: evidence not recorded at close (annotated 2026-07-20)
- [x] `npm --prefix librarian run kit:validate` — runs and reports 7 real kit issues (expected; tracked separately) — unproven: evidence not recorded at close (annotated 2026-07-20)

## Close Checklist

- [x] All steps completed
- [x] Validation passes
- [x] Backlog updated
- [x] Plan file moved to `closed/2026/07/02/`
