# Bug / Task: Make generated install/update guides one-shot executable

**Backlog ID:** TASK-040
**Type:** task
**Reported:** 2026-05-14
**Risk Level:** medium

---

### Input (what was reported)

> we can still not run the auto update from the generated branch file for update or install in this repo maybe add something at the top of the generated files from the source files that it is a new work item for installing or updating whatver the file is for the current repo and for update to check if .hawp already exisit on the repo being install and for install make sure there is no .hawp root folder and if there is just run the update command instead keep it simple and straightforward for any physical or digital agent to review and run in one shot without assistance update any docs or rules in the core .hawp and .github and then let me know what the changers are or were in a summarized way keep it simple scalabel and extensible

---

### Context

Install/update generated guides are branch-specific user entry points. Current scripts run but do not enforce simple preflight mode checks for install-vs-update and can produce shell glob warnings in zsh cleanup paths.

---

### Analysis

**Root cause (or most likely cause):**
The distribution scripts do not include explicit mode guards:

- install path does not signal that existing `.hawp/` means update mode should be used
- update path does not fail fast when `.hawp/` is missing
- generated guides do not include a top execution preflight checklist framing each run as a work item
- cleanup glob lines can emit zsh "no matches found" noise

**Directly verified:**

- Source scripts and generated guides currently lack the requested preflight section and guardrails.
- Running update previously produced zsh no-match warnings for cleanup globs.

**Inferred (not yet proven):**

- Adding explicit preflight checks and nullglob-safe cleanup should reduce user confusion and increase one-shot reliability.

**Scope — what else is affected:**

- Distribution source docs for shared install/update concepts
- Install and update shell script source blocks
- Regenerated guides in `core/distribution/generated/`

---

### Work Coordination

**Owner:** agent
**Implementation status:** in-progress
**Overlapping files:**

- `core/distribution/sources/shared/install.md`
- `core/distribution/sources/shared/update.md`
- `core/distribution/sources/install/script.md`
- `core/distribution/sources/update/script.md`
- `core/.github/instructions/intake.instructions.md`

**Parallel work risk:** medium
**Can implement now:** only after approval

**Coordination note:**
There may be overlap with existing install/update doc audit tasks. User explicitly requested this implementation now.

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
```

---

### Options

#### Option A — Documentation-only clarification

Add guidance text only in shared docs and generated guides.
Trade-off: scripts can still run in ambiguous modes.

#### Option B — Add script preflight guards plus docs contract

Add install/update preflight checks directly in scripts and concise top-level execution checklist in shared docs. Regenerate guides.
Trade-off: behavior changes are explicit and safer; small script logic increase.

---

### Recommended Fix

**Option chosen:** B
**Rationale:**
This directly enforces expected behavior in one-shot execution while keeping generated guides simple and consistent.

**Files to change:**

- `core/distribution/sources/shared/install.md` — add execution preflight checklist and install/update redirect rule
- `core/distribution/sources/shared/update.md` — add execution preflight checklist and required `.hawp` presence check
- `core/distribution/sources/install/script.md` — add preflight guard for existing `.hawp/` and safer cleanup globs
- `core/distribution/sources/update/script.md` — add preflight guard for missing `.hawp/` and safer cleanup globs
- `core/.github/instructions/intake.instructions.md` — add install/update execution rule for preflight behavior

**What to verify after:**

- [ ] Generated install/update files include top preflight section
- [ ] Install script guard is present and clear
- [ ] Update script guard is present and clear
- [ ] Distribution sync passes and updates generated files

---

### Implementation Notes

Keep wording concise and one-shot oriented for humans and agents. Avoid adding extra workflow complexity.

---

## Outcome (filled at close)

Implemented preflight-first install/update behavior and regenerated all branch guides.

- Added execution-preflight sections to shared install/update source docs so generated guides clearly frame runs as work items and enforce install-vs-update mode selection.
- Added install script preflight messaging when `.hawp/` already exists and kept execution update-compatible in one shot.
- Added update script fail-fast guard when `.hawp/` is missing, with install-first guidance.
- Replaced wildcard cleanup removals with `find ... -delete` to avoid zsh no-match failures.
- Updated intake/safety rule docs in both root and core template trees.
- Rebuilt and validated distribution outputs.

---

## Verification (filled at close)

- [x] Generated guides include preflight sections: **Evidence:** `rg -n "Execution Preflight \(Run First\)" core/distribution/sources core/distribution/generated` returned matches in shared sources and all generated install/update guides.
- [x] Install mode with existing `.hawp/` is explicit and one-shot: **Evidence:** Executed bash block extracted from `core/distribution/generated/install-dev.md`; output included `Preflight: detected existing .hawp/.` and completed successfully.
- [x] Update mode without `.hawp/` fails fast with clear guidance: **Evidence:** Executed bash block extracted from `core/distribution/generated/update-dev.md` in a temporary directory; output included `Preflight: .hawp/ not found in this repository.` and exited with code 1.
- [x] Distribution outputs are synced: **Evidence:** `npm --prefix librarian run distribution:sync` completed with `distribution validation passed: generated outputs are current`.
- [x] Path-sensitive diff quality check: **Evidence:** `git diff --check` returned no output.

---

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled (all claims have direct evidence or "unproven" tag)
- [x] Evidence files created if large/complex
- [x] Plan file moved to closed/YYYY/MM/DD/
- [x] BACKLOG.md updated
- [ ] Status report written (if non-trivial / unproven / decision-bearing)
- [ ] Decision file created (if applicable)
- [ ] Staged-path proof captured before commit:
  - [ ] `git diff --name-status`
  - [ ] `git diff --check`
  - [ ] `git diff --cached --name-status`
  - [ ] `git diff --cached --check`
  - [ ] `git status --short`
- [ ] If basename-only paths appeared during path-sensitive work, mark report unsafe and restart after correction (fresh path-locked pass after two failures)

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
