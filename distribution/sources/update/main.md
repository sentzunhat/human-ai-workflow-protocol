# Update HAWP To Main Branch

Use this when you have HAWP already installed and want to upgrade to the latest stable improvements from the main branch.

## Prerequisites

- HAWP already installed in your repository (you have `.hawp/kit/` and `.hawp/work/`).
- `curl` and `tar` (standard Unix tools).
- Ability to run shell scripts or copy/paste a shell command block.

## Before You Update

1. **Check your `.hawp/work/BACKLOG.md`** to see if there are any active or parked work items.
   - The update script will automatically reconcile closed work items based on your backlog.
   - All your work will be preserved.

2. **Review any GitHub Copilot customizations** in `.github/copilot-instructions.md`.
   - If you have customized it, the update will preserve your customizations.
   - HAWP-provided instructions in `.github/instructions/` and `.github/prompts/` will be updated.

## Update Steps

1. **Open your repository in a terminal and move to its root directory.**

2. **Copy and run the command** from the "Update Command (Copy/Paste)" section below.
   - No edits are needed.
   - This guide already pins `REF="main"`.
   - The script output will show `Source: sentzunhat/human-ai-workflow-protocol@main`.
   - The command handles migration + update in one run.

   Optional local-source mode for branch testing:
   - If you need to test local or unpushed branch changes, run `export HAWP_LOCAL_CORE="/absolute/path/to/human-ai-workflow-protocol/core"` before running the command.
   - The script will print `Source mode: local core (...)` when that override is active.

3. **Review the output** to see what was updated.
   - Check that `.hawp/kit/` was refreshed with new protocol docs.
   - Check that stale legacy folders were removed (`.hawp/templates/`, `.hawp/patterns/`, etc. if present).
   - If the script printed `reconciled (link): <src> -> <dest>`, `reconciled (id-fallback): <src> -> <dest>`, or `retired (orphan): <src> -> <dest>` lines, those are expected — done work items were moved from `active/` to `closed/`.
   - The final output will show: `Refreshed: .hawp/LICENSE, .hawp/kit/**, .github/instructions/*, .github/prompts/*`

4. **Verify your work is intact**:
   - Check `.hawp/work/BACKLOG.md` to see if your items are still there.
   - Look for any automatically reconciled items moved to `closed/`.

## What Gets Updated

- `.hawp/LICENSE` — Apache 2.0 license (updated to latest).
- `.hawp/kit/**` — All HAWP protocol docs, templates, patterns, examples.
- `.github/instructions/*.instructions.md` — HAWP Copilot instructions.
- `.github/prompts/*.prompt.md` — HAWP Copilot prompt templates.
- Missing `.hawp/work/` scaffold files are seeded if they don't exist.

## What Is Preserved

- All of `.hawp/work/` — your backlog, active items, parked items, closed work, decisions, and evidence.
- `.github/copilot-instructions.md` — your custom instructions (not overwritten).
- Your project code and configuration.

## Automatic Work Reconciliation

On every run, the update script reads `.hawp/work/BACKLOG.md` and:

- Moves `.hawp/work/active/` files matching Done rows or `done`/`wont-fix` Active Work rows to `.hawp/work/closed/YYYY/MM/DD/`. Prints each move as `reconciled (link): <src> -> <dest>` or `reconciled (id-fallback): <src> -> <dest>`.
- Retires `.hawp/work/active/` files with no matching BACKLOG entry to `closed/YYYY/MM/DD/`. Prints each as `retired (orphan): <src> -> <dest>`.
- Both passes only run when the backlog has at least one data row.

All your work remains safe — no files are deleted, only moved.

## Troubleshooting

If something looks wrong after update:

- All files under `.hawp/work/` are unchanged, so your work is safe.
- If `.github/copilot-instructions.md` was not updated (expected if you customized it), you can manually merge in new guidance from the source.
- If stale folders were removed and you want them back, you can restore them from the HAWP source repository manually.
- Run the update command again; it is safe to re-run.
