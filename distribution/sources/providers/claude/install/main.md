# Install HAWP — Claude Code Provider (Main Branch)

Stable install of HAWP kit plus Claude Code overlays: `.claude/rules/hawp-*.md` and root `CLAUDE.md`.

**Source → target mapping:**

| `core/providers/.claude/` | Your repo |
|---------------------------|-----------|
| `rules/hawp-*.md` | `.claude/rules/` |
| `CLAUDE.md.seed` | `CLAUDE.md` (seeded once) |

## Prerequisites

- A repository where you use Claude Code.
- `curl` and `tar`.

## Installation Steps

1. Open your target repository root in a terminal.
2. Run the **Install Command (Copy/Paste)** block below (`REF="main"`, `PROVIDER="claude"`).
3. Confirm `.hawp/kit/`, `.claude/rules/hawp-*.md`, and `CLAUDE.md` exist.
4. Open Claude Code and verify HAWP guidance loads (status report / backlog prompts).

Optional: `export HAWP_LOCAL_CORE="/absolute/path/to/human-ai-workflow-protocol/core"` for local testing.

For a temporary slash-named provider branch, use the visible command block below and set `REF` to that branch name after review.

## What Was Added

- `.hawp/kit/**` — agent-neutral HAWP kit (always installed).
- `.claude/rules/hawp-*.md` — Claude Code rules from `core/providers/.claude/rules/`.
- `CLAUDE.md` — from `core/providers/.claude/CLAUDE.md.seed` (seeded only when missing).
- `.hawp/work/` scaffold — seeded once when missing.

## What Was NOT Changed

- Non-HAWP files under `.claude/rules/` (only `hawp-*.md` from the provider pack are copied).
- `.hawp/work/**` project records.
- Existing `CLAUDE.md` if one already exists.

## Customizing CLAUDE.md

After install, add project-specific content to `CLAUDE.md`: build commands, test commands, coding standards, and architecture notes. The seeded file includes a placeholder section.

## Other guides

- Development branch: `distribution/generated/claude/install/development.md`
- GitHub/Copilot: `distribution/generated/github/install/main.md`
- Cursor: `distribution/generated/cursor/install/main.md`
- Continue: `distribution/generated/continue/install/main.md`
