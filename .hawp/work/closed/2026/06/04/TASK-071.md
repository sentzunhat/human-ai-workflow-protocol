## Task: Mirror all public standards files from attached folders into core .hawp

**Backlog ID:** TASK-071
**Type:** standards-update
**Reported:** 2026-06-04
**Risk Level:** low

---

### Input (what was reported)

> Every file on the path file paths that I send on the folders. If not, check all the files on the folders that I sent the docs and the standards and whatever is public added here to the. Core dot. HAWP folder.

---

### Context

The request expands prior standards work to ensure exhaustive handling of attached folders. This task mirrors every file under `standards/public/**` into `core/.hawp`, and documents classification for `docs/docs/**` files where public status is not explicit.

### Analysis

**Root cause (or most likely cause):**
Prior work promoted selected public standards, but not an explicit full-file mirror of the entire public standards source tree.

**Directly verified:**

- Full file inventory was captured for:
  - `/Users/beltrd/Desktop/projects/sentzunhat/docs/standards/public/**`
  - `/Users/beltrd/Desktop/projects/sentzunhat/docs/docs/**`
- `standards/public/**` contains additional public assets (context and exports) not yet mirrored in core.
- `docs/docs/**` contains project docs without explicit public/private boundaries in path naming.

**Inferred (not yet proven):**

- A full mirror subtree under core is the most reliable way to satisfy "every public file" without destabilizing existing canonical standards index files.

**Scope - what else is affected:**

- `core/.hawp/kit/standards/public/`
- `core/.hawp/kit/standards/README.md`
- `.hawp/work/evidence/2026/06/04/`

---

### Work Coordination

**Owner:** agent
**Implementation status:** in-progress
**Overlapping files:**

- `core/.hawp/kit/standards/public/**`
- `core/.hawp/kit/standards/README.md`
- `.hawp/work/evidence/2026/06/04/TASK-071-public-sync-and-docs-classification.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
This task is additive and does not rewrite existing project code. It only mirrors public standards source files and records docs classification evidence.

**Path discipline:**

- Use exact repo-relative paths from repository root for all files in this plan.
- Basenames alone are unsafe for path-sensitive work unless the file is truly at repository root and explicitly marked as such.

**Repo-root proof (required before path-sensitive edits):**

```bash
pwd
git rev-parse --show-toplevel
git rev-parse --show-prefix
git status --short
```

Recorded output (machine-local prefix redacted):

```text
<repo-root-abs>
<repo-root-abs>

 M .hawp/work/BACKLOG.md
 M core/.hawp/kit/standards/README.md
?? .hawp/work/closed/2026/06/
?? .hawp/work/evidence/2026/06/
?? core/.hawp/kit/standards/docs/
```

---

### Options

#### Option A - Full mirror of `standards/public/**` into a dedicated core subtree

Copy all public standards files into `core/.hawp/kit/standards/public/` preserving source-relative structure.

#### Option B - Continue selective promotion into canonical category folders only

Merge file-by-file into existing category trees and skip source mirror.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
This guarantees every public source file is present in core HAWP and keeps existing canonical folders stable.

**Files to change:**

- `core/.hawp/kit/standards/public/**` - full mirror from source public tree
- `core/.hawp/kit/standards/README.md` - add mirror section and usage note
- `.hawp/work/evidence/2026/06/04/TASK-071-public-sync-and-docs-classification.md` - inventory and classification evidence

**What to verify after:**

- [ ] All `standards/public/**` source files exist under `core/.hawp/kit/standards/public/**`.
- [ ] Standards README documents the mirrored subtree.
- [ ] Docs folder inventory is classified and logged with public-boundary decision.

---

### Implementation Notes

Do not mirror `standards/private/**` or `standards/project-specific/**`. For `docs/docs/**`, classify as project-specific unless explicit public lane evidence is available.

---

## Outcome (filled at close)

- Mirrored all files from `/Users/beltrd/Desktop/projects/sentzunhat/docs/standards/public/**` into `core/.hawp/kit/standards/public/**` while preserving directory structure.
- Updated `core/.hawp/kit/standards/README.md` with a `public/` mirror section that distinguishes full-source mirror content from canonical promoted standards categories.
- Recorded a full audit artifact with source/destination counts and `docs/docs/**` classification.

## Verification (filled at close)


### Evidence Follow-Up

- [ ] Research evidence for: Source and destination file counts match for full public mirror.
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: Mirrored destination contains expected category paths and files.
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: Standards index documents the mirrored subtree.
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: Docs folder was reviewed and classified with public-boundary outcome.
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.

- [x] Source and destination file counts match for full public mirror.
  - **Evidence:** `source_count=32`, `dest_count=32` captured in terminal output and summarized in `.hawp/work/evidence/2026/06/04/TASK-071-public-sync-and-docs-classification.md`.
- [x] Mirrored destination contains expected category paths and files.
  - **Evidence:** sample destination list in `.hawp/work/evidence/2026/06/04/TASK-071-public-sync-and-docs-classification.md`.
- [x] Standards index documents the mirrored subtree.
  - **Evidence:** `core/.hawp/kit/standards/README.md` (`public/ (full source mirror)` section).
- [x] Docs folder was reviewed and classified with public-boundary outcome.
  - **Evidence:** `.hawp/work/evidence/2026/06/04/TASK-071-public-sync-and-docs-classification.md` (`docs_total_count=119`, `docs_public_path_count=0`).

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled (all claims have direct evidence or "unproven" tag)
- [x] Evidence files created if large/complex
- [x] Plan file moved to closed/YYYY/MM/DD/
- [x] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [ ] Status report written (optional)
- [ ] Decision file created if applicable
- [ ] Staged-path proof captured before commit

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
