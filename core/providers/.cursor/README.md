# HAWP Cursor Provider

Install pack for Cursor Agent. Source: `core/providers/.cursor/` → target repo `.cursor/rules/` and `AGENTS.md`.

| Source | Installs to |
|--------|-------------|
| `rules/*.mdc` | `.cursor/rules/` (all files refreshed on install/update) |
| `AGENTS.md.seed` | `AGENTS.md` at repo root (seeded on install if missing; refreshed on update) |

Run install/update from `distribution/generated/cursor/install|update/{main|dev}.md`.

Docs: [Cursor Rules](https://cursor.com/docs/context/rules), [AGENTS.md](https://cursor.com/docs/context/rules#agentsmd). See [doc-verification.md](../doc-verification.md).
