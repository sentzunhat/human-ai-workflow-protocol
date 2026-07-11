# Continue Provider Boundaries

This guide sets `PROVIDER=continue`. Only the Continue overlay is installed — not GitHub or Cursor paths.

## Source pack

`core/providers/.continue/`

## Install mapping

| Source | Installs to | Install | Update |
|--------|-------------|---------|--------|
| `rules/hawp-*.md` | `.continue/rules/` | refresh | refresh |

## Not touched by this guide

- `.github/**`
- `.cursor/**`, `AGENTS.md`
- Non-HAWP rules in `.continue/rules/` (only `hawp-*.md` from the provider pack are refreshed)

## Boundary model

```
core/providers/.continue/  →  .continue/rules/hawp-*.md
```

Local Continue rules load automatically; no Hub `config.yaml` is required for this pack.
