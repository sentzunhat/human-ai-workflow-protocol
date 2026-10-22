---
uuid: b152c56c
title: "HAWP provider verification continuation"
type: status
date: 2026-10-21
---

# 2026-10-21 — HAWP provider scope verification

## Status Report

### Intent

Resume step 6 of `.hawp/work/active/b152c56c/plan.md`: verify the generated
HAWP agent's tool scope in a session explicitly started with that agent.

### Current State

`b152c56c` remains in progress. The local generated agent definition can be
inspected, but this conversation is running as Master Builder and has no
VS Code session-control surface to select the HAWP agent or inspect a live
agent's tool picker. The live closeout gate is therefore still open.

The duplicate-config item `cb096b6a` is plan-ready, installer audit
`70d8112f` is parked pending reconciliation, and protocol/structured-output
item `9c9f9bed` remains plan-ready.

### What Was Inspected

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/b152c56c/plan.md`
- `.github/agents/hawp.agent.md`
- `.hawp/work/active/cb096b6a/plan.md`
- `.mcp.json` and `.vscode/mcp.json` (read-only comparison)
- HAWP MCP search and validation results

### What Changed

No provider code or local MCP configuration was changed. This status report
records the current verification gate and the read-only config preflight.

### What Was Directly Verified

- `.github/agents/hawp.agent.md` lists seven HAWP MCP tools:
  `hawp/hawp_search`, `hawp/hawp_usage`, `hawp/hawp_work_intake`,
  `hawp/hawp_work_new`, `hawp/hawp_work_validate`, `hawp/hawp_work_doc`, and
  `hawp/hawp_work_reshape`.
- `hawp-hawp_search` returned no indexed results for the four continuation
  items; source files were checked directly.
- `hawp-hawp_work_validate` passed all three validation groups: kit checks
  passed with zero issues, work validation reported zero issues and one
  warning, and local links were valid across 152 Markdown files.
- Both local MCP config files contain only the `hawp` server and parse as
  JSON. Their command and arguments match; `.vscode/mcp.json` uses the VS Code
  `servers` shape with `type: stdio`, while `.mcp.json` uses `mcpServers`.
- `.gitignore` excludes both config files, and neither is tracked by Git.
- The pre-existing working-tree changes were left untouched.

### What Remains Unproven

- The exact tools exposed by a VS Code session explicitly started with the
  generated HAWP agent, and whether any unrelated tools are exposed.
- Whether normal client approval prompts remain in effect for write-capable
  tools; no write-capable tool was invoked for this check.
- Which config file the intended VS Code session is loading, and how the
  migration preview handles the existing same-name `hawp` entry.

### Constraints

- Keep manual approval for write-capable tools.
- Do not edit either ignored config file until the migration behavior and
  intended client setup are confirmed.
- Preserve existing unrelated working-tree changes.

### Help Wanted

In a fresh VS Code Copilot session, explicitly select the generated `HAWP`
agent and inspect its Tools list. Confirm whether it contains exactly the
seven tool IDs listed above. Also inspect the MCP migration preview and report
whether it merges with or replaces the existing `hawp` entry in `.mcp.json`.
Do not invoke `hawp_work_new` or `hawp_work_doc` as part of this scope check.

### Suggested Next Step

Record the live tool inventory and migration-preview behavior, then update
`b152c56c` and reconcile the local duplicate under `cb096b6a`. Resume
`70d8112f` only after reconciliation. Keep the separate modern/legacy client
matrix and structured-output implementation in `9c9f9bed`.
