# Release And Architecture Checkpoint

## Current Product Target

The current `v0.1.0` vision defines a small, local context tool: index the
HAWP kit/work corpus, search it, retain provenance, and optionally reshape
context through offline ONNX or Ollama backends. The `none` backend remains a
zero-network fallback. This is a product milestone, not a requirement to add
every future provider.

## What Is Now Stronger

- Search uses a typed search-index port with a SQLite search adapter, and CLI
  and RAG retrieval share the application search path.
- Context acquisition uses typed corpus input with a filesystem context
  adapter; domain context no longer owns repository reads.
- Kit validation and normalization use typed workspace snapshots with a
  filesystem kit adapter; domain kit no longer owns reads, Markdown parsing,
  renames, or writes.
- Ports and adapters are grouped by capability rather than under global
  `ports` or `adapters` directories.

## Release Gates Still Needed

1. Decide and document the authoritative milestone: the vision limits the
   release to offline ONNX/Ollama/none, while `librarian/docs/BACKENDS.md`
   still presents cloud providers as planned for `v0.1.0`.
2. Prove the release matrix, especially whether the ONNX LLM/ORT path is
   shippable on each supported target or must remain an explicitly unavailable
   optional backend.
3. Add end-to-end acceptance fixtures for index, lexical/semantic fallback,
   context provenance, token budgets, and missing-model behavior.
4. Run the six-binary release verification and the update/install smoke once a
   release candidate is intentionally selected.

## Architecture Order

1. Audit and extract the `domain/work` boundary (`c1d2e3f8` then `c1d2e3f9`).
2. Audit/index typed persistence and SQLite capability groups.
3. Split mixed CLI command families after application seams are stable.
4. Introduce a small composition layer for provider selection only after the
   capability contracts are complete; do not add a general plugin framework.

## Provider Decision

Current runtime providers are ONNX, Ollama, and `none`. Agent/editor overlays
are manifest-driven for Claude, Codex, GitHub, Cursor, and Continue. Cloud
OpenAI/Anthropic providers remain future work, not a current release blocker.
