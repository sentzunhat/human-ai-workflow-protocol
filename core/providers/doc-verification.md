# Provider documentation alignment

Last verified: 2026-06-04 against vendor docs linked below.

## GitHub / Copilot

| Artifact | Our path | Vendor reference |
|----------|----------|------------------|
| Repository instructions | `.github/copilot-instructions.md` | [Add repository instructions](https://docs.github.com/en/copilot/how-tos/copilot-on-github/customize-copilot/add-custom-instructions/add-repository-instructions) |
| Path-specific instructions | `.github/instructions/*.instructions.md` with `applyTo` | [IDE repository instructions](https://docs.github.com/en/copilot/how-tos/configure-custom-instructions-in-your-ide/add-repository-instructions-in-your-ide) |
| Prompt files | `.github/prompts/*.prompt.md` with `name` / `description` | [Prompt files tutorial](https://docs.github.com/en/copilot/tutorials/customization-library/prompt-files/your-first-prompt-file) |
| VS Code cross-check | same paths | [VS Code custom instructions](https://code.visualstudio.com/docs/copilot/customization/custom-instructions) |

**HAWP behavior:** seed `copilot-instructions.md` on install if missing; refresh on update. Instructions use YAML `applyTo` globs.

## Cursor

| Artifact | Our path | Vendor reference |
|----------|----------|------------------|
| Project rules | `.cursor/rules/*.mdc` | [Cursor Rules](https://cursor.com/docs/context/rules) |
| Agent instructions | `AGENTS.md` (repo root) | [Cursor AGENTS.md](https://cursor.com/docs/context/rules#agentsmd) |

**HAWP behavior:** refresh all provider-pack `.mdc` rules on install/update; seed `AGENTS.md` on install only when missing and refresh it on update. Plain `.md` in `.cursor/rules/` is ignored by Cursor — use `.mdc` or `AGENTS.md`.

Frontmatter: `description`, `globs`, `alwaysApply` per rule type table in Cursor docs.

## Continue

| Artifact | Our path | Vendor reference |
|----------|----------|------------------|
| Local workspace rules | `.continue/rules/*.md` | [Continue Rules](https://docs.continue.dev/customize/deep-dives/rules) |

**HAWP behavior:** refresh provider-pack rules only (`hawp-*.md`). Local rules load automatically; Hub `config.yaml` is not required for this pack.

Frontmatter: `name`, `description`, `globs` (string or array), `alwaysApply`. Rules apply in lexicographic order — use `hawp-01-*` … `hawp-04-*` prefixes.

## Shared behaviors (source of truth)

Four HAWP integration behaviors are authored once under `core/providers/shared/behaviors/` and materialized into GitHub, Cursor, and Continue packs. See [shared/README.md](shared/README.md).

Hand-maintained per provider: GitHub `copilot-instructions.md` and prompts; Cursor `AGENTS.md.seed`.

## Install guides

Per-provider bash blocks install kit + one overlay only. See `distribution/generated/<provider>/`.
