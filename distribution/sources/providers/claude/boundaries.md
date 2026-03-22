# Claude Code Provider Boundaries

This guide sets `PROVIDER=claude`. Only the Claude Code overlay is installed — not GitHub, Cursor, or Continue paths.

## Source pack

`core/providers/.claude/`

## Install mapping

| Source | Installs to | Install | Update |
|--------|-------------|---------|--------|
| `rules/hawp-*.md` | `.claude/rules/` | refresh | refresh |
| `CLAUDE.md.seed` | `CLAUDE.md` (repo root) | seed if missing | skip |

## Not touched by this guide

- `.github/**`
- `.cursor/**`
- `.continue/**`
- Non-HAWP files in `.claude/rules/` (only files copied from the provider pack are refreshed)
- `CLAUDE.md` on update (user-editable; never overwritten after first install)

## Boundary model

```
core/providers/.claude/  →  .claude/rules/hawp-*.md
                         →  CLAUDE.md  (seed only)
```
