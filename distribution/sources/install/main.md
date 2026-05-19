# Install HAWP From Main Branch

This is the standard, stable installation of HAWP.

## Prerequisites

- A repository where you want to use HAWP (GitHub or any Git repository).
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
   - If the script printed `reconciled (link): <src> -> <dest>`, `reconciled (id-fallback): <src> -> <dest>`, or `retired (orphan): <src> -> <dest>` lines, those are expected — it moved done work items from `active/` to `closed/`.

4. **Next steps after install:**
   - Read `.hawp/kit/README.md` to understand the HAWP protocol.
   - Read `.hawp/kit/start-here.md` to shape your first task.
   - Use `.hawp/work/BACKLOG.md` to track your work.

## What Was Added

- `.hawp/LICENSE` — Apache 2.0 license for the HAWP kit content.
- `.hawp/kit/**` — HAWP protocol docs, templates, patterns, examples, references.
- `.github/instructions/*.instructions.md` — HAWP Copilot instructions.
- `.github/prompts/*.prompt.md` — HAWP Copilot prompt templates.
- `.github/copilot-instructions.md` — seeded only if not already present.
- `.hawp/work/` scaffold files — `BACKLOG.md`, `STATUS.md`, `README.md`, and sub-folder `README.md` files (seeded once; never overwritten).

## What Was NOT Changed

- Your existing `.github/copilot-instructions.md` (if you have one).
- Any existing files in `.hawp/work/` (if you already had a partial `.hawp/` setup).
- Your project code or any other repository files.
