---
work-item: b6c4e8a2
type: task
title: "Audit and align librarian architecture with domain port/adapter boundaries"
status: in-progress
owner: unassigned
created: 2026-08-10
updated: 2026-08-11
---

# Architecture Port/Adapter Audit

## Mission

Compare the Go librarian architecture with the repository-local area and
port/adapter patterns used by the Mictlan and Tekit reference codebases, then
apply a bounded search/context refactor that improves seams without a broad
folder migration or public CLI break.

## Direct Findings

- The librarian already separates domain, application, infrastructure, and
  platform packages, and already has embedding/LLM interfaces.
- Search behavior is duplicated: the CLI performs lexical search, hybrid
  ranking, result conversion, and context preparation separately from the
  application search service.
- The application search service imports SQLite directly and exchanges
  untyped `map[string]interface{}` rows with the infrastructure layer.
- Mictlan/Tekit provide a useful reference pattern: self-contained domain
  areas, ports beside capabilities, infrastructure adapters beside external
  systems, and explicit composition at the application boundary.

## Bounded Implementation Slice

- Add typed search-domain port and candidate contracts.
- Add a SQLite search adapter implementing the port.
- Route the CLI and RAG retrieval through one application search service.
- Preserve the existing SQLite schema, command names, output modes, and
  lexical-only fallback.

## Deferred Follow-ups

- Full dense-retrieval candidate union so semantic search can recover lexical
  misses.
- Extract remaining search/index command handlers from the general CLI router
  into the search platform capability group.
- Composition-root construction for all providers and repositories.
- Broader package migration for non-search capabilities.
- Large-scale index/storage changes.

## Progress

- Added `domain/search/index/port.go` with typed index candidates and an
  `IndexPort` kept with the search-index capability group.
- Added `infrastructure/sqlite/search/adapter.go` to translate SQLite rows into
  the typed domain contract within the SQLite search capability group.
- Refactored application search to accept the port and carry real relevance,
  folder context, and work metadata into context formatting.
- Routed the CLI through the application search service, removing its second
  lexical/hybrid implementation.
- Extracted context JSON/reference rendering and LLM-reshape fallback into
  `platform/cli/search_output.go`; `run.go` is no longer responsible for that
  output concern.

## Verification

- Targeted tests pass for application search, domain search, SQLite, and CLI.
- `make build` passes with the default `CGO_ENABLED=0` release build.
- `git diff --check` passes.
- Full repository tests remain unproven because the model/integration-heavy
  suite did not complete promptly in this environment and was stopped.

## Ranked Next Work

1. Add semantic candidate retrieval and union it with lexical candidates so
   paraphrase queries can be found.
2. Extract the remaining search/index CLI handlers into a search platform
   capability file.
3. Move provider/repository construction to explicit composition roots and
   remove remaining application/domain imports of concrete infrastructure.
4. Add end-to-end acceptance fixtures for index, search, context, provenance,
   token budgets, and model-missing fallback.

## Folder Category Audit: `platform/cli`

### Confirmed Structure

| File/group | Capability | Finding | Next action |
| --- | --- | --- | --- |
| `registry.go`, `commands.go` | command discovery | coherent registry/output pair | keep together |
| `search_output.go` | context presentation | correctly isolated from routing | keep together; add output tests later |
| `benchmark.go` | search benchmark presentation | separate file, but still coupled to raw SQLite rows | port after semantic retrieval |
| `run.go` lines 44-1103 | command routing plus maintenance, update, work, index, model, and search handlers | mixed capability groups in one file | split by capability into command files |
| `*_test.go` | CLI behavior | routing tests are centralized and useful | add capability-local tests as files split |

### Category Decision

`platform/cli` should be a thin composition and transport category. Each
capability gets a colocated command file, while application services own the
use case and domain ports/adapters own the capability boundary. The next
mechanical split is `search_commands.go` for index/embed/search routing, then
`update_commands.go` and `work_commands.go`; no behavior rewrite is needed for
the file moves.

### Category Matrix

| Category | Status | Next compoundable slice |
| --- | --- | --- |
| `platform/cli` | audited | split remaining command families |
| `domain/search` | audited | semantic candidate union |
| `application/search` | audited | keep orchestration typed; add retrieval tests |
| `infrastructure/sqlite` | partially audited | isolate vector/search repository methods |
| `application/context` | partially audited | separate retrieval, formatting, and reshape acceptance seams |
| remaining `domain/*` | pending | audit after search/context stabilizes |
| composition roots | pending | add explicit provider/repository construction |

## Recursive Capability Audit

The audit was extended one level below each category. The key result is that
folder names alone do not yet guarantee the intended dependency direction.
The most important issue is not file size: several packages under `domain/*`
perform filesystem, Markdown, and repository I/O directly.

