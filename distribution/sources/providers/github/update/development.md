# Update HAWP — GitHub/Copilot Provider (Development Branch)

Upgrade to the latest unreleased HAWP kit and GitHub Copilot overlays from the `development` branch.

## When to Use Development Update

- Testing latest HAWP features before `main`.
- Contributing to HAWP and validating changes in a downstream repo.

## Prerequisites

- HAWP already installed.
- Willingness to use potentially unstable Development branch content.

## Before You Update

1. Back up `.hawp/work/` if needed (update does not overwrite it).
2. Check `.hawp/work/BACKLOG.md` for active items.

## Update Steps

1. Repository root in terminal.
2. Run the **Update Command (Copy/Paste)** block (`REF="development"`).
3. Review output and test new kit content.
4. Verify work files intact.

## What Gets Updated

- `.hawp/LICENSE`, `.hawp/kit/**` (Development branch)
- `.github/instructions/`, `.github/prompts/`, `.github/copilot-instructions.md`

## What Is Preserved

- `.hawp/work/**`, project code

## Reverting to Main

Use `distribution/generated/github/update/main.md` and run that script.

## Reporting Issues

Open an issue on the HAWP repository with repro steps and branch (`development`).
