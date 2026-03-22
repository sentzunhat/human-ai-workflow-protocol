# Bug / Task: Sanitize machine-local absolute paths in HAWP work artifacts

**Backlog ID:** TASK-034
**Type:** task
**Reported:** 2026-05-13
**Risk Level:** low

---

### Input (what was reported)

> can you do me a favour of going in the work folder or any references that have the full working file path of my machine and remove them please
>
> <user-home>/Desktop/projects/r-and-d/personal-projects/
>
> we can leave it from the repo folder name i dont want to share my personal full path do an audit and if anything on the workf low needs to change make those changes please this can be a new work item please make sure it is a new work item and that the intake is working well if not then fix the workflow and maybe make it easier for small to huge ai model to understand only focus on the working items of this work item and fiiles too

---

### Context

The user wants all machine-local absolute paths removed from HAWP work artifacts and related workflow references, keeping only repo-name anchored or repo-relative forms.

---

### Analysis

**Root cause (or most likely cause):**
Older plan/evidence files captured terminal commands and outputs with host-local absolute paths (for example `<user-home>/...`). Workflow guidance enforced repo-relative file evidence but did not explicitly ban machine-local path disclosure in prose/examples.

**Directly verified:**

- Absolute path matches were found under `.hawp/work/active/`, `.hawp/work/closed/`, and `.hawp/work/evidence/`.
- The reported local prefix `<user-home>/Desktop/projects/r-and-d/personal-projects/` appears in multiple artifacts.
- Additional absolute paths outside the reported prefix (for external repos) also appear under `.hawp/work/`.

**Inferred (not yet proven):**
Adding explicit privacy-safe path guidance to workflow docs/templates will reduce recurrence when humans or models paste command transcripts.

