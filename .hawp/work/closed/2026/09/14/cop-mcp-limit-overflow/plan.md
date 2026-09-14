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
