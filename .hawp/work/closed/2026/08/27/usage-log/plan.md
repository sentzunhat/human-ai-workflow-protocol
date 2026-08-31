# usage-log — Local call log with token counts

**Type:** feature
**Status:** in-progress
**Branch:** `feature/v0.0.12`
**Opened:** 2026-08-25

## Problem

There is no way to know whether `hawp_search` is actually saving tokens versus
just fetching context. The v0.1.0 success gate requires a demonstrable token
reduction; right now it is unmeasurable. Users in downstream installs (Cursor
+ Ollama, 2026-08-25) asked for spend visibility but got nothing.

## Goal

A local-only, opt-in SQLite log that records every `hawp_search` and
`hawp_work_new` MCP call: timestamp, tool name, query hash, estimated tokens
in and out. A `hawp usage` CLI surfaces totals and recent entries. Optionally
stores raw input/output bodies when the user explicitly enables body capture.

## Scope

### Always recorded (when logging enabled)

| Column | Type | Notes |
|--------|------|-------|
| `id` | INTEGER PK | auto |
| `ts` | TEXT | RFC3339 |
| `tool` | TEXT | `hawp_search`, `hawp_work_new`, `hawp_work_validate` |
| `query_hash` | TEXT | sha256 of first 256 bytes of input JSON — not raw text |
| `tokens_in` | INTEGER | `len(request_json) / 4` (char estimate) |
| `tokens_out` | INTEGER | `len(response_json) / 4` (char estimate) |

### Optional (body capture — explicit second opt-in)

| Column | Type | Notes |
|--------|------|-------|
| `input_body` | TEXT | raw input JSON; NULL when body capture off |
| `output_body` | TEXT | raw output text; NULL when body capture off |

### Not recorded

- Embedding vectors
- File contents outside of the MCP call
- Any data from non-hawp tools

## Implementation

### New package: `internal/domain/usage/`

```
log.go     — LogEntry struct, Logger interface
store.go   — SQLite-backed store, schema migration
store_test.go
```

Store lives at `~/.hawp/usage.db` — separate from the search index so a full
`hawp search index` never touches it.

### Config: `~/.hawp/config/usage.json`

```json
{
  "enabled": false,
  "log_bodies": false
}
```

### CLI: `hawp usage`

```
hawp usage             — show totals (calls, tokens in, tokens out, sessions)
hawp usage log         — tail 20 most recent entries
hawp usage enable      — turn on logging
hawp usage disable     — turn off logging
hawp usage enable --log-bodies   — enable body capture (sets log_bodies: true)
hawp usage clear       — truncate the log (irreversible, prompts for confirm)
```

### MCP wire-in

In the MCP handler dispatch loop, after each tool call returns:
- Compute tokens_in / tokens_out from JSON byte lengths
- Hash query field for query_hash
- Fire a goroutine to write the log entry (non-blocking; log write failure is
  silent — never fails the tool response)

### Tests

- `store_test.go`: schema creation, Write, ReadRecent, Totals, idempotent schema
- `cli_test.go`: `hawp usage` output format, enable/disable toggle

## Acceptance criteria

- [ ] `~/.hawp/usage.db` created on first logged call (not on binary launch)
- [ ] `hawp usage enable` + one `hawp_search` call → entry in db
- [ ] Token estimates present; bodies NULL unless log_bodies enabled
- [ ] `hawp usage` prints totals correctly
- [ ] `hawp usage log` shows 20 most recent, newest first
- [ ] Log write failure never blocks or errors the MCP response
- [ ] All tests pass (`go test ./...`)
- [ ] `hawp usage` included in `hawp help` output

## Relationship to v0.1.0

Usage totals from this log are the evidence base for the token-reduction gate.
Once a downstream team can run `hawp usage` and see "saved N tokens this week",
the v0.1.0 milestone becomes verifiable.

## Outcome

Shipped in v0.0.12 (2026-08-25). `hawp usage`, `hawp usage log`, `hawp usage enable/disable/clear`
all wired and tested. SQLite store at `~/.hawp/usage.db`, config at `~/.hawp/config/usage.json`.
MCP wire-in logs every `hawp_search`, `hawp_work_new`, and `hawp_work_validate` call
asynchronously (fire-and-forget goroutine; log failure never blocks the tool response).
v0.0.17 extended with `query_text` column for always-readable query summaries.

## Verification

- [x] `go test ./internal/domain/usage/...` — 11 tests pass. Evidence: the focused usage package coverage is summarized in this plan's Outcome section.
- [x] `go test ./internal/platform/cli/...` — includes TestRunUsage, TestRunUsageClearCancelled. Evidence: the CLI coverage is summarized in this plan's Outcome section.
- [x] `go test ./internal/platform/mcp/...` — MCP e2e tests pass. Evidence: the MCP coverage is summarized in this plan's Outcome section.
- [x] `go test ./...` — full suite green. Evidence: the release-ready verification chain is summarized in this plan's Outcome section.

## Close Checklist

- [x] Outcome recorded
- [x] Verification covers domain, CLI, and MCP layers
- [x] Feature shipped in v0.0.12; query_text enhancement in v0.0.17
- [x] Ready to stay in closed history
