# GitHub / Copilot Update Contract

## Work item goal

Refresh HAWP kit plus **GitHub Copilot overlays** from the selected branch. This refreshes `core/providers/.github/` into `.github/instructions/`, `.github/prompts/`, and `.github/copilot-instructions.md`.

## Agent execution

- Requires existing `.hawp/`.
- Report proof: `Source:`, `Provider: github`, `Source mode:`.
- File proof: `git status --short .hawp/kit .github/instructions .github/prompts .github/copilot-instructions.md`

## Provider-specific rules

- Refreshes `.github/copilot-instructions.md` on every update.
- Does not modify `.cursor/` or `AGENTS.md`.

## Auto-dispatch

Use the install-contract **Guide fetch (review-first)** block with `PROVIDER="github"`; it selects update when `.hawp/` exists.
