# Bug / Task: update flow leaves legacy .hawp/work/adrs and status trees

**Backlog ID:** BUG-002
**Type:** bug
**Reported:** 2026-05-01
**Risk Level:** low

---

### Input (what was reported)

> Update run is not quite right. After running the update flow, target repos still keep both `.hawp/work/active` and legacy `.hawp/work/status`, and both `.hawp/work/decisions/...` and legacy `.hawp/work/adrs`.

---

### Context

The issue is in the published one-shot update command in `core/update.md`. The command migrates legacy folders into current destinations but does not clean up source legacy folders, leaving duplicate trees in updated downstream repos.

---

### Analysis

**Root cause (or most likely cause):**
The migration steps for `.hawp/work/adrs` and `.hawp/work/status` use no-clobber copy helpers but never remove the old folders after copy succeeds.

**Directly verified:**

- `core/update.md` migration step 3b copies `.hawp/work/adrs` into `.hawp/work/decisions/$MDATE`.
- `core/update.md` migration step 3c copies `.hawp/work/status` into `.hawp/work/active`.
- Neither step removes the legacy source folder.

**Inferred (not yet proven):**
Downstream repos that had these legacy folders will show duplicate trees after update unless they manually remove legacy folders.

**Scope — what else is affected:**

- `core/update.md` command block and surrounding migration checklist/notes wording.

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `core/update.md`
- `.work/BACKLOG.md`
- `.work/active/2026-05-01-bug-002-update-flow-clean-legacy-work-folders.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
No active items are currently open in `.work/BACKLOG.md`.

---

### Options

#### Option A — keep copy-only migration

Keep current behavior and document manual cleanup.

#### Option B — copy then remove only migrated legacy folders

After successful copy for `.hawp/work/adrs` and `.hawp/work/status`, remove those source folders to complete migration while preserving no-clobber protections.

---

### Recommended Fix

**Option chosen:** B
**Rationale:**
Matches user expectation of migration semantics and prevents stale duplicate trees while keeping data-preserving behavior.

**Files to change:**

- `core/update.md` — migration steps, requirements bullets, post-update checklist, and notes wording
- `.work/BACKLOG.md` — status transitions and plan link
- `.work/active/2026-05-01-bug-002-update-flow-clean-legacy-work-folders.md` — implementation notes and closure status

**What to verify after:**

- [ ] Migration steps remove `.hawp/work/adrs` and `.hawp/work/status` after copy.
- [ ] Requirement and checklist text reflect removal of legacy work folders.
- [ ] No broader update boundary is changed.

---

### Implementation Notes

Implemented in `core/update.md` by adding guarded cleanup after migration copies:

- `.hawp/work/adrs` is removed after copy into `.hawp/work/decisions/$MDATE`.
- `.hawp/work/status` is removed after copy into `.hawp/work/active`.

Verification performed:

- Confirmed the script text now includes both cleanup lines.
- Ran a shell simulation with sample files showing copied files remain in destination trees while legacy source folders are removed.
- Confirmed requirements/checklist/notes sections were updated to match behavior.

---

### Status

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
