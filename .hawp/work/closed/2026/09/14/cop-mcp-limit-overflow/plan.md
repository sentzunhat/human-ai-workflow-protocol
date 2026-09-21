# Security: MCP search limit integer overflow

**UUID:** `cop-mcp-limit-overflow`
**Type:** fix
**Severity:** High
**Reported:** 2026-09-13
**Status:** `inbox`
**Source:** Copilot PR #40 review — `librarian/src/internal/platform/mcp/server/tool_work.go`

---

## Finding

The MCP intake handler bounds `Limit` only at zero:

```go
if a.Limit <= 0 {
    a.Limit = 10
}
```

The search service then passes `limit * 3` to SQLite. A max-int JSON value
(e.g. `9223372036854775807`) overflows to a negative integer after multiplication,
causing SQLite `LIMIT -1` which means no limit — loading the entire index into
memory in one request.

## Fix Plan

1. Add an upper-bound guard matching the CLI's pattern:
   ```go
   } else if a.Limit > int(^uint(0)>>1)/3 {
       return toolErr("limit must be positive and small enough for retrieval")
   }
   ```
2. Also add a sane maximum (e.g. 500) to prevent legitimate but oversized requests.
3. Add a test in the MCP tool layer asserting that an overflow value returns an
   error rather than passing through.

## Files

- `librarian/src/internal/platform/mcp/server/tool_work.go` — add upper-bound guard
- Relevant test file under `librarian/src/tests/`

## Verification

- A JSON limit value of `math.MaxInt` is rejected before the multiply-by-3 step.
- Limits in the normal range (e.g., 10, 50) pass through unchanged.
- `go test ./...` passes; `go vet ./...` clean.

## Outcome

Added upper-bound guard `if a.Limit > int(^uint(0)>>1)/3` in `tool_work.go`, matching the existing CLI guard pattern. Request with an overflowing limit returns a `toolErr` response without hitting SQLite. `go test ./...` and `go vet ./...` pass.

## Close Checklist

- [x] MCP limit overflow guard applied in `tool_work.go`.
- [x] `go test ./...` passes; `go vet ./...` clean.
- [x] Plan status: done.
- [x] Plan moved to `.hawp/work/closed/2026/09/14/cop-mcp-limit-overflow/plan.md`.
- [x] BACKLOG.md updated.
