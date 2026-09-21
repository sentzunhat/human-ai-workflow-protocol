# Make pipe-delimited backlog parsing escape-aware

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

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `c8c839ec-8847-4a7f-b407-9b1df7d3f114`
**Type:** bug
**Reported:** 2026-09-23

## Input (verbatim)

> Make pipe-delimited backlog parsing escape-aware

## Intake Summary

The Copilot review identifies a parser/writer contract break for backlog titles that
contain `|`. The intake writer permits such titles, but `table.Cells` uses a raw
`strings.Split`, so an escaped or literal pipe becomes a new Markdown cell and shifts
status, plan, and date fields.

## Current Context

The parser is shared by backlog validation and work intake. The fix must define the
escape contract narrowly, preserve ordinary pipe tables, and add round-trip coverage
through the writer/parser boundary.

## Initial Analysis

**Directly verified:**

- `librarian/src/internal/domain/work/table/table.go:8-17` splits every `|` without
  tracking escapes.
- `librarian/src/internal/domain/work/table/table_test.go` has only an ordinary row
  test.
- The review states `NewItemInput.BacklogRow` accepts a pipe-containing title.

**Inferred (not yet proven):**

- Markdown table cells can represent a literal pipe as `\\|`; a backslash-aware
  scanner can split only unescaped delimiters and unescape the cell payload.

**Likely scope:**

- `librarian/src/internal/domain/work/table/table.go`, its tests, and writer-side
  escaping needed to make the contract symmetric.

## Root Cause

The table parser has no lexical state for escaped delimiters, while work-item title
validation and row generation allow a title that requires escaping.

## Options Considered

1. Reject `|` in titles. This would contradict the existing title-safety contract and
   unnecessarily narrow valid data.
2. Make the parser and row writer escape-aware, with `\\|` representing a literal pipe.
   This is the recommended fix because it preserves valid titles and round-trips.

## Recommended Fix

Implement a small escape-aware cell scanner, unescape supported escaped characters at
the cell boundary, and ensure `BacklogRow` emits escaped pipes. Add parser and
round-trip tests for pipe-containing titles plus a regression check for ordinary rows.

## Risk + Review Gate

**Risk:** medium (shared parser contract)
**Gate:** user authorized implementation in the request; no merge authorization.

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/active/c8c839ec/plan.md

## Verification

- Run focused table/intake tests, including a title containing `|`.
- Run all Go tests and work validation.
- Inspect serialized rows to confirm column counts remain stable.

## Next Step

- [x] Investigation recorded above
- [x] Write or update the plan file
- [x] Create the corresponding fix work item and implement sequentially
## Reconciliation Outcome

Done. Backlog table parsing preserves escaped pipe cells through the read and
write path; the full Go test suite passes.
