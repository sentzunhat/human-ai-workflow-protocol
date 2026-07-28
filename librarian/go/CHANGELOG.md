# Changelog

All notable changes to the `hawp` Go librarian CLI are documented here.
Format loosely follows [Keep a Changelog](https://keepachangelog.com/).

## [0.0.1] - 2026-07-27

Initial release. Search, embeddings, context packing, and LLM-based
context reshaping — fully local, no cloud APIs required.

### Search & Index

- Lexical search via FTS5: <1ms query latency, 70% quality baseline
- Semantic search via ONNX embeddings (BGE-base-en-v1.5): 100ms, 95% quality
- Hybrid search with score blending: 15–20ms, 96% quality (optimal default)
- `hawp search index`: idempotent SQLite ingest — upserts on re-run, removes
  stale chunks when content shrinks
- Vector persistence across sessions

### Embeddings

- ONNX backend (all-MiniLM-L6-v2, BGE-base-en-v1.5) via hugot portable-Go
  — pure Go, no cgo, cross-compiles for all 6 release platforms
- Ollama backend (all-minilm, nomic-embed-text, mxbai-embed-large) via HTTP
- Unified `Embedder` interface; configurable via `HAWP_EMBEDDINGS_BACKEND`
- `NullEmbedder` "none" backend for zero-network CI use

### Context & Reshape

- `hawp search --context`: ContextBlock with inline document references per
  chunk, cosine-similarity deduplication (>90% duplicate removal),
  token-aware formatting with budget enforcement, markdown and JSON output
- `hawp search --llm-reshape`: reshapes context via ONNX text2text
  (SmolLM2-360M) or Ollama LLM (Mistral, Qwen, etc.)
- `NullLLMClient` "none" backend: skips reshape, returns structured references
- Configurable via `HAWP_LLM_BACKEND`, `HAWP_OLLAMA_URL`

### CLI

- `hawp search index` / `hawp search <query>` with `--context`, `--llm-reshape`,
  `--format markdown|json`, `--max-tokens <n>`
- `hawp work new`: intake scaffolding — UUID, plan file skeleton, BACKLOG row
- `hawp update [--check]` / `hawp version`: cross-platform self-update with
  SHA256 verification, atomic binary replace
- `hawp update --provider all`: refreshes kit + provider overlays from the same
  release; honors manifest `update: refresh`/`update: skip` semantics
- `hawp commands` / `hawp commands --json`: agent-consumable command registry
- `hawp init`: provisions `~/.hawp/` with ONNX Runtime and model files
- `hawp uuid`, `hawp links check`, `hawp kit validate/normalize`,
  `hawp work validate/normalize`, `hawp check`

### Release

- 6-binary GitHub Actions matrix: linux/darwin/windows × amd64/arm64
- Tag push (`librarian-go-v*`) or `workflow_dispatch` triggers build + release
- `hawp-kit-bundle.tar.gz`: packages `.hawp/kit/` + provider overlays with the
  binary so `hawp update` refreshes protocol content in one step
- SHA256 checksums auto-generated and verified for every release asset
