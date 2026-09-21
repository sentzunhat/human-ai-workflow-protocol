# `hawp_work_intake` MCP tool: compound search + reshape in one call

**UUID:** `a8797f44-...` (see BACKLOG.md)
**Type:** feature
**Reported:** 2026-09-12
**Status:** `done`
**Target Version:** `v0.0.24`

---

## Input (verbatim)

> the search and reshape should be one compound MCP call so agents dont drift by skipping either step. Hawp is a tool for the harness of ai agent providers for getting less drift and then we need to do a compact or I guess the reshape is the compact and how do we help ai agents like yourself for context and tokens savings with this hawp tool

## Intake Summary

HAWP's value loop for AI agents:

1. **Index** — `hawp search index` + embed: build the local knowledge base
2. **Retrieve** — `hawp_search`: get relevant kit + work context (~600–2000 tokens)
3. **Shape/Compact** — `hawp_work_reshape`: compress (request + context) → Mission/Constraints/Output/Checkpoint (~30–211 tokens, 95% savings)
4. **Create** — `hawp_work_new`: scaffold the work item with the shaped fields
5. **Document** — `hawp_work_doc --work-item <id>`: evidence/status/decision tied to the item

The gap: steps 2 and 3 are two separate MCP calls. An agent must remember to call search first, capture the context string, then pass it to reshape. Agents that skip search always pass empty context to reshape — defeating the token-reduction benefit entirely.

**`hawp_work_intake`** combines search + reshape into one call: agent passes the raw request, the tool retrieves context via hybrid search, runs reshape, and returns the shaped Draft fields ready for `hawp_work_new`. One call, no skipped steps.

## Current State (verified 2026-09-12)

- `hawp_search` returns a context block string (markdown)
- `hawp_work_reshape` accepts `{input, context, model, url}` — context is a free string
- No compound tool exists; agents must orchestrate manually

## 2026-09-13 MCP Architecture Lesson Queue

User shared the video `MCP Just Got a Whole Lot Better` by Neon Postgres as an
MCP architecture lesson. Transcript retrieval was blocked during intake, so the
lesson queue is grounded in confirmed video metadata plus current MCP
architecture themes rather than a transcript summary.

Follow-on plan-ready items created from the lesson pass:

- `2eea565c` — MCP capability catalog for HAWP resources and prompts
- `429e075e` — Remote MCP transport and authorization readiness audit
- `c8d9e1cb` — Interactive MCP intake refinement with structured results

This item remains the first concrete implementation slice: make search +
reshape one compound HAWP MCP call so agents cannot skip the context step.

## Release Scope

This is the main `v0.0.24` implementation candidate. It may absorb the smallest
structured-result subset from `c8d9e1cb` if that improves intake safety without
expanding the release into resources, prompts, remote MCP, or usage metering.

Version-scope note: see
`.hawp/work/status/2026/09/13/6059dd4b/status.md`.

## Plan

### Slice A — `hawp_work_intake` MCP tool

Input: `{input: string, model?: string, url?: string, limit?: int}`

Internally:
1. Run hybrid search with `input` as query (default limit 10)
2. Format results as context block string via `ContextBlock.String()`
3. Run `DraftIntake` (existing use case) with `{Input: input, Context: contextBlock}`
4. Return shaped Draft fields + token accounting

Output:
```json
{
  "mission": "...",
  "constraints": "...",
  "output": "...",
  "checkpoint": "...",
  "context_tokens": 1240,
  "shaped_tokens": 89,
  "savings_pct": 93
}
```

Agents can then pass `mission`/`constraints`/`output`/`checkpoint` directly to `hawp_work_new`.

Implementation:
- [x] Add `toolWorkIntake` in `platform/mcp/server/tool_work.go`
- [x] Add `hawp_work_intake` to `toolDefs()` in `tools.go`
- [x] Add case to `callTool()` dispatch
- [x] Tests: mock search index + mock shaper; verify savings accounting

### Slice B — Token accounting surface

The `context_tokens` / `shaped_tokens` / `savings_pct` fields in the response give agents a self-reported token budget summary. This surfaces HAWP's value inline, every call — agents can log or surface this to users.

## Implementation Notes (2026-09-13)

Implemented `hawp_work_intake` as a compound MCP tool that:

- Runs indexed search and formats the result through the same context block path
  used by `hawp_search` with `context:true`.
- Runs the existing `DraftIntake` use case with the retrieved context.
- Returns structured JSON with `state`, `draft`, `retrieval`, `questions`,
  `warnings`, and `token_accounting`.
- Uses `ready_for_work_new` when context is found and shaped.
- Uses `needs_user_input` when no indexed context matches.
- Uses `blocked_missing_index` when the index is missing or not initialized.
- Uses `blocked_reshape_failed` when retrieval succeeds but the local model
  returns an incomplete or invalid draft.

The implementation intentionally did not add MCP resources, prompts, remote
transport, or usage-metering scope.

## Risk

Medium: requires wiring the search index (needs a live SQLite index) and Ollama shaper together inside a single MCP handler. No new domain logic — reuses existing `ContextBlock`, `DraftIntake`, and `OllamaIntakeShaper`.

## Verification

- [x] Unit/integration coverage verifies `ready_for_work_new`,
  `needs_user_input`, `blocked_missing_index`, and `blocked_reshape_failed`
  responses.
- [x] Tool discovery smoke confirms `hawp_work_intake` appears in MCP
  `tools/list`.
- [x] Live MCP/Ollama smoke confirms `needs_user_input` response when the index
  has no matching chunks.
- [x] Live MCP/Ollama smoke confirms reshape failures are reported as structured
  blocked states instead of confident drafts.
- [x] `go test ./...` passes from `librarian/src`.
- [x] `go vet ./...` passes from `librarian/src`.
- [x] Live MCP stdio smoke returns valid initialization and an honest
  `needs_user_input` response when the current repository index has no match.
- [x] Ready-state behavior is proven by repository-backed integration coverage;
  a real local model producing a ready draft remains unproven in this checkout.
- [x] `hawp work validate` and `hawp check` pass with the current repository
  state; historical clarity warnings remain explicitly reported where present.

## Outcome

`hawp_work_intake` is complete for the v0.0.24 scope. It compounds indexed
search and local reshaping, exposes token accounting, and returns explicit
`ready_for_work_new`, `needs_user_input`, `blocked_missing_index`, and
`blocked_reshape_failed` states without fabricating a draft when context or
local model readiness is absent.

## Close Checklist

- [x] Compound MCP tool implemented and discoverable.
- [x] Structured states, retrieval metadata, and token accounting verified.
- [x] Repository-backed ready-state and blocked-state tests pass.
- [x] Live stdio smoke verified honest no-match behavior.
- [x] `go test ./...`, `go vet ./...`, and `hawp check` pass.
- [x] Plan moved to `closed/2026/09/19/a8797f44/`.
- [x] `BACKLOG.md` updated.
