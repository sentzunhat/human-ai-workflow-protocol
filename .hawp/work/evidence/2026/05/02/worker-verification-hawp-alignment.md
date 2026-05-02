# Worker Verification — HAWP alignment work

## Summary

Partial.

## Checks

### Three-folder boundary

Status:
Pass.
Evidence:

- README states source boundary explicitly: repo work at .hawp/work, reusable package source at core/.hawp/kit, downstream scaffold source at core/.hawp/work.
- core/.hawp/work/README.md marks core/.hawp/work as downstream scaffold source only and points live source-repo work to repo-root .hawp/work.
- .github and core/.github intake guidance both include: write repo work to .hawp/work and do not write repo work to core/.hawp/work.

### Install/update safety

Status:
Pass.
Evidence:

- core/install.md trust-boundary notes explicitly say install copies from core/.hawp, seeds target .hawp/work only from core/.hawp/work scaffold, never copies source-repo .hawp/work, and does not overwrite project-owned target .hawp/work.
- core/update.md update boundary says source-repo .hawp/work is never copied to downstream repos and existing .hawp/work/\*\* is never overwritten.
- Install/update script sections use no-clobber helpers and scaffold seeding only when files are missing.

### Backlog alignment

Status:
Partial.
Evidence:

- core/BACKLOG_ALIGNMENT.md defines compact active index behavior, capped Recently Closed, archive-by-date, and preserve-history rules.
- core/.hawp/work/BACKLOG.md scaffold follows compact model (Active Work + Recently Closed + Archive).
- Live repo .hawp/work/BACKLOG.md still contains a long Done section with historical accumulation instead of capped Recently Closed. This is policy drift, not fixed here by design.

### Source repo dogfooding

Status:
Partial.
Evidence:

- Repo-root .hawp/work exists and contains source-repo state (backlog, active/closed records, decisions, notes, evidence, status).
- Migration ADR and migration-note style records exist in repo work history under `.hawp/work/decisions/**`, `.hawp/work/status/**`, and `.hawp/work/notes/**`.
- Root .work is not present as active guidance surface.
- Follow-up inconsistency: TASK-006 exists in both .hawp/work/active/TASK-006.md and .hawp/work/closed/2026/05/02/TASK-006.md.

### Scaffold cleanliness

Status:
Pass.
Evidence:

- core/.hawp/work contains scaffold artifacts only (README.md, STATUS.md, starter BACKLOG.md, and folder README anchors).
- No source-repo closed items, decisions history, or project-specific evidence files were found in core/.hawp/work.

### Instruction/prompt alignment

Status:
Partial.
Evidence:

- Root instruction/prompt files include .hawp/work destination, core/.hawp/work prohibition, and compact-backlog rules.
- core instruction/prompt files now include the same compact-backlog and boundary guidance.
- Remaining gap: requested reusable compaction prompt file is still missing (`.github/prompts/hawp-compact-backlog.prompt.md`, and core scaffold equivalent).

### Stale references and links

Status:
Partial.
Evidence:

- Stale .work references in active policy docs are largely cleared.
- Remaining .work references appear in historical archived records under .hawp/work/closed and migration notes; treated as historical context, not active guidance.
- Path/link contradiction fixed in this pass: status scaffold references existed in docs/layout but status/README seeding/file was missing.

## Tiny corrections made

- Added missing scaffold file: core/.hawp/work/status/README.md.
- Updated `.github/copilot-instructions.md` with explicit compact-backlog/compaction-item guardrails.
- Updated `.github/instructions/intake.instructions.md` and `core/.github/instructions/intake.instructions.md` with compact-backlog + compaction-item rules.
- Updated `.github/prompts/intake.prompt.md` and `core/.github/prompts/intake.prompt.md` with compact-backlog constraints.

## Concerns for manager

- Live repo backlog model is still full-history Done style at .hawp/work/BACKLOG.md, which conflicts with compact-index policy in core/BACKLOG_ALIGNMENT.md.
- TASK-006 exists in both active and closed paths; manager should decide whether to close cleanly (single canonical location) or keep dual state temporarily with explicit note.
- Requested reusable compaction prompt file is still missing: `.github/prompts/hawp-compact-backlog.prompt.md` (and core scaffold equivalent).

## Recommended manager decision

Needs follow-up work.

## Verification commands run

- git status --short; git diff --stat; git diff --name-status
  - Result: expected .work deletions, .hawp additions, and docs/instructions/install-update edits.
- rg -n --hidden -uu "\.work/|root \.work|single source of truth|Done forever|permanent history" README.md core .github .hawp
  - Result: .work hits primarily in archived historical records; active policy docs aligned to .hawp/work.
- rg -n --hidden -uu "core/\.hawp/work|\.hawp/work|preserve|overwrite|Recently Closed|Done" README.md core .github .hawp
  - Result: trust-boundary and no-overwrite language present across README/install/update/instructions; compact guidance present in core alignment docs and intake overlays.
- find .hawp/work -maxdepth 5 -type f | sort
  - Result: repo-root work state present, including decisions/status/evidence/closed history.
- find core/.hawp/work -maxdepth 5 -type f | sort
  - Result: scaffold-only files, including newly present core/.hawp/work/status/README.md.
