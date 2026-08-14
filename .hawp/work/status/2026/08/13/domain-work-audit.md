# Domain Work Audit Checkpoint

## Work Item

Completed `c1d2e3f8`, the recursive audit of `domain/work`.

## Confirmed Architecture Finding

`domain/work` currently mixes pure HAWP record rules with filesystem reads,
dated-tree traversal, Markdown link parsing, path/existence checks, stderr
warnings, record moves, and content writes. It is the largest remaining
domain-to-infrastructure coupling in the librarian.

## Ordered Improvements

1. `c1d2e3f9`: extract a typed read-only work validation source and filesystem
   adapter for backlog, records, evidence, and dead-link inputs.
2. `c1d2e402`: after the read-only seam is stable, extract normalization scan
   and mutation operations separately.

This keeps HAWP validation semantics, UUID matching, legacy cutoff rules,
evidence containment, and dry-run safety intact while making the architecture
testable and extensible.

## Verification

`go test ./internal/domain/work ./internal/application/work ./internal/application/check ./internal/platform/cli` passes.
