# Bug / Task: Run and operationalize repo-root .hawp validation

**Backlog ID:** TASK-048
**Type:** task
**Reported:** 2026-05-14
**Risk Level:** low

---

### Input (what was reported)

> maybe also an item to run the validation on this repo root .hawp

---

### Context

This repository needs a dedicated operational task for running and documenting validator checks against the repo-root `.hawp/work` state, with clear output interpretation and follow-up handling.

---

### Analysis

**Root cause (or most likely cause):**
Validation has been run ad hoc during other tasks, but there is no dedicated work item that defines repeatable execution and interpretation for repo-root `.hawp` health checks.

**Directly verified:**

- Validation command exists and targets repo-root `.hawp/work`: `npm --prefix librarian run validate:workflow`.
- Current runs surface both consistency and historical completeness findings, which need structured triage.

**Inferred (not yet proven):**

- A dedicated runbook-style task for repo-root validation execution/reporting will reduce ambiguity and improve future unblocking cadence.

**Scope — what else is affected:**

- `.hawp/work/active/TASK-048.md`
- `.hawp/work/status/YYYY/MM/DD/TASK-048-status.md` (optional, if generated)
- Potentially `.hawp/kit/usage/` docs only if command interpretation guidance is missing

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `.hawp/work/active/TASK-048.md`
- `.hawp/work/status/` (if status artifact is produced)

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
This is an operational execution/reporting item, not a broad remediation item. Remediation work belongs in TASK-047.

Dependency note:

- Execute with TASK-047 policy classification (fail vs warn) so repo-root validation reporting is consistent and non-destructive.

**Path discipline:**

- Use exact repo-relative paths from repository root for all files in this plan.
- Basenames alone are unsafe for path-sensitive work unless the file is truly at repository root and explicitly marked as such.

**Repo-root proof (required before path-sensitive edits):**

```bash
pwd
<repo-root-abs>
git rev-parse --show-toplevel
<repo-root-abs>
git rev-parse --show-prefix


git status --short
```

---

### Options

#### Option A — One-time validation execution with status artifact

Run validator, summarize findings, and save a dated status report for operator handoff.

#### Option B — Add automation gate changes

Modify CI or scripts to enforce validation continuously.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
Matches requested scope (“an item to run validation”) without introducing policy or CI behavior changes.

**Files to change:**

- `.hawp/work/active/TASK-048.md` (this plan)
- Optional: `.hawp/work/status/YYYY/MM/DD/TASK-048-status.md`

**What to verify after:**

- [ ] Validation command executes from repo root successfully
- [ ] Findings are categorized as pass/fail/warn with explicit affected files
- [ ] Follow-up ownership points to TASK-047 for historical closed-file remediation

---

### Implementation Notes

Treat this task as an execution/reporting checkpoint to make repo-root `.hawp` health visible and repeatable.

---

## Outcome (filled at close)

Completed repo-root `.hawp/work` validation execution and operational reporting.

- Ran `npm --prefix librarian run validate:workflow` from repository root.
- Captured sanitized execution evidence and repo-root proof in `.hawp/work/evidence/2026/05/14/TASK-047-validate-workflow-output.txt`.
- Published status summary to `.hawp/work/status/2026/05/14/TASK-048-status.md` with lane-based follow-up ownership.
- Routed remediation ownership to TASK-047 and kept this item execution/reporting-only.

---

## Verification (filled at close)

- [x] Validation command executes from repo root successfully: **Evidence:** `.hawp/work/evidence/2026/05/14/TASK-047-validate-workflow-output.txt` includes command run and `VALIDATE_EXIT_CODE=1`.
- [x] Findings are categorized as pass/fail/warn with affected files: **Evidence:** `.hawp/work/status/2026/05/14/TASK-048-status.md`.
- [x] Follow-up ownership points to historical remediation lane: **Evidence:** `.hawp/work/status/2026/05/14/TASK-048-status.md` references TASK-047.

---

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled (all claims have direct evidence or "unproven" tag)
- [x] Evidence files created if large/complex
- [x] Plan file moved to closed/YYYY/MM/DD/
- [x] BACKLOG.md updated
- [x] Status report written (if non-trivial / unproven / decision-bearing)
- [ ] Decision file created if applicable (only if this task resolves a design question)

**Status:**

- [x] Plan written
- [x] Approved (explicit user request)
- [x] Implemented
- [x] Verified
- [x] Closed
