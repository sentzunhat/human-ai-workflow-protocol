# Install HAWP — GitHub/Copilot Provider (Main Branch)

Stable install of HAWP kit plus GitHub Copilot overlays.

**Source → target mapping:**

| `core/providers/.github/` | Your repo |
|---------------------------|-----------|
| `instructions/*.instructions.md` | `.github/instructions/` |
| `prompts/*.prompt.md` | `.github/prompts/` |
| `copilot-instructions.md` | `.github/copilot-instructions.md` (seed if missing on install) |

## Prerequisites

- A repository where you want to use HAWP with GitHub Copilot.
- `curl` and `tar` (standard Unix tools).
- Ability to run shell scripts or copy/paste a shell command block.

## Installation Steps

1. **Open your target repository in a terminal and move to its root directory.**

2. **Copy and run the command** from the "Install Command (Copy/Paste)" section below.
   - No edits are needed.
   - This guide already pins `REF="main"`.
   - The script output will show `Source: sentzunhat/human-ai-workflow-protocol@main`.
   - The command handles migration + install in one run.

   Optional local-source mode for branch testing:
   - If you need to test local or unpushed branch changes, run `export HAWP_LOCAL_CORE="/absolute/path/to/human-ai-workflow-protocol/core"` before running the command.
   - The script will print `Source mode: local core (...)` when that override is active.

3. **Review the output** to confirm files were added.
   - Check that `.hawp/kit/` was created with HAWP documentation.
   - Check that `.github/instructions/` and `.github/prompts/` now have HAWP files.
   - Check that `.hawp/work/BACKLOG.md` exists (your project backlog).
   - If the script printed `reconciled (link):` or `reconciled (id-fallback):` lines, those are expected — done work items moved from `active/` to `closed/`.

4. **Next steps after install:**
   - Read `.hawp/kit/README.md` to understand the HAWP protocol.
   - Read `.hawp/kit/start-here.md` to shape your first task.
   - Use `.hawp/work/BACKLOG.md` to track your work.

## What Was Added

- `.hawp/LICENSE` — Apache 2.0 license for the HAWP kit content.
- `.hawp/kit/**` — HAWP protocol docs, templates, patterns, examples, references.
- `.github/instructions/*.instructions.md` — HAWP Copilot instructions.
- `.github/prompts/*.prompt.md` — HAWP Copilot prompt templates.
- `.github/copilot-instructions.md` — seeded only if not already present on install.
- `.hawp/work/` scaffold files — seeded once; never overwritten.

## What Was NOT Changed

- `.github/copilot-instructions.md` — your existing file on install (seed skipped if present).
- Any existing files in `.hawp/work/` (partial `.hawp/` setups).
- Your project code or any other repository files.

## Other guides

- Development branch: `distribution/generated/github/install/development.md`
- Cursor: `distribution/generated/cursor/install/main.md`
- Continue: `distribution/generated/continue/install/main.md`