**Scope - what else is affected:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-034.md`
- `.hawp/work/active/TASK-030.md`
- `.hawp/work/active/TASK-032.md`
- `.hawp/work/closed/2026/05/13/TASK-033.md`
- `.hawp/work/closed/2026/05/10/TASK-019.md`
- `.hawp/work/closed/2026/05/11/BUG-006.md`
- `.hawp/work/evidence/2026/05/10/TASK-012-implementation-results.md`
- `.hawp/work/evidence/2026/05/10/TASK-014-validation-output.md`
- `.hawp/work/evidence/2026/05/10/TASK-019-tekit-validator-compatibility.md`
- `.hawp/work/evidence/2026/05/10/TASK-019-mictlan-validator-compatibility.md`
- `.hawp/work/evidence/2026/05/10/TASK-021-mictlan-validator-rerun-after-drift-cleanup.md`
- `.hawp/work/evidence/2026/05/10/TASK-021-tekit-validator-rerun-after-drift-cleanup.md`
- `.hawp/work/evidence/2026/05/11/TASK-020-backlog-consistency-repair.md`
- `.hawp/work/evidence/2026/05/11/TASK-020-cli-help-output.md`
- `.hawp/work/evidence/2026/05/11/BUG-006-validator-external-root-diagnostics.md`
- `.hawp/work/evidence/2026/05/11/mictlan-validator-rerun.md`
- `.hawp/work/evidence/2026/05/11/tekit-validator-rerun.md`
- `.hawp/work/evidence/2026/05/11/validator-rerun-hawp-tekit-mictlan.md`
- `.hawp/kit/usage/intake-workflow.md`
- `.hawp/kit/templates/intake-plan.md`

---

### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-034.md`
- `.hawp/work/active/TASK-030.md`
- `.hawp/work/active/TASK-032.md`
- `.hawp/work/closed/2026/05/13/TASK-033.md`
- `.hawp/work/closed/2026/05/10/TASK-019.md`
- `.hawp/work/closed/2026/05/11/BUG-006.md`
- `.hawp/work/evidence/2026/05/10/TASK-012-implementation-results.md`
- `.hawp/work/evidence/2026/05/10/TASK-014-validation-output.md`
- `.hawp/work/evidence/2026/05/10/TASK-019-tekit-validator-compatibility.md`
- `.hawp/work/evidence/2026/05/10/TASK-019-mictlan-validator-compatibility.md`
- `.hawp/work/evidence/2026/05/10/TASK-021-mictlan-validator-rerun-after-drift-cleanup.md`
- `.hawp/work/evidence/2026/05/10/TASK-021-tekit-validator-rerun-after-drift-cleanup.md`
- `.hawp/work/evidence/2026/05/11/TASK-020-backlog-consistency-repair.md`
- `.hawp/work/evidence/2026/05/11/TASK-020-cli-help-output.md`
- `.hawp/work/evidence/2026/05/11/BUG-006-validator-external-root-diagnostics.md`
- `.hawp/work/evidence/2026/05/11/mictlan-validator-rerun.md`
- `.hawp/work/evidence/2026/05/11/tekit-validator-rerun.md`
- `.hawp/work/evidence/2026/05/11/validator-rerun-hawp-tekit-mictlan.md`
- `.hawp/kit/usage/intake-workflow.md`
- `.hawp/kit/templates/intake-plan.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
Changes are confined to this task's plan, backlog row, and explicitly scoped work/workflow files. Unrelated working-tree changes remain untouched.

**Path discipline:**

- Keep all file-path evidence repo-relative.
- Replace machine-local absolute path prefixes with privacy-safe placeholders (for example `<repo-root>/...` or `<workspace>/...`).

**Repo-root proof (required before path-sensitive edits):**

```bash
pwd
git rev-parse --show-toplevel
git rev-parse --show-prefix
git status --short
```

Outputs were captured during analysis and intentionally redacted in this artifact to avoid storing machine-local prefixes.

---

### Options

#### Option A - Targeted replacement of reported prefix only

Replace only `<user-home>/Desktop/projects/r-and-d/personal-projects/` occurrences. Lower edit count but leaves other local absolute paths in historical files.

#### Option B - Comprehensive `.hawp/work` local-path sanitization + workflow guardrails

Replace all `<user-home>/...` absolute paths in `.hawp/work` with privacy-safe forms and update workflow docs/templates to prevent recurrence.

---

### Recommended Fix

**Option chosen:** B
**Rationale:**
Meets the immediate privacy request and hardens future intake execution for human and model contributors.

**Files to change:**

- `.hawp/work/**` scoped files listed above - sanitize host-local absolute paths
- `.hawp/kit/usage/intake-workflow.md` - add explicit privacy-safe path logging guidance
- `.hawp/kit/templates/intake-plan.md` - add explicit privacy-safe path logging guidance

**What to verify after:**

- [ ] No `<user-home>/` matches remain in `.hawp/work/**`
- [ ] No `<user-home>/Desktop/projects/r-and-d/personal-projects/` matches remain anywhere in repo
- [ ] Workflow docs include explicit privacy-safe path examples
- [ ] Only scoped files for TASK-034 are changed by this task

---

### Implementation Notes

Keep historical meaning intact while sanitizing only machine-local prefixes.

---

## Outcome (filled at close)

- Sanitized machine-local absolute path references across scoped `.hawp/work/**` artifacts to privacy-safe forms.
- Updated workflow guidance to reinforce privacy-safe path handling during intake and verification logging.
- Closed the backlog row and moved the plan from `active/` to `closed/2026/05/13/`.

## Verification (filled at close)

- Verified no active machine-local absolute path prefix remains in scoped `.hawp/work/**` artifacts outside this task plan's preserved quoted input/context.
- Verified no machine-local absolute user-home path prefix remains in scoped workflow/work artifacts.
- Confirmed modified files align with TASK-034 scoped paths; unrelated lanes remain unstaged.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [x] Staged-path proof captured before commit
- [x] Only approved TASK-034 paths are staged
- [x] Commit created with focused message

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
