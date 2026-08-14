# Context Source Boundary Checkpoint

## Work Item

Completed `c1d2e3f5`, the implementation follow-up from the domain context
audit.

## What Changed

- Added a typed context corpus/source contract under
  `librarian/go/internal/domain/context/source`.
- Moved kit/work enrichment to typed corpus inputs without concrete
  filesystem or Markdown imports in `domain/context`.
- Added a filesystem context adapter for discovery, reads, relative paths, and
  backlog acquisition.
- Wired adapter composition into the application index build service.
- Added in-memory enrichment tests and filesystem adapter policy tests.

## Preserved Behavior

- Kit README files remain included.
- Work README files remain excluded.
- Only the existing work-folder allow-list is indexed.
- Read failures retain the previous skip behavior.
- Existing build, search, and CLI behavior remains covered by targeted tests.

## Verification

- Domain context, filesystem context adapter, and application index tests pass.
- CLI and application search tests pass.
- CGO-free `make build` passes from `librarian/go`.
- `git diff --check` passes.

## Next Work

The next compoundable item is `c1d2e3f6`, the recursive audit of the domain
kit capability, followed by `c1d2e3f7` if the audit confirms a bounded seam.
