# Install HAWP — Continue Provider (Dev Branch)

Dev-branch install of HAWP kit plus Continue overlays: `.continue/rules/hawp-*.md`.

Use this when testing unreleased HAWP changes on the `dev` branch.

**Source → target mapping:**

| `core/providers/.continue/` | Your repo |
|-----------------------------|-----------|
| `rules/hawp-*.md` | `.continue/rules/` |

## Prerequisites

- A repository where you use Continue.
- `curl` and `tar`.

## Installation Steps

1. Open your target repository root in a terminal.
2. Run the **Install Command (Copy/Paste)** block below (`REF="dev"`, `PROVIDER="continue"`).
3. Confirm `.hawp/kit/` and `.continue/rules/hawp-*.md` exist.

Optional: `export HAWP_LOCAL_CORE="/absolute/path/to/human-ai-workflow-protocol/core"` for local testing.

## What Was Added

- `.hawp/kit/**`, `.continue/rules/hawp-*.md`, `.hawp/work/` scaffold (when missing).

## What Was NOT Changed

- Non-HAWP rules under `.continue/rules/`, `.hawp/work/**` project records.

## Other guides

- Stable branch: `distribution/generated/continue/install/main.md`
