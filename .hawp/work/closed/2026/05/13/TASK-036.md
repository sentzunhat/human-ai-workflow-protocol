# Bug / Task: Automate distribution sync and simplify install/update work-item guidance

**Backlog ID:** TASK-036
**Type:** task
**Reported:** 2026-05-13
**Risk Level:** medium

---

### Input (what was reported)

> the install and update generated branches files are not auto updating or auto installing without intervension can you help me make it more simple and less about the project with more instruction on this is an update or install work item please update the hawp lib or install hawp of the repo and to update the github root folder with the new workflow protocol

---

### Context

Generated install/update branch guides under `core/distribution/generated/` are built from source fragments and a composition script, but keeping them current still depends on a manual build step. The user also wants the guides to read less like project prose and more like direct execution instructions for an install/update work item.

---

### Analysis

**Root cause (or most likely cause):**
The repo has manual `distribution:build` and `distribution:validate` commands but no repo automation that regenerates `core/distribution/generated/*.md` when source fragments change. The shared install/update source content is also verbose and does not start with a direct work-item contract telling an agent to execute install/update to refresh `.hawp` and GitHub-root workflow files.

**Directly verified:**

- `librarian/package.json` exposed `distribution:build` and `distribution:validate` only before this change.
- There was no `.github/workflows/` automation in this repo before this change.
- `core/distribution/sources/shared/install.md` and `core/distribution/sources/shared/update.md` included execution guidance, but it was buried below project-oriented explanatory text.

**Inferred (not yet proven):**
Adding a path-filtered GitHub Actions sync workflow plus a one-step local `distribution:sync` command will keep generated guides current with much less manual intervention. Putting an explicit work-item contract at the top of the shared install/update guides will make agent execution more reliable.

**Scope — what else is affected:**

- `.hawp/work/BACKLOG.md`
- `.github/workflows/sync-distribution-generated.yml`
- `librarian/package.json`
- `librarian/scripts/distribution/shared/composition.ts`
- `core/distribution/sources/shared/install.md`
- `core/distribution/sources/shared/update.md`
- `core/distribution/generated/install-main.md`
- `core/distribution/generated/install-dev.md`
- `core/distribution/generated/update-main.md`
- `core/distribution/generated/update-dev.md`

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `.hawp/work/BACKLOG.md`
- `.github/workflows/sync-distribution-generated.yml`
- `librarian/package.json`
- `librarian/scripts/distribution/shared/composition.ts`
- `core/distribution/sources/shared/install.md`
- `core/distribution/sources/shared/update.md`
- `core/distribution/generated/install-main.md`
- `core/distribution/generated/install-dev.md`
- `core/distribution/generated/update-main.md`
- `core/distribution/generated/update-dev.md`

**Parallel work risk:** medium
**Can implement now:** yes (explicit user approval)

**Coordination note:**
This overlapped the install/update distribution lane already active in `TASK-032`, but stayed on the same area with the same agent and narrowed scope to automation plus execution-first guide wording.

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
 M .github/instructions/code-style.instructions.md
 M .hawp/kit/references/standards-precedence.md
 M .hawp/kit/templates/code-style.md
 M .hawp/kit/templates/documentation.md
