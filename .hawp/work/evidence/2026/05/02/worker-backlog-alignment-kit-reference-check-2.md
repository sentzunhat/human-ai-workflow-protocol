# Worker Check — Backlog alignment kit reference

Date: 2026-05-02
Worker: alignment-pass-2

## Summary

Pass — no changes needed.

## What I checked

- `core/BACKLOG_ALIGNMENT.md` — existence
- `core/.hawp/kit/references/backlog-alignment.md` — existence and is canonical
- `README.md`, `core/install.md`, `core/update.md` — path references
- `.github/copilot-instructions.md` — path references
- `.github/instructions/hawp-backlog-alignment.instructions.md` — path references
- `.github/prompts/hawp-backlog-alignment.prompt.md` — path references
- `core/.github/instructions/**` — path references
- `core/.github/prompts/**` — path references
- `.hawp/work/BACKLOG.md` — path references
- `.hawp/work/status/2026/05/02/**` — stale reference check
- `.hawp/work/closed/2026/05/02/**` — stale reference check
- `.hawp/work/evidence/2026/05/02/**` — stale reference check

Search patterns used:

- `core/BACKLOG_ALIGNMENT.md`, `BACKLOG_ALIGNMENT.md`, `backlog_alignment`
- `backlog-alignment.md`

## Changes made

None.

## Findings

- old path cleanup: `core/BACKLOG_ALIGNMENT.md` is absent. Move already completed by a prior worker.
- new canonical path: `core/.hawp/kit/references/backlog-alignment.md` exists and is the only copy.
- prompt/instruction references: `.github/instructions/hawp-backlog-alignment.instructions.md` and `.github/prompts/hawp-backlog-alignment.prompt.md` both reference `core/.hawp/kit/references/backlog-alignment.md`. Correct.
- core overlay references: `core/.github/instructions/` and `core/.github/prompts/` contain no stale path. Correct.
- install/update safety: `core/install.md` and `core/update.md` contain no stale path. Correct.
- copilot-instructions: `.github/copilot-instructions.md` and `core/.github/copilot-instructions.md` contain no stale path. Correct.
- work/scaffold boundary: `.hawp/work/` and `core/.hawp/work/` were not touched.
- Stale references found only in historical/archived records (`.hawp/work/status/`, `.hawp/work/closed/`, `.hawp/work/evidence/`). Per rules, archived records are not rewritten.

## Concerns for manager

None.

## Verification commands

```
test -f core/.hawp/kit/references/backlog-alignment.md       → success (EXISTS)
test ! -f core/BACKLOG_ALIGNMENT.md                          → success (ABSENT)
rg -n "core/BACKLOG_ALIGNMENT.md|BACKLOG_ALIGNMENT.md|backlog_alignment" README.md core .github .hawp
  → hits in .hawp/work/status/, .hawp/work/closed/, .hawp/work/evidence/ (archived only, expected)
  → no hits in active guidance files
rg -n "backlog-alignment.md" README.md core .github .hawp
  → .github/instructions/hawp-backlog-alignment.instructions.md:22 (correct new path)
  → .github/prompts/hawp-backlog-alignment.prompt.md:21 (correct new path)
```
