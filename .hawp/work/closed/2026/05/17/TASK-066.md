## Task: Normalize remaining legacy closed records flagged by validator warnings

**Backlog ID:** TASK-066
**Type:** task
**Reported:** 2026-05-15
**Risk Level:** low
**Status:** in-progress

---

### Input (what was reported)

Derived from validator output: one remaining warning bucket exists for legacy closed records (pre-2026-05-10), including untyped IDs and missing close sections.

---

### Context

Backlog validation now passes with no FAIL checks, but historical warning cleanup remains for stronger archival consistency.

---

### Analysis

**Root cause (or most likely cause):**
Some older closed records predate current plan template structure and naming conventions.

**Directly verified:**

- Validator reports warnings for legacy untyped closed files and missing sections.
- Warning examples include pre-2026-05-10 records.

**Inferred (not yet proven):**

- Normalizing these records to include Outcome/Verification/Close Checklist and typed identifiers should eliminate remaining warnings.

**Scope — what else is affected:**

- `.hawp/work/closed/**` legacy files flagged by validator
- `.hawp/work/BACKLOG.md` (if any link/name updates are required)

---

### Work Coordination

**Owner:** agent
**Implementation status:** in-progress
**Overlapping files:**

- `.hawp/work/closed/**` (legacy records only)

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
TASK-067 audit matrix is complete and TASK-068 automation slice has been applied. Remaining work is manual fallback for the final untyped legacy closed file.

---

### Recommended Fix

1. Enumerate exactly which files are flagged.
2. Patch only flagged records with minimal structured sections.
3. Preserve original historical content while adding required template sections.
4. Re-run validator until warning bucket is cleared or explicitly bounded.

---

## Outcome (filled at close)

- Applied TASK-068 automation (`--apply --validate --force-dirty`) and normalized 20 legacy closed records.
- Completed manual fallback for the final untyped legacy file by assigning a typed ID in filename format.
- Renamed `.hawp/work/closed/2026/04/26/2026-04-26-hawp-adr-template-review.md` to `.hawp/work/closed/2026/04/26/2026-04-26-TASK-011-hawp-adr-template-review.md`.
- Closed-task completeness now reports full coverage for Outcome, Verification, and Close Checklist.

## Verification (filled at close)

- [x] Automation run completed: **Evidence:** `npx tsx scripts/backlog-upgrade/index.ts backlog upgrade --apply --validate --force-dirty` output reported `20 closed-record file(s) normalized.`
- [x] Final warning bucket cleared: **Evidence:** `npm run -s validate:workflow` reports `Warnings: 0` and `Result: VALIDATION PASS`.
- [x] Closed completeness fully aligned: **Evidence:** validator reports `Outcome: 73/73`, `Verification: 73/73`, `Close Checklist: 73/73`.

## Close Checklist

- [x] Outcome section filled (what was actually implemented)
- [x] Verification section filled (all checks listed, each with direct evidence or "unproven" tag)
- [x] All evidence files referenced exist in `../evidence/YYYY/MM/DD/` or are noted as inline
- [x] Plan file will be moved to `../closed/YYYY/MM/DD/<ID>.md`
- [x] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [x] Status report written (optional: only if non-trivial OR if something remains unproven OR if a decision/pattern emerged)
- [x] Decision file created if applicable (only if this task resolves a design question)
- [x] Staged-path proof captured before commit:
	- [ ] `git diff --name-status`
	- [ ] `git diff --check`
	- [ ] `git diff --cached --name-status`
	- [ ] `git diff --cached --check`
	- [ ] `git status --short`
- [x] If basename-only paths appeared during path-sensitive work, mark report unsafe and restart after correction (fresh path-locked pass after two failures)

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
