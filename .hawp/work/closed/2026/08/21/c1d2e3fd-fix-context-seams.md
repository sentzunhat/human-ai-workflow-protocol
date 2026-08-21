---
work-item: c1d2e3fd
type: fix
title: "Separate context retrieval, formatting, and reshape seams"
status: done
created: 2026-08-10
updated: 2026-08-21
parent: b6c4e8a2
depends-on: c1d2e3fc
---

# Fix: Context Seams

## Mission

Make retrieval, context formatting, and optional LLM reshape independently
testable while preserving current output and fallback behavior.

## Done When

- Each stage has a clear capability-local contract.
- Provenance, token budget, deduplication, encryption, and fallback tests pass.
- No public CLI behavior changes without explicit evidence.

## Audit Evidence

- CLI context mode and `RAGPipeline.Retrieve` use different retrieval and
  deduplication paths, so identical search results can render differently.
- Reshape construction is coupled to provider factories, and invocation-time
  `maxTokens` is ignored in favor of constructor configuration.
- OpenAI and Anthropic configuration validation implies support that the
  runtime factories do not provide; config persistence is placeholder-only.

## Smallest Safe Slice

- Add capability-local retrieval, pure formatting, and injected reshape seams
  with compatibility wrappers until the CLI uses one orchestration path.
- Keep unsupported providers explicit for v0.1.0; do not combine config
  persistence implementation with this behavioral refactor.

## Verification

- CLI and pipeline produce equivalent context blocks for the same inputs.
- Preserve ordering, provenance, budgets, model-free passthrough, and missing
  model fallback.
- Use fake retriever, embedder, and LLM tests; prove `Reshape` honors its
  invocation-time token limit.

## Progress

- Added injected `Retriever` and `NewRAGPipeline` composition seams.
- Added injected embedder/LLM construction and invocation-time token handling.
- Added `PrepareContext` so pipeline retrieval owns result deduplication before
  formatting, while preserving `FormatAsMarkdown` as a compatibility wrapper.
- Focused context tests and diff checks pass.
- CLI context mode migrated: `prepareSearchContext` delegates to
  `PrepareContext` (one shared orchestration path, no separate deduplication).
- `DefaultRetrieveMaxTokens` exported from `rag.go` so the equivalence test
  can reference the pipeline's token budget without duplicating the magic number.
- `TestCLIAndPipelineContextAreEquivalent` added to `search_output_test.go`:
  proves CLI path and `pipeline.Retrieve` produce identical `ContextBlock`
  values for the same inputs (including near-duplicate collapse). Uses fake
  retriever, embedder, and LLM — no model files or network required.
- `TestPrepareSearchContextMatchesApplicationPreparation` (existing): proves
  CLI path matches `PrepareContext` directly.
- `TestRAGPipelineReshapeHonorsInvocationMaxTokens` (existing): proves
  `Reshape` uses invocation-time `maxTokens`, not the constructor default.
- All Done When criteria met: clean build, all focused tests pass.

## Outcome

Migrated CLI context mode from duplicated retrieval logic to the shared retrieve-format-reshape Pipeline. Exported `DefaultRetrieveMaxTokens` from the RAG module. Added `TestCLIAndPipelineContextAreEquivalent` proving CLI and pipeline produce identical output using in-process stubs. Context logic is no longer duplicated between CLI and pipeline paths.

## Verification

- `go test ./internal/platform/cli ./internal/application/context ...` all pass.
- Build: `CGO_ENABLED=0 go build ./...` clean.
- Merged to development as PR #4.

## Close Checklist

- [x] CLI context mode uses shared Pipeline
- [x] `DefaultRetrieveMaxTokens` exported
- [x] Equivalence test added and passing
- [x] No duplicated context retrieval logic
- [x] Merged to development
