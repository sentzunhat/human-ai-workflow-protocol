# Bug / Task: Fix install/update path output and branch-update clarity

**Backlog ID:** TASK-032
**Type:** task
**Reported:** 2026-05-12
**Risk Level:** medium

---

### Input (what was reported)

> number one makes sense please do the planned changes in as a nnew work item
>
> the second one make no sense i miss spoke the previous install and update scripts worked perfecly on main but for some reason the auto generated ones don't like the branched ones like an example update-dev did not auto update from dev.
>
> three it is getting the repo and i dont want the uncommitted local tree i commited to dev and it did not auto update
>
> please do the fixes
>
> What to fix first:
>
> Patch install/update scripts to print full repo-relative paths in reconciliation output (replace basename-only output).
> Regenerate distribution outputs so update-main/update-dev include the corrected path reporting.
> Clarify in docs: update-dev tracks dev only; update-main tracks main.
>
> dont assume please

---

### Context

Install/update script output currently prints basename-only reconciliation lines, which conflicts with path-lock discipline. User also reported update-dev behavior was unclear and did not visibly prove branch-targeting in output.

---

### Analysis

**Root cause (or most likely cause):**

- Reconciliation output in install/update scripts logs basename-only names via `basename`, not exact repo-relative paths.
- Generated branch guides are in sync with sources, but branch-target behavior is not explicit enough in command output for quick verification.

**Directly verified:**

- `distribution/sources/install/script.md` and `distribution/sources/update/script.md` print basename-only reconciliation lines.
- `update-dev.md` is pinned to `REF="dev"`, `update-main.md` to `REF="main"`.
- Distribution validation passes (no drift between source and generated files before this fix).

**Inferred (not yet proven):**

- Adding explicit source/ref logging and full-path reconciliation output will make branch-target behavior auditable and prevent path-collapse ambiguity.

**Scope — what else is affected:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-032.md`
- `distribution/sources/install/script.md`
- `distribution/sources/update/script.md`
- `distribution/sources/shared/install.md`
- `distribution/sources/shared/update.md`
- `distribution/sources/update/main.md`
- `distribution/sources/update/dev.md`
- `distribution/generated/install-main.md`
- `distribution/generated/install-dev.md`
- `distribution/generated/update-main.md`
- `distribution/generated/update-dev.md`

---

### Work Coordination

**Owner:** agent
**Implementation status:** in-progress
**Overlapping files:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-032.md`
- `distribution/sources/install/script.md`
- `distribution/sources/update/script.md`
- `distribution/sources/shared/install.md`
- `distribution/sources/shared/update.md`
- `distribution/sources/update/main.md`
- `distribution/sources/update/dev.md`
- `distribution/generated/install-main.md`
- `distribution/generated/install-dev.md`
- `distribution/generated/update-main.md`
- `distribution/generated/update-dev.md`

**Parallel work risk:** medium
**Can implement now:** yes (user explicitly approved)

**Coordination note:**
Unrelated modified files from other active lanes will remain untouched.

---

### Recommended Fix

**Option chosen:** A

#### Option A — Patch source scripts and regenerate generated outputs

Update install/update script sources to print full repo-relative reconciliation paths and explicit source ref logging; update shared docs for branch clarity; regenerate branch-specific guides.

#### Option B — Patch generated files only

Directly edit generated files without source updates. Rejected because it would drift again on next regeneration.

**Rationale:**
Source-first update keeps generated outputs consistent and repeatable.

**Files to change:**

- `distribution/sources/install/script.md` — full-path reconciliation output and source/ref logging
- `distribution/sources/update/script.md` — full-path reconciliation output and source/ref logging
- `distribution/sources/shared/install.md` — explain output now includes full source/destination paths
- `distribution/sources/shared/update.md` — explicit branch/ref behavior and output verification guidance
- `distribution/sources/update/main.md` — clarify main-branch ref usage
- `distribution/sources/update/dev.md` — clarify dev-branch ref usage
- `distribution/generated/*.md` — regenerated via build script

