# HAWP Codex Provider

Install pack for Codex. Source: `core/providers/.codex/` -> target repo `AGENTS.md`.

| Source | Installs to |
|--------|-------------|
| `AGENTS.md.seed` | `AGENTS.md` at repo root (seeded on install if missing; refreshed on update) |

Run install/update from `distribution/generated/codex/install|update/{main|dev}.md`.

Codex reads repo-local instructions from `AGENTS.md`. This provider does not create a `.codex/` folder or runtime CLI adapter.

