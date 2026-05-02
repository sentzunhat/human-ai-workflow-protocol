# Work Intake - Plan Template

## Bug / Task: core backlog alignment guidance for compact active backlog management

**Backlog ID:** TASK-005
**Type:** improvement
**Reported:** 2026-05-02
**Risk Level:** low

### Input (what was reported)

> Add a new core alignment document that defines HAWP compact backlog management so BACKLOG.md remains an active index and not a permanent history database.

### Context

The repository currently keeps all done items in `.work/BACKLOG.md`, which grows quickly and adds noise for both humans and agents.

### Analysis

**Root cause (or most likely cause):**
There is no explicit, dedicated core guidance that defines backlog compaction rules, recently closed caps, and archive boundaries.

**Directly verified:**
- `.work/BACKLOG.md` includes a long `Done` section with all historical closed rows.
- The repo already has date-partitioned folders under `.work/closed/`, `.work/evidence/`, and `.work/decisions/`.

**Inferred (not yet proven):**
- Without a compact-backlog policy, downstream users and agents are likely to keep appending done history in the main backlog.

**Scope - what else is affected:**
- `core/BACKLOG_ALIGNMENT.md` (new guidance document)
- `.work/BACKLOG.md` status row updates for this task lifecycle

### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- `core/BACKLOG_ALIGNMENT.md`
- `.work/BACKLOG.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
No active plan files currently overlap this new documentation file.

### Options

#### Option A - Add dedicated alignment document

Create `core/BACKLOG_ALIGNMENT.md` with concise, markdown-first rules and a compact backlog example.
Trade-off: one more core doc to maintain.

#### Option B - Update existing README/install docs only

Embed guidance in current docs without a dedicated policy file.
Trade-off: lower discoverability and weaker norm-setting for agents.

### Recommended Fix

**Option chosen:** A
**Rationale:**
A single dedicated alignment file is explicit, scoped, and easy for agents to load quickly.

**Files to change:**

- `core/BACKLOG_ALIGNMENT.md` - add compact backlog model, rules, example, and migration guidance

**What to verify after:**

- [ ] New file includes all required sections from request
- [ ] Guidance explicitly states BACKLOG is active index, not permanent history
- [ ] Example compact backlog shape is present and practical

### Implementation Notes

Keep wording practical and short for agent readability. Do not add automation or runtime behavior.

### Outcome

Created `core/BACKLOG_ALIGNMENT.md` with a compact backlog policy that defines active index boundaries, archive folders, recently closed caps, repeated update-task handling, migration steps, agent rules, and an acceptance checklist.

### Verification

- Confirmed the new file exists at `core/BACKLOG_ALIGNMENT.md`.
- Confirmed explicit language that `BACKLOG.md` is an active index and not permanent history.
- Confirmed required folder model includes active, closed, status, evidence, and decisions paths.
- Confirmed compact backlog example and long-backlog migration guidance are included.

### Evidence Links

- Documentation artifact: `core/BACKLOG_ALIGNMENT.md`

### Final Status

- Final status: done
- Close date: 2026-05-02

### Status

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
