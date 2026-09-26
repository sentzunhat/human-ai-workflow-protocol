# Fix installer regular-file destination validation

## Outcome

The PR-review finding was implemented and reconciled with the current branch.

## Verification

The focused regression coverage and repository-wide Go, HAWP, distribution,
formatting, and diff-hygiene checks passed before close.

## Close Checklist

- [x] Outcome recorded.
- [x] Verification evidence recorded or referenced.
- [x] Backlog row removed from active coordination.
- [x] Plan archived under the close date.

**UUID:** `d57a02e6-f79d-463f-906e-5d0b5fb36924`
**Type:** fix
**Reported:** 2026-09-23
**Source plan:** `b8dfd620`

## Context

The Makefile installer at `librarian/src/Makefile` guards symlinked `.hawp` paths but does not reject an existing directory or other non-regular destination at `.hawp/bin/hawp`.

## Fix Work

Update both pre- and pre-rename destination checks so an existing destination must be a regular file. Preserve the temporary sibling copy and final rename behavior. Add focused coverage that verifies a directory destination is rejected and does not receive the temporary binary.

## Verification

Run the installer-focused regression, `go test ./...`, `go vet ./...`, `go run ./cmd/hawp check`, formatting, and diff hygiene checks.

## Outcome

Done. Both pre-copy and pre-rename Makefile guards reject symlinked or
non-regular `.hawp/bin/hawp` destinations before `mv -f`.
