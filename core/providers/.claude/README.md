# Claude Code Provider Pack

Source pack for the Claude Code provider overlay.

## What's here

| Path | Purpose |
|------|---------|
| `CLAUDE.md.seed` | Seeded as `CLAUDE.md` on first install (not overwritten on update) |
| `rules/hawp-*.md` | Generated from `core/providers/shared/behaviors/` — do not edit here |

## Edit rules

Edit `core/providers/shared/behaviors/<behavior>.md`, then run:

```bash
npm --prefix librarian run providers:sync
```

## Install mapping

| Source | Target | Install | Update |
|--------|--------|---------|--------|
| `rules/hawp-*.md` | `.claude/rules/` | refresh | refresh |
| `CLAUDE.md.seed` | `CLAUDE.md` (repo root) | seed if missing | skip |

## Claude Code rules format

Rules in `.claude/rules/` use YAML frontmatter:
- `paths:` — path glob array; rule loads only when matching files are opened
- No `paths:` field — rule loads at every session start (always-on)

See [Claude Code memory docs](https://code.claude.com/docs/en/memory) for full specification.
