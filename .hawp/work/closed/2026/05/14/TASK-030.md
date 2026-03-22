# Bug / Task: TASK-030 Fresh Path-Locked Pass

**Backlog ID:** TASK-030
**Type:** task
**Reported:** 2026-05-12
**Risk Level:** low

---

### Input (what was reported)

> Add base path discipline to HAWP workflow guidance so agents use exact repo-relative paths from repository root in plans, reports, staging summaries, evidence, and handoffs.

---

### Context

Previous TASK-030 attempts had path-discipline failures caused by basename/path collapse in path-sensitive reporting.

---

### Analysis

**Root cause (or most likely cause):**
Path-sensitive guidance did not explicitly enforce repository-root-relative path evidence and path-proof gates before edit/commit.

**Directly verified:**

- Required repo-root proof commands were run.
- All approved TASK-030 paths exist.
- Existing workflow and instruction files lacked complete path-locked gates for repo-root proof, staged-path proof, basename-unsafety, loop-breaker restart, and lane-boundary handling for unrelated changes.

**Inferred (not yet proven):**
Adding explicit path-locked rules in both root and core scaffold guidance will reduce repeated basename/path-collapse failures.

**Scope - what else is affected:**

- `.hawp/kit/usage/intake-workflow.md`
- `.hawp/kit/templates/intake-plan.md`
- `.hawp/kit/patterns/parallel-work-guardrails.md`
- `.hawp/kit/start-here.md`
- `.github/instructions/intake.instructions.md`
- `core/.hawp/kit/usage/intake-workflow.md`
- `core/.hawp/kit/templates/intake-plan.md`
- `core/.hawp/kit/patterns/parallel-work-guardrails.md`
- `core/.hawp/kit/start-here.md`
- `core/.github/instructions/intake.instructions.md`
- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-030.md`

---

### Required Path Table

| Purpose                    | Exact repo-relative path                            | Exists? | Action        |
| -------------------------- | --------------------------------------------------- | ------- | ------------- |
| Active work item           | .hawp/work/active/TASK-030.md                       | Yes     | Create/update |
| Backlog                    | .hawp/work/BACKLOG.md                               | Yes     | Update        |
| Workflow usage guide       | .hawp/kit/usage/intake-workflow.md                  | Yes     | Update        |
| Workflow template          | .hawp/kit/templates/intake-plan.md                  | Yes     | Update        |
| Parallel work guardrails   | .hawp/kit/patterns/parallel-work-guardrails.md      | Yes     | Update        |
| Start-here guide           | .hawp/kit/start-here.md                             | Yes     | Update        |
| GitHub instruction overlay | .github/instructions/intake.instructions.md         | Yes     | Update        |
| Core scaffold usage guide  | core/.hawp/kit/usage/intake-workflow.md             | Yes     | Update        |
| Core scaffold template     | core/.hawp/kit/templates/intake-plan.md             | Yes     | Update        |
| Core scaffold guardrails   | core/.hawp/kit/patterns/parallel-work-guardrails.md | Yes     | Update        |
| Core scaffold start-here   | core/.hawp/kit/start-here.md                        | Yes     | Update        |
| Core scaffold instructions | core/.github/instructions/intake.instructions.md    | Yes     | Update        |

---

### Work Coordination

**Owner:** agent
**Implementation status:** in-progress
**Overlapping files:**

- .hawp/work/active/TASK-030.md
- .hawp/work/BACKLOG.md
- .hawp/kit/usage/intake-workflow.md
- .hawp/kit/templates/intake-plan.md
- .hawp/kit/patterns/parallel-work-guardrails.md
- .hawp/kit/start-here.md
- .github/instructions/intake.instructions.md
- core/.hawp/kit/usage/intake-workflow.md
- core/.hawp/kit/templates/intake-plan.md
- core/.hawp/kit/patterns/parallel-work-guardrails.md
- core/.hawp/kit/start-here.md
- core/.github/instructions/intake.instructions.md

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
Unrelated working-tree changes are treated as lane boundaries and remain untouched.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
Add explicit path-locked rules and proof gates in workflow-first docs and overlays at both root and core scaffold paths.

**Files to change:**

- .hawp/work/active/TASK-030.md - reset plan for fresh path-locked pass and capture required table/proofs
- .hawp/work/BACKLOG.md - align TASK-030 title/update metadata
- .hawp/kit/usage/intake-workflow.md - enforce exact repo-relative path evidence and path-proof gates
- .hawp/kit/templates/intake-plan.md - require repo-root proof and staged-path proof in plan workflow
- .hawp/kit/patterns/parallel-work-guardrails.md - enforce lane boundaries and exact-path overlap lists
- .hawp/kit/start-here.md - add quick path-discipline startup rules
- .github/instructions/intake.instructions.md - enforce path-locked operating rules
- core/.hawp/kit/usage/intake-workflow.md - mirror workflow rules in core scaffold
- core/.hawp/kit/templates/intake-plan.md - mirror template rules in core scaffold
- core/.hawp/kit/patterns/parallel-work-guardrails.md - mirror guardrails in core scaffold
- core/.hawp/kit/start-here.md - mirror startup rules in core scaffold
- core/.github/instructions/intake.instructions.md - mirror overlay rules in core scaffold

**What to verify after:**

- [ ] Required diff/status checks are clean and path-locked
- [ ] Optional npm verification checks run if available
- [ ] Only approved TASK-030 paths are staged and committed

---

### Repo-Root Proof

```bash
pwd
human-ai-workflow-protocol

