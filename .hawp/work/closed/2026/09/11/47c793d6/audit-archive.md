# CLI architecture continuation audit

Scope: `librarian/src/internal/platform/cli`, initially inspected 2026-09-05; model-parser follow-up 2026-09-07. This is a
bounded architecture review, not a repository-wide security certification.

## Implemented

Normalization safety follow-up: active-row and duplicate cleanup now have
separate files in `internal/domain/work`. Duplicate linking preserves all files
and supporting artifacts, adds reciprocal references for a unique archived
counterpart, and reports ambiguous matches for review. No duplicate removal.
Focused tests verify preview/apply agreement, preservation, and idempotency.

- Confirmed: command routing shared a 1,674-line file with unrelated command
  implementations. It now occupies 203 lines; feature handlers, help, and
  corpus assembly have separate files within the same private Go package.
- Confirmed: `strconv.ParseFloat` accepts `NaN`, which bypassed ordinary range
  comparisons for `--hybrid-ratio`. Non-finite values are now rejected before
  retrieval; black-box regression coverage includes NaN, infinities, malformed
  values, and out-of-range numbers.

## 2026-09-10 CLI input contracts (partial)

`parseProviderFlags` in `cli/init/command.go` and `cli/update/commands.go` was a
hand-rolled loop that silently ignored unknown flags, a trailing `--provider` with
no value, and extra positional arguments. Replaced with typed `FlagSet`-based parsers
(`parseInitArgs`, `parseUpdateSyncArgs`) for the `init` and `update sync` handlers.
22 table-driven cases added in `cli/init/args_test.go` and `cli/update/args_test.go`.
`RunFull` retains the legacy `parseProviderFlags` helper; its `--no-providers` flag
makes typed migration a separate follow-up.

Note: agent worktree for this slice branched from `main` instead of the feature branch.
Changes were re-applied manually to the decomposed `cli/init/` and `cli/update/`
subpackages (not the old monolithic `run.go`). Commit `35107e6`.

## Remaining Queue

Repository-port continuation: `domain/search.Index` now owns the retrieval
contract and embedding metadata. SQLite implements it with a compile-time
assertion; its old metadata name remains a compatibility alias. The application
accepts injected repositories using this contract. Vector-state read failures
now propagate instead of being silently treated as a vector-free index.
Fake-repository tests cover success, open failure, vector failure, query failure,
and closure of successfully opened repositories.

Current architecture is partial: `Embedder` and `LLMClient` are provider ports,
but provider implementations and factories remain in domain packages. Search's
default constructor still wires filesystem and SQLite adapters inside the
application package. CLI handlers are inbound adapters, with direct storage
access still present in indexing and usage. Follow-up extraction should preserve
the existing defaults while making composition explicit and testing failures.

1. Corpus correctness: read errors now abort collection with their original error
   preserved; corpus lifecycle derives from active/parked/closed folders, and
   supporting files keep unknown status. Paths use `filepath.Rel` and slash
   normalization. Tests cover the lifecycle cases and dangling-file read failures.
   Remaining: resolve exact backlog status and UUID metadata before ingestion,
   and move corpus collection behind an application-owned port. The existing
   ingest path only persists status when WorkUUID is supplied, so this change
   alone does not establish persisted backlog metadata parity. The separate
   `index build` enrichment path also needs a coordinated follow-up review.
2. CLI input contracts: several handlers silently accept missing flag values or
   invalid integer options. Define compatibility expectations and add table-driven
   tests before adopting shared parsing. Avoid a new CLI framework.
3. Storage boundaries: usage handlers directly open the usage database; corpus
   handlers directly walk files. Extract use cases when they reduce duplication
   or enable meaningful error-path tests; a file split alone does not invert
   these dependencies.
4. Search shaping: token capping and dedup orchestration remain in the CLI.
   Compare MCP and CLI contracts before extracting shared context policy, with
   sparse/dense benchmark evidence kept separate from structural claims.
5. Work scaffolding: resolved in this continuation. `work new` now renders
   against the actual Active Work header, including reordered and legacy columns.
   It inserts within the table before trailing prose, encodes title pipes and
   newlines, and refuses unsupported tables before writing a plan or backlog.
   Regression fixtures cover UUID-only, UUID-plus-legacy, reordered, and numeric
   ID headers, as well as refusal with no writes.

