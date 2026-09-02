# Install HAWP — Continue Provider (Main Branch)

Stable install of HAWP kit plus Continue overlays: `.continue/rules/hawp-*.md`.

**Source → target mapping:**

| `core/providers/.continue/` | Your repo |
|-----------------------------|-----------|
| `rules/hawp-*.md` | `.continue/rules/` |

## Prerequisites

- A repository where you use Continue (VS Code extension or IDE integration).
- `curl` and `tar`.

## Installation Steps

1. Open your target repository root in a terminal.
2. Run the **Install Command (Copy/Paste)** block below (`REF="main"`, `PROVIDER="continue"`).
3. Confirm `.hawp/kit/` and `.continue/rules/hawp-*.md` exist.
4. Open Continue and verify HAWP rules appear in the rules toolbar.

Optional: `export HAWP_LOCAL_CORE="/absolute/path/to/human-ai-workflow-protocol/core"` for local testing.

## What Was Added

- `.hawp/kit/**` — agent-neutral HAWP kit (always installed).
- `.continue/rules/hawp-*.md` — Continue rules from `core/providers/.continue/rules/`.
- `.hawp/work/` scaffold — seeded once when missing.

## What Was NOT Changed

- Non-HAWP files under `.continue/rules/` (only `hawp-*.md` from the provider pack are copied).
- `.hawp/work/**` project records.

## Other guides

- Development branch: `distribution/generated/continue/install/development.md`
- GitHub/Copilot: `distribution/generated/github/install/main.md`
- Cursor: `distribution/generated/cursor/install/main.md`
