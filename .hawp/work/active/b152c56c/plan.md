# Provider MCP setup and tool availability allowlists

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `b152c56c-253e-4a48-a970-231b1764fa84`
**Type:** improvement
**Reported:** 2026-10-16

---

## Input (verbatim)

> Update HAWP MCP provider configuration so HAWP can set up every supported provider and configure per-tool availability allowlists where the client supports them. Preserve manual approval for write-capable tools. Verify Codex and VS Code GitHub Copilot only; document unsupported or client-managed allowlists for other providers.

## Intake Summary

The provider configuration command already creates repository-scoped MCP server
entries for Claude Code, Cursor, Codex, and VS Code GitHub Copilot; Continue
prints a snippet for the user's profile config. Codex's `enabled_tools` field
sets availability. VS Code supports a custom agent with an explicit `tools`
list, which can scope HAWP tools without changing approval behavior. Similar
allowlist terminology in Claude and Cursor refers to auto-approval, so HAWP must
not write it as a substitute for tool visibility.

## Current Context

- `librarian/src/internal/platform/mcp/configure/config.go` dispatches the
  supported provider configurations. `--provider all` expands to Claude,
  Cursor, Continue, Codex, and GitHub Copilot.
- Codex config already enables all seven HAWP tools through `enabled_tools`.
- VS Code's provider config is `.vscode/mcp.json`. The current VS Code guide
  documents custom-agent frontmatter `tools` entries and the `<server>/*`
  notation; tool-call approval is a separate client decision.
- Claude and Cursor have approval rules with allowlist names. Automatically
  populating them could preapprove write-capable HAWP tools such as
  `hawp_work_new` or `hawp_work_doc`.
- Continue config is user-scoped (`~/.continue/config.yaml`) and is not
  modified by the repository command.
- The user limited live client verification to VS Code GitHub Copilot and
  Codex.

## Initial Analysis

**Directly verified:**

- `./.hawp/bin/hawp mcp configure --provider all` ran successfully and emitted
  setup for all five named providers; Continue remains a printed user-config
  block.
- `codex mcp get hawp` reports the repository-scoped HAWP server enabled and
  lists all seven HAWP tools in `enabled_tools`.
- A Codex-hosted `hawp_work_validate` MCP call passed for this repository:
  kit/work checks passed and local links were valid across 149 Markdown files.
- The HAWP VS Code workspace has `.vscode/mcp.json`; after the user approved
  the per-session confirmation, Copilot called `hawp_work_validate` through
  MCP and validation passed. The generated custom-agent tool selection was
  not independently exercised in a HAWP agent session.
- Official VS Code documentation supports a custom agent's `tools` field for
  MCP tool selection and separately describes server enablement and tool-call
  trust/approval.

**Inferred (not yet proven):**

- A generated, repo-local VS Code custom agent can make HAWP's tool visibility
  repeatable without changing workspace-wide approval settings.
- Exact HAWP tool identifiers in the custom-agent allowlist still require
  client verification after generating the file.

**Likely scope:**

- Generate a Copilot custom agent in `.github/agents/` with only HAWP's seven
  MCP tools listed and instructions that preserve call approvals.
- Preflight the target path and refuse to overwrite an unrelated existing
  custom agent file.
- Keep Codex's existing tool allowlist intact; clarify other providers'
  client-managed behavior in MCP setup documentation.
- Do not write auto-approval rules or mutate user-level Continue config.
- Verify Codex and Copilot only, with read-only HAWP validation calls.

## Risk + Review Gate

**Risk:** medium
**Gate:** User requested provider setup and verification. Keep changes
repo-scoped, preserve unrelated custom agent files, and do not auto-approve
write-capable tools.

## Backlog + Plan Link

**Status now:** in-progress
**Plan file:** work/active/b152c56c/plan.md

## Next Step

- [x] Investigation recorded above (required before planning)
- [x] Write or update the plan file
- [x] Move backlog status accordingly

## Plan

1. [x] Add a generated HAWP-only VS Code agent allowlist with safe path checks,
   preserving manual tool approval.
2. [x] Keep the existing Codex `enabled_tools` allowlist and update provider
   setup guidance with the difference between availability and auto-approval.
3. [x] Verify Codex with its repo config and a read-only HAWP MCP validation.
4. [x] Run focused configuration tests and HAWP consistency checks.
5. [x] Verify Copilot's HAWP tool call after the user approves the per-session
   permission prompt.
6. [ ] In the already-selected fresh HAWP chat, inspect the live Tools list
   and approval behavior after refreshing the generated agent. Do not infer
   them from the file or ask the user to select the agent again.

## Verification

- `go test ./internal/platform/mcp/configure` passed.
- `go run ./cmd/hawp mcp configure --provider all` completed. It wrote Claude,
  Cursor, Codex, and VS Code config, generated the VS Code HAWP agent, and
  printed Continue's user-scoped config block.
