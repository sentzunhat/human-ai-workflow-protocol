# Bug / Task: Compact BACKLOG.md and archive closed work.

**Backlog ID:** TASK-008
**Type:** task
**Reported:** 2026-05-02
**Risk Level:** low

---

### Input (what was reported)

> work on the active task and aligned active works with the backlog and clean up anything else on the working and active and backlog do the cap backlog task until completion

---

### Context

The backlog had an open compaction task and still carried a long historical Done table. Active plan state was out of sync because TASK-006 still existed in active even though it was already closed.

---

### Analysis

**Root cause (or most likely cause):**
Backlog compaction policy existed but had not yet been applied to the repo-local BACKLOG.md, and an old active plan file remained after closure.

**Directly verified:**

- `.hawp/work/BACKLOG.md` contained a long Done table with historical rows.
- `.hawp/work/active/TASK-006.md` existed while `.hawp/work/closed/2026/05/02/TASK-006.md` also existed.
- Backlog Active Work listed TASK-008 as open with pending plan.

**Inferred (not yet proven):**

- The stale active file likely remained from a prior close step that did not remove the active copy.

**Scope - what else is affected:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/`
- `.hawp/work/closed/2026/05/02/`
- `.hawp/work/evidence/2026/05/02/`

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-006.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
No other active plan files overlapped this cleanup scope.

---

### Options

#### Option A - Keep full historical Done table in BACKLOG.md

Preserves all history in one file but violates compact-index guidance and increases triage noise.

#### Option B - Convert BACKLOG.md to compact active index and cap Recently Closed

Aligns with backlog-alignment policy and keeps historical detail in closed archives.

---

### Recommended Fix

**Option chosen:** B
**Rationale:**
Matches canonical HAWP backlog-alignment policy and improves day-to-day coordination clarity.

**Files to change:**

- `.hawp/work/BACKLOG.md` - replace long Done section with capped Recently Closed and archive links
- `.hawp/work/active/TASK-006.md` - remove stale duplicate active plan file

**What to verify after:**

- [x] BACKLOG.md is compact and operational
- [x] Recently Closed is capped
- [x] Active folder matches active backlog rows

---

### Implementation Notes

Applied the compact-backlog shape and preserved project history via archive links and closed plan files.

---

### Outcome

Completed TASK-008 by compacting BACKLOG.md, capping Recently Closed, aligning active plans with backlog state, and recording verification evidence.

### Verification

- Confirmed compact backlog structure and archive links in `.hawp/work/BACKLOG.md`.
- Confirmed stale `.hawp/work/active/TASK-006.md` removed.
- Confirmed evidence captured in `.hawp/work/evidence/2026/05/02/task-008-backlog-compaction-verification.md`.

### Evidence Links

- `.hawp/work/evidence/2026/05/02/task-008-backlog-compaction-verification.md`

### Final Status

- Final status: done
- Close date: 2026-05-02

### Status

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
