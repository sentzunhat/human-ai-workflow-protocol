# Status Report

## Intent

Capture the current checkpoint for the `feature/v0.0.23-refactor-v010` branch so
the next cleanup pass can continue the architecture audit, shaping work, and
benchmark evidence without losing direct proof or widening scope.

## Current State

- Branch: `feature/v0.0.23-refactor-v010`
- Base: `main` at `5b693d8`
- Current local commit on branch: `082e8fc`
- Working tree is not clean:
  - `librarian/src/internal/application/context/format.go`
  - `librarian/src/internal/application/context/format_test.go`
  - `.hawp/work/evidence/2026/08/29/`

## What Was Inspected

- HAWP active work items and backlog:
  - `.hawp/work/BACKLOG.md`
  - `.hawp/work/active/387a37b7-v0-0-23-refactor-search-pipeline-and-layout-unification.md`
  - `.hawp/work/active/c804eec0-v0-0-23-benchmark-evidence-and-release-artifact-refresh.md`
  - `.hawp/work/active/5957aaf4-v0-0-23-follow-up-architecture-audit-and-simplification-queu.md`
- Refactor/code paths:
  - `librarian/src/internal/application/search/service.go`
  - `librarian/src/internal/application/context/format.go`
  - `librarian/src/internal/infrastructure/filesystem/hawp_project.go`
  - `librarian/src/internal/infrastructure/filesystem/readme_generator.go`
- Verification and benchmark surfaces:
  - `./.hawp/bin/hawp search index`
  - `./.hawp/bin/hawp search embed --backend ollama --model nomic-embed-text`
  - `./.hawp/bin/hawp search benchmark`
  - `./.hawp/bin/hawp search benchmark --token`
  - `/tmp/hawp-v0023 search benchmark --token`
  - SQLite checks on `.hawp/db/index.sqlite`

## What Changed

- Created a non-merge patch branch and split work into three HAWP lanes:
  - main refactor lane
  - benchmark/evidence lane
  - follow-up read-only audit lane
- Committed the first architecture slice as `082e8fc`:
  - aligned reusable `.hawp` runtime-path helpers with the active
    `.hawp/db` and `.hawp/config` contract
  - preserved richer search provenance through the typed formatting path
  - added focused tests for runtime layout and context/reference provenance
- Rebuilt the branch-local search index and restored it to a usable state.
- Ran a full local embed pass against the rebuilt index with Ollama
  `nomic-embed-text`.
- Added benchmark evidence artifacts under `.hawp/work/evidence/2026/08/29/`.
- Started a second shaping tune in `format.go` to reduce overhead on sparse
  result sets; this is still uncommitted and under evaluation.

## What Was Directly Verified

- `go test ./internal/application/search ./internal/infrastructure/filesystem`
  passed during the first refactor slice.
- Focused formatter tests passed after the shaping tune:
  - `go test ./internal/application/context -run 'TestFormatAsMarkdown|TestFormatAsMarkdownRespectsBudget|TestFormatAsMarkdownEmpty|TestContextBlockString|TestContextBlockReferences|TestContextBlockStringInterleavesReferences'`
- `go test ./internal/platform/cli -run TestRunSearchWithContext` passed for the
  committed refactor slice.
- `go build ./internal/application/search ./internal/application/context ./internal/infrastructure/filesystem ./internal/platform/cli` passed for the committed refactor slice.
- `git diff --check` passed before the first local commit.
- `./.hawp/bin/hawp work validate` returned `VALIDATION PASS` with the existing
  repo warning:
  - `VERIFICATION CLARITY [WARN]` (`Ambiguous: 7`)
- Branch-local index was rebuilt successfully:
  - documents: `434`
  - chunks: `3382`
- Full embed coverage is now complete:
  - embedded chunks: `3382`
  - remaining unembedded chunks: `0`
- Benchmark results after full embed using the checked-in branch-local binary:
  - lexical: `0.7ms`, `10/10` high
  - semantic: `522.2ms`, `9/10` high
  - hybrid: `72.7ms`, `10/10` high
- Token-savings benchmark after full embed using the checked-in branch-local
  binary:
  - total raw tokens: `22866`
  - total shaped tokens: `18527`
  - total saved: `4339`
  - overall savings: `19%`
- Token-savings benchmark after the shaping tune using a fresh source-built
  binary at `/tmp/hawp-v0023`:
  - total shaped tokens: `18278`
  - total saved: `4588`
  - overall savings: `20%`
  - previously negative sparse-query cases improved:
    - intake ordering: `-11%` -> `-5%`
    - provider materialization: `-10%` -> `-7%`
    - core shape fields: `-20%` -> `-9%`

## What Remains Unproven

- The current shaping tune has not yet been committed or rerun as the checked-in
  `.hawp/bin/hawp` binary on this branch.
- The sparse-query negative cases are improved but not eliminated yet.
- The next architecture slices are still only audit-recommended, not yet
  implemented:
  - unify `runSearchIndex` with the canonical corpus-builder service
  - unify `runSearch` with the shared typed search service
  - continue reducing formatting overhead on sparse result sets
- No new main-branch binary chore update or version bump has been performed in
  this checkpoint.

## Constraints

- Keep work on this branch unmerged until the next cleanup and evidence pass is
  reviewed.
- Use other repositories only as design references; do not copy private paths,
  names, or project-specific content into this repo.
- Keep changes minimal, capability-local, and scalable; avoid broad rewrites.
- Do not commit machine-local absolute paths or personal data.

## Help Wanted

- Continue the audit-driven cleanup using the same bounded, compoundable slices.
- Prioritize removing the remaining negative token-savings cases without
  sacrificing provenance or making the formatter harder to reason about.
- Keep usage/benchmark evidence updated as shaping changes land.

## Suggested Next Step

1. Finish the current sparse-result shaping cleanup and rerun the token benchmark
   with the checked-in branch-local binary, not only `/tmp/hawp-v0023`.
2. If the negatives are still present, add an explicit sparse-result passthrough
   rule or lower-overhead render mode for under-budget result sets.
3. Then continue the architecture cleanup with the next high-leverage slice:
   consolidate `runSearchIndex` and `runSearch` onto shared typed services.
