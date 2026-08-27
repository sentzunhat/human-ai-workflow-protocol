# usage-tracking — `hawp usage` CLI + `hawp_usage` MCP tool

**Type:** feature  
**Status:** parked  
**Opened:** 2026-08-25  
**Parked reason:** Too large for v0.0.11 incremental patch strategy. Prerequisite for
v0.1.0 token-reduction gate (actual measured savings, not just retrieval quality).

## Input

Downstream install of hawp 0.0.9. The install team wanted to verify that
hawp MCP search actually reduced Cursor cloud token spend. No tool exists to
measure this. `--max-tokens` counts packed context size (chars/4 estimate), not
Cursor subscription billing.

Evidence source: downstream install evidence 2026-08-25.

## Goal

An optional, local-only usage log that records:
- prompt token count, completion token count (when available from MCP caller)
- model, backend (cursor-cloud / ollama / onnx)
- tool name, query hash (not raw prompt text)
- wall clock

Expose via:
- CLI: `hawp usage` (show totals), `hawp usage log` (tail recent entries)
- MCP: `hawp_usage` tool (structured query)
- Auto-record: optional hook on `hawp_search`, `hawp_work_new` calls

## Constraints

- Default off or opt-in; never log raw prompt/response bodies unless opted in.
- No secrets; no external telemetry.
- Must distinguish local Ollama tokens from cloud tokens.
- Do not expand HAWP into a runtime or billing system.

## Why parked

- Token metering requires the MCP caller to pass token counts (not all callers do).
- `hawp_usage` MCP tool adds scope that belongs in a minor version bump.
- The v0.1.0 gate (demonstrable token reduction) is the right milestone for this.

## Resume when

- v0.1.0 planning begins
- A downstream team asks for this as a blocking requirement
- Token count pass-through from MCP callers is standardized

## Relationship

- Extends [[v010-cost]] (parked cloud-API cost tracking) with a local/Ollama dimension.
- Feeds the v0.1.0 success metric: "tokens in vs tokens saved".
