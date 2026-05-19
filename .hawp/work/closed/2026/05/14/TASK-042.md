# Bug / Task: Add explicit branch execution command for install/update

**Backlog ID:** TASK-042
**Type:** task
**Reported:** 2026-05-14
**Risk Level:** low

---

### Input (what was reported)

> the biggest issue is not installing the update or installing a new hawp on this repo how can you fix that add an explisit command where treat it as a work item to update/install the hawp library from the branch

---

### Context

Install/update guidance already says to treat execution as a work item, but users can still miss the concrete action and stop at reading docs. We need one explicit branch-based command that performs install vs update work-item execution directly.

---

### Analysis

**Root cause (or most likely cause):**
The docs rely on conceptual language and separate generated guides, but do not include a single explicit dispatcher command that tells users and agents exactly how to run install vs update from a selected branch.

**Directly verified:**
- `core/distribution/sources/shared/install.md` and `core/distribution/sources/shared/update.md` contain execution-contract text but no explicit branch dispatcher command.
- Generated guides already include branch-pinned script blocks, but users can still skip command execution and remain in analysis-only mode.

**Inferred (not yet proven):**
A single command that selects install/update by `.hawp/` presence and branch will reduce execution misses.

**Scope — what else is affected:**
- `core/distribution/sources/shared/install.md`
- `core/distribution/sources/shared/update.md`
- Regenerated output under `core/distribution/generated/`

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `core/distribution/sources/shared/install.md`
- `core/distribution/sources/shared/update.md`

**Parallel work risk:** medium
**Can implement now:** yes

**Coordination note:**
No currently active task explicitly targets these two shared concept files in this pass; scope is narrow and documentation-only.

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
(no output)
```

Machine-local absolute path prefixes were redacted as `<repo-root-abs>` before persistence.

---

### Options

#### Option A — Add explicit dispatcher command in shared install/update docs

Add a single branch-aware command users can paste to execute install/update work item directly, plus brief verification notes.

#### Option B — Only strengthen prose instructions

Clarify wording without adding executable command; lower maintenance but may not solve execution misses.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
User request explicitly asks for an explicit command. A branch-aware command removes ambiguity and gives a direct action path.

**Files to change:**

- `core/distribution/sources/shared/install.md` — add explicit work-item command section.
- `core/distribution/sources/shared/update.md` — add matching explicit command section.
- `core/distribution/generated/*` — regenerate distribution output.

**What to verify after:**

- [x] New explicit command appears in generated install/update guides.
- [x] Distribution sync runs successfully and only expected files change.

---

### Implementation Notes

Keep command copy/paste ready, branch-selectable via `REF`, and safe by checking `.hawp/` existence to dispatch to install or update mode.

---

## Outcome (filled at close)

Added an explicit branch-aware execution command section to both shared install/update concept docs. The command now dispatches to install or update based on `.hawp/` presence and executes the selected branch guide (`main` or `dev`) as a real work item.

Regenerated distribution outputs so the new section is present in all four generated guides:

- `core/distribution/generated/install-main.md`
- `core/distribution/generated/install-dev.md`
- `core/distribution/generated/update-main.md`
- `core/distribution/generated/update-dev.md`

---

## Verification (filled at close)

- [x] Claim 1: Explicit branch-aware command exists in source shared docs. **Evidence:** `rg -n` shows `## Explicit Work Item Command (Branch-Aware)` in `core/distribution/sources/shared/install.md` and `core/distribution/sources/shared/update.md`.
- [x] Claim 2: Generated install/update guides include the new command section. **Evidence:** `rg -n` matches in `core/distribution/generated/install-main.md`, `core/distribution/generated/install-dev.md`, `core/distribution/generated/update-main.md`, and `core/distribution/generated/update-dev.md`.
- [x] Claim 3: Distribution generation and validation succeeded. **Evidence:** `npm --prefix librarian run distribution:sync` reported `distribution build complete: 4/4 file(s) updated` and `distribution validation passed: generated outputs are current`.
- [x] Claim 4: Staged-path proof captured for close readiness. **Evidence:**
  - `git diff --name-status` lists only expected files.
  - `git diff --check` produced no output.
  - `git diff --cached --name-status` produced no output.
  - `git diff --cached --check` produced no output.
  - `git status --short` matches expected modified/untracked files.

---

## Close Checklist

**Before marking this task done in BACKLOG.md, verify ALL:**

- [x] Outcome section filled (what was actually implemented)
- [x] Verification section filled (all checks listed, each with direct evidence or "unproven" tag)
- [x] All evidence files referenced exist in `../evidence/YYYY/MM/DD/` or are noted as inline
- [x] Plan file moved to `../closed/YYYY/MM/DD/<ID>.md`
- [x] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [ ] Status report written (optional: only if non-trivial OR if something remains unproven OR if a decision/pattern emerged)
- [x] Decision file created if applicable (not applicable for this task)
- [x] Staged-path proof captured before commit:
  - [x] `git diff --name-status`
  - [x] `git diff --check`
  - [x] `git diff --cached --name-status`
  - [x] `git diff --cached --check`
  - [x] `git status --short`
- [x] If basename-only paths appeared during path-sensitive work, mark report unsafe and restart after correction (fresh path-locked pass after two failures)

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
