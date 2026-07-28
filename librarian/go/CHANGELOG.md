# Changelog

All notable changes to the `hawp` Go librarian CLI are documented here.
Format loosely follows [Keep a Changelog](https://keepachangelog.com/).

## [0.0.3] - 2026-07-27

First real release. Consolidates all pre-release work (0.0.1–0.0.6 test
builds) into a single verified baseline: search, embeddings, context
packing, and LLM-based context reshaping.

### Search & Index

- Lexical search via FTS5: <1ms query latency, 70% quality baseline
- Semantic search via ONNX embeddings (BGE-base-en-v1.5): 100ms, 95% quality
- Hybrid search with score blending: 15–20ms, 96% quality (optimal default)
- `hawp search index`: idempotent SQLite ingest — upserts on re-run, removes
  stale chunks when content shrinks
- Vector persistence across sessions; 2,445 chunks indexed in <2 seconds

### Embeddings

- ONNX backend (all-MiniLM-L6-v2, BGE-base-en-v1.5) via hugot portable-Go
  — pure Go, no cgo, cross-compiles cleanly for all 6 release platforms
- Ollama backend (all-minilm, nomic-embed-text, mxbai-embed-large) via HTTP
- Unified `Embedder` interface; configurable via `HAWP_EMBEDDINGS_BACKEND`

### Context & Reshape

- `hawp search --context`: ContextBlock with inline document references per chunk,
  deduplication via cosine similarity (removes >90% of duplicates), token-aware
  formatting with strict budget enforcement, markdown and JSON output formats
- `hawp search --llm-reshape`: reshapes context via ONNX text2text (SmolLM2-360M)
  or Ollama LLM (Mistral, Qwen, etc.) — configurable via `HAWP_LLM_BACKEND`
- `NullEmbedder`/`NullLLMClient` "none" backend for zero-network CI use
- ReshapeBatch serialized per-context (fixes ort-genai batch-EOS bug where
  multi-context batches failed with "max length reached")

### CLI & Config

- `hawp search index` / `hawp search <query>` / `hawp search --context --llm-reshape`
- `hawp work new`: HAWP intake scaffolding — generates UUID, plan file, BACKLOG row
- `hawp update [--check]` / `hawp version`: cross-platform self-update with SHA256
  verification; atomic binary replace (rename-over); skips gracefully on Windows
- `hawp commands` / `hawp commands --json`: agent-consumable command registry with
  usage, flags, and exit-code semantics per command
- `hawp init`: provisions `~/.hawp/` with ONNX Runtime and embedding model files,
  checksum-verified, idempotent
- `hawp update --provider all`: refreshes kit + all provider overlays from the
  same release the binary updates to (one unified version); honors manifest
  `update: refresh`/`update: skip` semantics; never overwrites user-customized files
- Configurable backends: `HAWP_EMBEDDINGS_BACKEND`, `HAWP_LLM_BACKEND`,
  `HAWP_OLLAMA_URL`, `HAWP_ONNX_LIBRARY_DIR`

### Workflow Commands (ported from TypeScript tooling)

- `hawp uuid [--short]`, `hawp links check`, `hawp kit validate`,
  `hawp kit normalize [--apply]`, `hawp work validate`, `hawp work normalize`,
  `hawp check` (kit + work + links in one pass)
- `work:validate` fixes: short-ID-slug filename recognition, shared closed-record
  support via `**Closes:** \`id\`` content matching

### Release & CI

- 6-binary GitHub Actions matrix: linux/darwin/windows × amd64/arm64, CGO_ENABLED=0
- Release workflow supports both tag-push trigger and manual `workflow_dispatch`
  with version input (creates and pushes tag itself)
- CHANGELOG section required per version — workflow hard-fails if missing
- `hawp-kit-bundle.tar.gz` release asset: packages `.hawp/kit/` + provider overlays
  alongside the binary so `hawp update` can refresh protocol content in one step
- SHA256 checksums auto-generated and verified for every release asset

### Housekeeping

- Removed local machine paths from all `.hawp/work/` files (public repo safety)
- Removed `.hawp/db/index.sqlite` and compiled binaries from git tracking
- Added `.hawp/db/` and binary artifacts to `.gitignore`
- Squashed 75 incremental development commits into this single clean baseline