```

---

### Options

#### Option A — Add sync automation and sharpen shared guidance

Add a branch-filtered GitHub Actions workflow that regenerates committed distribution files when source fragments change, add a one-step local sync script, and simplify the shared install/update wording so the generated guides lead with an execution contract.

#### Option B — Rewrite docs only

Improve wording in the generated guides but keep regeneration manual. This helps readers but does not fix the stale-generated-file problem.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
This addresses both failure modes the user described: stale branch guides and unclear install/update execution instructions.

**Files to change:**

- `.github/workflows/sync-distribution-generated.yml` — auto-regenerate committed branch guides on `main`/`dev` source changes
- `librarian/package.json` — add one-step local `distribution:sync`
- `librarian/scripts/distribution/shared/composition.ts` — clarify generated footer with auto-sync/local-sync guidance
- `core/distribution/sources/shared/install.md` — add direct install work-item contract
- `core/distribution/sources/shared/update.md` — add direct update work-item contract
- `core/distribution/generated/*.md` — regenerate from updated sources

**What to verify after:**

- [x] `distribution:sync` succeeds locally
- [x] distribution validation still passes after regeneration
- [x] generated install/update guides start with a direct execution contract
- [x] generated guides mention auto-sync or local sync instead of implying manual branch-file editing

---

### Implementation Notes

Kept workflow triggers narrow so generated-file commits do not loop on `core/distribution/generated/**` changes. Reused the existing build/validate commands instead of adding new TypeScript sync logic.

---

## Outcome (filled at close)

Implemented a repo-side sync path plus simpler install/update guide framing.

Completed changes:

- Added `.github/workflows/sync-distribution-generated.yml` to rebuild and commit `core/distribution/generated/**` on `main` or `dev` pushes when distribution sources change.
- Added `distribution:sync` to `librarian/package.json` so local refresh is one command instead of separate build and validate steps.
- Updated `librarian/scripts/distribution/shared/composition.ts` so every generated guide footer says the file is generated and points to automatic sync and local sync.
- Added direct work-item contract sections to `core/distribution/sources/shared/install.md` and `core/distribution/sources/shared/update.md` so the guides tell an agent to execute install/update rather than review content.
- Regenerated all four branch-specific generated guides.

---

## Verification (filled at close)

- [x] Local one-step sync works. **Evidence:** `npm --prefix librarian run distribution:sync` rebuilt all four generated guides and then reported `distribution validation passed: generated outputs are current`.
- [x] Generated guides now expose the execution-first wording. **Evidence:** `core/distribution/generated/install-main.md`, `core/distribution/generated/install-dev.md`, `core/distribution/generated/update-main.md`, and `core/distribution/generated/update-dev.md` each contain `Install Work Item Contract` or `Update Work Item Contract`.
- [x] Generated guides now explain how sync happens. **Evidence:** the same generated guides contain `Automatic sync: push source changes on main or dev...` and `Local sync: run npm --prefix librarian run distribution:sync...`.
- [x] Touched implementation files are structurally clean. **Evidence:** diagnostics reported no errors for `librarian/scripts/distribution/shared/composition.ts` and `.github/workflows/sync-distribution-generated.yml`.
- [x] GitHub-hosted auto-sync run verified on `dev`. **Evidence:** `.github/workflows/sync-distribution-generated.yml` run `25897172755` succeeded for commit `0da2ffd7e7fc65b461ad273c4173c7a10edd36c1` (`https://github.com/sentzunhat/human-ai-workflow-protocol/actions/runs/25897172755`). Evidence: ../evidence/../../../../evidence/2026/05/15/TASK-065-github-hosted-distribution-sync-proof.md

---

## Close Checklist

**Before marking this task done in BACKLOG.md, verify ALL:**

- [x] Outcome section filled (what was actually implemented)
- [x] Verification section filled (all checks listed, each with direct evidence or `unproven` tag)
- [x] All evidence files referenced exist in `../evidence/YYYY/MM/DD/` or are noted as inline
- [x] Plan file moved to `../closed/YYYY/MM/DD/<ID>.md`
- [x] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [x] Status report written (non-trivial change with one unproven GitHub-hosted step)
- [x] Decision file not needed (no separate design decision beyond the implementation plan)
- [x] Staged-path proof captured before commit:
  - [x] `git diff --name-status`
  - [x] `git diff --check`
  - [x] `git diff --cached --name-status`
  - [x] `git diff --cached --check`
  - [x] `git status --short`
- [x] If basename-only paths appeared during path-sensitive work, mark report unsafe and restart after correction (not needed; path discipline remained intact)

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (user requested implementation)
- [x] Implemented
- [x] Verified
- [x] Closed
