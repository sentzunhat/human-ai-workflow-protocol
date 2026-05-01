## Bug / Task: reconcile closed plans during update/install flow and docs sync

**Backlog ID:** BUG-003
**Type:** bug
**Reported:** 2026-05-01
**Risk Level:** low

### Input (what was reported)

> can we fix it and maybe tell it to move the active files to closed if they are closed update the update files and install files or any references in the docs too please

### Context

The HAWP update/install docs include one-shot shell flows for migrating legacy layouts. A downstream run left many files in `.hawp/work/active/` even when they were already closed in the backlog, and the user asked to fix this plus update related references.

### Analysis

**Root cause (or most likely cause):**
The documented update/install flow migrates legacy folders into `.hawp/work/active/` but does not include a post-migration reconciliation pass that moves already-closed plans into `.hawp/work/closed/YYYY/MM/DD/` based on backlog links.

**Directly verified:**
- `core/update.md` and `core/install.md` include migration steps for `usage/status` and `work/status` into `work/active`.
- `core/install.md` currently documents copies into `work/active` but does not remove legacy `work/adrs` and `work/status` folders after migration.
- Neither file contains a reconciliation step for closed plan files.

**Inferred (not yet proven):**
- In downstream repos with historical plan files in legacy status folders, this can leave closed files in `active/` until manual cleanup.

**Scope — what else is affected:**
- `core/update.md`
- `core/install.md`
- Any sections in those docs that describe migration behavior, checklists, and notes.

### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- `core/update.md`
- `core/install.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
No active work items are open in `.work/BACKLOG.md`, so there is no overlap conflict.

### Options

#### Option A — Add deterministic backlog-based reconciliation

Add a post-migration shell step that parses `.hawp/work/BACKLOG.md` `Done` links and moves matching files from `active/` to `closed/` when needed. Update docs/checklists/notes accordingly.

#### Option B — Keep script unchanged and document manual cleanup

Add only post-update instructions telling users to manually move closed files.

### Recommended Fix

**Option chosen:** A
**Rationale:**
It removes a repeat manual step, keeps active/closed consistent, and is safe because it only moves files when a backlog `Done` plan link explicitly points to `closed/...` and the source exists in `active/...`.

**Files to change:**

- `core/update.md` — add reconciliation step, update migration notes/checklist/prompt requirements.
- `core/install.md` — mirror reconciliation behavior and align migration descriptions/checklist/notes.

**What to verify after:**

- [ ] Both one-shot scripts include the same reconciliation behavior.
- [ ] Install/update docs explicitly mention legacy `work/adrs` and `work/status` removal behavior consistently.
- [ ] References and checklist text are aligned with the new behavior.

### Implementation Notes

Use a conservative parser: read markdown link targets from the Done table in `.hawp/work/BACKLOG.md`, only move when the target path is `closed/...`, source exists under `active/`, and destination does not already exist.

### Status

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
