---
description: "HAWP workflow core - kit entry, status reports, backlog discipline"
---

<!-- Generated from core/providers/shared/behaviors - edit shared source and run `hawp providers sync` -->

# HAWP Core

This repo uses HAWP. Read `.hawp/kit/start-here.md` for task shaping, `.hawp/kit/usage/status-report.md` for handoffs, and `.hawp/kit/usage/workflow-loop.md` for multi-iteration work.

- Track work in `.hawp/work/BACKLOG.md`; plans in `.hawp/work/active/`; close to `.hawp/work/closed/YYYY/MM/DD/`.
- Do not append completed work endlessly to BACKLOG.md — keep Recently Closed capped.
- HAWP is not a runtime, validator, orchestrator, or memory system.
- Do not invent per-field runtime folders (input/, context/, mission/, constraints/, output/, checkpoint/) unless repo policy requires them.
- Prefer compact, decision-useful outputs. Separate direct evidence from inference.

## HAWP MCP

When HAWP MCP tools are available, use the existing tool groups as a workflow
index: **Discover** (`hawp_search`, `hawp_usage`), **Shape**
(`hawp_work_reshape`, `hawp_work_intake`), **Create** (`hawp_work_new`,
`hawp_work_doc`), and **Verify** (`hawp_work_validate`). These are guidance
categories; tool names and boundaries stay unchanged. Check the installed
`.hawp/kit/usage/mcp/README.md` for tool effects, local usage logging privacy,
and manual workflows.

HAWP's local server may expose the fixed read-only resources and
`hawp_intake_request` prompt. Use them only when the connected client actually
lists them; client documentation is not proof that a given HAWP connection
exposes them. If resources or prompts are unavailable, use the MCP guide and
existing tools. Intake prompts and shaping do not create work; `hawp_work_new`
and `hawp_work_doc` are explicit writes.
