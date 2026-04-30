# Work Intake — Plan Template

Use this template for plan files saved to `work/status/<date>-<id>-<slug>.md`.

---

## Bug / Task: [Short title]

**Backlog ID:** TASK-XXX
**Type:** bug | task | improvement
**Reported:** YYYY-MM-DD
**Risk Level:** low | medium | high

---

### Input (what was reported)

> Paste the original bug report or task description verbatim.

---

### Context

What area of the codebase is affected? What is the user-visible symptom?

---

### Analysis

**Root cause (or most likely cause):**
_What is broken and why._

**Directly verified:**
_What I confirmed by reading files or running checks._

**Inferred (not yet proven):**
_What I believe to be true but cannot confirm without a live environment._

**Scope — what else is affected:**
_Other files, routes, or modules touched by this._

---

### Work Coordination

**Owner:** human | agent | unassigned
**Implementation status:** not-started | in-progress | blocked | done
**Overlapping files:**

- `path/to/file`

**Parallel work risk:** low | medium | high
**Can implement now:** yes | no | only after approval

**Coordination note:**
_Explain whether another active item touches the same files._

---

### Options

#### Option A — [Name]

_Description. Trade-offs._

#### Option B — [Name]

_Description. Trade-offs._

---

### Recommended Fix

**Option chosen:** A | B
**Rationale:**

**Files to change:**

- `path/to/file` — what changes

**What to verify after:**

- [ ] Verify item 1
- [ ] Verify item 2

---

### Implementation Notes

_Any gotchas, ordering concerns, or things to watch for during implementation._

---

### Status

- [ ] Plan written
- [ ] Approved / auto-approved (low risk)
- [ ] Implemented
- [ ] Verified
- [ ] Closed
