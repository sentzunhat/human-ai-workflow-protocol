# Fix escaped pipe backlog table round-trip

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

**UUID:** `e3e2b3f0`
**Type:** fix
**Source plan:** `work/active/c8c839ec/plan.md`
**Status:** done

## Context

Backlog writers accept titles containing `|`, but raw table splitting changes the
column count. The shared table contract must represent literal pipes as escaped
delimiters.

## Fix Work

Implement escape-aware table-cell parsing and symmetric writer escaping, then add
ordinary-row and pipe-containing-title round-trip tests.

## Verification

- Run focused table and intake tests.
- Run all Go tests and work validation.
- Confirm serialized rows retain their original column count.
+## Outcome

Implemented the planned fix and added the scoped regression coverage where applicable.

## Direct Verification

- Full `go test ./...`: PASS.
- `go vet ./...`: PASS.
- Tracked Go files checked with `gofmt -l`: no output.
- `git diff --check`: PASS.
- `hawp distribution validate`: PASS.
- `hawp work validate`: PASS with the repository's pre-existing verification-clarity warning.
- `hawp check`: all three validations passed.

## Review State

The corresponding Copilot finding remains to be answered and resolved after the commit
is pushed and the live PR thread is re-read. No merge authorization is granted.
## Reconciliation Outcome

Done. Escaped-pipe backlog round trips are implemented and covered by parser
and intake regressions.
