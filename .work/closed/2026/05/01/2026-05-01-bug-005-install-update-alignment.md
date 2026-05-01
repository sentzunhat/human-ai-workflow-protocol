## Bug / Task: align install/update flow with validated migration and reconciliation

**Backlog ID:** BUG-005
**Type:** bug
**Reported:** 2026-05-01
**Risk Level:** low

### Input (what was reported)

> Update the HAWP install/update flow to match the current migration and reconciliation behavior validated in this repository.

### Context

The source docs and one-shot scripts in `core/install.md` and `core/update.md` define downstream migration and refresh behavior. The current repository has validated behavior that must remain the source of truth, including strict update boundaries, no-overwrite handling in `.hawp/work/**`, backlog-driven reconciliation, and legacy cleanup.

### Analysis

**Root cause (or most likely cause):**
`core/update.md` is already aligned for Active Work row reconciliation (`done`/`wont-fix`) and closed-date fallback logic, but `core/install.md` still documents and scripts the older Done-only reconciliation behavior.

**Directly verified:**

- `core/update.md` reconciliation parses both `## Done` and `## Active Work` rows for eligible closures.
- `core/install.md` reconciliation currently scans only `## Done` and requires Closed date fallback without filename/today fallback.
- Both files already preserve `.hawp/work/**` no-overwrite semantics, preserve parked content, avoid overwriting `.github/copilot-instructions.md`, remove stale overlay names, and refresh only HAWP-managed files.

**Inferred (not yet proven):**

- Keeping install/update out of sync can cause divergent behavior between first install and later updates in downstream repos.

**Scope — what else is affected:**

- `core/install.md`
- `core/update.md`

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `core/install.md`
- `core/update.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
No active overlapping work in `.work/BACKLOG.md`.

### Options

#### Option A — Align install script/docs to update flow

Update `core/install.md` reconciliation logic, requirement bullets, prompt text, and verification checklist to match `core/update.md` behavior.

#### Option B — Leave install script unchanged and only clarify docs

Document the mismatch and ask users to run update after install.

### Recommended Fix

**Option chosen:** A
**Rationale:**
Keeps install and update deterministic and idempotent with the same migration and reconciliation semantics.

**Files to change:**

- `core/install.md` — reconcile function and narrative/checklist updates.
- `core/update.md` — light consistency pass for wording and checklist where needed.

**What to verify after:**

- [ ] Install and update both describe and implement backlog Done + Active Work (`done`/`wont-fix`) reconciliation.
- [ ] Closed-date fallback order is consistent (Closed column, filename prefix, today).
- [ ] Required boundaries and no-overwrite semantics remain unchanged.

### Implementation Notes

Keep changes conservative and idempotent. Preserve no-clobber behavior and only adjust reconciliation/wording where drift exists.

### Status

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