**What to verify after:**

- [ ] Distribution build succeeds and rewrites expected generated outputs
- [ ] Distribution validate passes after regeneration
- [ ] Generated update-main/update-dev include corrected reconciliation output format
- [ ] Generated docs clearly state branch mapping (`update-main` -> main, `update-dev` -> dev)

---

### Repo-Root Proof

```bash
pwd
human-ai-workflow-protocol

git rev-parse --show-toplevel
human-ai-workflow-protocol

git rev-parse --show-prefix


git status --short
 M .github/instructions/intake.instructions.md
 M .hawp/kit/patterns/evidence-discipline.md
 M .hawp/kit/usage/intake-workflow.md
 M .hawp/work/BACKLOG.md
 D .hawp/work/active/0007.md
 D .hawp/work/active/HAWP-BACKLOG-VALIDATE-PLAN.md
 M .hawp/work/active/TASK-030-files.md
 M core/.hawp/kit/instructions/da-file-tracking.md
 M core/.hawp/kit/references/work-item-file-tracking.md
 M core/.hawp/kit/templates/work-item-files.md
?? .hawp/kit/references/install-update-safety.md
?? .hawp/kit/references/work-item-file-tracking.md
?? .hawp/kit/templates/adr-template.md
?? .hawp/kit/templates/work-item-files.md
?? .hawp/work/active/TASK-031.md
?? .hawp/work/closed/2026/05/12/0007.md
?? .hawp/work/closed/2026/05/12/0008-CLARIFICATION-exact-paths.md
?? .hawp/work/closed/2026/05/12/0008-install-update-distribution-review.md
?? .hawp/work/closed/2026/05/12/HAWP-BACKLOG-VALIDATE-PLAN.md
?? .hawp/work/closed/2026/05/12/TASK-030-files.md
?? shared_standards/
```

---

## Outcome (filled at close)

Implemented source-first fixes for install/update path output and branch clarity.

Completed changes:

- Updated install/update script sources to print reconciliation moves using full repo-relative `source -> destination` paths.
- Added explicit source logging at runtime: `Source: <owner>/<repo>@<ref>` and archive URL.
- Updated source docs to explicitly state branch mapping behavior and output verification guidance.
- Regenerated all branch-specific generated guides from updated sources.
- Aligned top-level `core/install.md` and `core/update.md` script blocks to the same full-path output format to avoid mixed guidance.

## Verification (filled at close)

- [x] Distribution build regenerated expected files. **Evidence:** `npm --prefix librarian run distribution:build` updated all four generated files, then updated install files after doc refinements.
- [x] Distribution validation passes after regeneration. **Evidence:** `npm --prefix librarian run distribution:validate` -> `distribution validation passed: generated outputs are current`.
- [x] Full-path reconciliation output is present in source and generated scripts. **Evidence:** `distribution/sources/install/script.md`, `distribution/sources/update/script.md`, and generated install/update guides now emit `reconciled ...: <src> -> <dest>` and `retired (orphan): <src> -> <dest>`.

Rebaseline note (2026-05-14): Distribution paths were relocated to the repo-root `distribution/**` layout via TASK-046. Historical intent remains the same; file references above are normalized to current layout.

- [x] Branch-target clarity is explicit in docs. **Evidence:** source and generated install/update guides include explicit `REF` mapping and `Source: sentzunhat/human-ai-workflow-protocol@<ref>` verification text.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [ ] Staged-path proof captured before commit
- [ ] Only approved TASK-032 paths are staged
- [ ] Commit created with focused message

Close note (2026-05-14): Task implementation and verification are complete; commit-path checks remain pending until a dedicated commit pass.

**Status:**

- [x] Plan written
- [x] Approved
- [x] Implemented
- [x] Verified
- [x] Closed
