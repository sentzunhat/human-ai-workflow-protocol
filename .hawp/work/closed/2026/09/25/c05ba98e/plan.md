# Fix nested canonical work-plan dead-link validation

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

**UUID:** `c05ba98e`
**Type:** fix
**Source plan:** `work/active/a78f1048/plan.md`
**Status:** done

## Context

Dead-link validation scans only files directly under `active/` and `parked/`, so
canonical `<uuid>/plan.md` records are skipped.

## Fix Work

Add scoped recursive collection for Markdown files under active and parked work roots,
preserving flat files and excluding archives. Add a nested broken-link regression.

## Verification

- Run focused dead-link tests for nested and flat plans.
- Run all Go tests and workflow validation.
- Record any unproven filesystem adapter behavior.
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

Done. Nested canonical work-plan dead-link handling is implemented and
covered by validation tests.
