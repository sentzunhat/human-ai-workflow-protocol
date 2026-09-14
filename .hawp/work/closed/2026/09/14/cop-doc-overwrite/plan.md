# Bug: silent overwrite of existing work docs

**UUID:** `cop-doc-overwrite`
**Type:** fix
**Severity:** Medium
**Reported:** 2026-09-13
**Status:** `inbox`
**Source:** Copilot PR #40 review — `librarian/src/internal/application/work/doc/doc.go`

---

## Finding

`os.WriteFile` truncates an existing `{type}.md` without warning:

```go
if err := os.WriteFile(filePath, []byte(template(...)), 0o644); err != nil {
```

Repeating a `hawp work doc` or `hawp_work_doc` call for the same date and
work-item ID silently destroys the previously written document.

## Fix Plan

Option A (preferred): check for existence first; if the file exists, return the
existing path without overwriting it. The caller can then append content.

Option B: use `os.OpenFile` with `O_EXCL | O_CREATE` for exclusive creation;
return an error with the existing path so the caller knows where to write.

Implement Option A as the default MCP/CLI behavior. Document it in the tool
description so callers understand repeated calls are idempotent.

## Files

- `librarian/src/internal/application/work/doc/doc.go` — guard before WriteFile
- `librarian/src/internal/platform/mcp/server/tool_work.go` — propagate the
  "already exists, returning path" result cleanly
