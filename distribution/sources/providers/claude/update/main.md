# Update HAWP — Claude Code Provider (Main Branch)

Refresh HAWP kit and Claude Code provider rules from `main`. `CLAUDE.md` is preserved.

**Refreshed paths:**

| Source | Target |
|--------|--------|
| `core/providers/.claude/rules/hawp-*.md` | `.claude/rules/` |

## Prerequisites

- HAWP already installed (`.hawp/` present).

## Steps

1. Repository root in terminal.
2. Run **Update Command (Copy/Paste)** (`REF="main"`, `PROVIDER="claude"`).
3. Verify `.hawp/work/` intact and Claude Code rules updated.

## Preserved

- `.hawp/work/**` — never overwritten.
- `CLAUDE.md` — user-editable; never overwritten on update.

## Other guides

- Dev update: `distribution/generated/claude/update/dev.md`