| Nested group | Direct evidence | Assessment | Compoundable next slice |
| --- | --- | --- | --- |
| `domain/context` | `context/kit.go` and `context/work.go` import concrete `infrastructure/markdown` and `infrastructure/repo`; both walk/read repository files | confirmed domain-to-infrastructure coupling; context enrichment is currently an I/O workflow, not pure domain logic | introduce a context corpus/source port and a filesystem-backed context adapter; keep enrichment rules in the domain context group |
| `domain/kit` | `kit/normalize.go` and `kit/validate.go` import Markdown/repository infrastructure | confirmed parsing/validation and filesystem concerns are mixed | separate kit content input behind a capability-specific port before moving behavior |
| `domain/work` | `work/deadlinks.go`, `work/completeness.go`, and normalization rules import repository/Markdown infrastructure | confirmed work rules are coupled to repository layout and files | extract work-source and link-resolution ports; preserve existing validation APIs through wrappers |
| `application/index` | ingest/embed services open SQLite directly and construct storage records | application owns part of persistence composition | move document/chunk/embedding persistence behind capability-local contracts, starting with index ingestion |
| `application/context` | retrieval, formatting, deduplication, encryption, and reshape are separate files, but the package still spans several use-case stages | healthy internal grouping, incomplete seam boundaries | add explicit retrieval/format/reshape contracts and acceptance fixtures; avoid arbitrary file splitting |
| `infrastructure/sqlite` | root SQLite package still contains broad index operations; search adapter is now isolated, while benchmark still consumes raw maps | partial capability extraction | group document/chunk/vector/search persistence by capability and replace raw result maps |
| `platform/cli` nested command groups | output is isolated, but `run.go` still routes maintenance, work, update, index, model, and search | transport category remains too broad | split command families mechanically, starting with search/index |

### Dependency-Direction Rule

For the remaining migrations, a domain capability may define a port beside
its capability, and an infrastructure capability may implement that port. A
domain or application package should not import a concrete infrastructure
package merely to read files, walk folders, open SQLite, or invoke a provider.
Composition roots may wire those concrete adapters. This preserves the user
decision that ports and adapters stay together by capability rather than in
generic top-level folders.

### Recursive Work Order

1. Extract the `domain/context` repository-reading seam with compatibility
   wrappers and focused tests.
2. Apply the same pattern to `domain/work` and `domain/kit`, keeping parsing
   and validation capability-local.
3. Split `application/index` persistence contracts into document, chunk, and
   embedding capability groups.
4. Split remaining CLI command families after the application seams are
   stable.
5. Add explicit composition roots and end-to-end acceptance fixtures.

## Completed Child Slice: `c1d2e3f5`

- Added `domain/context/source/port.go` with typed corpus inputs and a source
  contract.
- Added `infrastructure/filesystem/context/adapter.go` for repository
  acquisition and backlog loading.
- Removed concrete filesystem/Markdown acquisition from `domain/context` and
  wired the default adapter through `application/index`.
- Preserved the existing README, work-folder, and read-error policies with
  focused tests.
- Targeted tests, CLI/search tests, CGO-free build, and diff checks pass.

## Completed Child Slice: `c1d2e3f7`

- Added a typed kit workspace contract under `domain/kit/source`.
- Added a filesystem kit adapter for discovery, Markdown link extraction, and
  controlled mutation.
- Removed concrete filesystem/Markdown/repository access from `domain/kit`.
- Preserved application-owned dry-run/apply reporting and dirty-worktree
  protection with focused pure-rule and adapter tests.

## Child Work Items

Each recursive audit now has a linked implementation item. Audits produce
confirmed findings, boundary proposals, and verification requirements; the
paired fix item turns those findings into a bounded code change.

| Audit | Improvement item | Order |
| --- | --- | --- |
| `c1d2e3f4` domain context | `c1d2e3f5` context corpus/source boundary | complete; next |
| `c1d2e3f6` domain kit | `c1d2e3f7` kit content boundary | complete; next |
| `c1d2e3f8` domain work | `c1d2e3f9` work source/link boundary | complete; next |
| `c1d2e3fa` application index | `c1d2e3fb` typed index persistence | after domain seams |
| `c1d2e3fc` application context | `c1d2e3fd` retrieval/format/reshape seams | after search baseline |
| `c1d2e3fe` SQLite infrastructure | `c1d2e3ff` capability-local persistence | after index audit |
| `c1d2e400` CLI capabilities | `c1d2e401` command-family splits | after application seams |

The work audit identified a separate normalization scan/mutation boundary;
`c1d2e402` follows `c1d2e3f9` rather than expanding its validation scope.
