# usage-tracking — `hawp usage` CLI + `hawp_usage` MCP tool

**UUID:** `e1dc28ac-58eb-42a7-a9f2-3cb8cd624764`
**Legacy ID:** `usage-tracking` (preserved from the parked pre-UUID record)
**Type:** feature  
**Status:** in-progress  
**Opened:** 2026-08-25  
**Target:** Complete the token-reduction evidence gate and resolve remaining
accuracy/privacy gaps; exact release target should be aligned with current
release planning.

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

## Resume investigation (2026-10-16)

### Directly verified complete / already covered

- The feature is not wholly missing. The prior closed item
  `.hawp/work/closed/2026/08/27/usage-log/plan.md` records the local usage log
  and CLI as shipped in v0.0.12, with query summary added in v0.0.17.
- Current MCP dispatch logs calls when usage logging is enabled. The local
  SQLite store records timestamp, tool, query hash, token estimates based on
  serialized request/response size, and a query/title summary. Raw input/output
  bodies are stored only when `log_bodies` is enabled.
- `hawp_usage` exposes local totals and a report; CLI commands support enable,
  disable, report/log, and clear. These current functions must not be
  reimplemented as a duplicate.
- Current source: `librarian/src/internal/platform/mcp/server/tools.go`,
  `tool_usage.go`, `librarian/src/internal/infrastructure/repositories/usage/store.go`,
  `librarian/src/internal/domain/usage/types.go`, and
  `librarian/src/internal/platform/cli/usage/commands.go`.

### Remaining gap against the original goal

- Stored token counts are estimates from JSON byte length divided by four; they
  are not caller/provider-reported token counts and do not isolate model prompt
  and completion tokens.
- A timestamp is recorded, but call duration/elapsed wall-clock time is not.
- Model/backend identity is not recorded with each MCP usage entry.
- Query/title text is stored up to 256 characters even with raw-body logging
  disabled. Full bodies are separately opt-in. Retention and deletion are
  manual (`hawp usage clear`); there is no expiry policy in the inspected code.
- These gaps are distinct from the MCP capability catalog (`2eea565c`): the
  catalog documents the existing usage capability and its privacy boundary but
  does not implement accurate metering.

### Updated scope and plan

1. Preserve the existing opt-in local usage logger and CLI/MCP interfaces.
2. Decide which additional measurements are technically available from each
   supported client/backend: elapsed duration, model/backend labels, and actual
   provider token counts. Do not claim cloud billing or actual token usage when
   only payload-size estimates are available.
3. Define a privacy-minimal default schema and explicit query/body retention
   controls. Keep telemetry local; no external reporting.
4. Add only measurements with reliable sources, label estimates versus actual
   counters, and keep logging failures from breaking MCP calls.
5. Verify migration compatibility, opt-in defaults, clear/retention behavior,
   reporting semantics, and end-to-end metering evidence for the v0.1.0 gate.

**Status:** `in-progress` — existing base implementation is complete; accurate
duration/provider-token/model metadata and retention-policy review remain
unresolved. This work item is resumed, not duplicated or replaced by the
capability catalog.
