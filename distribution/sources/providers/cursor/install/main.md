# Install HAWP — Cursor Provider (Main Branch)

Stable install of HAWP kit plus Cursor overlays: `.cursor/rules/hawp-*.mdc` and root `AGENTS.md`.

**Source → target mapping:**

| `core/providers/.cursor/` | Your repo |
|---------------------------|-----------|
| `rules/*.mdc` | `.cursor/rules/` |
| `AGENTS.md.seed` | `AGENTS.md` |

## Prerequisites

- A repository where you use Cursor Agent.
- `curl` and `tar`.

## Installation Steps

1. Open your target repository root in a terminal.
2. Run the **Install Command (Copy/Paste)** block below (`REF="main"`, `PROVIDER="cursor"`).
3. Confirm `.hawp/kit/`, `.cursor/rules/hawp-*.mdc`, and `AGENTS.md` exist.
4. Open Cursor Agent and verify HAWP guidance loads (status report / backlog prompts).

Optional: `export HAWP_LOCAL_CORE="/absolute/path/to/human-ai-workflow-protocol/core"` for local testing.

## What Was Added

- `.hawp/kit/**` — agent-neutral HAWP kit (always installed).
- `.cursor/rules/hawp-*.mdc` — Cursor rules from `core/providers/.cursor/rules/`.
- `AGENTS.md` — from `core/providers/.cursor/AGENTS.md.seed`.
- `.hawp/work/` scaffold — seeded once when missing.

If your repo already has `AGENTS.md`, HAWP preserves it. Manually blend in any
HAWP guidance you want from the provider seed instead of overwriting your file.

## What Was NOT Changed

- Non-HAWP files under `.cursor/rules/` (only `hawp-*.mdc` from the provider pack are copied).
- `.hawp/work/**` project records.

## Other guides

- Dev branch: `distribution/generated/cursor/install/dev.md`
- GitHub/Copilot: `distribution/generated/github/install/main.md`
- Continue: `distribution/generated/continue/install/main.md`
