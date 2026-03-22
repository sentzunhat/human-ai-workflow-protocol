# Claude Code Update Contract

## Work item goal

Refresh HAWP kit plus **Claude Code overlays** from the selected branch. This refreshes `core/providers/.claude/rules/hawp-*.md` into `.claude/rules/`. `CLAUDE.md` is never overwritten on update.

## Agent execution

- Requires existing `.hawp/`.
- Report proof: `Source:`, `Provider: claude`, `Source mode:`.
- File proof: `git status --short .hawp/kit .claude/rules`

## Provider-specific rules

- Refreshes all provider-pack `hawp-*.md` rules on every update.
- Never overwrites `CLAUDE.md` (user-editable).
- Does not modify `.github/`, `.cursor/`, or `.continue/`.

## Auto-dispatch

Use the install-contract **Guide fetch (review-first)** block with `PROVIDER="claude"`; it selects update when `.hawp/` exists.
