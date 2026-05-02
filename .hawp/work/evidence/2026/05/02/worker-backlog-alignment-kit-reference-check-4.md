# Worker Check — Backlog alignment kit reference

## Summary

Pass. Another worker had already completed the primary move and most reference updates before this worker ran. This worker verified full alignment and performed no changes (all were already done).

## What I checked

- `core/BACKLOG_ALIGNMENT.md` — existence check
- `core/.hawp/kit/references/backlog-alignment.md` — existence and content
- `.github/instructions/hawp-backlog-alignment.instructions.md` — reference line
- `.github/prompts/hawp-backlog-alignment.prompt.md` — reference line
- `core/.github/instructions/hawp-backlog-alignment.instructions.md` — reference line
- `core/.github/prompts/hawp-backlog-alignment.prompt.md` — reference line
- `README.md` — old path mention
- `.github/copilot-instructions.md` — old path mention
- `core/.github/copilot-instructions.md` — old path mention
- `.hawp/work/**` historical records — reviewed, not modified

## Changes made

None. Previous worker(s) already completed the migration and reference updates.

## Findings

- old path cleanup: `core/BACKLOG_ALIGNMENT.md` does not exist. Move already done.
- new canonical path: `core/.hawp/kit/references/backlog-alignment.md` — exists with full policy content.
- prompt/instruction references:
  - Root `.github/instructions` and `.github/prompts` → `core/.hawp/kit/references/backlog-alignment.md` ✅
  - Core `core/.github/instructions` and `core/.github/prompts` → `.hawp/kit/references/backlog-alignment.md` (downstream installed path) ✅
  - `README.md` — no stale reference ✅
  - `copilot-instructions.md` (root and core) — reference the instruction file, not the policy directly ✅
- install/update safety: no install/update reference to old path found. `core/install.md` and `core/update.md` do not reference the policy file by path.
- work/scaffold boundary: no regression. `core/.hawp/work/` and `.hawp/work/` boundaries unchanged.

## Concerns for manager

- `worker-verification-hawp-alignment.md` (evidence from a prior worker) references `core/BACKLOG_ALIGNMENT.md` at lines 35 and 88. This is an archived evidence record, not active guidance. No change needed.
- Closed plan records under `.hawp/work/closed/` reference the old path. These are historical records, not active guidance. No change needed.
- Status reports under `.hawp/work/status/` reference the old path. These are archived artifacts. No change needed.

## Verification commands

- `test -f core/.hawp/kit/references/backlog-alignment.md` → exit 0 (exists)
- `test ! -f core/BACKLOG_ALIGNMENT.md` → exit 0 (not present)
- grep for `BACKLOG_ALIGNMENT` in active docs (README.md, .github, core/.github) → no matches in active guidance
- grep for `backlog-alignment.md` in `.github` and `core/.github` → all 4 instruction/prompt files point to new path
