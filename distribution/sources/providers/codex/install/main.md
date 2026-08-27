# Install HAWP — Codex Provider (Main Branch)

Stable install of HAWP kit plus Codex `AGENTS.md` instructions.

**Source -> target mapping:**

| `core/providers/.codex/` | Your repo |
|--------------------------|-----------|
| `AGENTS.md.seed` | `AGENTS.md` |

## Prerequisites

- A repository where you use Codex.
- `curl` and `tar`.

## Installation Steps

1. Open your target repository root in a terminal.
2. Run the **Install Command (Copy/Paste)** block below (`REF="main"`, `PROVIDER="codex"`).
3. Confirm `.hawp/kit/` and `AGENTS.md` exist.

Optional: `export HAWP_LOCAL_CORE="/absolute/path/to/human-ai-workflow-protocol/core"` for local testing.

## What Was Added

- `.hawp/kit/**` — agent-neutral HAWP kit (always installed).
- `AGENTS.md` — Codex repo-local instructions, seeded only when missing.
- `.hawp/work/` scaffold — seeded once when missing.

If your repo already has `AGENTS.md`, HAWP preserves it. Manually blend in any
HAWP guidance you want from the provider seed instead of overwriting your file.

## What Was NOT Changed

- `.github/**`
- `.cursor/**`
- `.continue/**`
- `.claude/**`
- Runtime CLI participant adapters.
- `.hawp/work/**` project records.

## Other guides

- Dev branch: `distribution/generated/codex/install/dev.md`
- GitHub/Copilot: `distribution/generated/github/install/main.md`
- Cursor: `distribution/generated/cursor/install/main.md`
- Continue: `distribution/generated/continue/install/main.md`
