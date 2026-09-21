# Preserve byte offsets when masking Unicode fenced content

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

**UUID:** `f793efc9-c59c-4ddc-8041-05f72bab7396`
**Type:** bug
**Reported:** 2026-09-23

## Intake Summary

Copilot reports that `librarian/src/internal/infrastructure/markdown/markdown.go` converts fenced content to runes while `librarian/src/internal/domain/kit/normalize.go` applies extracted offsets to original byte strings. Unicode within an earlier fence can therefore shift a later link update and corrupt `hawp kit normalize --apply`.

## Current Context

`markdown.BlankFences` is the production `KitSource.BlankFences` adapter. `PlanLinkUpdates` derives byte offsets from the masked string, and `ApplyLinkUpdates` slices the unmasked original content at those offsets.

## Initial Analysis

**Directly verified:** `BlankFences` uses `[]rune`, while `ExtractLinks` returns regexp byte offsets and `ApplyLinkUpdates` uses Go string byte slicing.

**Inferred:** A multibyte character in a masked fence shortens the masked byte sequence, causing every following extracted link offset to point before its original byte position.

**Likely scope:** `librarian/src/internal/infrastructure/markdown/markdown.go` and kit-normalization regression coverage.

## Root Cause

The masking representation does not preserve the byte-coordinate system used by the link extractor and mutation planner.

## Options

1. Mask individual bytes, preserving newline bytes and replacing every other byte with ASCII space. This keeps all later byte offsets stable.
2. Convert all consumers to rune offsets. This broadens the coordinate contract and risks unrelated parsing behavior.

## Recommended Fix

Replace the rune conversion with byte-wise masking and add a `PlanLinkUpdates` plus apply regression where a Unicode fence precedes a renamed link.

## Risk + Review Gate

**Risk:** high — incorrect offsets can corrupt user-authored kit Markdown.
**Gate:** user explicitly authorized sequential remediation in this PR loop.

## Verification

Pending linked fix item `84892d9d`.
## Reconciliation Outcome

Done. Unicode fenced-content masking preserves byte offsets for link
analysis.
