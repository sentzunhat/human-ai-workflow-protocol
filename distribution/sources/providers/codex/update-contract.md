# Codex Update Contract

## Work item goal

Refresh HAWP kit plus the **Codex overlay** from the selected branch. This seeds `core/providers/.codex/AGENTS.md.seed` into `AGENTS.md` only when the target repo does not already have one.

## Agent execution

- Requires existing `.hawp/`.
- Report proof: `Source:`, `Provider: codex`, `Source mode:`.
- File proof: `git status --short .hawp/kit AGENTS.md`

## Provider-specific rules

- Seeds `AGENTS.md` only when it is missing; existing custom `AGENTS.md` content is preserved.
- If you want new HAWP instruction wording in an already-customized `AGENTS.md`, blend it manually instead of expecting update to overwrite the file.
- Does not modify `.github/`, `.cursor/`, `.continue/`, or `.claude/`.
- Does not create runtime CLI participant adapters.

## Auto-dispatch

Use the install-contract **Guide fetch (review-first)** block with `PROVIDER="codex"`; it selects update when `.hawp/` exists.
