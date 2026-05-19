# Install HAWP From Dev Branch

Use this installation path when you want to test the latest unreleased or in-progress HAWP changes and improvements.

## When to Use Dev Install

- You are actively contributing to HAWP development.
- You want early access to new features or bug fixes in progress.
- You are testing HAWP improvements before they land on `main`.
- You want to provide feedback on new HAWP guidance or tools.

## Prerequisites

- A repository where you want to test HAWP.
- `curl` and `tar` (standard Unix tools).
- Ability to run shell scripts or copy/paste a shell command block.
- Understanding that dev branch changes may not be fully stable.

## Installation Steps

1. **Open your target repository in a terminal and move to its root directory.**

2. **Copy and run the command** from the "Install Command (Copy/Paste)" section below.
   - No edits are needed.
   - This guide already pins `REF="dev"`.
   - The script output will show `Source: sentzunhat/human-ai-workflow-protocol@dev`.
   - The command handles migration + install in one run.

   Optional local-source mode for branch testing:
   - If you need to test local or unpushed branch changes, run `export HAWP_LOCAL_CORE="/absolute/path/to/human-ai-workflow-protocol/core"` before running the command.
   - The script will print `Source mode: local core (...)` when that override is active.

3. **Review the output** to confirm files were added from the dev branch.
   - Check that `.hawp/kit/` was created with dev-branch HAWP documentation.
   - Check that `.github/instructions/` and `.github/prompts/` have dev-branch files.
   - Check that `.hawp/work/BACKLOG.md` exists.

4. **Test the new features or changes**:
   - Read the HAWP kit to understand any new concepts or improvements.
   - Test the new features in your workflow.
   - Open issues or feedback on the HAWP repository if you find problems.

## What Was Added

Same as main branch install, but with dev-branch versions of:

- `.hawp/kit/**` (dev-branch protocol docs, templates, patterns).
- `.github/instructions/` and `.github/prompts/` (dev-branch Copilot overlays).
- Any other HAWP-managed files.

## Important: Dev Branch Is Unstable

- Dev branch changes are in progress and may not be fully tested.
- You may encounter bugs or incomplete features.
- Do not use dev branch for production unless you are actively testing and providing feedback.
- Switch back to `main` branch if you need stability.

## Switching Back to Main

To return to the stable `main` branch:

1. Use the main branch install guide (`install-main.md`) and run that script.
2. All HAWP-managed files will be updated to main branch versions.
3. Your `.hawp/work/` and `.github/copilot-instructions.md` will be preserved.
