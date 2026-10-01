# Fix quality workflow all-zero base SHA handling

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

**UUID:** `2c588680`
**Type:** fix
**Source plan:** `work/active/06fd7e9b/plan.md`
**Status:** done

## Context

PR #41's Copilot finding reports that non-PR pushes use the all-zero
`github.event.before` sentinel during branch creation. The quality workflow must
retain diff hygiene without passing that sentinel to Git.

## Fix Work

Update `.github/workflows/quality.yml` to replace an all-zero base with the first
root commit before `git diff --check`. Keep the pull-request merge-base branch and
ordinary push behavior unchanged.

## Verification

- Inspect the three event paths.
- Run workflow/configuration validation and `git diff --check`.
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

## Outcome

Done. The workflow contains the explicit all-zero SHA fallback and keeps the
pull-request merge-base path unchanged. Workflow inspection and repository
diff-hygiene validation passed.

## Review State

The corresponding Copilot finding remains to be answered and resolved after the commit
is pushed and the live PR thread is re-read. No merge authorization is granted.
