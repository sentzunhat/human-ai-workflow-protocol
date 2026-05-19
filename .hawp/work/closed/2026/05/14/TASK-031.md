# Bug / Task: Enforce path-locked self-validation for DA file references

**Backlog ID:** TASK-031
**Type:** task
**Reported:** 2026-05-12
**Risk Level:** low

---

### Input (what was reported)

> Repeated DA outputs claimed "exact repo-relative paths only" while still emitting basename-only references. Add deterministic self-validation and reject output when path-lock rules fail.

---

### Context

Path-sensitive DA reports currently rely on manual checking and can still emit ambiguous basename-only paths. This creates operational risk and directly conflicts with path-lock discipline.

---

### Analysis

**Root cause (or most likely cause):**
Path rules exist, but report-level self-validation and hard reject semantics are not explicit enough in core scaffold guidance.

**Directly verified:**

- `core/.hawp/kit/instructions/da-file-tracking.md` defines repo-relative path rules but no explicit violation taxonomy.
- `core/.hawp/kit/references/work-item-file-tracking.md` marks v0.2 as future but does not define deterministic rejection categories.
- `core/.hawp/kit/templates/work-item-files.md` lacks a report self-validation gate section.
- `core/.hawp/kit/usage/intake-workflow.md` and `core/.github/instructions/intake.instructions.md` include path discipline but no explicit path-gate failure categories.

**Inferred (not yet proven):**
Adding explicit v0.2 violation categories and output-reject semantics in core scaffold docs will reduce repeated basename/path-collapse failures in downstream DA outputs.

**Scope — what else is affected:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-031.md`
- `core/.hawp/kit/instructions/da-file-tracking.md`
- `core/.hawp/kit/references/work-item-file-tracking.md`
- `core/.hawp/kit/templates/work-item-files.md`
- `core/.hawp/kit/usage/intake-workflow.md`
- `core/.github/instructions/intake.instructions.md`

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-031.md`
- `core/.hawp/kit/instructions/da-file-tracking.md`
- `core/.hawp/kit/references/work-item-file-tracking.md`
- `core/.hawp/kit/templates/work-item-files.md`
- `core/.hawp/kit/usage/intake-workflow.md`
- `core/.github/instructions/intake.instructions.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
Scoped to core scaffold guidance and task bookkeeping only. No CLI, automation, queue, or SQLite implementation in this task.

**Block reason:**
Resolved — historical validator blockers were normalized under TASK-047.

Unblock dependency:

- Historical closed-file normalization tracked under `.hawp/work/active/TASK-047.md`.

Re-evidence (2026-05-14):

- `npm --prefix librarian run validate:workflow` now returns `VALIDATION PASS` with warnings only after TASK-047 apply-phase normalization.
- Historical fail-class findings (post-cutoff untyped and missing required sections) were resolved in TASK-047 evidence/status artifacts.

**Path discipline:**

- Use exact repo-relative paths from repository root for all files in this plan.
- Basenames alone are unsafe for path-sensitive work unless the file is truly at repository root and explicitly marked as such.

**Repo-root proof (required before path-sensitive edits):**

```bash
pwd
/Users/beltrd/Desktop/projects/r-and-d/personal-projects/human-ai-workflow-protocol

git rev-parse --show-toplevel
/Users/beltrd/Desktop/projects/r-and-d/personal-projects/human-ai-workflow-protocol

git rev-parse --show-prefix


git status --short
 M .hawp/work/active/TASK-030-files.md
 M core/.hawp/kit/instructions/da-file-tracking.md
 M core/.hawp/kit/references/work-item-file-tracking.md
 M core/.hawp/kit/templates/work-item-files.md
