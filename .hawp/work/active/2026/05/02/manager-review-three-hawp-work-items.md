# Manager Review — HAWP backlog alignment, install/update wiring, and repo-root .hawp dogfooding

## Summary

Partial. The three work streams are mostly aligned and trust-boundary intent is strong, but compact-backlog adoption is not yet fully applied to this repo's live `.hawp/work/BACKLOG.md` and a few historical/stale references remain in archived work artifacts. Tiny documentation contradictions found during review were corrected.

## Work items reviewed

- Core backlog alignment guidance
- Install/update/docs/prompts wiring
- Repo-root .hawp dogfooding

## Findings

### 1. Three-folder boundary

Status: partial
Evidence:

- Boundary model is explicit in `README.md`, `core/update.md`, and ADR ` .hawp/work/decisions/2026/05/02/adr-root-hawp-work-for-source-repo.md`.
- `core/.hawp/work/` remains scaffold-only (README/BACKLOG/README anchors), not live source-repo history.
- Source-repo live state is now under `.hawp/work/**` and tracked as project-owned.
- Historical references to root `.work` remain inside archived and migration-note artifacts (expected as historical record, but still noisy in broad text scans).

### 2. Install/update safety

Status: pass
Evidence:

- `core/install.md` and `core/update.md` both describe no-overwrite behavior for target `.hawp/work/**`.
- Scripts use no-clobber copy helpers and seed scaffold files only when missing.
- Text and script behavior state downstream install/update sources are from `core/.hawp/*` only.
- Tiny contradiction in `core/install.md` scope-clarification boundary text (root `.work`) was corrected to repo-root `.hawp/work`.

### 3. Backlog compaction

Status: partial
Evidence:

- `core/BACKLOG_ALIGNMENT.md` clearly defines active-index model, Recently Closed cap, archive-by-date, and anti-forever-Done behavior.
- Scaffold backlog template `core/.hawp/work/BACKLOG.md` was corrected to use Recently Closed + Archive sections (no permanent Done section).
- Live source-repo backlog `.hawp/work/BACKLOG.md` still has a long permanent `## Done` table; compaction policy is documented but not fully applied to live repo state.

### 4. Source repo dogfooding

Status: partial
Evidence:

- Repo-root `.hawp/work/` exists with backlog/active/closed/decisions/evidence/notes/status content.
- `.work/` appears removed from working tree (deletions staged in git status output), indicating migration in progress.
- Migration ADR and migration note exist and preserve historical records intent.
- Historical files still reference old `.work` paths; these appear archival, not active policy.

### 5. Scaffold cleanliness

Status: pass
Evidence:

- `core/.hawp/work/**` contains generic scaffold artifacts only:
  - `core/.hawp/work/README.md`
  - `core/.hawp/work/BACKLOG.md`
  - README anchors under `active/`, `parked/`, `closed/`, `decisions/`, `evidence/`, `notes/`
- No source-repo closed histories, source decisions, or project-specific evidence were found in scaffold source.

### 6. Prompt/instruction alignment

Status: partial
Evidence:

- Root and core intake instructions/prompts now point to `.hawp/work/` and explicitly forbid writing repo work into `core/.hawp/work/`.
- Boundary messaging in copilot instructions is mostly aligned.
- Compact-backlog behavior is not yet consistently encoded in all instruction surfaces (for example, wording about tracking open and completed work can still be interpreted as ever-growing backlog tables).

### 7. Links and stale references

Status: partial
Evidence:

- No critical broken path found in key active docs reviewed (`README.md`, `core/install.md`, `core/update.md`, core/root instructions/prompts).
- Corrected stale/incorrect link casing in `core/.hawp/work/README.md` usage links.
- Historical artifacts under `.hawp/work/active` and `.hawp/work/closed` intentionally contain legacy `.work` references; these are not active policy but create noisy grep matches.

## Corrections made

- Updated `.hawp/work/README.md`:
  - corrected heading from `.work/` to `.hawp/work/`
  - updated backlog description to compact active index wording
  - added `status/YYYY/MM/DD/` layout entry
- Updated `core/.hawp/work/README.md`:
  - updated backlog description to compact active index wording
  - corrected usage guide link casing/paths (`init.md`, `intake-workflow.md`, `status-report.md`)
  - added `status/YYYY/MM/DD/` layout entry
- Updated `core/.hawp/work/BACKLOG.md`:
  - switched from permanent Done model to Recently Closed + Archive model
  - added explicit note to keep Recently Closed capped
- Updated `core/install.md`:
  - fixed stale boundary sentence that still described source-repo state as root `.work/`

## Risks / follow-ups

- Create a dedicated backlog-compaction work item to migrate live `.hawp/work/BACKLOG.md` from long `Done` history to capped Recently Closed with archive links.
- Decide and document explicit handling for legacy-path mentions in historical artifacts (keep as historical-is vs annotate once to reduce scan noise).
- Complete and verify `.work/` retirement in git history (ensure migration commit preserves all records and removes legacy path as active surface).

## Final decision

Ready after listed tiny fixes

## Verification performed

- `git status --short`
  - showed migration-state changes, including `.work/**` deletions and `.hawp/**` additions, plus doc updates.
- `git diff --name-status`
  - confirmed modified boundary docs and newly introduced repo-root `.hawp/work` records.
- `rg -n --hidden -uu "\.work|\.hawp/work|core/\.hawp/work|single source of truth|Recently Closed|Done" ...`
  - identified stale references and policy drift locations.
- `find`/file inventory via workspace tools:
  - verified `.hawp/work/**` contains source-repo records.
  - verified `core/.hawp/work/**` remains scaffold-only.
- Manual inspection:
  - `README.md`
  - `core/install.md`
  - `core/update.md`
  - `core/BACKLOG_ALIGNMENT.md`
  - `core/.hawp/work/BACKLOG.md`
  - `core/.hawp/work/README.md`
  - `.hawp/work/README.md`
  - `.hawp/work/BACKLOG.md`
  - `.github/copilot-instructions.md`
  - `.github/instructions/intake.instructions.md`
  - `.github/prompts/intake.prompt.md`
  - `core/.github/instructions/intake.instructions.md`
  - `core/.github/prompts/intake.prompt.md`

Key outcomes from verification:

- Trust boundary language is now consistent in install/update docs after tiny fix.
- Scaffold source is clean and generic.
- Backlog compaction policy exists and is scaffolded, but live source backlog still needs compaction follow-through.
