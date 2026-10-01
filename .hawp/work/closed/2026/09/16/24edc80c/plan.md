# Wrap SQLite "no such table" as IndexNotFoundError in search infra

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `24edc80c-0de6-4cb4-af04-5b4663a4323a`
**Type:** fix
**Reported:** 2026-09-16

---

## Input (verbatim)

> Wrap SQLite "no such table" as IndexNotFoundError in search infra

## Intake Summary

SQLite returns a "no such table" error when the database file exists but the schema was never applied (the file was created by `os.Create` but the index build step did not run). This error was not wrapped as `IndexNotFoundError`, so callers received a generic error and could not distinguish a missing index from a query failure.

## Current Context

The `IndexNotFoundError` sentinel existed in `search/service.go` for the case where the database file was absent, but the "no such table" path (schema not applied) fell through to a generic `fmt.Errorf` wrap.

A string heuristic in `runWorkIntake` (MCP layer) duplicated detection logic.

## Initial Analysis

**Directly verified:**

- `service.Execute()` called `db.HasVectors()` and `db.QueryChunksLexical()` without wrapping schema errors.
- The string heuristic in `runWorkIntake` was redundant and fragile.

**Inferred (not yet proven):**

- All callers that handle `IndexNotFoundError` would benefit from consistent wrapping.

**Likely scope:**

- `librarian/src/internal/application/search/service.go` — two call sites.
- `librarian/src/internal/application/mcp/tools.go` — remove string heuristic.

## Risk + Review Gate

**Risk:** low
**Gate:** auto-implement on low

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/closed/2026/09/16/24edc80c/plan.md

## Verification

- `go test ./...` passes; `go vet ./...` clean.
- `isSchemaError()` helper added; `IndexNotFoundError` is now returned from `Execute()` at both `HasVectors()` and `QueryChunksLexical()` call sites when "no such table" is encountered.
- String heuristic removed from `runWorkIntake`; callers receive the typed sentinel.

## Outcome

Added `isSchemaError()` helper in `search/service.go` that checks for the "no such table" substring. Wrapped the result in `IndexNotFoundError` at two call sites in `Execute()` (`HasVectors` and `QueryChunksLexical`). Removed the fragile string heuristic from `runWorkIntake` in the MCP layer; all schema-missing detection is now centralised in the service.

## Close Checklist

- [x] Schema error wrapped as `IndexNotFoundError` at both call sites in `Execute()`.
- [x] String heuristic removed from MCP layer.
- [x] `go test ./...` passes; `go vet ./...` clean.
- [x] Plan moved to `closed/2026/09/16/24edc80c/`.
- [x] BACKLOG.md updated.
