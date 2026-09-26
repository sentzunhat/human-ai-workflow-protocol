# Fix missing SQLite index preflight

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

**UUID:** `c3371bf5`
**Type:** fix
**Source plan:** `work/active/80c677fb/plan.md`
**Status:** done

## Context

`hawp search embed` opens SQLite before checking existence, allowing a missing
database to be created and collapsing permission/corruption errors into a successful
missing-index message.

## Fix Work

Check the resolved index path before `sqlite.Open`; provide first-run guidance only
for a confirmed missing path and propagate other open errors with context. Preserve
successful embedding and close handling.

## Outcome

Done. The missing-index preflight and successful-open path are implemented;
the command closes the database on all completed preflight paths and focused
search/embed validation passed.

## Verification

- Cover missing, usable, and failing-open paths.
- Run focused search-embed tests, all Go tests, and vet.
- Confirm missing-index preflight creates no SQLite file.
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
