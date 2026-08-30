# v0.0.23 compact render follow-up

type: benchmark
date: 2026-08-29
backlog_uuid: 387a37b7-e8e4-4c6b-9972-3c9712bc9774

## Scope

Verify whether a smaller context renderer reduces the remaining sparse-result
token negatives without widening the search/context contract.

## Code change under test

- `librarian/src/internal/application/context/format.go`
  - sparse markdown output skips the title header when result count is 5 or less
  - inline markdown provenance is reduced to `Ref: <source[:line[-line]]>`
  - chunk title and relevance remain in typed result/reference metadata and JSON
- `librarian/src/internal/application/context/format_test.go`
  - coverage added for sparse header skipping and compact provenance rendering

## Commands

```bash
cd librarian/src && go test ./internal/application/context ./internal/application/search
cd librarian/src && go run ./cmd/hawp search benchmark --token
```

## Direct evidence

- Focused Go packages passed:
  - `ok github.com/sentzunhat/hawp/librarian/src/internal/application/context`
  - `ok github.com/sentzunhat/hawp/librarian/src/internal/application/search`
- Token benchmark after the compact render change:
  - Intake process and investigation ordering: `1707 -> 1753` (`-46`, `-3%`)
  - Provider distribution and materialization: `1796 -> 1870` (`-74`, `-4%`)
  - Core HAWP protocol shape fields: `591 -> 613` (`-22`, `-4%`)
  - Total: `22866 -> 18169` (`4697`, `21%`)

## Interpretation

- The renderer is now simpler and the structured/markdown boundary is cleaner:
  markdown carries only inline source-line provenance while richer ranking/title
  data stays in the typed shapes.
- Sparse-result negatives were reduced but not eliminated.
- No evidence from this pass justifies a larger behavioral change inside the
  same bounded slice; the next step should be an explicit sparse passthrough
  policy if removing those negatives remains a release requirement.
