---
work-item: c1d2e3fc
type: audit
title: "Recursive audit: application context capability"
status: done
created: 2026-08-10
updated: 2026-08-13
parent: b6c4e8a2
follow-up: c1d2e3fd
---

# Audit: Application Context Capability

## Outcome

Context has solid building blocks, but CLI context mode and the RAG pipeline
are separate paths. Retrieval, formatting, and reshape need explicit seams
before further behavior is added.

## Findings

- CLI retrieval/deduplication differs from `RAGPipeline.Retrieve`.
- Formatting deduplicates references but not chunks, producing output drift.
- Reshape construction hardwires provider factories and ignores invocation-time
  `maxTokens`.
- Provider configuration accepts unsupported OpenAI and Anthropic choices.
- Config persistence currently reports success without writing configuration.

## Follow-up

`c1d2e3fd` will establish retrieval, pure formatting, and injected reshape
seams while retaining compatibility wrappers until the CLI has one
orchestration path. Config persistence is explicitly out of scope.

## Verification

- Equivalent context output through CLI and pipeline for identical inputs.
- Fake retriever, embedder, and LLM coverage plus token-limit proof.
- Preserve provenance, budgets, model-free passthrough, and fallback behavior.
