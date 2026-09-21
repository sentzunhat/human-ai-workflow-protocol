# Add symmetric ampersand table-cell decoding

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

**UUID:** `7e73bcb8-ac73-4494-ae2a-1622d0f1558d`
**Type:** fix
**Reported:** 2026-09-23
**Source plan:** `f23c9369`

## Context

Backlog-row generation escapes ampersands as `&amp;`, but `table.Cells` only decodes `&#124;` and structural backslash escapes.

## Fix Work

Decode `&amp;` symmetrically in parsed cells and add a writer/parser round-trip test for `A & B`. Keep the decode narrow rather than applying a general HTML unescape that could alter literal entities.

## Verification

Run table/intake focused tests, full Go tests, vet, HAWP checks, formatting, and diff hygiene.

## Outcome

Completed implementation; see the reconciliation outcome below.
## Reconciliation Outcome

Done. Symmetric ampersand table-cell decoding is implemented and covered by
the work-table tests.
