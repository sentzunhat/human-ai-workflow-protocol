## Bug / Task: reconcile done active plans by id and date

**Backlog ID:** BUG-004
**Type:** bug
**Reported:** 2026-05-01
**Risk Level:** low

### Input (what was reported)

> we are still not moving the active done work to the other folders

### Context

A downstream repo run of the update flow kept many done plan files under `.hawp/work/active/`.

### Analysis

**Root cause (or most likely cause):**
The current reconciliation function only moves files when Done rows in `.hawp/work/BACKLOG.md` already link to `closed/...` targets. If Done rows still point at `active/...` or use inconsistent links, no move happens.

**Directly verified:**
- `core/update.md` and `core/install.md` use `reconcile_closed_plans_from_backlog` that filters only `closed/*`, `work/closed/*`, `.hawp/work/closed/*` links.
- No fallback path exists that uses Done row ID + Closed date to move matching `active` files.

**Inferred (not yet proven):**
- Downstream backlogs that mark items done but still link plan paths to active files will retain stale done plans in `active/`.

**Scope — what else is affected:**
- `core/update.md`
- `core/install.md`

### Work Coordination

**Owner:** agent
**Implementation status:** in-progress
**Overlapping files:**

- `core/update.md`
- `core/install.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
No other active work items in this repo overlap these files.

### Options

#### Option A — Backlog row-aware reconciliation

Parse Done table rows (ID, Closed date, Plan link). Move matching `active/*-<ID>-*.md` files into `closed/YYYY/MM/DD/` even if Plan links are missing or still point to active paths.

#### Option B — Keep strict link-only behavior

Continue requiring closed links and add only documentation warnings.

### Recommended Fix

**Option chosen:** A
**Rationale:**
It preserves no-overwrite guarantees while making reconciliation resilient to common backlog drift.

**Files to change:**

- `core/update.md` — update function and docs/checklist language.
- `core/install.md` — mirror function and docs/checklist language.

**What to verify after:**

- [ ] Reconciliation supports closed links and active-link fallback via ID + Closed date.
- [ ] Existing behavior for no-overwrite and legacy migration remains intact.
- [ ] Install/update docs remain aligned.

### Implementation Notes

Keep logic conservative: only move files from `active/`, never overwrite destination files, and require a valid Closed date for fallback destination generation.

### Status

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
