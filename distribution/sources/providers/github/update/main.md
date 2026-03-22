# Update HAWP — GitHub/Copilot Provider (Main Branch)

Upgrade an existing HAWP installation to the latest stable kit and GitHub Copilot overlays from `main`.

## Prerequisites

- HAWP already installed (`.hawp/kit/` and `.hawp/work/` present).
- `curl` and `tar`.

## Before You Update

1. **Check `.hawp/work/BACKLOG.md`** — reconciliation may move done items from `active/` to `closed/`.
2. **Review `.github/copilot-instructions.md`** — update refreshes this file from the HAWP GitHub provider pack.

## Update Steps

1. **Open your repository root in a terminal.**

2. **Copy and run the command** from the "Update Command (Copy/Paste)" section below (`REF="main"`).

   Optional: `export HAWP_LOCAL_CORE="/absolute/path/to/human-ai-workflow-protocol/core"`.

3. **Review output** — `.hawp/kit/` refreshed; reconciliation lines if any; final `Refreshed:` summary.

4. **Verify `.hawp/work/`** is intact.

## What Gets Updated

- `.hawp/LICENSE`, `.hawp/kit/**`
- `.github/instructions/*.instructions.md`, `.github/prompts/*.prompt.md`
- `.github/copilot-instructions.md` (refreshed on update)
- Missing `.hawp/work/` scaffold files (seeded only when absent)

## What Is Preserved

- All of `.hawp/work/**`
- Your project code and configuration

## Troubleshooting

- `.hawp/work/` is never overwritten.
- Re-run update safely if needed.
- Source pack: `core/providers/.github/` in the HAWP repository.

## Other branches

- Dev update: `distribution/generated/github/update/dev.md`