Search argument parsing now uses a pure typed parser built on Go `flag.FlagSet`.
It supports separate and equals-form values, existing `-v`, and rejects unknown
flags, missing values, invalid formats/numbers, and extra positional arguments.
The query-first command shape remains. Other command parsers remain in the queue.
`work normalize` now also uses typed flag parsing and rejects conflicting modes,
competing roots, and missing/empty path arguments before mutation. Its parser
tests run without filesystem writes. Local/CI installation now writes the
native build to `.hawp/bin/hawp` via temporary file plus rename, as confirmed
in `librarian/src/Makefile`; the earlier launcher/sibling layout is historical.
Existing CI runs the new tests through `go test ./...` and now also runs `go vet`.

## 2026-09-07 Verified Follow-Up

The existing domain model consolidation retains nullable folder context and
line ranges, with SQLite aliases and round-trip coverage. The embed parser
now keeps option values separate from text and returns parse errors.
`model pull` now accepts the documented trailing options as well as leading
options; 16 pure cases cover compatibility and rejected inputs. Full Go tests,
vet, distribution validation, kit parity, and connected MCP validation pass.
No model inference or cross-platform runtime claim follows from parser tests.

Request-to-intake reshaping is tracked separately in
[the reshape plan](../a3df8a9c/plan.md). Worker guidance uses existing tools;
a callable reshape implementation remains an explicit follow-up.

## Mutation Boundary Follow-Up

`kit normalize` no longer accepts both enabled apply/dry-run modes. Omitted
paths still discover defaults; explicit empty/whitespace paths are rejected by
work creation and both kit parsers. The prior no-op work-new validation callback
and incorrect FlagSet comment are removed. Parser regressions and handler-level
fixture snapshots verify rejection without file/directory changes. The source
README and command registry describe these distinctions.

## 2026-09-10 CLI input contracts (RunFull + usage boundary)

**Slice A — RunFull typed parser (done).** `RunFull` previously used
`parseProviderFlags` (a hand-rolled loop that silently ignored unknown flags,
a trailing `--provider` with no value, and extra positionals) and
`containsArg` to detect `--no-providers`. Both helpers are now removed.
`parseUpdateFullArgs` uses `flag.FlagSet` with `--provider` (repeatable),
`--no-providers`, and `--no-update-check`; it rejects unknown flags, trailing
`--provider`, and extra positionals before any side effect. Nine table-driven
cases (5 valid, 4 invalid) appended to `cli/update/args_test.go`. Commit
`8a0a790` on `agent/47c793d6-cli-validation`.

**Slice B — usage DB storage boundary (done).** `RunLog`, `RunReport`,
`RunClear`, and `RunTotals` previously called `usageinfra.Open(h.UsageDB)`
directly in the CLI layer. Extracted into `internal/application/usage` package
with four use-case functions (`RecentLog`, `GetReport`, `ClearLog`, `GetTotals`)
that each open the store, operate, and close, returning domain types only.
`commands.go` now calls through the application layer; it retains only
`usageinfra.LoadConfig`/`SaveConfig` for config-only operations (no store open).
`home()` helper stays in the CLI layer (path resolution is platform-layer
responsibility). 12 table-driven tests in `application/usage/service_test.go`
cover success, empty-store, and store-open-failure for each operation. `go test
./internal/application/usage/...` and `go test ./...` both pass; `go vet ./...`
clean. Branch `agent/47c793d6-storage-boundary`.

## Applied Guidance

- `.hawp/kit/instructions/clean-code-and-structure.md`: split mixed ownership,
  preserve bounded scope, validate the first structural edit.
- `.hawp/kit/standards/service-design/handler-responsibilities.md`: handlers
  translate and validate; services own policy.
- `.hawp/kit/standards/guidelines/security.md`: validate boundary input.
- `.hawp/kit/standards/nodejs/project-structure.md`: use responsibility-based
  organization where applicable. TypeScript framework and file conventions do
  not replace the existing Go module/package conventions.
- Go toolchain and architecture are defined by `librarian/src/go.mod` and the
  existing application/domain/infrastructure/platform packages. No dependencies
  or public contracts were introduced by this pass.
