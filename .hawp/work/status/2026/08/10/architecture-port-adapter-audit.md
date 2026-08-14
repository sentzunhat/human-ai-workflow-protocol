# Architecture Audit Checkpoint

## Scope

Compared the librarian Go search/context structure with the local Mictlan and
Tekit area patterns: self-contained domain slices, ports beside capabilities,
infrastructure adapters beside external systems, and explicit application
composition.

## What Changed

- Added a typed `domain/search/index/IndexPort` and candidate contract.
- Added an `infrastructure/sqlite/search/Adapter` that owns storage-row mapping
  inside the SQLite search capability group.
- Routed CLI search through the application search service shared by RAG.
- Preserved existing commands, SQLite schema, output modes, and lexical-only
  fallback behavior.
- Preserved real relevance and indexed folder/work metadata instead of
  assigning a fixed context score in the CLI.
- Extracted context output and reshape fallback into
  `platform/cli/search_output.go`, reducing the general CLI router by that
  presentation responsibility.

## Direct Verification

- Targeted Go tests pass for search, domain search, SQLite, and CLI packages.
- `make build` passes with `CGO_ENABLED=0`.
- `git diff --check` passes.
- No Mictlan or Tekit files were modified.

## Remaining Work

- Semantic retrieval still reranks lexical candidates; it does not yet recover
  lexical misses.
- Full provider/repository composition is not part of this first slice.
- Full repository tests remain unproven because the integration-heavy suite did
  not complete promptly and was stopped.

## Next Order

1. Semantic candidate retrieval plus lexical/semantic result union.
2. Extract remaining search/index handlers from `run.go`.
3. Establish explicit composition roots for providers and repositories.
4. Add end-to-end context acceptance fixtures.

## Folder Category Audit: `platform/cli`

- `registry.go` and `commands.go` are a coherent command-discovery group.
- `search_output.go` is now a coherent context-presentation group.
- `benchmark.go` is separate but still reads raw SQLite result maps.
- `run.go` remains mixed: routing, maintenance, update, work, indexing,
  model, and search command handlers share one file.

Decision: keep the CLI as a thin transport/composition category and split
`run.go` by capability next. The first split should be search/index command
routing, followed by update and work command routing. No behavior rewrite is
needed for those moves.

## Category Matrix

| Category | Status | Next slice |
| --- | --- | --- |
| `platform/cli` | audited | split command families |
| `domain/search` | audited | semantic candidate union |
| `application/search` | audited | typed retrieval tests |
| `infrastructure/sqlite` | partial | isolate vector/search repository methods |
| `application/context` | partial | retrieval/format/reshape seams |
| remaining domains | pending | audit after search/context |
| composition roots | pending | explicit provider/repository wiring |

## Deeper Folder Audit

The recursive pass went inside the category folders and found a real
dependency-direction issue, not just oversized files:

- `domain/context/kit.go` and `domain/context/work.go` directly import the
  concrete Markdown and repository infrastructure packages while walking and
  reading files.
- `domain/kit/*` and several `domain/work/*` files do the same for kit
  validation, normalization, completeness, and dead-link rules.
- `application/index/*` opens SQLite and constructs persistence records
  directly; `application/context` is internally grouped but still combines
  retrieval, formatting, deduplication, encryption, and reshape stages.
- `infrastructure/sqlite` is partly grouped by capability after the search
  adapter change, but benchmark and older index paths still expose raw maps.
- `platform/cli` has a good registry/output grouping, but `run.go` remains a
  broad command router across unrelated capabilities.

## Deeper Category Matrix

| Nested category | Status | Next compoundable work item |
| --- | --- | --- |
| `domain/context` | finding confirmed | extract a context corpus/source port and filesystem adapter |
| `domain/kit` | finding confirmed | isolate content input from pure normalization/validation |
| `domain/work` | finding confirmed | isolate repository/link resolution from work rules |
| `application/index` | partial | introduce typed document/chunk/embedding persistence seams |
| `application/context` | partial | formalize retrieval, format, and reshape boundaries |
| `infrastructure/sqlite` | partial | group persistence by capability and remove raw result maps |
| `platform/cli` | partial | split search/index, then update and work command families |

## Architecture Decision

Do not split files by arbitrary line count. Split when a folder contains a
distinct capability, external boundary, or composition responsibility. Keep
each port beside the capability that owns it and each adapter beside the
infrastructure capability that implements it; do not create generic global
`ports/` or `adapters/` buckets.

## Ordered Follow-Up

The next implementation slice should be `domain/context`, because it has the
clearest bounded seam and directly affects search context quality. Preserve
the existing `EnrichKit` and `EnrichWork` entry points while moving concrete
file access behind capability-local contracts. Then apply the same pattern to
`domain/work`, `domain/kit`, and index persistence.

## Audit-to-Improvement Work Items

The recursive audit is now decomposed into paired backlog items. Each audit
must produce confirmed findings, a smallest safe boundary, and tests/evidence
for its linked improvement item.

| Audit item | Follow-up improvement | Status |
| --- | --- | --- |
| `c1d2e3f4` domain context | `c1d2e3f5` extract corpus/source boundary | plan-ready |
| `c1d2e3f6` domain kit | `c1d2e3f7` isolate content input | plan-ready |
| `c1d2e3f8` domain work | `c1d2e3f9` extract source/link boundaries | plan-ready |
| `c1d2e3fa` application index | `c1d2e3fb` typed persistence contracts | plan-ready |
| `c1d2e3fc` application context | `c1d2e3fd` separate context seams | plan-ready |
| `c1d2e3fe` SQLite infrastructure | `c1d2e3ff` group persistence capabilities | plan-ready |
| `c1d2e400` CLI capabilities | `c1d2e401` split command routing | plan-ready |

Execution should begin with `c1d2e3f4`, then its follow-up `c1d2e3f5`; the
remaining pairs can be audited in parallel only if they do not overlap files.

## Current Work Item: `c1d2e3f4`

The domain context audit is complete for the inspected files. Confirmed:

- `domain/context` owns the correct `Document` and enrichment rules, but also
  performs Markdown discovery, file reads, path conversion, and backlog file
  loading.
- Kit and work enrichment are coupled to concrete infrastructure and to the
  work domain parser.
- Kit read failures are skipped and work read failures are ignored; this
  behavior must be preserved explicitly or changed through a separate
  decision.
- `EnrichWork` intentionally scans a fixed allow-list of work folders, so the
  ignored-folder policy needs a focused test.

Verification passed: `go test ./internal/domain/context ./internal/application/index`.

The follow-up `c1d2e3f5` is implementation-ready: add a context-specific
source port, a filesystem context adapter, application-level composition, and
compatibility wrappers while preserving the observed policies.
