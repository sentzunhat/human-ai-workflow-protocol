# Preserve ampersands in backlog table round trips

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

**UUID:** `f23c9369-8ee7-4d86-b02d-60115ac4a83c`
**Type:** bug
**Reported:** 2026-09-23

## Intake Summary

Copilot found that `escapeTableCell` converts every ampersand to `&amp;`, while `table.Cells` decodes escaped pipes but not HTML entities. A title such as `A & B` therefore returns as `A &amp; B` after writing and parsing a backlog row.

## Current Context

Pipe and backslash escaping was added in the prior review cycle. The parser has legacy `&#124;` compatibility but no symmetric ampersand handling.

## Initial Analysis

**Directly verified:** `librarian/src/internal/domain/work/intake/intake.go` replaces `&` with `&amp;`; `librarian/src/internal/domain/work/table/table.go` does not decode `&amp;`.

**Inferred:** Work titles containing ampersands fail the promised round-trip contract and can be displayed or validated with altered content.

**Likely scope:** table codec, intake writer, and round-trip tests.

## Root Cause

Writer and parser implement asymmetric entity handling.

## Options

1. Stop ampersand escaping and keep only structural pipe/backslash escaping. This risks ambiguity with existing HTML-like content.
2. Decode `&amp;` symmetrically after structural parsing and add a round-trip test. This preserves the current writer contract.

## Recommended Fix

Keep the writer’s escaping and decode only the entity it emits, after escaped structural characters are processed. Add an ampersand title round-trip test and retain legacy `&#124;` behavior.

## Risk + Review Gate

**Risk:** medium — shared backlog serialization.
**Gate:** user explicitly requested implementation and sequential remediation.

## Verification

Pending linked fix item `7e73bcb8`.
## Reconciliation Outcome

Done. Ampersands survive backlog table round trips through the current cell
encoding/decoding path.
