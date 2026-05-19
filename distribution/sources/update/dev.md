# Update HAWP To Dev Branch

Use this when you want to upgrade to the latest unreleased or in-progress HAWP improvements and test new features.

## When to Use Dev Update

- You are testing the latest HAWP features or bug fixes in development.
- You want to provide feedback on upcoming HAWP improvements before they reach `main`.
- You are actively contributing to HAWP and want to validate your changes in a downstream repository.
- You are willing to work with potentially unstable code.

## Prerequisites

- HAWP already installed in your repository (you have `.hawp/kit/` and `.hawp/work/`).
- `curl` and `tar` (standard Unix tools).
- Ability to run shell scripts or copy/paste a shell command block.
- Understanding that dev branch changes may not be fully stable.

## Before You Update

1. **Back up your `.hawp/work/` if you have sensitive decisions or evidence files**.
   - The update will not touch them, but having a backup is good practice when testing.
   - Specifically backup: `.hawp/work/BACKLOG.md`, `.hawp/work/decisions/`, `.hawp/work/evidence/`.

2. **Check your `.hawp/work/BACKLOG.md`** to see what active and parked items you have.
   - The update script will reconcile closed work based on your backlog.
   - All your work will be preserved.

3. **Note your current HAWP version** (if you know it) so you can easily revert if needed.

## Update Steps

1. **Open your repository in a terminal and move to its root directory.**

2. **Copy and run the command** from the "Update Command (Copy/Paste)" section below.
   - No edits are needed.
   - This guide already pins `REF="dev"`.
   - The script output will show `Source: sentzunhat/human-ai-workflow-protocol@dev`.
   - The command handles migration + update in one run.

   Optional local-source mode for branch testing:
   - If you need to test local or unpushed branch changes, run `export HAWP_LOCAL_CORE="/absolute/path/to/human-ai-workflow-protocol/core"` before running the command.
   - The script will print `Source mode: local core (...)` when that override is active.

3. **Review the output carefully** for any warnings or new features.
   - Look for migration messages or new HAWP features being applied.
   - Check that `.hawp/kit/` was refreshed with dev-branch improvements.
   - Look for any stale folders being removed (expected).

4. **Test the new features**:
   - Read the updated `.hawp/kit/README.md` to understand any new concepts.
   - Try new templates or patterns from `.hawp/kit/templates/` or `.hawp/kit/patterns/`.
   - Provide feedback on the HAWP repository if you find issues.

5. **Verify your work is intact**:
   - Check `.hawp/work/BACKLOG.md` and all your work files.
   - Verify nothing was lost.

## What Gets Updated

- `.hawp/LICENSE` — Apache 2.0 license from dev branch.
- `.hawp/kit/**` — Dev-branch HAWP protocol docs, templates, patterns, examples.
- `.github/instructions/*.instructions.md` — Dev-branch Copilot instructions.
- `.github/prompts/*.prompt.md` — Dev-branch Copilot prompt templates.
- Missing `.hawp/work/` scaffold files from dev branch.

## What Is Preserved

- All of `.hawp/work/` — your backlog, work items, decisions, evidence.
- `.github/copilot-instructions.md` — your custom instructions (not overwritten).
- Your project code and configuration.

## Important: Dev Branch Is Unstable

- Dev branch HAWP improvements may be incomplete or have bugs.
- New features may change or be removed before reaching `main`.
- Do not rely on dev branch changes for production work unless you are actively testing.
- Keep notes on what you like and don't like so you can provide feedback.

## Reverting to Main Branch

If dev branch breaks something or you want to go back to stable HAWP:

1. Use the main branch update guide (`update-main.md`) and run that script.
2. All HAWP-managed files will revert to main branch versions.
3. Your `.hawp/work/` will be unchanged and preserved.
4. Your project remains safe and unaffected.

## Reporting Issues

If you find problems with dev branch HAWP:

1. Note the specific issue (what you did, what went wrong).
2. Check if the problem exists on `main` branch too (run update back to `main` and retry).
3. Open an issue on the HAWP repository with your findings.
4. Include the steps to reproduce and what you expected to happen.
