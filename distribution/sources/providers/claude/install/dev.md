# Install HAWP — Claude Code Provider (Dev Branch)

Install HAWP kit plus Claude Code overlays from the `dev` branch.

Same source → target mapping as main: `core/providers/.claude/` → `.claude/rules/` + `CLAUDE.md` (seed).

## When to Use

- Testing unreleased HAWP Claude Code provider changes.
- Contributing to HAWP development.

## Steps

1. Repository root in terminal.
2. Run install command block (`REF="dev"`, `PROVIDER="claude"`).
3. Verify `.hawp/kit/`, `.claude/rules/`, and `CLAUDE.md`.

## Reverting to Main

Use `distribution/generated/claude/install/main.md`.
