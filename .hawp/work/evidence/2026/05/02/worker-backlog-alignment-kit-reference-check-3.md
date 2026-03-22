# Worker Check — Backlog alignment kit reference

## Summary

Pass — no changes needed.

## What I checked

- `core/BACKLOG_ALIGNMENT.md` — existence check
- `core/.hawp/kit/references/backlog-alignment.md` — existence and content check
- `README.md` — searched for stale path references
- `core/install.md` — searched for stale path references
- `core/update.md` — searched for stale path references
- `.github/copilot-instructions.md` — checked backlog alignment references
- `.github/instructions/intake.instructions.md` — checked backlog alignment references
- `.github/instructions/hawp-backlog-alignment.instructions.md` — verified canonical path
- `.github/prompts/intake.prompt.md` — checked backlog alignment references
- `.github/prompts/hawp-backlog-alignment.prompt.md` — verified canonical path
- `core/.github/copilot-instructions.md` — checked backlog alignment references
- `core/.github/instructions/intake.instructions.md` — checked backlog alignment references
- `core/.github/prompts/intake.prompt.md` — checked backlog alignment references
- `.hawp/work/**` — searched for stale references

## Changes made

None. The move was already completed by a prior worker.

## Findings

- **old path cleanup**: `core/BACKLOG_ALIGNMENT.md` no longer exists (`test -f` returned false).
- **new canonical path**: `core/.hawp/kit/references/backlog-alignment.md` exists, 149 lines, content well-formed.
- **prompt/instruction references**: All active guidance files reference the new path via `hawp-backlog-alignment.instructions.md` and `hawp-backlog-alignment.prompt.md`, which themselves cite `core/.hawp/kit/references/backlog-alignment.md` correctly.
- **install/update safety**: `core/install.md` and `core/update.md` reference backlog alignment via the overlay files (instructions/prompts), not by a direct path to the policy document. No stale path in either.
- **work/scaffold boundary**: No regressions. `.hawp/work/`, `core/.hawp/work/`, and `core/.hawp/kit/` boundaries are intact.
- **stale references in historical records**: `core/BACKLOG_ALIGNMENT.md` appears in `.hawp/work/closed/`, `.hawp/work/status/`, and `.hawp/work/evidence/` files. These are archived historical records and are intentionally not updated.

## Concerns for manager

- None. State is consistent.

## Verification commands

```
test -f core/.hawp/kit/references/backlog-alignment.md → exit 0 (NEW_EXISTS)
test -f core/BACKLOG_ALIGNMENT.md → exit 1 (OLD_GONE)
rg "core/BACKLOG_ALIGNMENT.md|BACKLOG_ALIGNMENT.md" README.md core/install.md core/update.md .github/copilot-instructions.md → no matches
rg "backlog-alignment.md" .github/prompts/hawp-backlog-alignment.prompt.md .github/instructions/hawp-backlog-alignment.instructions.md → matches at correct new path
```
