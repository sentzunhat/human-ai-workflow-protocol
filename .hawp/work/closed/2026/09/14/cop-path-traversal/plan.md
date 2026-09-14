# Security: path traversal in work doc creation

**UUID:** `cop-path-traversal`
**Type:** fix
**Severity:** High
**Reported:** 2026-09-13
**Status:** `inbox`
**Source:** Copilot PR #40 review — `librarian/src/internal/application/work/doc/doc.go`

---

## Finding

`folderID` is inserted into the filesystem path without any format validation:

```go
docDir := filepath.Join(workDir, dirName, date[:4], date[5:7], date[8:10], folderID)
```

A value such as `../../../../../../tmp` escapes `.hawp/work` and writes the
document outside the work directory. The same risk applies to any caller that
passes an unvalidated work-item ID from MCP tool input.

## Fix Plan

1. Validate `folderID` before path join: accept only short UUIDs (8 hex chars)
   or full UUIDs (36-char `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`).
2. Return an error if the ID does not match — do not create the directory.
3. Add a unit test covering the traversal case and an invalid ID.
4. Review all callers (`hawp_work_doc` MCP handler, CLI `work doc` command) to
   confirm they pass raw user input through and cannot bypass this check.

## Files

- `librarian/src/internal/application/work/doc/doc.go` — add validation helper
- `librarian/src/internal/platform/mcp/server/tool_work.go` — confirm input is passed through
- New test in `librarian/src/tests/application/work/` or inline in `doc/`
