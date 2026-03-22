# Generated distribution guides

Do not edit directly. Regenerate:

```bash
npm --prefix librarian run distribution:sync
```

Edits to provider overlay rules: change `core/providers/shared/behaviors/` (materialized into provider packs automatically).

## Layout

```text
generated/
└── <provider>/
    ├── install/main.md | install/dev.md
    └── update/main.md  | update/dev.md
```

Each guide embeds `PROVIDER=<name>` and `REF=main|dev` in the bash block. Scripts install `.hawp/kit/**` plus **only that provider's overlay** from `core/providers/.<provider>/`.

## Claude Code

| | Stable (`main`) | Dev (`dev`) |
|--|-----------------|-------------|
| **Install** | [install/main.md](claude/install/main.md) | [install/dev.md](claude/install/dev.md) |
| **Update** | [update/main.md](claude/update/main.md) | [update/dev.md](claude/update/dev.md) |

Installs: `core/providers/.claude/` → `.claude/rules/hawp-*.md`, `CLAUDE.md`

## Codex

| | Stable (`main`) | Dev (`dev`) |
|--|-----------------|-------------|
| **Install** | [install/main.md](codex/install/main.md) | [install/dev.md](codex/install/dev.md) |
| **Update** | [update/main.md](codex/update/main.md) | [update/dev.md](codex/update/dev.md) |

Installs: `core/providers/.codex/` → `AGENTS.md`

## GitHub / Copilot

| | Stable (`main`) | Dev (`dev`) |
|--|-----------------|-------------|
| **Install** | [install/main.md](github/install/main.md) | [install/dev.md](github/install/dev.md) |
| **Update** | [update/main.md](github/update/main.md) | [update/dev.md](github/update/dev.md) |

Installs: `core/providers/.github/` → `.github/`

## Cursor

| | Stable (`main`) | Dev (`dev`) |
|--|-----------------|-------------|
| **Install** | [install/main.md](cursor/install/main.md) | [install/dev.md](cursor/install/dev.md) |
| **Update** | [update/main.md](cursor/update/main.md) | [update/dev.md](cursor/update/dev.md) |

Installs: `core/providers/.cursor/` → `.cursor/rules/`, `AGENTS.md`

## Continue

| | Stable (`main`) | Dev (`dev`) |
|--|-----------------|-------------|
| **Install** | [install/main.md](continue/install/main.md) | [install/dev.md](continue/install/dev.md) |
| **Update** | [update/main.md](continue/update/main.md) | [update/dev.md](continue/update/dev.md) |

Installs: `core/providers/.continue/` → `.continue/rules/hawp-*.md`

## Legacy

Root-level `install-main.md` / `update-dev.md` were removed. Use provider folders above.
