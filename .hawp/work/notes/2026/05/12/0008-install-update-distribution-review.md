# Install/Update Distribution Final Review

**Date:** 2026-05-12
**Scope:** HAWP install/update documentation and distribution system
**Status:** READY FOR IMPLEMENTATION WITH SOURCE-MODEL CHANGES FIRST

---

## Clarification: Exact File Paths & Roles

**CRITICAL NOTE:** Earlier sections used imprecise language ("root files"). This section proves exact paths.

### Files That EXIST and Are TRACKED

| Path                                          | Exists | Tracked | Role                                                         | Size |
| --------------------------------------------- | ------ | ------- | ------------------------------------------------------------ | ---- |
| `core/install.md`                             | ✅ YES | ✅ YES  | Generator script source (scriptSourceFile in composition.ts) | 29k  |
| `core/update.md`                              | ✅ YES | ✅ YES  | Generator script source (scriptSourceFile in composition.ts) | 21k  |
| `core/distribution/sources/shared/install.md` | ✅ YES | ✅ YES  | Source fragment (assembled into guides)                      | N/A  |
| `core/distribution/sources/shared/update.md`  | ✅ YES | ✅ YES  | Source fragment (assembled into guides)                      | N/A  |
| `install.md` (repo root)                      | ❌ NO  | ❌ NO   | Does not exist                                               | N/A  |
| `update.md` (repo root)                       | ❌ NO  | ❌ NO   | Does not exist                                               | N/A  |

**PROOF:**

```bash
# Files reported by git ls-files | grep:
core/distribution/sources/shared/install.md
core/install.md
core/distribution/sources/shared/update.md
core/update.md

# Files found by ls -la:
.rw-r--r--@ 29k <owner> 12 May 14:49 core/install.md
.rw-r--r--  21k <owner> 12 May 13:37 core/update.md

# No files at repo root:
ls: install.md: No such file or directory
ls: update.md: No such file or directory
```

### Generator Dependencies (scriptSourceFile in composition.ts)

Exact lines from `librarian/scripts/distribution/shared/composition.ts`:

```
Line 19:  scriptSourceFile: "core/install.md",  ← for install-main variant
Line 30:  scriptSourceFile: "core/install.md",  ← for install-dev variant
Line 41:  scriptSourceFile: "core/update.md",   ← for update-main variant
Line 52:  scriptSourceFile: "core/update.md",   ← for update-dev variant
```

**Generator currently reads: `core/install.md` and `core/update.md`**

### User-Facing Documentation References (README.md)

Exact lines from `README.md`:

```
Line 9:  | **Install HAWP** into a repo (copy/paste) | → [core/install.md](../../../../../../core/install.md)         |
Line 10: | **Update HAWP** to latest `main`          | → [core/update.md](../../../../../../core/update.md)           |
Line 27: - Authoritative scripts live in [core/install.md](../../../../../../core/install.md) and [core/update.md](../../../../../../core/update.md).
```

**README currently directs users to: `core/install.md` and `core/update.md`**

---

## Executive Summary

The generated branch-specific install/update guides are **complete and self-contained**, but the distribution system still depends on root `core/install.md` and `core/update.md` files for bash script extraction. To achieve the desired final state (generated guides as only user-facing entry points), two sequential steps are required:

1. **Phase 1 (Source Model):** Move bash script blocks to distribution sources and update the build system
2. **Phase 2 (User-Facing):** Update README and remove root files (after Phase 1 verification)

---

## Current Files Reviewed

**Root-level install/update files:**

- `core/install.md` — 542 lines, contains full install bash script + narrative
- `core/update.md` — 385 lines, contains full update bash script + narrative

**Generated guides (complete, verified self-contained):**

- `core/distribution/generated/install-main.md` — 509 lines
- `core/distribution/generated/install-dev.md` — 524 lines
- `core/distribution/generated/update-main.md` — 546 lines
- `core/distribution/generated/update-dev.md` — 566 lines

**Distribution source fragments:**

- `core/distribution/sources/shared/safety.md`
- `core/distribution/sources/shared/repo-boundaries.md`
- `core/distribution/sources/shared/install.md`
- `core/distribution/sources/shared/update.md`
- `core/distribution/sources/install/main.md`
- `core/distribution/sources/install/dev.md`
- `core/distribution/sources/update/main.md`
- `core/distribution/sources/update/dev.md`

**Build system:**

