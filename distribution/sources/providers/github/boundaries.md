# GitHub / Copilot Provider Boundaries

This guide sets `PROVIDER=github`. Only the GitHub overlay is installed — not Cursor or Continue paths.

## Source pack

`core/providers/.github/`

## Install mapping

| Source | Installs to | Install | Update |
|--------|-------------|---------|--------|
| `instructions/*.instructions.md` | `.github/instructions/` | refresh | refresh |
| `prompts/*.prompt.md` | `.github/prompts/` | refresh | refresh |
| `copilot-instructions.md` | `.github/copilot-instructions.md` | seed if missing | refresh |

## Not touched by this guide

- `.cursor/**`, `AGENTS.md`, `.continue/**`
- Custom files under `.github/` outside the HAWP-managed paths above

## Boundary model

```
core/providers/.github/  →  .github/instructions/
                         →  .github/prompts/
                         →  .github/copilot-instructions.md
```
