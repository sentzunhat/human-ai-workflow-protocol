# Changelog

All notable changes to the `hawp` Go librarian CLI are documented here.
Format loosely follows [Keep a Changelog](https://keepachangelog.com/).

## [0.0.2] - 2026-08-24

Maintenance release: validator improvements, source layout cleanup,
release pipeline simplification, and platform-aware binary install.

### Work Validator

- `hawp work validate` now recognizes folder-per-item layouts:
  `active/{uuid}/plan.md` and `closed/YYYY/MM/DD/{uuid}/plan.md`
- Shared batch-close records (`**Closes:** \`id1\`, \`id2\`…`) are matched
  by content — one closed record can cover multiple backlog rows
- Short-ID-slug (`b7e2a4f9-title`) and bare 8-char alphanumeric directory
  names are now recognized as work item IDs

### Internal

- Source moved from `librarian/go/` to `librarian/src/`; Go module path
  updated from `github.com/sentzunhat/hawp/librarian/go` to
  `github.com/sentzunhat/hawp/librarian/src`
- npm hawp scripts (`kit:validate`, `work:validate`, `hawp:check`, etc.)
  now invoke the compiled Go binary directly — TypeScript fallbacks retired

### Install

- Provider install scripts now detect OS and arch at runtime and download
  the matching `hawp` binary from the GitHub release (macOS arm64/amd64,
  Linux amd64/arm64, Windows amd64/arm64). SHA256 verification runs when
  `checksums.txt` is available. Falls back gracefully when the platform is
  unsupported or the network is unavailable.

### Release

- Tag format changed from `librarian-go-vX.Y.Z` to plain `X.Y.Z`
  (e.g. `0.0.2`). `hawp update` handles both formats transparently.

---

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
