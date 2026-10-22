# Reconcile VS Code HAWP MCP config duplicates

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `cb096b6a-99b3-4042-b4ed-cad5a02a8788`
**Type:** improvement
**Reported:** 2026-10-21

---

## Input (verbatim)

> Reconciling the duplicate first and seeing what the diff is on both and reconcile into a good one.

## Intake Summary

Compare the two local HAWP MCP definitions, determine which format should be
canonical for the current VS Code Copilot setup, and reconcile the duplicate
without losing other server entries or settings.

## Current Context

The VS Code Agent Customizations screen offers to migrate the workspace entry
from `.vscode/mcp.json` to `.mcp.json`. Both files currently contain a `hawp`
server entry. They are ignored by Git and contain machine-local paths.

## Initial Analysis

**Directly verified:**

- `.vscode/mcp.json` uses the VS Code `servers` shape with `type: stdio`.
- `.mcp.json` uses the portable `mcpServers` shape. Both entries currently
  specify the same HAWP executable and repository root.
- `.gitignore` excludes both config files, so this is local setup rather than
  a team-shared configuration.
- Current VS Code docs describe `.vscode/mcp.json` as the VS Code format and
  `.mcp.json` as the portable Copilot format. The migration guide says
  migration can drop VS Code-specific metadata and advises resolving duplicate
  server names.
- The read-only `hawp-hawp_work_validate` MCP check passed all three validation
  groups; work validation reported zero issues and one warning.

**Inferred (not yet proven):**

- The migration notice is caused by the two equivalent `hawp` entries using
  different configuration formats. The existing `.mcp.json` entry may collide
  with the proposed migration; the migration's behavior for this exact
  same-name destination has not been tested.

**Likely scope:**

- `.mcp.json` and `.vscode/mcp.json` only, unless live VS Code verification
  shows that the existing provider guide also needs a narrowly scoped update.

## Risk + Review Gate

**Risk:** low
**Gate:** This request is for planning only. Do not change either config until
the user asks to execute the plan.

## Backlog + Plan Link

**Status now:** plan-ready
**Plan file:** work/active/cb096b6a/plan.md

## Work Coordination

- **Owner:** unassigned
- **Implementation status:** not-started
- **Overlapping files:** `.mcp.json`, `.vscode/mcp.json`
- **Parallel work risk:** low
- **Can implement now:** no — current request is planning only.
- **Coordination note:** Complete this config reconciliation before item
  `70d8112f` reviews the provider implementation. Re-capture repository-root
  proof and the worktree status before implementation; the current checkout
  already contains unrelated in-progress changes.

## Options

### Option A — Use the portable workspace config as canonical

Keep the single `hawp` definition in `.mcp.json`, preserve any other entries,
and remove the duplicate from `.vscode/mcp.json` only after confirming that the
current VS Code Copilot session loads the portable entry and HAWP still starts.

### Option B — Keep the VS Code-specific config

Retain `.vscode/mcp.json` and leave `.mcp.json` for its existing provider use.
Ignore the migration if live verification shows the VS Code-specific entry is
still required. This avoids changing runtime behavior but leaves the migration
notice unresolved.

## Recommended Approach

Compare the full files and the migration preview before choosing. Prefer Option
A if current VS Code Copilot and the existing `.mcp.json` consumer both work
from the portable format. Otherwise keep Option B and record why migration is
not appropriate. Do not overwrite or discard unrelated server entries.

## Plan

1. Capture a focused diff of `.vscode/mcp.json` and `.mcp.json`; confirm the
   complete contents and any other MCP server entries before editing.
2. Re-check the current VS Code MCP configuration and migration documentation,
   including the changes shown by the migration UI.
3. Choose one canonical HAWP definition based on observed client behavior; do
   not assume the migration merges an existing same-name server.
4. Apply only the minimal local config change. Keep the machine-specific files
   ignored and preserve the existing HAWP command, repository scope, and other
   server settings.
5. Validate JSON, start HAWP in VS Code, and call the read-only
   `hawp_work_validate` MCP tool from the intended Copilot workspace.
6. Confirm the migration no longer represents a duplicate and normal tool
   approval prompts remain enabled.

## Verification

- [ ] Record the before/after diff and the chosen canonical file.
- [ ] Both retained config files parse as JSON.
- [ ] VS Code starts HAWP from the intended repository and Copilot can call
  `hawp_work_validate`.
- [ ] Any non-HAWP entries and normal tool-call approvals are preserved.
- [ ] Confirm `.mcp.json` and `.vscode/mcp.json` remain ignored by Git.

## Next Step

- [x] Investigation recorded above (required before planning)
- [x] Write or update the plan file
- [x] Move backlog status accordingly
