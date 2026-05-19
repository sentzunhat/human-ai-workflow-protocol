# Task: Align parked-work folder rules and validator checks

**Backlog ID:** TASK-014
**Type:** task
**Reported:** 2026-05-10
**Risk Level:** low

---

### Input (what was reported)

> TASK-012 added a TypeScript validator. During UUID planning, TASK-013 was parked, exposing a drift point: parked items may live in `work/active/` while BACKLOG.md shows them as `parked`. Clarify and enforce the parked-work rule, update docs, and add a validator check.

---

### Context

After TASK-012 closed, TASK-013 was correctly placed in `work/parked/` and BACKLOG.md links to `parked/TASK-013.md`. However:

1. `backlog-alignment.md` Folder Model does not list `work/parked/` — it is an undocumented folder.
2. The validator (`backlog-consistency.ts`) parses only Active Work and Recently Closed sections. It ignores the `## Blocked / Parked` section entirely, so parked items with broken file links are silently missed.
3. `TASK-013.md` contains a stale link to `../active/TASK-012.md` (TASK-012 was moved to closed/).
4. `work/active/2026/` contains old status review files (not plan files) that belong in `work/status/` — noted as a follow-up finding, not addressed here.

---

### Analysis

**Directly verified:**

- `work/parked/TASK-013.md` exists — TASK-013 is correctly placed
- `BACKLOG.md` links to `parked/TASK-013.md` — link is correct
- `backlog-alignment.md` Folder Model has no `parked/` entry
- `intake-workflow.md` already has a "Parked / Deferred" section — no change needed
- Validator ignores `## Blocked / Parked` section (section = null)
- TASK-013.md has stale internal link `../active/TASK-012.md`

**Inferred:**

- Adding `parked/` to the Folder Model is sufficient doc coverage
- Validator parked check should: verify linked file exists, detect orphaned parked files

---

### Recommended Fix

1. Add `work/parked/` to Folder Model in `backlog-alignment.md`
2. Fix stale link in `TASK-013.md`
3. Update validator:
   - Parse `## Blocked / Parked` section in `parseBacklog()`
   - Add `parkedWork` sub-check to `BacklogCheck` type
   - Check each parked row's Detail link resolves to a real file
   - Detect orphaned files in `parked/` with no matching backlog row
   - Update reporter output to show parked check

---

### Scope

**In scope:**

- `core/.hawp/kit/references/backlog-alignment.md` — add `parked/` to Folder Model
- `.hawp/work/parked/TASK-013.md` — fix stale link
- `librarian/scripts/validate-hawp-workflow/types.ts` — add `parkedWork` to BacklogCheck
- `librarian/scripts/validate-hawp-workflow/index.ts` — parse parked section
- `librarian/scripts/validate-hawp-workflow/validations/backlog-consistency.ts` — add parked check
- `librarian/scripts/validate-hawp-workflow/reporter.ts` — display parked results

**Out of scope:** UUID migration, SQLite, indexing, search, queueing, rewriting old closed tasks, moving `active/2026/` status files.

---

### Follow-up Finding

`work/active/2026/` contains 5 dated status review files that are not plan files and likely belong in `work/status/`. Recommend a separate cleanup item.

---

## Outcome (filled at close)

All changes implemented as planned:

1. **`core/.hawp/kit/references/backlog-alignment.md`** — Added `work/parked/` entry to Folder Model with description: "intentionally paused work item detail files (move back to `active/` to resume)"
2. **`.hawp/work/parked/TASK-013.md`** — Fixed stale link from `../active/TASK-012.md` → `../closed/2026/05/10/TASK-012.md`
3. **Validator — `types.ts`** — Added `parkedWork: { total, found, missing }` and `orphanedParked: string[]` to `BacklogCheck`
4. **Validator — `index.ts`** — `parseBacklog()` now parses `## Blocked / Parked` section into a `parked: BacklogRow[]` array; added catch-all `section = null` for unknown `##` headers
5. **Validator — `backlog-consistency.ts`** — Added parked file verification (extracts path from Detail markdown link via `extractLinkPath()`); added orphaned parked file detection
6. **Validator — `reporter.ts`** — Report now shows `Blocked / Parked (N items): Found M/N` and `Orphaned Files (in parked/)` sections

**Final parked-work rule:**
- `work/active/` — plan files for work actively moving
- `work/parked/` — plan files for intentionally paused work; BACKLOG.md Detail link must point here
- `work/closed/YYYY/MM/DD/` — plan files for completed work
- BACKLOG.md must link to the correct folder for each status; validator now enforces this for all three

---

## Verification (filled at close)

- [x] `backlog-alignment.md` Folder Model includes `parked/` — **Evidence:** file edited, line reads "intentionally paused work item detail files"
- [x] TASK-013.md stale link fixed — **Evidence:** link now points to `../closed/2026/05/10/TASK-012.md`
- [x] `npm run typecheck` passes with 0 errors — **Evidence:** see work/evidence/2026/05/10/TASK-014-validation-output.md
- [x] Parked check in validator output shows `Blocked / Parked (1 items): Found: 1/1` — **Evidence:** see TASK-014-validation-output.md
- [x] Orphaned parked section shows `(none)` — **Evidence:** see TASK-014-validation-output.md
- [x] `extractLinkPath()` parses `[plan](parked/TASK-013.md)` correctly — **Evidence:** validator found TASK-013 without error

---

## Close Checklist

- [x] `backlog-alignment.md` Folder Model includes `parked/`
- [x] TASK-013.md stale link fixed
- [x] `BacklogCheck` type has `parkedWork` sub-check
- [x] `parseBacklog()` parses `## Blocked / Parked` section
- [x] `checkBacklogConsistency()` verifies parked file links
- [x] Orphaned parked file detection added
- [x] Reporter shows parked check output
- [x] `npm run typecheck` passes with 0 errors
- [x] `npm run validate:workflow` shows parked section in output
- [x] Evidence saved
- [x] TASK-014 closed and moved to `closed/2026/05/10/`
- [x] BACKLOG.md updated

---

### Status

- [x] Plan written
- [x] Approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
