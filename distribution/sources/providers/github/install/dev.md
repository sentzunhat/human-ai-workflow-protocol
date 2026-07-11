# Install HAWP — GitHub/Copilot Provider (Dev Branch)

Install HAWP kit plus GitHub Copilot overlays from the unreleased `dev` branch.

## When to Use Dev Install

- You are actively contributing to HAWP development.
- You want early access to new features or bug fixes in progress.
- You are testing HAWP improvements before they land on `main`.

## Prerequisites

- A repository where you want to test HAWP with GitHub Copilot.
- `curl` and `tar` (standard Unix tools).
- Understanding that dev branch changes may not be fully stable.

## Installation Steps

1. **Open your target repository in a terminal and move to its root directory.**

2. **Copy and run the command** from the "Install Command (Copy/Paste)" section below.
   - This guide pins `REF="dev"`.
   - Optional: `export HAWP_LOCAL_CORE="/absolute/path/to/human-ai-workflow-protocol/core"` for local core testing.

3. **Review the output** — confirm `.hawp/kit/`, `.github/instructions/`, `.github/prompts/`, and `.hawp/work/BACKLOG.md`.

4. **Test and provide feedback** on the HAWP repository if you find problems.

## What Was Added

Same as main-branch GitHub install, but dev-branch versions of `.hawp/kit/**` and `.github/**` overlays.

## Important: Dev Branch Is Unstable

Do not use dev branch for production unless you are actively testing.

## Switching Back to Main

1. Use [install/main.md](main.md) (`distribution/generated/github/install/main.md`).
2. HAWP-managed files update to main branch versions; `.hawp/work/` is preserved.
