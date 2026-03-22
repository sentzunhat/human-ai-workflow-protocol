# Task: Move misplaced status review files out of active work

**Backlog ID:** TASK-015
**Type:** task
**Reported:** 2026-05-10
**Risk Level:** low

---

### Input (what was reported)

> TASK-014 validation found that `work/active/2026/` contains 5 dated status review files. These are status/review artifacts, not active work plan files, and must be moved so `active/` contains only active work item plan files.

---

### Analysis

**Files found under `work/active/2026/`:**

| File | Date | Content type |
|------|------|-------------|
| `active/2026/05/02/hawp-self-update-review-summary.md` | 2026-05-02 | Status review summary |
| `active/2026/05/02/manager-final-verification-hawp-alignment.md` | 2026-05-02 | Manager verification review |
| `active/2026/05/02/manager-review-backlog-alignment-kit-reference.md` | 2026-05-02 | Manager review |
| `active/2026/05/02/manager-review-three-hawp-work-items.md` | 2026-05-02 | Manager review |
| `active/2026/05/03/remove-legacy-root-work-cleanup.md` | 2026-05-03 | Status summary |

All five are status/review documents. None match the plan file naming pattern (TASK-NNN.md, BUG-NNN.md). Correct destination: `work/status/YYYY/MM/DD/`.

**Validator impact:** The orphaned file detector in `backlog-consistency.ts` only flags `.md` files whose names match the ID pattern (`extractIdFromFilename`). These files have no ID prefix so they were silently ignored — the validator did not flag them.

**Note:** `hasSection()` prefix-match fix also applied as part of this session (prerequisite).

---

### Scope

- Move 5 files from `active/2026/` to `status/2026/` at the same YYYY/MM/DD path
- Scan for any links pointing to the old paths and update them
- No content edits to moved files

---

## Outcome (filled at close)

1. **`hasSection()` prefix-match fix** — `closed-task-completeness.ts` now uses `startsWith(header)` instead of `=== header`. TASK-012 and TASK-014 no longer appear as false-positives. Verification Clarity rose from 9/9 to 15/15.

2. **5 status/review files moved:**

| Old path | New path |
|----------|----------|
| `active/2026/05/02/hawp-self-update-review-summary.md` | `status/2026/05/02/hawp-self-update-review-summary.md` |
| `active/2026/05/02/manager-final-verification-hawp-alignment.md` | `status/2026/05/02/manager-final-verification-hawp-alignment.md` |
| `active/2026/05/02/manager-review-backlog-alignment-kit-reference.md` | `status/2026/05/02/manager-review-backlog-alignment-kit-reference.md` |
| `active/2026/05/02/manager-review-three-hawp-work-items.md` | `status/2026/05/02/manager-review-three-hawp-work-items.md` |
| `active/2026/05/03/remove-legacy-root-work-cleanup.md` | `status/2026/05/03/remove-legacy-root-work-cleanup.md` |

3. **Empty directories removed:** `active/2026/05/02/`, `active/2026/05/03/`, `active/2026/05/`, `active/2026/`

4. **No stale links found** — grep for `active/2026/` across `.hawp/` returned no matches.

5. **TASK-016 logged** to backlog for legacy tolerance mode implementation.

---

## Verification (filled at close)

- [x] `hasSection()` uses `startsWith` — **Evidence:** validator output shows TASK-012/TASK-014 absent from missing list; Proven count 15/15
- [x] All 5 files at correct destination — **Evidence:** `find .hawp/work/status/2026 -name "*.md"` returns all 5
- [x] `active/2026/` directory gone — **Evidence:** `ls .hawp/work/active/` shows only README.md and TASK-015.md, TASK-016.md
- [x] No stale links — **Evidence:** `grep -r "active/2026/"` returned no results
- [x] `npm run typecheck` 0 errors — **Evidence:** see TASK-015-validation-output.md
- [x] Backlog Consistency PASS, Parked PASS — **Evidence:** validator output section 1 all green

---

## Close Checklist

- [x] All 5 files moved from `active/2026/` to `status/2026/`
- [x] `active/2026/` directory empty/removed
- [x] No stale links to old paths remain
- [x] `hasSection()` prefix-match fix applied
- [x] `npm run typecheck` passes
- [x] `npm run validate:workflow` passes backlog consistency check
- [x] Evidence saved
- [x] TASK-015 closed and moved to `closed/2026/05/10/`
- [x] BACKLOG.md updated

---

### Status

- [x] Plan written
- [x] Approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
