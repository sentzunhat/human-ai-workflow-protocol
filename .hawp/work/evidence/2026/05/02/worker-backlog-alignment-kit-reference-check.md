# Worker Check - Backlog alignment kit reference

## Summary

Pass.

## What I checked

- `core/BACKLOG_ALIGNMENT.md`
- `core/.hawp/kit/references/backlog-alignment.md`
- `README.md`
- `core/install.md`
- `core/update.md`
- `.github/copilot-instructions.md`
- `.github/instructions/hawp-backlog-alignment.instructions.md`
- `.github/prompts/hawp-backlog-alignment.prompt.md`
- `core/.github/instructions/hawp-backlog-alignment.instructions.md`
- `core/.github/prompts/hawp-backlog-alignment.prompt.md`
- `core/.hawp/work/**`
- `.hawp/work/**`
- search results for old/new path strings and compact-backlog terms

## Changes made

- Moved:
  - `core/BACKLOG_ALIGNMENT.md` -> `core/.hawp/kit/references/backlog-alignment.md`
- Updated references:
  - `.github/instructions/hawp-backlog-alignment.instructions.md`
  - `.github/prompts/hawp-backlog-alignment.prompt.md`
  - `core/.github/instructions/hawp-backlog-alignment.instructions.md`
  - `core/.github/prompts/hawp-backlog-alignment.prompt.md`

## Findings

- old path cleanup:
  - Active canonical policy file no longer at `core/BACKLOG_ALIGNMENT.md`.
  - Remaining old-path mentions are in historical records under `.hawp/work/closed/`, `.hawp/work/status/`, and prior evidence; treated as archive history, not active guidance.
- new canonical path:
  - Canonical policy now exists at `core/.hawp/kit/references/backlog-alignment.md` (lowercase kebab-case).
- prompt/instruction references:
  - Source-repo live overlays now reference `core/.hawp/kit/references/backlog-alignment.md`.
  - Downstream-shipped overlays now reference `.hawp/kit/references/backlog-alignment.md`.
- install/update safety:
  - `core/install.md` and `core/update.md` still state preserve/no-overwrite behavior for project-owned `.hawp/work/` content.
- work/scaffold boundary:
  - No boundary changes were made to `.hawp/work/` or `core/.hawp/work/`.
  - No source-repo work records were copied into `core/.hawp/work/`.

## Concerns for manager

- None.

## Verification commands

- `git status --short`
  - Confirmed file move/edit set and no unexpected scope expansion by this pass.
- `git diff --stat`
  - Shows focused move/reference updates within existing wider working tree.
- `test -f core/.hawp/kit/references/backlog-alignment.md`
  - `NEW_FILE_OK`
- `test ! -f core/BACKLOG_ALIGNMENT.md`
  - `OLD_FILE_REMOVED_OK`
- `rg -n "core/BACKLOG_ALIGNMENT.md|BACKLOG_ALIGNMENT.md|backlog_alignment" README.md core .github .hawp`
  - No active-guidance hits in README/core/.github; hits are historical archive/evidence records.
- `rg -n "backlog-alignment.md|compact backlog|active index|Recently Closed" README.md core .github .hawp`
  - Confirmed updated active references and compact-backlog language.
- `rg -n "preserve.*\.hawp/work|overwrite.*\.hawp/work|project-owned" README.md core/install.md core/update.md .github core/.github`
  - Confirmed install/update safety language remains intact.
