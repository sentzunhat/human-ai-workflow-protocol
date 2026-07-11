# TASK-084 — Implement kit:normalize script

**Backlog ID:** TASK-084

**Type:** feature  
**Status:** done  
**Created:** 2026-07-02  
**Closed:** 2026-07-02

## Goal

New TypeScript script at `librarian/scripts/kit-normalize/` that fixes structural drift in `.hawp/kit/`. Defaults to dry-run; applies with `--apply`. Wired as `kit:normalize` npm script.

## What it fixes

1. **File naming** — renames files that don't follow lowercase-hyphen convention
2. **Internal links** — updates links in other kit files that pointed to renamed files

## Steps

- [x] Open TASK-084 in BACKLOG and plan file
- [x] Create `index.ts`, `cli.ts`, `script.ts`
- [x] Create `mutations/file-naming.ts` and `mutations/internal-links.ts`
- [x] Add `__tests__/kit-normalize.test.ts` (dry-run + apply tests)
- [x] Wire `kit:normalize` in `librarian/package.json`
- [x] Update `librarian/scripts/README.md` domain map
- [x] Run `npm --prefix librarian run validate`
- [x] Close TASK-084

## Outcome

Script implemented with `planFileRenames` + `planLinkUpdates` + `applyLinkUpdates`. Link updates apply patches in reverse-index order to preserve offsets. Test required committing files to a temp git repo before calling `--apply` (dirty-tree guard). Full validate: 4 checks passed, 0 issues.

## Verification


### Evidence Follow-Up

- [ ] Research evidence for: `npm --prefix librarian run typecheck` — clean
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: `npm --prefix librarian run validate` — PASS, 0 issues
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: `npm --prefix librarian run kit:normalize` — reports no normalization needed (kit is already clean)
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: All kit-normalize tests pass
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.

- [x] `npm --prefix librarian run typecheck` — clean
- [x] `npm --prefix librarian run validate` — PASS, 0 issues
- [x] `npm --prefix librarian run kit:normalize` — reports no normalization needed (kit is already clean)
- [x] All kit-normalize tests pass

## Close Checklist

- [x] All steps completed
- [x] Validation passes
- [x] Backlog updated
- [x] Plan file moved to `closed/2026/07/02/`
