# TASK-085 — Fix 7 kit issues found by kit:validate

**Backlog ID:** TASK-085

**Type:** bug  
**Status:** done  
**Created:** 2026-07-02  
**Closed:** 2026-07-02

## Goal

Fix 7 issues found by `npm --prefix librarian run kit:validate`: 1 bad filename and 6 broken links.

## Steps

- [x] Fix bug in `internal-links.ts` — links inside fenced code blocks were being checked; now blanked before scanning
- [x] Rename `standards/public/templates/ADR.template.md` → `adr-template.md`
- [x] Fix `standards/public/guidelines/README.md` — removed dead `security.md` link, updated `adr-template.md` link to `../templates/adr-template.md`
- [x] Fix `standards/public/standards/docs/hawp-install-update-safety.md` — corrected link from `../../../../docs/` to `../../../docs/`
- [x] Fix `standards/public/standards/nodejs/git-workflow.md` — corrected relative path from `../guidelines/` to `../../guidelines/`
- [x] Run `npm --prefix librarian run kit:validate` — 0 issues
- [x] Close TASK-085

## Outcome

All 7 issues resolved. `kit:validate` now reports: ✓ 3 checks passed, 0 issues.

## Verification


### Evidence Follow-Up

- [ ] Research evidence for: `npm --prefix librarian run kit:validate` → 0 issues
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: `npm --prefix librarian run validate` → PASS
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.

- [x] `npm --prefix librarian run kit:validate` → 0 issues
- [x] `npm --prefix librarian run validate` → PASS

## Close Checklist

- [x] All steps completed
- [x] Validation passes
- [x] Backlog updated
- [x] Plan file moved to `closed/2026/07/02/`