- `codex mcp get hawp` confirms the project server is enabled and its seven
  tools are allowlisted. A read-only HAWP `hawp_work_validate` MCP call passed:
  kit/work checks passed and links were valid across 149 Markdown files.
- `.hawp/bin/hawp work validate` passed with zero issues and one existing
  verification-clarity warning; `links check` passed for 150 Markdown files;
  `git diff --check` passed.
- VS Code Copilot called `hawp-hawp_work_validate` through MCP after the user
  approved the session prompt. Kit validation passed (3 checks, 0 issues); work
  validation passed (0 issues, 1 warning); link validation passed (150
  Markdown files, all local links valid). The generated HAWP custom-agent scope
  still needs a separate session selected explicitly with that agent.

## Scope extension: VS Code workspace tools and HAWP-aware docs (2026-10-22)

### Directly verified

- `librarian/src/internal/platform/mcp/configure/config_agent.go` and the
  generated `.github/agents/hawp.agent.md` enumerate the seven HAWP MCP tools
  but no VS Code workspace tools. `configure_test.go` asserts that allowlist
  and preserves the client confirmation wording.
- `.github/agents/master-builder.agent.md` already has broad VS Code tool
  aliases and `hawp/*`; its HAWP workflow section does not explain when to use
  the MCP tools or how to interpret empty search results.
- `.github/agents/docs-alignment.agent.md` has VS Code workspace tools but no
  HAWP MCP tools.
- `.github/skills/hawp-change-review-and-reference-sync/SKILL.md` and
  `core/providers/.github/prompts/hawp-change-review-and-reference-sync.prompt.md`
  do not specify HAWP MCP usage. The provider prompt can carry a tool allowlist.
- `.github/copilot-instructions.md` points to a prompt path that is absent in
  this workspace; the matching local docs-sync skill exists.
- The HAWP Copilot usage guides describe the generated agent as exposing only
  the seven HAWP MCP tools.

### Decision

- Give the generated HAWP agent the native VS Code `read`, `search`, `edit`,
  and `execute` tools plus the seven explicit HAWP MCP tools. Keep the HAWP MCP
  allowlist explicit and preserve normal client confirmation prompts.
- Keep Master Builder's existing tool scope and add HAWP-first usage guidance.
- Give Docs Alignment and docs sync only the read-oriented HAWP
  `hawp_search` and `hawp_work_validate` tools; their existing workspace tools
  remain available. Do not grant HAWP work-creation tools to a docs audit.
- Update the local skill, provider prompt source, Copilot guidance, and
  affected HAWP usage guides. Do not alter ignored MCP configuration files.

### Plan

- [x] Update the generated-agent source, output message, and tests for the
  seven HAWP tools plus four VS Code workspace tools.
- [x] Add read-oriented HAWP tool guidance to Master Builder, Docs Alignment,
  and docs sync; retain client-managed approval behavior.
- [x] Align the Copilot reference and HAWP setup guides with the generated
  tool allowlist.
- [x] Run the configure tests and canonical HAWP checks; record known
  unrelated validation failures without expanding scope.
- [ ] Verify the actual Tools list and approval behavior in the selected
  Copilot session after the agent definition refresh.
- [ ] Keep the same-name MCP migration preview unresolved until its exact
  effect is visible; make no MCP config changes meanwhile.

### Verification results (2026-10-22)

- `go test ./internal/platform/mcp/configure` and `go test ./...` passed.
- `go run ./cmd/hawp distribution sync` updated 0 of 40 generated files and
  reported generated outputs current.
- `hawp check` did not pass: kit validation passed (3 checks, 0 issues), work
  validation reported 2 issues and 1 warning, and link checking found two
  pre-existing missing backlog targets: `active/2eea565c/plan.md` and
  `active/429e075e/plan.md`.
- Detailed `go run ./cmd/hawp work validate` confirmed those two missing-plan
  issues. Verification clarity reported 309/392 evidence-backed, 73 explicitly
  unproven, and 10 ambiguous. These unrelated backlog records were not changed.
- YAML frontmatter parsing passed for the updated agent, skill, and prompt
  files; `git diff --check` passed.

### Deferred live checks

- The VS Code session's actual post-change Tools list and client approval
  behavior remain unverified; generated frontmatter is not live-session proof.
- The migration preview for the existing `hawp` entry remains unverified.
- No ignored MCP configuration file or downstream consumer was changed.

### Repository proof before edits

```text
pwd: <repo-root-abs>/human-ai-workflow-protocol
git rev-parse --show-toplevel: <repo-root-abs>/human-ai-workflow-protocol
git rev-parse --show-prefix: (empty)
git status --short:
 M .hawp/work/BACKLOG.md
?? .hawp/work/active/09ee9986/
```
