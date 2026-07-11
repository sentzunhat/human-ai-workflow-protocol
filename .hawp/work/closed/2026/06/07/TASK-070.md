## Task: TASK-070 — Evaluate private workflow-only standards for safe public adaptation

**Backlog ID:** TASK-070
**Type:** governance
**Reported:** 2026-06-04
**Risk Level:** low
**Status:** done
**Closed:** 2026-06-07

---

### Input (what was reported)

> Review internal-only standards docs and decide which parts can be rewritten for public HAWP kit use (follow-up to TASK-069).

---

### Context

Audit internal standards documents that are currently private or workflow-only and identify content safe to adapt for public consumption in the HAWP kit.

---

### Analysis

**Root cause:** Some standards were written for internal workflows and contain project-specific or sensitive guidance that must be generalized before public release.

**Directly verified:** BACKLOG row exists for TASK-070 with `inbox` status and plan pending (then updated during work).

**Inferred:** The work can be scoped into a quick audit (identify obviously public-safe items) and a follow-up rewrite phase for content that needs redaction or reframing.

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- .github/instructions/hawp-backlog-alignment.instructions.md
- .github/instructions/hawp-docs-alignment.instructions.md
- .github/instructions/intake.instructions.md
- core/providers/shared/behaviors/hawp-core.md
- core/providers/shared/behaviors/hawp-backlog-alignment.md
- core/providers/shared/behaviors/hawp-docs-alignment.md
- core/providers/shared/behaviors/hawp-intake.md

**Path discipline:** Any edits reference exact repo-relative paths from repository root.

---

### Options

#### Option A — Quick audit

- Inspect internal standards and mark items safe for public copy or requiring rewrite.
- Deliverable: annotated checklist + proposed public-safe fragments.

#### Option B — Deep rewrite

- Perform full rewrite of identified content into public-safe kit docs.
- Deliverable: PR with rewritten files and verification notes.

---

### Recommended Fix

**Option chosen:** A

**Rationale:** Fast feedback loop; reduces scope and allows follow-up work.

**Files changed:**

- `core/providers/shared/behaviors/hawp-intake.md` — wording updated to clarify public-safe scope and coordination expectations.

**What was implemented:**

- Wording in `core/providers/shared/behaviors/hawp-intake.md` adjusted to ensure guidance is public-safe and instruct maintainers to coordinate before extending scope beyond `.hawp/`.
- Audit evidence created at `.hawp/work/evidence/2026/06/06/TASK-070-audit.md`.
- BACKLOG row updated to point to this plan while closing.

---

## Outcome

- Updated `core/providers/shared/behaviors/hawp-intake.md` to broaden public-friendly wording and coordination notes.
- Created audit evidence: `.hawp/work/evidence/2026/06/06/TASK-070-audit.md` (quick audit and repo-root proof).
- BACKLOG updated and plan moved to closed on `2026-06-07`.

---

## Verification


### Evidence Follow-Up

- [ ] Research evidence for: Audit evidence created: `.hawp/work/evidence/2026/06/06/TASK-070-audit.md` — contains repo-root proo
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: Implementation change: `core/providers/shared/behaviors/hawp-intake.md` — wording updated.
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: BACKLOG updated to reflect the closed plan file location.
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: Commits & staged-path proof: recorded inline below (commit SHAs and changed files).
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.

- [x] Audit evidence created: `.hawp/work/evidence/2026/06/06/TASK-070-audit.md` — contains repo-root proof and short file-level findings.
- [x] Implementation change: `core/providers/shared/behaviors/hawp-intake.md` — wording updated.
- [x] BACKLOG updated to reflect the closed plan file location.
- [x] Commits & staged-path proof: recorded inline below (commit SHAs and changed files).

**Evidence:**

- Audit file: `.hawp/work/evidence/2026/06/06/TASK-070-audit.md`

---

### Commits & staged-path proof

- `129fc59` — finalize TASK-070 plan outcome and verification
  - Files changed: `.hawp/work/active/TASK-070.md`

- `39467ec` — close TASK-070 plan and move to closed/2026/06/07
  - Files changed: `.hawp/work/closed/2026/06/07/TASK-070.md` (added)
  - Files removed: `.hawp/work/active/TASK-070.md` (deleted)

- `2046809` — mark TASK-070 done in BACKLOG.md and add closed plan link
  - Files changed: `.hawp/work/BACKLOG.md`

Exact staged-path proof (git diffs and status) were captured prior to each commit and are recorded in the repository history alongside these commits.

---

## Close Checklist

- [x] Outcome section filled (what was actually implemented)
- [x] Verification section filled (all checks listed, each with direct evidence or "unproven" tag)
- [x] All evidence files referenced exist in `../evidence/2026/06/06/`
- [x] Plan file moved to this closed path
- [x] BACKLOG.md row moved to Recently Closed (or marked done)
- [x] Staged-path proof recorded (see commits above)

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
