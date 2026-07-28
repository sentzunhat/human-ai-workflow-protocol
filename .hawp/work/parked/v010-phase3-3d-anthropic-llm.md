---
work-item: v010-3-3d
type: feature
title: "v0.1.0 Phase 3.3d: Anthropic LLM Backend"
status: parked
owner: unassigned
created: 2026-07-25
updated: 2026-07-27
---

# Phase 3.3d: Anthropic LLM (parked)

## Status

Parked indefinitely, no v0.1.0 milestone — all future work ships as
incremental v0.0.4+ patches. This detail doc was never written before the
strategy changed; recorded here purely so the backlog row and its file are
consistent (see `librarian/docs/CONTEXT_AUDIT_v003.md` for the current
capability/gap picture).

## Mission (if picked up later)

Implement Anthropic Messages API (`claude-3-sonnet`, `claude-3-opus`) as a
configurable `LLMClient` backend, alongside the existing ONNX and Ollama
backends (see `internal/domain/llm/llm_client.go`).

## Related Work

- `v010-3-2d` (Anthropic Embeddings Stub) — same parked status, same rationale
- `v010-cost` (Cost Tracking) — would wrap this backend if ever built
