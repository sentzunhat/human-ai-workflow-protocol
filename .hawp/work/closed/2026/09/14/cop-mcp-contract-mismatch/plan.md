# Bug: MCP response contract mismatches (state names, JSON tags, draft-on-no-match)

**UUID:** `cop-mcp-contract-mismatch`
**Type:** fix
**Severity:** Medium
**Reported:** 2026-09-13
**Status:** `inbox`
**Source:** Copilot PR #40 review — multiple files

---

## Finding

Three related contract inconsistencies found across `types.go`, `tool_work.go`,
and `CHANGELOG.md`:

### 1. JSON tag mismatch (`types.go:101-102`)
Struct fields emit `output` and `checkpoint`, but the release contract
documents `output_spec` and `done_signal`:
```go
Output      string `json:"output"`
Checkpoint  string `json:"checkpoint,omitempty"`
```
Consumers using the advertised schema cannot find either field.

### 2. State name mismatch (`tool_work.go` no-index path)
Missing-index path emits `blocked_missing_index`; the documented contract says
`missing_index`. Clients matching the advertised value misclassify the result.

### 3. Draft returned on no-match (`tool_work.go:208-212`)
When `ChunksUsed == 0` the handler sets `State = "needs_user_input"` but still
returns the populated `Draft`. The contract says uncertain states should not
return a usable draft.

## Fix Plan

1. **JSON tags** — change `json:"output"` → `json:"output_spec"` and
   `json:"checkpoint,omitempty"` → `json:"done_signal,omitempty"` in `types.go`.
   Update any test assertions that check these field names.

2. **State name** — change `"blocked_missing_index"` → `"missing_index"` in
   `tool_work.go`. Verify `tools_e2e_test.go` assertions and update if needed.

3. **Draft on no-match** — when `ChunksUsed == 0`, set `response.Draft = WorkIntakeDraft{}`
   (zero value) before returning. Add a test asserting draft is empty in this state.

4. **CHANGELOG** — correct the v0.0.24 notes to match the shipped JSON field
   names and state values.

## Files

- `librarian/src/internal/platform/mcp/server/types.go`
- `librarian/src/internal/platform/mcp/server/tool_work.go`
- `librarian/src/internal/platform/mcp/server/tools_e2e_test.go` (if it exists)
- `librarian/src/CHANGELOG.md`

## Verification

- `WorkIntakeDraft` JSON tags emit `output_spec` and `done_signal` (not `output`/`checkpoint`).
- `missing_index` state is returned when no index exists (not `blocked_missing_index`).
- `response.Draft` is nil/zero when `ChunksUsed == 0` (no indexed context matched).
- `go test ./...` passes including MCP e2e tests; `go vet ./...` clean.

## Outcome

Fixed `types.go` JSON tags to `output_spec`/`done_signal`. Corrected state string to `missing_index` in `tool_work.go`. Cleared `response.Draft` on the no-match path. All MCP contract fields now match the documented schema. `go test ./...` and `go vet ./...` pass.

## Close Checklist

- [x] JSON tags `output_spec`/`done_signal` match documented contract.
- [x] State `missing_index` matches documented contract.
- [x] Draft cleared on no-match path.
- [x] `go test ./...` passes; `go vet ./...` clean.
- [x] Plan status: done.
- [x] Plan moved to `.hawp/work/closed/2026/09/14/cop-mcp-contract-mismatch/plan.md`.
- [x] BACKLOG.md updated.
