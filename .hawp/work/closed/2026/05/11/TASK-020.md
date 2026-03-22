# Bug / Task: Add validator CLI help output

**Backlog ID:** TASK-020
**Type:** task
**Reported:** 2026-05-10
**Risk Level:** low

---

### Input (what was reported)

> Add validator CLI help output.

---

### Context

The local validator reports backlog inconsistency because TASK-020 exists in backlog but had no active plan file.
This task should remain scoped to CLI help output only and is not being implemented in this backlog-consistency repair pass.

---

### Analysis

**Root cause (or most likely cause):**
TASK-020 was logged in `BACKLOG.md` but its corresponding active plan file was missing.

**Directly verified:**

- `BACKLOG.md` includes TASK-020 in Active Work.
- `.hawp/work/active/` did not contain `TASK-020.md`.

**Inferred (not yet proven):**

- None.

**Scope — what else is affected:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-020.md`

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-020.md`

**Parallel work risk:** low
**Can implement now:** only after approval

**Coordination note:**
This plan creation repairs backlog consistency only; no TASK-020 feature implementation is included.

---

### Options

#### Option A — Create missing active plan file

Keep TASK-020 active and add the missing plan file to satisfy validator backlog consistency checks.

#### Option B — Remove TASK-020 from backlog

If task intent is canceled, remove TASK-020 row instead of creating a plan.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
TASK-020 is still present in active backlog and appears intended; creating the plan is the minimal consistency repair.

**Files to change:**

- `.hawp/work/BACKLOG.md` — ensure TASK-020 row points to active plan and correct status.
- `.hawp/work/active/TASK-020.md` — create intake plan file.

**What to verify after:**

- [ ] Local validator no longer reports TASK-020 missing plan file.
- [ ] No TASK-020 implementation changes were made.

---

### Implementation Notes

Do not implement CLI help output in this pass. Stop after backlog consistency repair, validation rerun, and evidence capture.

### Outcome

- Added validator CLI help output in the entrypoint with support for --help and -h.
- Help output now documents default local behavior, --hawp-root, --work-root, --debug-closed-task, exit code behavior, and FAIL/WARN/INFO meaning.
- Existing validator behavior and validation rules remain unchanged.

### Verification

- npm run validate:workflow -- --help prints the documented help content successfully.
- npm run typecheck passes in librarian.
- npm run validate:workflow (local) still returns VALIDATION PASS.
- Evidence saved to .hawp/work/evidence/2026/05/11/TASK-020-cli-help-output.md.

### Close Checklist

- [x] Outcome added
- [x] Verification added
- [x] Evidence saved
- [x] Validation rules unchanged
- [x] Scope limited to TASK-020 CLI help output

---

### Status

- [x] Plan written
- [x] Approved by user
- [x] Implemented
- [x] Verified
- [x] Closed
