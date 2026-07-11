# Cursor Update Contract

## Work item goal

Refresh HAWP kit plus **Cursor overlays** from the selected branch. This refreshes `core/providers/.cursor/` into `.cursor/rules/` and `AGENTS.md`.

## Agent execution

- Requires existing `.hawp/`.
- Report proof: `Source:`, `Provider: cursor`, `Source mode:`.
- File proof: `git status --short .hawp/kit .cursor/rules AGENTS.md`

## Provider-specific rules

- Refreshes all provider-pack `.mdc` rules and `AGENTS.md` on every update.
- Does not modify `.github/`.

## Auto-dispatch

Use the install-contract **Guide fetch (review-first)** block with `PROVIDER="cursor"`; it selects update when `.hawp/` exists.
