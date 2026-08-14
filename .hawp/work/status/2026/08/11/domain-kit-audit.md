# Domain Kit Audit Checkpoint

## Work Item

Completed `c1d2e3f6`, the recursive audit of the domain kit capability.

## Confirmed Findings

- `domain/kit` contains pure filename, required-file, and link rules, but also
  performs directory traversal, file reads, existence checks, renames, and
  writes.
- Markdown parsing and repository path helpers are concrete infrastructure
  dependencies of both validation and normalization.
- `application/kit` correctly owns reporting, dry-run/apply mode, and the
  dirty-worktree safety gate, but delegates path-based I/O to the domain.
- Existing tests are filesystem-driven; no in-memory rule tests or focused
  adapter mutation tests exist.

## Next Implementation

`c1d2e3f7` should introduce a kit-specific workspace contract under
`domain/kit/source` and a filesystem adapter under
`infrastructure/filesystem/kit`. It should migrate read-only validation and
planning first, then rename/write mutation, while preserving all CLI behavior
and safety checks.

## Verification

`go test ./internal/domain/kit ./internal/application/kit ./internal/application/check ./internal/platform/cli` passes.