- `librarian/scripts/distribution/shared/composition.ts` — DISTRIBUTION_PLAN + build logic
- `librarian/scripts/distribution/build/index.ts` — main build script
- `librarian/scripts/distribution/validate/index.ts` — validation script

**User-facing documentation:**

- `README.md` — 139 lines, primary entry point
- `librarian/README.md` — 59 lines, tooling documentation
- `benchmark/README.md` — (not in scope)

**Other references:**

- `shared_standards/public/standards/docs/hawp-install-update-safety.md` — (excluded from scope per user instruction)

---

## Generated Guide Findings

### core/distribution/generated/install-main.md

**Status:** ✅ COMPLETE & SELF-CONTAINED

**Branch correctness:** ✅ YES

- Uses `REF="main"` in shell command
- Clear statement: "This is the standard, stable installation of HAWP"
- Prerequisites section is present and clear
- No ambiguous branch references

**Self-contained:** ✅ YES

- All required context present (prerequisites, steps, what gets added, what isn't changed)
- Verification steps provided
- Next steps after install provided
- No instruction to "see install.md for details"
- Source reference section provided but not required for following the guide

**Assumptions/interventions found:** NONE

- Copy/paste command block requires no edits
- No placeholder values for users to fill
- REF value is already set to "main"

**Required fixes:** NONE — guide is complete as-is

---

### core/distribution/generated/install-dev.md

**Status:** ✅ COMPLETE & SELF-CONTAINED

**Branch correctness:** ✅ YES

- Uses `REF="dev"` in shell command
- Clear section: "Use this installation path when you want to test the latest unreleased or in-progress HAWP changes"
- Explicitly states: "Dev branch changes may not be fully stable"
- Guidance on testing and reverting to main provided

**Self-contained:** ✅ YES

- All required context present
- Prerequisites include "Understanding that dev branch changes may not be fully stable"
- Clear distinction from main branch install
- Switching back to main instructions included

**Assumptions/interventions found:** NONE

**Required fixes:** NONE — guide is complete as-is

---

### core/distribution/generated/update-main.md

**Status:** ✅ COMPLETE & SELF-CONTAINED

**Branch correctness:** ✅ YES

- Uses `REF="main"` in shell command
- Clear: "Use this when you have HAWP already installed and want to upgrade to the latest stable improvements from the main branch"

**Self-contained:** ✅ YES

- Clear distinction from fresh install (prerequisites: "HAWP already installed")
- Before-you-update checklist provided
- Update steps provided
- What gets updated vs preserved sections present
- Automatic work reconciliation behavior documented
- Troubleshooting section provided

**Assumptions/interventions found:** NONE

**Required fixes:** NONE — guide is complete as-is

---

### core/distribution/generated/update-dev.md

**Status:** ✅ COMPLETE & SELF-CONTAINED

**Branch correctness:** ✅ YES

- Uses `REF="dev"` in shell command
- Clear: "Use this when you want to upgrade to the latest unreleased or in-progress HAWP improvements"
- Backup recommendation included
- Explicit warning about instability

**Self-contained:** ✅ YES

- Clear prerequisites and backup guidance
- Testing and feedback section provided
- Revert-to-main instructions included
- Issue reporting guidance provided

**Assumptions/interventions found:** NONE

**Required fixes:** NONE — guide is complete as-is

---

## Distribution Source Model

### How Generated Guides Are Produced

**Composition Engine:** `librarian/scripts/distribution/shared/composition.ts`

**DISTRIBUTION_PLAN Definition:**

```
For each of 4 variants (install-main, install-dev, update-main, update-dev):
  - outputFile: target path in generated/
  - sectionFiles: array of source fragments to compose
  - scriptSourceFile: authoritative bash script source (core/install.md or core/update.md)
  - ref: branch identifier ("main" or "dev")
```

**Build Process Flow:**

1. **Variant Configuration** (`DISTRIBUTION_PLAN`)
   - install-main/dev: scriptSourceFile = `core/install.md`, ref = "main"/"dev"
   - update-main/dev: scriptSourceFile = `core/update.md`, ref = "main"/"dev"

2. **Section Assembly**
   - Read and concatenate all files in `sectionFiles` array
   - Example for install-main: safety.md + repo-boundaries.md + shared/install.md + install/main.md

3. **Bash Script Extraction** — `extractBashBlock()`
   - Read `core/install.md` or `core/update.md`
   - Find bash block containing `REF="main"`
   - Replace `REF="main"` with `REF="<ref>"` (main or dev)
   - Insert extracted/substituted block into guide

4. **Source Reference Section Generation**
   - Lists output file path
   - Lists scriptSourceFile as "authoritative source"
   - Lists all sectionFiles
   - Notes that changes to root source require regeneration

5. **Output**
   - Write assembled guide to `core/distribution/generated/<outputFile>`

### Current Dependency: core/install.md and core/update.md ARE REQUIRED

**Evidence:**

- `composition.ts:DISTRIBUTION_PLAN` explicitly maps `scriptSourceFile: "core/install.md"` and `"core/update.md"`
- `extractBashBlock()` function reads these files directly:
  ```typescript
  const filePath = join(repoRoot, scriptSourceFile);
  if (!existsSync(filePath)) {
    throw new Error(`Script source file not found: ${filePath}`);
  }
  const content = readFileSync(filePath, "utf-8");
  ```
- If either file is deleted, `npm run distribution:build` will fail with "Script source file not found"

---

## File Dependency & Reference Inventory

### core/install.md

**Currently exists:** ✅ YES (29k, tracked)

**Generator dependency:** ⚠️ YES — CRITICAL BLOCKER

- `librarian/scripts/distribution/shared/composition.ts` lines 19, 30: `scriptSourceFile: "core/install.md"`
- Generator calls `extractBashBlock()` which reads this file:
  ```typescript
  const filePath = join(repoRoot, scriptSourceFile); // = "core/install.md"
  const content = readFileSync(filePath, "utf-8");
  ```
- If deleted before Phase 1, `npm run distribution:build` fails with "Script source file not found"

**User-facing references in README.md:**

- Line 9: `[core/install.md](../../../../../../core/install.md)` in "Get started" table
- Line 27: "Authoritative scripts live in [core/install.md](../../../../../../core/install.md)..."

**References in librarian/README.md:**

- Line 12: "Extracts bash blocks from authoritative scripts (`core/install.md`, `core/update.md`)"

**References in generated guides:**

- All 4 guides' "Source Reference" sections list: `core/install.md` — shell script (authoritative source)

**Safe to delete now:** ❌ NO

- Phase 1 required first (extract bash block to sources/)
- Phase 2 required (update README)
- Phase 3: Then can be deleted

---

### core/update.md

**Currently exists:** ✅ YES (21k, tracked)

**Generator dependency:** ⚠️ YES — CRITICAL BLOCKER

- `librarian/scripts/distribution/shared/composition.ts` lines 41, 52: `scriptSourceFile: "core/update.md"`
- Generator calls `extractBashBlock()` which reads this file
- If deleted before Phase 1, `npm run distribution:build` fails with "Script source file not found"

**User-facing references in README.md:**

- Line 10: `[core/update.md](../../../../../../core/update.md)` in "Get started" table
- Line 27: "Authoritative scripts live in ... [core/update.md](../../../../../../core/update.md)"

**References in librarian/README.md:**

- Line 12: "Extracts bash blocks from authoritative scripts (`core/install.md`, `core/update.md`)"

**References in generated guides:**

- All 4 guides' "Source Reference" sections list: `core/update.md` — shell script (authoritative source)

**Safe to delete now:** ❌ NO

- Phase 1 required first (extract bash block to sources/)
- Phase 2 required (update README)
- Phase 3: Then can be deleted

---

### install.md (repo root)

**Currently exists:** ❌ NO

**Tracked:** ❌ NO

**Role:** NONE — This file does not exist and never has in this repository

---

### update.md (repo root)

**Currently exists:** ❌ NO

**Tracked:** ❌ NO

**Role:** NONE — This file does not exist and never has in this repository

---

## Reference Cleanup Plan

### User-Facing Entry Points (README.md)

**Current state:**

```markdown
| **Install HAWP** into a repo (copy/paste) | → [core/install.md](../../../../../../core/install.md) |
| **Update HAWP** to latest `main` | → [core/update.md](../../../../../../core/update.md) |
...

- Stable/default install: [core/distribution/generated/install-main.md](...)
- Stable/default update: [core/distribution/generated/update-main.md](...)
- Dev/test install only: [core/distribution/generated/install-dev.md](...)
- Dev/test update only: [core/distribution/generated/update-dev.md](...)
```

**Problem:** Users see `core/install.md` and `core/update.md` as primary entry points in "Get started" table, then see generated guides as secondary "quick links"

**Proposed change (PHASE 2):**

- Swap the priority: Generated branch-specific guides become primary in "Get started" table
- core/install.md and core/update.md become "For contributors" section only
- Link in contributors section becomes: "Authoritative scripts live in [core/distribution/sources/](../../../../../../core/distribution/sources/) for composition; see [librarian/README.md](../../../../../../librarian/README.md) for how guides are generated"

**Before:**

```
| **Install HAWP** into a repo (copy/paste) | → [core/install.md](../../../../../../core/install.md)         |
| **Update HAWP** to latest `main`          | → [core/update.md](../../../../../../core/update.md)           |
```

**After:**

```
| **Install HAWP** into a repo (copy/paste) | → [core/distribution/generated/install-main.md](core/distribution/generated/install-main.md) |
| **Update HAWP** to latest `main`          | → [core/distribution/generated/update-main.md](core/distribution/generated/update-main.md) |
```

---

### Contributor Documentation (README.md)

**Current state:**

```markdown
### For contributors

- Authoritative scripts live in [core/install.md](../../../../../../core/install.md) and [core/update.md](../../../../../../core/update.md).
- Editable distribution fragments live under [core/distribution/sources/](../../../../../../core/distribution/sources/).
```

**Proposed change (PHASE 2, after generator is updated):**

```markdown
### For contributors

- Distribution fragments live under [core/distribution/sources/](../../../../../../core/distribution/sources/) and are composed into user-facing guides by [librarian/scripts/distribution/](librarian/scripts/distribution/).
- Shell script blocks are maintained in [core/distribution/sources/install/script.md](core/distribution/sources/install/script.md) and [core/distribution/sources/update/script.md](core/distribution/sources/update/script.md).
```

---

### Tooling Documentation (librarian/README.md)

**Current state:**

```
- **`build/index.ts`** — ... Extracts bash blocks from authoritative scripts (`core/install.md`, `core/update.md`) ...
```

**Proposed change (PHASE 1):**

```
- **`build/index.ts`** — ... Extracts bash blocks from authoritative scripts in `core/distribution/sources/` ...
```

---

### Generated Guide Source References

**Current:** Each generated guide includes section:

```
This generated guide is built from:
- `core/install.md` — shell script (authoritative source)
- `core/distribution/sources/shared/safety.md`
...
```

**Proposed (PHASE 2, after files are moved):**

```
This generated guide is built from:
- `core/distribution/sources/install/script.md` — shell script (authoritative source)
- `core/distribution/sources/shared/safety.md`
...
```

---

## Implementation Plan

### ⚠️ CRITICAL: Two Phases Required

**Do NOT delete `core/install.md` or `core/update.md` until Phase 1 is complete.**

---

### PHASE 1: Update Distribution Source Model

**Objective:** Move bash script blocks to sources/ directory, update generator, verify build succeeds

**1.1 Extract bash blocks to new source files**

Create two new files with only the bash script blocks (no narrative):

- Create: `core/distribution/sources/install/script.md`
  - Extract bash block from `core/install.md` with comment explaining it's the install script
  - NO narrative, NO markdown headings except maybe a comment
  - Content: just the bash script block as-is

- Create: `core/distribution/sources/update/script.md`
  - Extract bash block from `core/update.md`
  - NO narrative, NO markdown headings except maybe a comment
  - Content: just the bash script block as-is

**Evidence files:** Save extracted blocks to `.hawp/work/evidence/2026/05/12/phase1-extracted-blocks.md`

**1.2 Update distribution plan configuration**

File: `librarian/scripts/distribution/shared/composition.ts`

- Change line 19 from `scriptSourceFile: "core/install.md"` to `scriptSourceFile: "core/distribution/sources/install/script.md"`
- Change line 30 from `scriptSourceFile: "core/install.md"` to `scriptSourceFile: "core/distribution/sources/install/script.md"`
- Change line 41 from `scriptSourceFile: "core/update.md"` to `scriptSourceFile: "core/distribution/sources/update/script.md"`
- Change line 52 from `scriptSourceFile: "core/update.md"` to `scriptSourceFile: "core/distribution/sources/update/script.md"`

**1.3 Verify build succeeds**

```bash
cd librarian
npm run distribution:build
npm run distribution:validate
npm run typecheck
```

Expected: All passes, generated files updated with new source paths

**1.4 Verify generated guides still work**

- Spot-check: each generated guide still has valid bash command block
- Spot-check: source reference section now lists `core/distribution/sources/install/script.md` and `core/distribution/sources/update/script.md` (not `core/install.md` or `core/update.md`)

**1.5 Commit Phase 1**

```bash
git add librarian/scripts/distribution/shared/composition.ts
git add core/distribution/sources/install/script.md
git add core/distribution/sources/update/script.md
git commit -m "extract install/update scripts to sources for composition"
```

**Result:** Generator no longer depends on `core/install.md` or `core/update.md`; build system is self-contained.

**1.5 Commit Phase 1**

```bash
git add librarian/scripts/distribution/shared/composition.ts
git add core/distribution/sources/install/script.md
git add core/distribution/sources/update/script.md
git commit -m "extract install/update scripts to sources for composition"
```

---

### PHASE 2: Update User-Facing References

**Objective:** Point users to generated branch-specific guides instead of root files

**2.1 Update README.md "Get started" table**

File: `README.md`

Replace:

```markdown
| **Install HAWP** into a repo (copy/paste) | → [core/install.md](../../../../../../core/install.md) |
| **Update HAWP** to latest `main` | → [core/update.md](../../../../../../core/update.md) |
```

With:

```markdown
| **Install HAWP** into a repo (choose your branch) | → [core/distribution/generated/install-main.md](core/distribution/generated/install-main.md) (stable) or [install-dev.md](core/distribution/generated/install-dev.md) (testing) |
| **Update HAWP** to latest | → [core/distribution/generated/update-main.md](core/distribution/generated/update-main.md) (stable) or [update-dev.md](core/distribution/generated/update-dev.md) (testing) |
```

**2.2 Update README.md "For contributors" section**

File: `README.md`

Replace:

```markdown
### For contributors

- Authoritative scripts live in [core/install.md](../../../../../../core/install.md) and [core/update.md](../../../../../../core/update.md).
- Editable distribution fragments live under [core/distribution/sources/](../../../../../../core/distribution/sources/).
- After editing either authoritative scripts or distribution sources, regenerate and validate:
```

With:

```markdown
### For contributors

- Distribution sources live under [core/distribution/sources/](../../../../../../core/distribution/sources/) (fragments and install/update scripts).
- Build system: [librarian/scripts/distribution/](librarian/scripts/distribution/) composes fragments into user-facing guides.
- After editing distribution sources, regenerate and validate:
```

**2.3 Update README.md "Repo layout" section**

File: `README.md`

Replace:

```markdown
- `core/install.md` — install script + boundary notes
- `core/update.md` — update script + migration notes
```

With:

```markdown
- `core/distribution/sources/` — installation/update script blocks and reusable documentation fragments
```

**2.4 Update librarian/README.md**

File: `librarian/README.md`

Replace:

```markdown
- **`build/index.ts`** — ... Extracts bash blocks from authoritative scripts (`core/install.md`, `core/update.md`) ...
```

With:

```markdown
- **`build/index.ts`** — ... Extracts bash blocks from authoritative scripts in `core/distribution/sources/` ...
```

**2.5 Commit Phase 2**

```bash
git add README.md
git add librarian/README.md
git commit -m "point users to branch-specific install and update guides"
```

---

### PHASE 3: Remove or Archive Root Files

**Objective:** Clean up root files after Phases 1 & 2 are verified

**Prerequisites for Phase 3:**

- Phase 1 build verification complete
- Phase 2 README updates in place
- All tests pass
- User can verify by reading Phase 2 README that generated guides are now primary

**3.1 Remove root files**

```bash
rm core/install.md core/update.md
git add -u
git commit -m "remove root install and update files; source moved to distribution/sources"
```

**3.2 Verify**

```bash
npm run distribution:build
npm run distribution:validate
npm run typecheck
npm run validate:workflow
grep -r "core/install.md\|core/update.md" . --include="*.md" --include="*.ts" 2>/dev/null | grep -v generated | grep -v ".git"
```

Expected:

- Build and validate succeed
- No remaining references to core/install.md or core/update.md (except in generated guides' source reference sections)

**3.3 Final commit verification**

```bash
git status --short
git log --oneline -n 3
```

---

## Verification Plan

After **EACH** phase, run:

```bash
# Verify build succeeds
cd librarian
npm run distribution:build

# Verify generated outputs are correct
npm run distribution:validate

# Verify TypeScript compiles
npm run typecheck

# Verify workflow integrity (if available)
npm run validate:workflow

# Verify no dangling references (for Phase 2+)
grep -RIn --exclude-dir=.git --exclude-dir=node_modules \
  'core/install\.md|core/update\.md' . | grep -v "distribution/generated" | grep -v "core/distribution/sources"
```

**Phase 1 success criteria:**

- ✅ `npm run distribution:build` completes with 0 errors
- ✅ `npm run distribution:validate` completes with 0 errors
- ✅ 4 generated guides are updated (git status shows modified generated/\*.md)
- ✅ Generated guide source references now point to `core/distribution/sources/install/script.md` and `core/distribution/sources/update/script.md`

**Phase 2 success criteria:**

- ✅ README.md updated to point users to `core/distribution/generated/` guides
- ✅ `librarian/README.md` updated to reflect new source locations
- ✅ All tests pass
- ✅ git diff shows README changes only (no code changes)

**Phase 3 success criteria:**

- ✅ core/install.md and core/update.md are deleted
- ✅ `npm run distribution:build` still passes (uses sources/ files)
- ✅ No remaining active references to deleted files
- ✅ Users can follow generated guides without seeing root file links

---

## Risks & Human Decisions

### Risk 1: Generator Failure If Phase 1 Not Done First

**Risk:** If root files are deleted before Phase 1, `npm run distribution:build` will fail
**Mitigation:** Phases are sequential; Phase 1 must complete with passing tests before Phase 2
**Decision required:** APPROVED — this report recommends Phase 1 first

---

### Risk 2: Breaking External Links

**Risk:** If external documentation or downstream projects link to `core/install.md` or `core/update.md`, removal breaks them
**Mitigation:** Generated guides include source reference sections, providing a path forward for external readers
**Decision required:** APPROVED — generated guides maintain traceability

---

### Risk 3: Contributor Confusion During Transition

**Risk:** Contributors may not find install/update script sources during the transition
**Mitigation:** Phase 2 README updates explicitly document new source locations in "For contributors" section
**Decision required:** APPROVED — Phase 2 includes contributor documentation updates

---

### Risk 4: Validator/Build Dependency on Old Files

**Risk:** Other scripts might have hardcoded references to core/install.md or core/update.md
**Mitigation:** Comprehensive grep search conducted; no other references found except as documented
**Decision required:** APPROVED — only identified references are in composition.ts (to be updated in Phase 1) and docs

---

## Summary Table

| Item                             | Status | Notes                                                                         |
| -------------------------------- | ------ | ----------------------------------------------------------------------------- |
| Generated guides complete?       | ✅ YES | All 4 guides are self-contained, branch-specific, no manual edits required    |
| Generated guides self-contained? | ✅ YES | No "see core/install.md" cross-references; source sections are reference-only |
| Generator depends on root files? | ⚠️ YES | BLOCKER — Phase 1 required to extract bash to sources/                        |
| README points to root files?     | ⚠️ YES | ISSUE — Phase 2 required to redirect to generated guides                      |
| Safe to remove root files now?   | ❌ NO  | Must complete Phase 1 & 2 first                                               |
| Ready for implementation?        | ✅ YES | If phases are done in order: 1 → 2 → 3                                        |

---

## Final Recommendation

**Status: READY FOR IMPLEMENTATION (phases required)**

The distribution system is well-designed. The generated branch-specific guides are complete and accurate. However, to achieve the desired final state where generated guides are the **only** user-facing entry points:

1. **Execute Phase 1** to move bash blocks to `core/distribution/sources/` and update the generator
2. **Execute Phase 2** to update README and contributor docs to point to generated guides
3. **Execute Phase 3** to remove the now-unused root files

Each phase is small, well-defined, and can be verified independently. **Do not skip Phase 1** — it is the critical prerequisite for removing the root files safely.

After all phases complete:

- ✅ Generated branch-specific guides are the only user-facing install/update entry points
- ✅ Users choose between `main` (stable) and `dev` (testing) directly from README
- ✅ Generated guides are self-contained with no manual intervention needed
- ✅ Distribution sources are consolidated under `core/distribution/sources/`
- ✅ Generator no longer depends on external root files
- ✅ Contributors see new source locations clearly documented
