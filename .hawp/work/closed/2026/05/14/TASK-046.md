# Bug / Task: Move distribution folder to repo root and update references

**Backlog ID:** TASK-046
**Type:** task
**Reported:** 2026-05-14
**Risk Level:** medium

---

### Input (what was reported)

> what we need now is to move the distribution folder from the core folder to the root folder where we can still run the auto gen script and references the files from the new paths, please update all the refences of the files to the new root location of the repo

---

### Context

Distribution guide sources and generated outputs currently live under `core/distribution/`. The repository uses librarian scripts and CI to regenerate and validate those outputs. Paths are hardcoded in docs, workflow filters, and composition code.

---

### Analysis

**Root cause (or most likely cause):**
The distribution layout and references were originally anchored to `core/distribution/`. Moving the folder to repo root requires synchronized updates across source discovery, composition metadata, generated file references, CI drift checks, and docs.

**Directly verified:**

- `librarian/scripts/distribution/shared/composition.ts` resolves and validates using `core/distribution/...` paths.
- `.github/workflows/sync-distribution-generated.yml` tracks `core/distribution/sources/**` and diffs `core/distribution/generated`.
- `README.md` links and contributor guidance point to `core/distribution/...`.
- Source/generated distribution docs include explicit `core/distribution/...` references in implementation/source-reference sections.

**Inferred (not yet proven):**

- A mechanical move plus path updates should preserve generation behavior if build/validate scripts and workflow gate are updated together.

**Scope — what else is affected:**

- `distribution/sources/**` and `distribution/generated/**` after move
- `librarian/scripts/distribution/**`
- `.github/workflows/sync-distribution-generated.yml`
- `README.md`

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `README.md`
- `.github/workflows/sync-distribution-generated.yml`
- `librarian/scripts/distribution/shared/composition.ts`

**Parallel work risk:** medium
**Can implement now:** yes

**Coordination note:**
User explicitly requested immediate implementation. This change overlaps historical distribution-path work items but is a new bounded migration.

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

#### Option A — Alias-only compatibility layer

Keep `core/distribution/` and add tooling aliases to also support `distribution/`.
Trade-off: split source of truth and ongoing maintenance cost.

#### Option B — Full move and path rewrite

Move folder to `distribution/` at repo root and update all references in tooling, workflow, and docs.
Trade-off: larger one-time diff, cleaner long-term structure.

---

### Recommended Fix

**Option chosen:** B
**Rationale:**
Matches requested structure and keeps one canonical location with no compatibility shims.

**Files to change:**

- `distribution/**` (moved from `core/distribution/**`) — new canonical distribution location
- `librarian/scripts/distribution/shared/composition.ts` — source discovery, output roots, and source-reference text
- `.github/workflows/sync-distribution-generated.yml` — trigger paths and drift-check target
- `README.md` — links and contributor guidance
- Any file containing `core/distribution/` references related to distribution generation

**What to verify after:**

- [ ] `npm --prefix librarian run distribution:sync` passes
- [ ] No remaining `core/distribution/` references for distribution paths
- [ ] Generated guide source references point to `distribution/sources/...` and `distribution/generated/...`

---

### Implementation Notes

Do the folder move first, then update references and run distribution sync to normalize generated files.

---

## Outcome (filled at close)

Moved the distribution system from `core/distribution/` to top-level `distribution/` and updated all active references to the new location.

- Relocated directory: `core/distribution/**` -> `distribution/**`
- Updated composition logic in `librarian/scripts/distribution/shared/composition.ts` to read/write from `distribution/` paths and detect repo root using `distribution/sources/shared`.
- Updated docs and workflow references in `README.md`, `librarian/README.md`, and `.github/workflows/sync-distribution-generated.yml`.
- Updated source-fragment references in `distribution/sources/shared/install.md` and `distribution/sources/shared/update.md`.
- Regenerated distribution guides in `distribution/generated/` using the updated source location.

---

## Verification (filled at close)

- [x] Distribution generation still works from new root location: **Evidence:** `npm --prefix librarian run distribution:sync` succeeded with `distribution build complete: 4/4 file(s) updated` and `distribution validation passed: generated outputs are current`.
- [x] TypeScript tooling remains valid after path migration: **Evidence:** `npm --prefix librarian run typecheck` succeeded.
- [x] Old path references removed from active project files: **Evidence:** `rg -n "core/distribution" --glob '!**/.hawp/**' --glob '!librarian/node_modules/**'` returned no matches.

---

## Close Checklist

- [ ] Outcome section filled
- [ ] Verification section filled (all claims have direct evidence or "unproven" tag)
- [ ] Evidence files created if large/complex
- [ ] Plan file moved to closed/YYYY/MM/DD/
- [ ] BACKLOG.md updated
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
- [x] Approved (explicit user request)
- [x] Implemented
- [x] Verified
- [x] Closed
