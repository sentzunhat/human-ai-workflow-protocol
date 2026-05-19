# Bug / Task: Add distribution source folders and repo-local build/validate scripts

**Backlog ID:** TASK-022
**Type:** task
**Reported:** 2026-05-11
**Risk Level:** low

---

### Input (what was reported)

> # HAWP Work Item: Add distribution source folders and repo-local build/validate scripts
>
> We want to add a scalable install/update distribution structure to the HAWP repo.
> Selected structure: core/distribution/** source, dist/distribution/** generated, librarian/scripts/distribution/{build,validate}/\*\* and npm scripts distribution:build, distribution:validate.

---

### Context

HAWP currently has install/update documents under core and one existing TypeScript maintenance script under librarian/scripts/validate-hawp-workflow/. The new request adds a source-fragment distribution model with deterministic generated outputs and CI-friendly validation.

---

### Analysis

**Root cause (or most likely cause):**
There is no dedicated distribution source tree, generator, or freshness validator for install/update bundles, so maintainability and CI enforceability are limited.

**Directly verified:**

- Existing script style is TypeScript entrypoint-first in librarian/scripts/validate-hawp-workflow/.
- npm scripts are defined in librarian/package.json and run with tsx.
- No .github/workflows files currently exist.
- No dist/ folder currently exists.
- Existing install/update docs are core/install.md and core/update.md.

**Inferred (not yet proven):**

- Contributors run maintenance commands from librarian/ package context.

**Scope — what else is affected:**

- core/distribution/\*\* (new source fragments)
- dist/distribution/\*\* (generated outputs)
- librarian/scripts/distribution/\*\* (new tooling)
- librarian/package.json (script wiring)

---

### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-022.md`
- `librarian/package.json`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
Current unstaged edits are in instruction/reference files and do not overlap implementation targets.

---

### Options

#### Option A — Deterministic composable generator

Create a shared composition module used by both build and validate scripts. Build writes files; validate compares expected vs on-disk content and fails with clear paths.

#### Option B — Duplicate composition logic in build and validate

Implement separate logic for each script. Faster initial coding, but higher drift risk and lower maintainability.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
Single composition source minimizes drift and guarantees validate checks the exact same generation logic used by build.

**Files to change:**

- `.hawp/work/BACKLOG.md` — status transitions for TASK-022
- `.hawp/work/active/TASK-022.md` — plan and close evidence
- `core/distribution/**` — source fragments
- `dist/distribution/**` — initial generated artifacts
- `librarian/scripts/distribution/build/**` — build entrypoint
- `librarian/scripts/distribution/validate/**` — validate entrypoint
- `librarian/scripts/distribution/shared/**` — composition helper/types
- `librarian/package.json` — npm scripts

**What to verify after:**

- [x] distribution build command regenerates expected files with deterministic content
- [x] distribution validate command fails on stale/missing outputs and passes when current
- [x] typecheck passes for new scripts
- [x] existing workflow validator command still runs

---

### Implementation Notes

Keep source and generated files explicitly separated. Validation must not write files. Keep HAWP CLI boundary intact by limiting all command wiring to librarian package scripts.

---

## Outcome (filled at close)

Implemented a minimal distribution pipeline with clear source/generated separation.

- Added source fragments under `core/distribution/shared`, `core/distribution/install`, and `core/distribution/update`.
- Added shared deterministic composition logic under `librarian/scripts/distribution/shared/composition.ts`.
- Added repo-local build and validate entrypoints under `librarian/scripts/distribution/build/index.ts` and `librarian/scripts/distribution/validate/index.ts`.
- Added npm scripts in `librarian/package.json`:
	- `distribution:build`
	- `distribution:validate`
- Generated initial outputs under `dist/distribution/` with a generated-file warning header.
- Preserved boundary: no public HAWP CLI changes.

---

## Verification (filled at close)

- [x] Build generated all expected files. **Evidence:** `npm run distribution:build` updated install-main.md, install-dev.md, update-main.md, update-dev.md under `dist/distribution/`.
- [x] Validate passes when outputs are current. **Evidence:** `npm run distribution:validate` returned "distribution validation passed: generated outputs are current".
- [x] TypeScript checks pass for new scripts. **Evidence:** `npm run typecheck` exited successfully.
- [x] Existing workflow validator still passes. **Evidence:** `npm run validate:workflow` returned `Result: VALIDATION PASS` (with existing tolerated legacy WARNs only).
- [x] No diagnostics in new distribution script files. **Evidence:** editor diagnostics for build/validate/composition files report no errors.

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
- [x] Outcome section filled (what was actually implemented)
- [x] Verification section filled (all checks listed, each with direct evidence or "unproven" tag)
- [x] All evidence files referenced exist in `../evidence/YYYY/MM/DD/` or are noted as inline
- [x] Plan file will be moved to `../closed/YYYY/MM/DD/<ID>.md`
- [x] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [ ] Status report written (optional: only if non-trivial OR if something remains unproven OR if a decision/pattern emerged)
- [ ] Decision file created if applicable (only if this task resolves a design question)

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
