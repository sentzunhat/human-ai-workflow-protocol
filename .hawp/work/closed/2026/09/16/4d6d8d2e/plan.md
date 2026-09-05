# Deduplicate idSet in duplicate_links.go

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `4d6d8d2e`
**Type:** fix
**Reported:** 2026-09-16

---

## Input (verbatim)

> Deduplicate idSet in duplicate_links.go

## Intake Summary

`duplicate_links.go` defined a private `idSet` type that duplicated an equivalent type in the `domain/work` package. Two separate implementations diverge over time and make the code harder to follow.

## Current Context

Both `domain/work` and the application layer (`duplicate_links.go`) maintained a `map[string]struct{}`-based set type for work item IDs. The domain package's version was private (`idSet`); the application file's version was also private.

## Initial Analysis

**Directly verified:**

- `duplicate_links.go` contained a local private `idSet` definition.
- `domain/work` had an equivalent but separate private type.

**Inferred (not yet proven):**

- Exporting `IDSet` from `domain/work` and removing the local copy unifies the implementation.

**Likely scope:**

- `librarian/src/internal/domain/work/` — export `IDSet`.
- `librarian/src/internal/application/work/duplicate_links.go` — remove local `idSet`, use `domainwork.IDSet`.

## Risk + Review Gate

**Risk:** low
**Gate:** auto-implement on low

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/closed/2026/09/16/4d6d8d2e/plan.md

## Verification

- `go test ./...` passes; `go vet ./...` clean.
- `IDSet` exported from `domain/work`; local `idSet` in `duplicate_links.go` removed; all callers use `domainwork.IDSet`.

## Outcome

Exported `IDSet` (capitalised) from the `domain/work` package. Removed the private `idSet` definition from `duplicate_links.go` and updated that file to use `domainwork.IDSet`. There is now a single canonical implementation.

## Close Checklist

- [x] `IDSet` exported from `domain/work`.
- [x] Duplicate private `idSet` removed from `duplicate_links.go`.
- [x] `go test ./...` passes; `go vet ./...` clean.
- [x] Plan moved to `closed/2026/09/16/4d6d8d2e/`.
- [x] BACKLOG.md updated.
