# Cursor Provider Boundaries

This guide sets `PROVIDER=cursor`. Only the Cursor overlay is installed — not GitHub or Continue paths.

## Source pack

`core/providers/.cursor/`

## Install mapping

| Source | Installs to | Install | Update |
|--------|-------------|---------|--------|
| `rules/*.mdc` | `.cursor/rules/` | refresh | refresh |
| `AGENTS.md.seed` | `AGENTS.md` (repo root) | seed if missing | refresh |

## Not touched by this guide

- `.github/**`
- `.continue/**`
- Non-HAWP rules in `.cursor/rules/` (only files copied from the provider pack are refreshed)

## Boundary model

```
core/providers/.cursor/  →  .cursor/rules/hawp-*.mdc
                       →  AGENTS.md
```
