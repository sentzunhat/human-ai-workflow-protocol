# Update HAWP — Cursor Provider (Main Branch)

Refresh HAWP kit and Cursor provider overlays from `main`.

**Refreshed paths:**

| Source | Target |
|--------|--------|
| `core/providers/.cursor/rules/*.mdc` | `.cursor/rules/` |
| `core/providers/.cursor/AGENTS.md.seed` | `AGENTS.md` |

## Prerequisites

- HAWP already installed (`.hawp/` present).

## Steps

1. Repository root in terminal.
2. Run **Update Command (Copy/Paste)** (`REF="main"`, `PROVIDER="cursor"`).
3. Verify `.hawp/work/` intact and Cursor rules updated.

## Preserved

- `.hawp/work/**` — never overwritten.

## Other guides

- Dev update: `distribution/generated/cursor/update/dev.md`
