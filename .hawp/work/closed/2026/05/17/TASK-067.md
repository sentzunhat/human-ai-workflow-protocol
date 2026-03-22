## Task: Audit and align all closed work items to current close-record standards

**Backlog ID:** TASK-067
**Type:** task
**Reported:** 2026-05-17
**Risk Level:** medium
**Status:** done

---

### Input (what was reported)

Review closed work items before legacy cleanup, ensure each record is aligned with current standards and best-intent implementation, and organize stray records into date-based folders.

---

### Context

Current workflow validation passes, but legacy warnings remain and closed-history consistency is mixed. Before normalization, the closed archive needs a rule-aligned audit pass so later automation uses confirmed patterns.

---

### Analysis

**Root cause (or most likely cause):**
Closed records were created across multiple workflow eras; some files predate current Outcome/Verification/Close Checklist conventions and some naming/date placement conventions.

**Directly verified:**

- Validator reports tolerated legacy warnings tied to missing sections and untyped IDs in `.hawp/work/closed/`.
- `BACKLOG.md` currently tracks only active visibility and does not yet include a completed closed-record audit pass.

**Inferred (not yet proven):**

- A complete closed-item alignment audit will reduce ambiguity and prevent accidental over-normalization.
- Some files may require date-folder placement reconciliation before section templating is applied.

**Scope — what else is affected:**

- `.hawp/work/closed/**`
- `.hawp/work/BACKLOG.md`
- `.hawp/work/status/**` (if checkpoint notes are needed)

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `.hawp/work/closed/**`
- `.hawp/work/closed/2026/05/17/TASK-066.md`
- `.hawp/work/active/TASK-068.md`

**Parallel work risk:** medium
**Can implement now:** yes

**Coordination note:**
This task defines the authoritative audit baseline. TASK-068 automation should consume this baseline; TASK-066 legacy normalization should execute after this task and TASK-068.

**Repo-root proof (redacted prefix):**

```bash
$ pwd
<repo-root-abs>

$ git rev-parse --show-toplevel
<repo-root-abs>

$ git rev-parse --show-prefix

$ git status --short
 D .hawp/kit/guidance/da-schema-planning.md
 D .hawp/kit/guidance/shared-standards-review-rubric.md
 M .hawp/kit/standards/README.md
 M .hawp/kit/templates/intake-plan.md
 M .hawp/kit/usage/intake-workflow.md
 M .hawp/work/BACKLOG.md
 M .hawp/work/active/TASK-066.md
?? .hawp/kit/standards/database/
?? .hawp/kit/standards/patterns/
?? .hawp/kit/standards/service-design/
?? .hawp/work/active/TASK-067.md
?? .hawp/work/active/TASK-068.md
?? .hawp/work/parked/TASK-066.md
```

---

### Recommended Fix

1. Build an inventory of all closed records and classify by:
   - naming/ID style,
   - folder date alignment,
   - required close-section completeness.
2. Produce an audit matrix listing exact file paths and required action type (none, move, scaffold, manual review).
3. Update backlog coordination state with findings and explicit handoff inputs for TASK-068.

Current execution snapshot (2026-05-17):

- Inventory and matrix completed for 73 closed files.
- 20 files classified as alignment candidates.
- No closed-plan folder/date relocation candidates detected.
- Evidence: `.hawp/work/evidence/2026/05/17/TASK-067-closed-audit-matrix.md`
- Checkpoint summary: `.hawp/work/status/2026/05/17/TASK-067-status.md`

---

## Outcome (filled at close)

- Produced a complete closed-record audit matrix for all closed markdown plans and classified required actions by type.
- Established exact baseline counts used by TASK-068 apply automation and TASK-066 manual fallback.
- Confirmed no date-folder relocation candidates existed at audit time; all relocations could be treated as explicit conditional logic in TASK-068.
- Delivered evidence and checkpoint artifacts that were directly consumed by downstream tasks.

## Verification (filled at close)

- [x] Audit matrix exists with exact path-level findings: **Evidence:** `.hawp/work/evidence/2026/05/17/TASK-067-closed-audit-matrix.md`.
- [x] Checkpoint summary exists with actionable handoff to TASK-068: **Evidence:** `.hawp/work/status/2026/05/17/TASK-067-status.md`.
- [x] Downstream execution consumed the matrix and cleared warning bucket: **Evidence:** `npm run -s validate:workflow` now reports `Warnings: 0`.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled (all claims have direct evidence or "unproven" tag)
- [x] Evidence files created if large/complex
- [x] Plan file moved to `closed/YYYY/MM/DD/`
- [x] BACKLOG.md updated
- [x] Status report written (non-trivial)
- [x] Decision file created (if applicable)

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
