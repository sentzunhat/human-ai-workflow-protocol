# Generated distribution guides

Do not edit directly. Regenerate:

```bash
cd librarian/src && go run ./cmd/hawp distribution sync
```

Edits to provider overlay rules: change `core/providers/shared/behaviors/` (materialized into provider packs automatically).

## Layout

```text
generated/
└── <provider>/
    ├── install/main.md | install/development.md
    └── update/main.md  | update/development.md
```

Each guide embeds `PROVIDER=<name>` and `REF=main|development` in the bash block. Scripts install `.hawp/kit/**` plus **only that provider's overlay** from `core/providers/.<provider>/`.

## Claude Code

| | Stable (`main`) | Development (`development`) |
|--|-----------------|-------------|
| **Install** | [install/main.md](claude/install/main.md) | [install/development.md](claude/install/development.md) |
| **Update** | [update/main.md](claude/update/main.md) | [update/development.md](claude/update/development.md) |

Installs: `core/providers/.claude/` → `.claude/rules/hawp-*.md`, `CLAUDE.md`

## Codex

| | Stable (`main`) | Development (`development`) |
|--|-----------------|-------------|
| **Install** | [install/main.md](codex/install/main.md) | [install/development.md](codex/install/development.md) |
| **Update** | [update/main.md](codex/update/main.md) | [update/development.md](codex/update/development.md) |

Installs: `core/providers/.codex/` → `AGENTS.md`

## GitHub / Copilot

| | Stable (`main`) | Development (`development`) |
|--|-----------------|-------------|
| **Install** | [install/main.md](github/install/main.md) | [install/development.md](github/install/development.md) |
| **Update** | [update/main.md](github/update/main.md) | [update/development.md](github/update/development.md) |

Installs: `core/providers/.github/` → `.github/`

## Cursor

| | Stable (`main`) | Development (`development`) |
|--|-----------------|-------------|
| **Install** | [install/main.md](cursor/install/main.md) | [install/development.md](cursor/install/development.md) |
| **Update** | [update/main.md](cursor/update/main.md) | [update/development.md](cursor/update/development.md) |

Installs: `core/providers/.cursor/` → `.cursor/rules/`, `AGENTS.md`

## Continue

| | Stable (`main`) | Development (`development`) |
|--|-----------------|-------------|
| **Install** | [install/main.md](continue/install/main.md) | [install/development.md](continue/install/development.md) |
| **Update** | [update/main.md](continue/update/main.md) | [update/development.md](continue/update/development.md) |

Installs: `core/providers/.continue/` → `.continue/rules/hawp-*.md`

## Legacy

Root-level `install-main.md` / `update-development.md` were removed. Use provider folders above.
