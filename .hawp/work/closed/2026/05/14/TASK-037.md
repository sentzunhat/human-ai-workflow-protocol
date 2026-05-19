# Bug / Task: Convert distribution workflow into a fail-on-drift gate and document it

**Backlog ID:** TASK-037
**Type:** task
**Reported:** 2026-05-13
**Risk Level:** low

---

### Input (what was reported)

> i love the distribution generate workflow can it be without commiting and have it done a fail if there is a difference please and add this workflow activity on the readme file of the root repo

---

### Context

The repo now has a GitHub Actions workflow for generated distribution guides, but it currently writes a commit back to the branch. The user wants a non-writing workflow that fails when generated outputs drift, plus a root README note so contributors know the workflow exists and what local command fixes drift.

---

### Analysis

**Root cause (or most likely cause):**
`.github/workflows/sync-distribution-generated.yml` is currently a sync-and-commit workflow. That behavior no longer matches the user's preferred workflow contract. The generated footer wording in `librarian/scripts/distribution/shared/composition.ts` also says "Automatic sync", which would become inaccurate once the workflow is converted into a drift gate.

**Directly verified:**

- `.github/workflows/sync-distribution-generated.yml` currently ends with a commit-and-push step.
- `README.md` does not currently mention the GitHub Actions distribution workflow.
- Generated footer text is sourced from `librarian/scripts/distribution/shared/composition.ts`.

**Inferred (not yet proven):**
Replacing the commit step with a diff failure and documenting the workflow plus the local recovery command will match contributor expectations without changing the local `distribution:sync` path.

**Scope — what else is affected:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-037.md`
- `.github/workflows/sync-distribution-generated.yml`
- `librarian/scripts/distribution/shared/composition.ts`
- `README.md`
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
- `.github/workflows/sync-distribution-generated.yml`
- `librarian/scripts/distribution/shared/composition.ts`
- `README.md`
- `distribution/generated/install-main.md`
- `distribution/generated/install-dev.md`
- `distribution/generated/update-main.md`
- `distribution/generated/update-dev.md`

**Parallel work risk:** low
**Can implement now:** yes (explicit user request)

**Coordination note:**
This is a direct follow-up to `TASK-036` in the same file area and stays within the distribution workflow lane.

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
 M .hawp/work/BACKLOG.md
 M .github/workflows/sync-distribution-generated.yml
 M README.md
 M librarian/scripts/distribution/shared/composition.ts
 M distribution/generated/install-dev.md
 M distribution/generated/install-main.md
 M distribution/generated/update-dev.md
 M distribution/generated/update-main.md
```

---

### Options

#### Option A — Convert workflow to drift gate and update docs

Keep the workflow trigger surface, remove write-back behavior, fail on generated drift, update generated footer wording, and document the gate in `README.md`.

#### Option B — Remove the workflow entirely and rely on local commands

Simpler, but loses repo-side feedback on stale generated outputs.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
This preserves the useful automation surface while matching the user's request not to have automation commit to the branch.

**Files to change:**

- `.github/workflows/sync-distribution-generated.yml` — switch from commit/push to fail-on-drift
- `librarian/scripts/distribution/shared/composition.ts` — change footer wording from auto-sync to workflow validation
- `README.md` — add a short note about the workflow and the local recovery command
- `distribution/generated/*.md` — regenerate from updated composition text

**What to verify after:**

- [ ] local distribution sync still succeeds
- [ ] generated footer text matches fail-on-drift behavior
- [ ] README mentions the workflow and the local recovery command
- [ ] workflow file has no write-back step and no write permission requirement

---

### Implementation Notes

Keep the existing workflow filename to avoid unnecessary churn. Use a terminal-visible `git diff` failure message so contributors know to run `npm --prefix librarian run distribution:sync`.

Rebaseline note (2026-05-14): Distribution-generated outputs were moved to `distribution/generated/**` (TASK-046). This plan now tracks the current root distribution layout.

---

## Outcome (filled at close)

Converted the distribution workflow from sync-and-commit to fail-on-drift and documented recovery flow for contributors.

Implemented changes:

- Updated `.github/workflows/sync-distribution-generated.yml` to run distribution sync and fail when `distribution/generated/**` has drift.
- Removed write-back behavior from the workflow (no commit/push step).
- Updated generated-source footer wording in `librarian/scripts/distribution/shared/composition.ts` to describe workflow gating and local sync recovery.
- Added README contributor guidance describing the workflow contract and local fix command.
- Regenerated branch guides under `distribution/generated/**` to match updated composition text.

---

## Verification (filled at close)

- [x] Workflow is fail-on-drift and does not push commits. **Evidence:** `.github/workflows/sync-distribution-generated.yml` has `Fail on generated drift` and no commit/push step.
- [x] Contributor docs mention workflow gate and local recovery. **Evidence:** `README.md` For contributors section references the workflow and `npm run distribution:sync`.
- [x] Generated guide source-reference text matches gate behavior. **Evidence:** `librarian/scripts/distribution/shared/composition.ts` emits "Workflow gate ... fail when generated guides drift".
- [x] Distribution workflow validation remains passing after close. **Evidence:** `npm --prefix librarian run validate:workflow` result `VALIDATION PASS` (warnings only).

---

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [x] Plan file moved to `../closed/YYYY/MM/DD/<ID>.md`
- [x] BACKLOG.md updated
- [ ] Status report written if useful
- [ ] Path-sensitive proof captured before commit if committing

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (user requested implementation)
- [x] Implemented
- [x] Verified
- [x] Closed
