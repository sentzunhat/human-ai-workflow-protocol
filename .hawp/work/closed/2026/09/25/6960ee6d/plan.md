# Fix UUID-folder closed-plan reconciliation

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

**UUID:** `6960ee6d`
**Type:** fix
**Source plan:** `work/active/e85e0a84/plan.md`
**Status:** done

## Context

The generated update helper currently reduces a closed link to `basename`, so
`closed/YYYY/MM/DD/<uuid>/plan.md` cannot reconcile its corresponding
`active/<uuid>/plan.md`.

## Fix Work

Change `distribution/sources/update/script-core.md` to preserve the UUID folder for
canonical plans while retaining flat legacy and ID-fallback paths. Regenerate all
provider artifacts with `hawp distribution sync`.

## Verification

- Exercise canonical UUID-folder and flat legacy fixtures.
- Run distribution validation and generated-output checks.
- Record the final commit and review-thread response after implementation.
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

Done. UUID-folder closed-plan reconciliation is implemented and covered by
the work normalization/validation test suite.
