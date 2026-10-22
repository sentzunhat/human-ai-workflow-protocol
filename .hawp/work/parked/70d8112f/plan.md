# Audit GitHub Copilot provider install against current VS Code MCP changes

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `70d8112f-5b33-47e9-a312-3116c9578a4d`
**Type:** improvement
**Reported:** 2026-10-21

---

## Input (verbatim)

> Review the code logic is doing the latest VS Code Copilot changes properly for installing that provider.

## Intake Summary

Audit HAWP's GitHub Copilot provider install and update path against current
VS Code MCP configuration, migration, and custom-agent behavior. Record whether
the implementation is current and make only evidence-backed changes if the
review gate is approved.

## Current Context

The existing provider implementation writes `.vscode/mcp.json` for GitHub
Copilot and writes `.mcp.json` for Claude Code. VS Code now offers migration
from the former to the portable Copilot `.mcp.json` format. The local HAWP
configuration has both files, so the config reconciliation item must complete
first.

## Initial Analysis

**Directly verified:**

- `librarian/src/internal/platform/mcp/configure/config.go` dispatches GitHub
  setup to `.vscode/mcp.json` with the VS Code `servers` schema and writes
  Claude's `mcpServers` config to root `.mcp.json`.
- `librarian/src/internal/platform/mcp/configure/configure.go` preflights the
  GitHub config and generated HAWP agent before writing.
- `librarian/src/internal/platform/mcp/configure/configure_test.go` covers
  repeated provider setup, preservation of existing settings, and refusal to
  overwrite an unmanaged HAWP agent.
- `.hawp/kit/usage/mcp/github-copilot.md` currently documents `.vscode/mcp.json`
  and says VS Code forwards supported server configurations to the Copilot
  Agent Host.
- The existing active item `b152c56c` covers provider MCP setup and the HAWP
  tool allowlist; its remaining item is a live session check of that custom
  agent's tool scope.
- Official references: [MCP configuration](https://code.visualstudio.com/docs/agents/reference/mcp-configuration),
  [customization migration](https://code.visualstudio.com/docs/agent-customization/migrate-customizations),
  and [custom agents](https://code.visualstudio.com/docs/agent-customization/custom-agents).

**Inferred (not yet proven):**

- The current `.vscode/mcp.json` implementation may remain valid for VS Code
  while the portable `.mcp.json` format is preferable for the newer Copilot
  harness. The migration notice alone does not establish that the provider
  installer is broken.
- The provider install/update behavior may need to merge GitHub Copilot's
  server into the existing root `.mcp.json` rather than keep two equivalent
  entries; this requires confirmation against current docs and a live client.

**Likely scope:**

- Provider configuration and tests under
  `librarian/src/internal/platform/mcp/configure/`, the generated HAWP agent,
  and `.hawp/kit/usage/mcp/github-copilot.md`.
- Keep this review focused on current VS Code migration/install behavior.
  Do not redo the already-completed Codex or generic tool-allowlist work from
  `b152c56c`.

## Risk + Review Gate

**Risk:** medium
**Gate:** Present confirmed findings and a minimal change plan before editing
shared provider code or documentation.

## Backlog + Plan Link

**Status now:** parked
**Plan file:** work/parked/70d8112f/plan.md

## Work Coordination

- **Owner:** unassigned
- **Implementation status:** not-started
- **Overlapping files:** `librarian/src/internal/platform/mcp/configure/`,
  `.github/agents/hawp.agent.md`, `.hawp/kit/usage/mcp/github-copilot.md`
- **Parallel work risk:** medium
- **Can implement now:** no — resume after `cb096b6a` reconciles the local
  config, then coordinate the review with active item `b152c56c`.
- **Coordination note:** `b152c56c` already owns the original provider setup
  implementation. This follow-up is limited to changes in current VS Code
  Copilot MCP migration/install behavior; do not duplicate its completed
  allowlist work.

## Options

### Option A — Keep the current provider format if still supported

Verify the existing `.vscode/mcp.json` install/update flow against current
official docs and a live Copilot session. If it remains supported, document
the migration choice and avoid unnecessary code changes.

### Option B — Update the installer to the portable Copilot format

Change the provider to merge into `.mcp.json` only if official docs and live
tests confirm that is the required current behavior. Account for the fact that
Claude Code already uses `.mcp.json`, and preserve existing entries.

## Recommended Approach

Run a bounded, evidence-led audit after `cb096b6a` completes. Do not assume the
migration notice means the existing provider logic is invalid. If a mismatch is
confirmed, present the smallest change and test plan for review before editing.

## Plan

1. After `cb096b6a` closes, re-check current official VS Code MCP, migration,
   and custom-agent documentation.
2. Trace `hawp init --provider github` and `hawp mcp configure --provider
   github` through provider selection, preflight, JSON merge, HAWP agent
   generation, and update behavior.
3. Compare the installed output with current VS Code requirements, including
   existing root `.mcp.json` entries, duplicate server names, settings
   preservation, workspace scope, and tool-call approval behavior.
4. Add or update focused tests for fresh install, repeat/update, existing
   unrelated servers, same-name HAWP entries, and unmanaged custom-agent
   protection as applicable.
5. Verify in VS Code Copilot that the HAWP server and custom agent work in the
   intended repository and that write-capable tools still require normal
   approval.
6. If a mismatch is confirmed, propose the minimal code/docs change and wait
   for approval before implementation.

## Verification

- [ ] Record current official docs and distinguish confirmed requirements
  from inferences.
- [ ] Tests cover the actual install/update contract, including existing
  configuration preservation.
- [ ] Live Copilot verification confirms repository scope, expected HAWP tools,
  and manual call approval.
- [ ] Any code or documentation changes are validated with focused Go tests
  and the appropriate HAWP checks.

## Next Step

- [x] Investigation recorded above (required before planning)
- [x] Write or update the plan file
- [x] Move backlog status accordingly
