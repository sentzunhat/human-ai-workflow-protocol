# Codex Update Contract

## Work item goal

Refresh HAWP kit plus the **Codex overlay** from the selected branch. This refreshes `core/providers/.codex/AGENTS.md.seed` into `AGENTS.md`.

## Agent execution

- Requires existing `.hawp/`.
- Report proof: `Source:`, `Provider: codex`, `Source mode:`.
- File proof: `git status --short .hawp/kit AGENTS.md`

## Provider-specific rules

- Refreshes `AGENTS.md` on every update.
- Does not modify `.github/`, `.cursor/`, `.continue/`, or `.claude/`.
- Does not create runtime CLI participant adapters.

## Auto-dispatch

Use the install-contract **Guide fetch (review-first)** block with `PROVIDER="codex"`; it selects update when `.hawp/` exists.

