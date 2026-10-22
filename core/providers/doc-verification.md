# Provider documentation alignment

Last verified: 2026-10-16 against the vendor docs linked below. This records
documented client features only; it is not a HAWP client smoke result.

## HAWP MCP Capability Guidance

HAWP's source server exposes the established tools and now also implements
`resources/list`, `resources/read`, `prompts/list`, and `prompts/get` over local
stdio. Client UX remains connection- and version-specific. The group names
Discover, Shape, Create, and Verify are HAWP documentation categories, not MCP
protocol namespaces.

| Client surface | Vendor documentation observed | HAWP guidance |
| --- | --- | --- |
| Claude Code | Documents MCP resources through `@` references and server prompts as slash commands. | Resources/prompts can be tried after connecting; no HAWP client smoke has been run. |
| GitHub Copilot in VS Code | Documents MCP prompts as slash commands and MCP resources through Add Context. | Resources/prompts can be tried in this VS Code surface; do not generalize to Copilot cloud agent or code review. No HAWP client smoke has been run. |
| Cursor | MCP documentation lists prompts as supported; it does not establish a resource reference workflow. | Prompt exposure is documented generically; resource UX remains unconfirmed. No HAWP client smoke has been run. |
| Continue | Documents MCP server configuration/context-provider support; the inspected docs do not specify HAWP resource/prompt discovery UX. | Use tools and the manual guide unless the client explicitly lists a resource or prompt. |
| Codex local clients | Documents local MCP server support and server instructions; the inspected supported-feature list does not establish MCP resources/prompts. | Use tools and the manual guide unless the client explicitly lists a resource or prompt. |

HAWP tool grouping and the fallback workflow are maintained in
`.hawp/kit/usage/mcp/README.md`. Avoid claiming client support from the server's
protocol handlers or an SDK feature list alone.

## GitHub / Copilot

| Artifact | Our path | Vendor reference |
|----------|----------|------------------|
| Repository instructions | `.github/copilot-instructions.md` | [Add repository instructions](https://docs.github.com/en/copilot/how-tos/copilot-on-github/customize-copilot/add-custom-instructions/add-repository-instructions) |
| Path-specific instructions | `.github/instructions/*.instructions.md` with `applyTo` | [IDE repository instructions](https://docs.github.com/en/copilot/how-tos/configure-custom-instructions-in-your-ide/add-repository-instructions-in-your-ide) |
| Prompt files | `.github/prompts/*.prompt.md` with `name` / `description` | [Prompt files tutorial](https://docs.github.com/en/copilot/tutorials/customization-library/prompt-files/your-first-prompt-file) |
| VS Code cross-check | same paths | [VS Code custom instructions](https://code.visualstudio.com/docs/copilot/customization/custom-instructions) |

**HAWP behavior:** seed `copilot-instructions.md` on install if missing; refresh on update. Instructions use YAML `applyTo` globs.

Vendor MCP references: [Copilot Chat MCP servers](https://docs.github.com/en/copilot/how-tos/copilot-in-your-ide/customize-copilot/extend-copilot-with-tools-and-context/extend-copilot-chat-with-mcp), [Claude Code MCP](https://code.claude.com/docs/en/mcp), [Cursor MCP](https://docs.cursor.com/context/model-context-protocol), [Continue MCP configuration](https://docs.continue.dev/reference), and [Codex MCP](https://developers.openai.com/codex/mcp).

## Cursor

| Artifact | Our path | Vendor reference |
|----------|----------|------------------|
| Project rules | `.cursor/rules/*.mdc` | [Cursor Rules](https://cursor.com/docs/context/rules) |
| Agent instructions | `AGENTS.md` (repo root) | [Cursor AGENTS.md](https://cursor.com/docs/context/rules#agentsmd) |

**HAWP behavior:** refresh all provider-pack `.mdc` rules on install/update; seed `AGENTS.md` only when missing on both install and update. Plain `.md` in `.cursor/rules/` is ignored by Cursor — use `.mdc` or `AGENTS.md`.

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