?? .hawp/work/active/0008-CLARIFICATION-exact-paths.md
?? .hawp/work/active/0008-install-update-distribution-review.md
?? shared_standards/
```

---

### Options

#### Option A — Core guidance-only enforcement

Add explicit v0.2 violation taxonomy and self-reject rules in core scaffold instructions, references, template, and workflow docs.

#### Option B — Start CLI enforcement now

Implement parser/validator categories in `librarian/scripts/backlog-upgrade/`.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
Directly addresses the failure mode with low-risk documentation enforcement in core scaffold, while preserving the no-automation boundary for this step.

**Files to change:**

- `core/.hawp/kit/instructions/da-file-tracking.md` — add explicit violation categories and hard report-gate rule
- `core/.hawp/kit/references/work-item-file-tracking.md` — define v0.2 enforcement semantics
- `core/.hawp/kit/templates/work-item-files.md` — add deterministic self-validation checklist
- `core/.hawp/kit/usage/intake-workflow.md` — require output rejection on invalid path references
- `core/.github/instructions/intake.instructions.md` — add ambient guardrails for path-gate failures

**What to verify after:**

- [ ] New categories exist: `INVALID_REPO_RELATIVE_PATH`, `BASENAME_ONLY_REFERENCE`, `SELF_VALIDATION_FAILURE`
- [ ] Core workflow text requires self-reject before emitting unsafe reports
- [ ] No automation/CLI implementation files were changed

---

### Implementation Notes

Keep scope constrained to core scaffold guidance and path-lock behavior. Do not implement validator code in this task.

---

## Outcome (filled at close)

Implemented core HAWP v0.2 path-lock enforcement semantics in scaffold guidance.

Changes completed:

- Added TASK-031 to `.hawp/work/BACKLOG.md` Active Work.
- Added task plan `.hawp/work/active/TASK-031.md`.
- Added explicit violation categories to core guidance:
  - `INVALID_REPO_RELATIVE_PATH`
  - `BASENAME_ONLY_REFERENCE`
  - `SELF_VALIDATION_FAILURE`
- Added deterministic pre-emit report self-validation gate in core instruction/workflow/template overlays.
- Added strict v0.2 path validator semantics in root workflow overlays:
  - path must contain `/`
  - path must start with repo-root directory prefix
  - path must be copy-pasteable in git commands from repo root
  - `INVALID_REPO_RELATIVE_PATH` explicitly includes basename-only references

Scope kept tight:

- No CLI/parser/automation implementation was added.
- No queue/SQLite/runtime systems were introduced.

---

## Verification (filled at close)

- [x] Core violation categories were added. **Evidence:** `core/.hawp/kit/instructions/da-file-tracking.md`, `core/.hawp/kit/references/work-item-file-tracking.md`
- [x] Deterministic report self-validation gate was added. **Evidence:** `core/.hawp/kit/instructions/da-file-tracking.md`, `core/.hawp/kit/usage/intake-workflow.md`, `core/.hawp/kit/templates/work-item-files.md`, `core/.github/instructions/intake.instructions.md`
- [x] Strict path validator rule was added in root overlays. **Evidence:** `.hawp/kit/usage/intake-workflow.md`, `.github/instructions/intake.instructions.md`
- [x] TypeScript checks pass. **Evidence:** `npm --prefix librarian run typecheck` passed.
- [x] Full workflow validation pass (warnings-only tolerated): **Evidence:** `.hawp/work/evidence/2026/05/14/TASK-047-validate-workflow-post-apply.txt`.

---

## Close Checklist

**Before marking this task done in BACKLOG.md, verify ALL:**

- [ ] Outcome section filled (what was actually implemented)
- [ ] Verification section filled (all checks listed, each with direct evidence or "unproven" tag)
- [ ] All evidence files referenced exist in `../evidence/YYYY/MM/DD/` or are noted as inline
- [ ] Plan file will be moved to `../closed/YYYY/MM/DD/<ID>.md`
- [ ] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [ ] Status report written (optional: only if non-trivial OR if something remains unproven OR if a decision/pattern emerged)
- [ ] Decision file created if applicable (only if this task resolves a design question)
- [ ] Staged-path proof captured before commit:
  - [ ] `git diff --name-status`
  - [ ] `git diff --check`
  - [ ] `git diff --cached --name-status`
  - [ ] `git diff --cached --check`
  - [ ] `git status --short`

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
