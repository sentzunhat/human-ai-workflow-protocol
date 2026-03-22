# HAWP Agent Providers

Install-pack sources under `core/providers/`. Each provider guide copies **kit + that provider's overlay only**.

## Layout

```text
providers/
├── manifest.yaml       # source → target routing reference
├── README.md
├── shared/             # canonical behaviors → materialized into provider packs
│   └── behaviors/
├── .claude/            → installs to .claude/rules/ + CLAUDE.md
├── .codex/             → installs to AGENTS.md
├── .github/            → installs to .github/
├── .cursor/            → installs to .cursor/ + AGENTS.md
└── .continue/          → installs to .continue/rules/
```

## Install mapping

| Provider | Source pack | Downstream targets |
|----------|-------------|-------------------|
| Claude Code | `providers/.claude/` | `.claude/rules/hawp-*.md`, `CLAUDE.md` |
| Codex | `providers/.codex/` | `AGENTS.md` |
| GitHub | `providers/.github/` | `.github/instructions/`, `.github/prompts/`, `.github/copilot-instructions.md` |
| Cursor | `providers/.cursor/` | `.cursor/rules/*.mdc`, `AGENTS.md` |
| Continue | `providers/.continue/` | `.continue/rules/hawp-*.md` |

## Distribution guides

```text
distribution/generated/<provider>/install|update/{main|dev}.md
```

See [distribution/generated/README.md](../../distribution/generated/README.md) for links.

Documentation alignment: [doc-verification.md](doc-verification.md).

Shared behaviors: [shared/README.md](shared/README.md) — run `npm --prefix librarian run providers:sync` after editing.

## Adding A Provider

1. Create `core/providers/.<provider>/` with the provider source pack.
2. Add install targets to `core/providers/manifest.yaml`.
3. Add provider-facing fragments under `distribution/sources/providers/<provider>/`.
4. Register the provider in `ACTIVE_PROVIDERS` in `librarian/scripts/librarian/distribution/shared/composition.ts`.
5. Add generated shared-behavior targets in `librarian/scripts/librarian/providers/materialize/composition.ts` when the provider uses materialized rule files.
6. Update `distribution/sources/providers/README.md` and `distribution/generated/README.md`.
7. Run `npm --prefix librarian run distribution:sync`.

Keep runtime adapters out of provider packs. Provider guides install the HAWP kit plus provider instructions only.
