# Codex Provider Boundaries

This guide sets `PROVIDER=codex`. Only the Codex overlay is installed — not GitHub, Cursor, Continue, or Claude paths.

## Source pack

`core/providers/.codex/`

## Install mapping

| Source | Installs to | Install | Update |
|--------|-------------|---------|--------|
| `AGENTS.md.seed` | `AGENTS.md` (repo root) | seed if missing | seed if missing |

## Not touched by this guide

- `.github/**`
- `.cursor/**`
- `.continue/**`
- `.claude/**`
- Runtime CLI participant adapters

## Boundary model

```text
core/providers/.codex/  ->  AGENTS.md
```
