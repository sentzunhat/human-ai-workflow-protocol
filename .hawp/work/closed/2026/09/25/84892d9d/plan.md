# Make fenced-code masking byte-length preserving

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

**UUID:** `84892d9d-c38e-4ac2-af0b-b900dbe0642a`
**Type:** fix
**Reported:** 2026-09-23
**Source plan:** `f793efc9`

## Context

`librarian/src/internal/domain/kit/normalize.go` applies byte offsets from masked Markdown to the original string. The current `librarian/src/internal/infrastructure/markdown/markdown.go` fence masking uses runes and changes byte length when fences contain Unicode.

## Fix Work

Mask fenced content byte-for-byte, retaining newline bytes. Add a regression that normalizes a renamed target after Unicode fenced content and proves only the intended live link changes.

## Verification

Run focused markdown and kit-normalization tests, then `go test ./...`, `go vet ./...`, HAWP validation, formatting, and diff-hygiene checks.

## Outcome

Completed implementation; see the reconciliation outcome below.
## Reconciliation Outcome

Done. Fenced-code masking is byte-length preserving and covered by the
Markdown tests.
