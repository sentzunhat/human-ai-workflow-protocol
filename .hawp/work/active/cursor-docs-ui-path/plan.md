# cursor-docs-ui-path — Fix stale Cursor MCP UI path in kit docs

**Type:** docs  
**Status:** plan-ready  
**Opened:** 2026-08-25  
**Target:** v0.0.11

## Input

Downstream install of hawp 0.0.9 (Cursor provider). Existing Cursor MCP docs
(kit and possibly generated install guide) reference Settings → Tools & MCP
as the UI path to enable a workspace MCP server. That path is stale in current
Cursor. The correct path is: sidebar **Customize → MCPs**, toggle the
**workspace** server toggle.

Additional: enabling an MCP server in Cursor may require opening a new chat
or reloading the window before the agent actually sees the tools.

Evidence source: confirmed by the downstream install session on 2026-08-25.

## Goal

Kit docs show the correct Cursor UI path and include the reload note.
No stale `Settings → Tools & MCP` references remain.

## Constraints

- Docs only. Do not change CLI behavior.
- Only update what was confirmed stale — do not invent other UI details.

## Plan

### Step 1 — Grep for stale references

```bash
grep -rn "Settings.*Tools.*MCP\|Tools.*MCP.*Settings\|Tools & MCP" .hawp/kit/ distribution/
```

### Step 2 — Replace with correct path

Replace stale references with:

> Enable in Cursor: sidebar **Customize → MCPs** → toggle the **workspace**
> server. After enabling, open a new chat or reload the window for the agent
> to see the tools. File presence of `.cursor/mcp.json` alone does not guarantee
> the agent has access.

### Step 3 — Validate links

Run `npm run check:markdown-links`.

## Verification

- No `Settings → Tools & MCP` references in kit or distribution docs
- Correct UI path present wherever Cursor MCP setup is described
