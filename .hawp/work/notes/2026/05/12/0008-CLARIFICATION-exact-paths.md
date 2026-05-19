# Clarification: Exact File Paths & Roles (Pre-Phase-1)

**Date:** 2026-05-12
**Purpose:** Prove which exact files exist, are tracked, and will be modified in Phases 1-3

---

## Executive Summary

**The review contained imprecise language.** This document fixes that with exact paths and proof.

**Key correction:**

- ❌ WRONG: "root files" (ambiguous)
- ✅ RIGHT: "`core/install.md` and `core/update.md`" (exact paths)
- ❌ CLARIFICATION: `install.md` and `update.md` at repo root **do not exist** and were never mentioned as targets

---

## File Existence Proof

**Command run:**

```bash
git ls-files | grep install.md
git ls-files | grep update.md
ls -la install.md update.md core/install.md core/update.md 2>&1
```

**Results:**

### FILES THAT EXIST & ARE TRACKED

| Path                                          | Exists | Tracked | Size | Role                                       |
| --------------------------------------------- | ------ | ------- | ---- | ------------------------------------------ |
| `core/install.md`                             | ✅ YES | ✅ YES  | 29k  | Generator script source (scriptSourceFile) |
| `core/update.md`                              | ✅ YES | ✅ YES  | 21k  | Generator script source (scriptSourceFile) |
| `core/distribution/sources/shared/install.md` | ✅ YES | ✅ YES  | N/A  | Source fragment (composed into guides)     |
| `core/distribution/sources/shared/update.md`  | ✅ YES | ✅ YES  | N/A  | Source fragment (composed into guides)     |

### FILES THAT DO NOT EXIST

| Path                     | Exists | Tracked | Role           |
| ------------------------ | ------ | ------- | -------------- |
| `install.md` (repo root) | ❌ NO  | ❌ NO   | DOES NOT EXIST |
| `update.md` (repo root)  | ❌ NO  | ❌ NO   | DOES NOT EXIST |

---

## Generator Dependency Proof

**File:** `librarian/scripts/distribution/shared/composition.ts`

**DISTRIBUTION_PLAN variant definitions:**

```
Line 19:  scriptSourceFile: "core/install.md",   ← for install-main
Line 30:  scriptSourceFile: "core/install.md",   ← for install-dev
Line 41:  scriptSourceFile: "core/update.md",    ← for update-main
Line 52:  scriptSourceFile: "core/update.md",    ← for update-dev
```

**Generator reads exactly:** `core/install.md` and `core/update.md`

**Generator does NOT read:** `install.md` or `update.md` at repo root (they don't exist)

---

## User-Facing Reference Proof

**File:** `README.md`

**User-facing table (lines 9-10):**

```markdown
| **Install HAWP** into a repo (copy/paste) | → [core/install.md](core/install.md) |
| **Update HAWP** to latest `main` | → [core/update.md](core/update.md) |
```

**Contributor section (line 27):**

```markdown
- Authoritative scripts live in [core/install.md](core/install.md) and [core/update.md](core/update.md).
```

**User currently accesses:** `core/install.md` and `core/update.md`

**User does NOT access:** `install.md` or `update.md` at repo root (links don't exist)

---

## Phase-by-Phase File Operations

### PHASE 1: Create sources, update generator, preserve `core/install.md` & `core/update.md`

**New files created:**

- `core/distribution/sources/install/script.md` ← extract bash from `core/install.md`
- `core/distribution/sources/update/script.md` ← extract bash from `core/update.md`

**Existing files modified:**

- `librarian/scripts/distribution/shared/composition.ts` ← update scriptSourceFile paths (4 lines)

**Existing files PRESERVED (not touched):**

- ✅ `core/install.md` — STILL NEEDED (generator fallback during Phase 1)
- ✅ `core/update.md` — STILL NEEDED (generator fallback during Phase 1)
- ✅ `README.md` — not modified in Phase 1
- ✅ `librarian/README.md` — not modified in Phase 1

**Result after Phase 1:** Generator reads from new `core/distribution/sources/install/script.md` and `core/distribution/sources/update/script.md`, but old files still exist (no deletions yet).

---

### PHASE 2: Update user-facing references in docs

**Files modified:**

- `README.md` ← update user-facing links to point to generated guides
- `librarian/README.md` ← update contributor notes

**Existing files PRESERVED (not touched):**

- ✅ `core/install.md`
- ✅ `core/update.md`

**Result after Phase 2:** README points users to generated branch-specific guides; old files still exist.

---

### PHASE 3: Remove proven-obsolete files

**Files deleted:**

- ❌ `core/install.md` — NOW SAFE (generator no longer reads it)
- ❌ `core/update.md` — NOW SAFE (generator no longer reads it)

**Files NOT deleted (only created in Phase 1):**

- ✅ `core/distribution/sources/install/script.md` — KEPT (generator still reads this)
- ✅ `core/distribution/sources/update/script.md` — KEPT (generator still reads this)

**Verification:**

- Run `npm run distribution:build` → must still pass (uses new source files)
- Grep for remaining references → should find none (except in generated guides' source sections)

**Result after Phase 3:** Old files deleted, build system self-contained in `core/distribution/sources/`.

---

## Critical Blockers to Check Before Phase 1

Phase 1 may **only** proceed after confirming:

1. ✅ `0008-install-update-distribution-review.md` has been updated with this clarification
2. ✅ User confirms exact paths are now clear
3. ✅ User confirms understanding that `core/install.md` and `core/update.md` must exist during Phase 1 (not deleted prematurely)
4. ✅ User confirms Phase 1 creates NEW files in `core/distribution/sources/` (not moves)
5. ✅ User confirms Phase 3 deletes `core/install.md` and `core/update.md` (only after Phase 1 & 2)
6. ✅ User confirms no untracked `shared_standards/` modifications are required

---

## Summary

**Before Phase 1:**

- Generator reads: `core/install.md`, `core/update.md`
- README points to: `core/install.md`, `core/update.md`
- Files at repo root: NONE (never existed)

**After Phase 1:**

- Generator reads: `core/distribution/sources/install/script.md`, `core/distribution/sources/update/script.md`
- README points to: `core/install.md`, `core/update.md` (unchanged in Phase 1)
- Old files: Still exist (safe to keep until Phase 3)

**After Phase 2:**

- Generator reads: `core/distribution/sources/install/script.md`, `core/distribution/sources/update/script.md`
- README points to: `core/distribution/generated/install-main.md`, etc. (users now see generated guides)
- Old files: Still exist (safe to delete in Phase 3)

**After Phase 3:**

- Generator reads: `core/distribution/sources/install/script.md`, `core/distribution/sources/update/script.md`
- README points to: `core/distribution/generated/install-main.md`, etc.
- Old files: DELETED (`core/install.md`, `core/update.md` removed)
- Build system: Self-contained in `core/distribution/sources/`
