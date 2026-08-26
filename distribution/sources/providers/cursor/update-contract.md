# Cursor Update Contract

## Work item goal

Refresh HAWP kit plus **Cursor overlays** from the selected branch. This refreshes `core/providers/.cursor/rules/` into `.cursor/rules/` and seeds `AGENTS.md` only when it is missing.

## Agent execution

- Requires existing `.hawp/`.
- Report proof: `Source:`, `Provider: cursor`, `Source mode:`.
- File proof: `git status --short .hawp/kit .cursor/rules AGENTS.md`

## Provider-specific rules

- Refreshes all provider-pack `.mdc` rules on every update.
- Seeds `AGENTS.md` only when it is missing; existing custom `AGENTS.md` content is preserved.
- If you want new HAWP instruction wording in an already-customized `AGENTS.md`, blend it manually instead of expecting update to overwrite the file.
- Does not modify `.github/`.

## Auto-dispatch

Use the install-contract **Guide fetch (review-first)** block with `PROVIDER="cursor"`; it selects update when `.hawp/` exists.