git rev-parse --show-toplevel
human-ai-workflow-protocol

git rev-parse --show-prefix


git status --short
 M .hawp/work/BACKLOG.md
 M .hawp/work/active/TASK-030-files.md
 M .hawp/work/active/TASK-030.md
 M core/.hawp/kit/instructions/da-file-tracking.md
 M core/.hawp/kit/references/work-item-file-tracking.md
 M core/.hawp/kit/templates/work-item-files.md
?? .hawp/work/active/0008-CLARIFICATION-exact-paths.md
?? .hawp/work/active/0008-install-update-distribution-review.md
?? shared_standards/
```

---

## Outcome (filled at close)

Completed a fresh path-locked pass and confirmed path-discipline guardrails are present in both root and core scaffold guidance.

Delivered outcomes:

- Path-sensitive operating rules now require exact repo-relative paths from repository root.
- Repo-root proof and staged-path proof requirements are explicitly documented for path-sensitive edits/commits.
- Basename/path-collapse handling now includes explicit unsafe classification and loop-breaker restart rule.
- Parallel lane boundaries are explicitly enforced in guidance.
- Root and core overlays are aligned on path-lock policy and error-class terminology.

## Verification (filled at close)

- [x] Root and core instruction overlays include exact-path and proof-gate rules. **Evidence:** `.github/instructions/intake.instructions.md` and `core/.github/instructions/intake.instructions.md` include rules 15-24 with repo-root and staged-path proof steps.
- [x] Root and core intake workflow docs include path-lock gates and failure classification. **Evidence:** `.hawp/kit/usage/intake-workflow.md` and `core/.hawp/kit/usage/intake-workflow.md` include `INVALID_REPO_RELATIVE_PATH`, `BASENAME_ONLY_REFERENCE`, and `SELF_VALIDATION_FAILURE` guidance.
- [x] Root and core templates/start-here/guardrails include path-sensitive exact-path language. **Evidence:** `.hawp/kit/templates/intake-plan.md`, `.hawp/kit/patterns/parallel-work-guardrails.md`, `.hawp/kit/start-here.md`, `core/.hawp/kit/templates/intake-plan.md`, `core/.hawp/kit/patterns/parallel-work-guardrails.md`, `core/.hawp/kit/start-here.md`.
- [x] Workflow validation remains passing after closeout. **Evidence:** `npm --prefix librarian run validate:workflow` returns `VALIDATION PASS` (warnings only).

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [ ] Staged-path proof captured before commit
- [ ] Only approved TASK-030 paths are staged
- [ ] Commit created with focused message

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
