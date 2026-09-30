# Fix non-atomic writes in ApplyDuplicateLinks plan files

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `a818b914`
**Type:** fix
**Reported:** 2026-09-16

---

## Input (verbatim)

> Fix non-atomic writes in ApplyDuplicateLinks plan files

## Intake Summary

`ApplyDuplicateLinks` wrote plan files using `os.WriteFile`, which is not atomic. An interrupted write could leave a truncated or partially-written plan file.

## Current Context

`os.WriteFile` on the plan path creates or truncates the file in-place. If the process is interrupted mid-write, the file is left in a corrupt state with no recovery path.

## Initial Analysis

**Directly verified:**

- `ApplyDuplicateLinks` used `os.WriteFile` to write each plan file.

**Inferred (not yet proven):**

- Atomic write via temp-file + rename eliminates the partial-write risk.

**Likely scope:**

- `librarian/src/internal/domain/work/apply_duplicate_links.go` — plan file write path.

## Risk + Review Gate

**Risk:** low
**Gate:** auto-implement on low

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/closed/2026/09/16/a818b914/plan.md

## Verification

- `go test ./...` passes; `go vet ./...` clean.
- Plan file writes now use `os.CreateTemp` + `os.Rename` for atomic replacement.

## Outcome

Changed the plan file write in `ApplyDuplicateLinks` to use a temp file in the same directory followed by `os.Rename`. The rename is atomic on POSIX systems, so an interrupted write leaves the original plan file intact.

## Close Checklist

- [x] Plan file writes are atomic (CreateTemp → Rename).
- [x] `go test ./...` passes; `go vet ./...` clean.
- [x] Plan moved to `closed/2026/09/16/a818b914/`.
- [x] BACKLOG.md updated.
