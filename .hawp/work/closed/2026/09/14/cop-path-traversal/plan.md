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

## Verification

- `folderID` values like `../../../../../../tmp` are rejected with an error before any directory is created.
- Valid short UUIDs (8 hex chars) and full UUIDs (36-char) are accepted.
- `go test ./...` passes; `go vet ./...` clean.

## Outcome

Added `validateFolderID` helper in `doc.go` that rejects any ID not matching short or full UUID format. MCP and CLI callers pass raw user input through the same validation gate. Unit tests cover the traversal case and invalid-ID rejection. `go test ./...` and `go vet ./...` pass.

## Close Checklist

- [x] Path traversal rejected at `doc.go` before directory creation.
- [x] `go test ./...` passes; `go vet ./...` clean.
- [x] Plan status: done.
- [x] Plan moved to `.hawp/work/closed/2026/09/14/cop-path-traversal/plan.md`.
- [x] BACKLOG.md updated.
