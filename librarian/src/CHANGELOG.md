# Changelog

All notable changes to the `hawp` Go librarian CLI are documented here.
Format loosely follows [Keep a Changelog](https://keepachangelog.com/).

## [0.0.7] - 2026-08-24

Real `--semantic` search mode and three-way benchmark.

### Added

- **`--semantic` flag** — pure-vector search: embeds the query using the same
  backend/model recorded in the index, ranks all stored chunk vectors by cosine
  similarity, and returns the top-N results. No FTS5 involved. Works with both
  ONNX and Ollama backends. Useful for conceptual queries with few keyword
  matches: in the v0.0.7 benchmark, queries that returned 1–4 FTS5 results
  returned 10 semantic results at 9/10 quality.
- **`QueryChunksByIDs`** in `infrastructure/sqlite/index.go` — fetches full
  document/chunk rows for a list of IDs in caller-supplied order, so semantic
  ranking survives the SQLite lookup.
- **Three-way benchmark** — `hawp search benchmark` now covers
  lexical / semantic / hybrid. v0.0.7 results: lexical 0.1ms 10/10,
  semantic 478.6ms 9/10, hybrid 72.0ms 10/10.

### Changed

- Benchmark header updated to "LEXICAL / SEMANTIC / HYBRID".
- Registry search description updated to document `--semantic`.
- `search.md` modes table updated with semantic row and corrected flag
  reference block.

## [0.0.6] - 2026-08-24

Configurable index paths, release title rename, and `--semantic` removal.

### Added

- **Configurable index paths** — `hawp search index` now reads
  `.hawp/config/search.json` (project) and `~/.hawp/config/search.json`
  (home) to determine which paths to walk. Project config overrides home;
  missing config files fall back to the default (`[".hawp/kit", ".hawp/work"]`).
  Extra paths (anything other than kit/work) are walked generically as
  `"custom"` corpus — directories recursed for `.md` files, single files
  indexed directly. Missing paths warn and skip rather than failing.
- Added `.hawp/config/search.json` to this repo, adding `librarian/docs/`
  and `README.md` to the search corpus (392 docs / 3064 chunks).

### Fixed

- **SQLite schema** — `documents.category` CHECK constraint extended to
  include `'custom'` so user-configured paths can be indexed. Existing
  indexes must be rebuilt (`hawp search index`).
- **Benchmark relative-performance display** — when lexical avg < 1ms,
  the division produced `+Inf`. Now prints absolute hybrid latency with a
  note that the ratio is not meaningful at sub-millisecond speeds.

### Changed

- **Release title** renamed from `librarian-go vX.Y.Z` to
  `HAWP librarian vX.Y.Z`.
- **`--semantic` flag removed** from registry and docs — it was never
  implemented as a separate code path (silently fell through to lexical).
  Search automatically uses hybrid re-ranking when vectors are present.
- **`search.md`** updated: removed phantom `--semantic` row, corrected
  hybrid latency range to 50–90ms, added configurable paths section.

## [0.0.5] - 2026-08-24

`hawp update --check` read-only fix, benchmark rewrite, and docs naming cleanup.

### Fixed

- **`hawp update --check`** — previously triggered the full update instead of being
  read-only. Now routes to `runUpdateVerify()` (same path as `hawp update verify`):
  exits 0 when up-to-date, exits 1 when an update is available, no side effects.
  `--check` is now documented in registry help text and `hawp update --help`.

- **`hawp search benchmark`** — benchmark was entirely fake: all three patterns called
  `QueryChunksLexical`; quality scores were hardcoded; semantic mode doesn't exist as
  a separate code path. Fixed:
  - Hybrid now calls real `appsearch.HybridRank` (lexical candidates → Ollama embed →
    cosine re-rank), matching the production search path.
  - Query set replaced: 15 Go/ML queries (8 zero-hit) → 10 HAWP corpus queries
    (all verified to return results against the kit + work index).
  - Quality scoring: hardcoded → keyword-matching against the actual top result text.
  - Patterns reduced to `["lexical", "hybrid"]` — no fake semantic-only mode.
  - Dead code removed: `benchmarkGetStr` (deduped with existing `getStr` in same
    package), `"semantic"` entries in the summary pattern loops.

### Changed

- `librarian/docs/` files renamed UPPERCASE → kebab-case to match the repo
  documentation naming standard. All internal cross-references updated.

### Documentation

- `search.md` — corrected Ollama warm embed time from ~110ms/chunk to ~32ms/chunk
  (measured on 2873-chunk batch with nomic-embed-text).
- Added `librarian/docs/benchmarks-v004.md` — v0.0.4 benchmark results, pre-fix and
  post-fix tables, comparison vs v0.0.3, FTS5 fix confirmation, corpus mismatch notes.

## [0.0.4] - 2026-08-24

stdio MCP server (`hawp mcp`) and expanded `hawp init`.

### Added

- **`hawp mcp`** — start a stdio JSON-RPC 2.0 MCP server. Exposes three tools
  to any connected AI agent (Claude Code, Cursor, Continue, etc.):
  - `hawp_search` — hybrid lexical+vector search over indexed kit and work docs.
  - `hawp_work_new` — scaffold a new work item (UUID, plan file, BACKLOG row).
  - `hawp_work_validate` — run composite kit + work + links validation.
- **`hawp init --provider <name>|all`** — new `--provider` flag writes the MCP
  server entry into the appropriate config file for each requested provider:
  - `claude` → `.mcp.json`
  - `cursor` → `.cursor/mcp.json`
  - `continue` → prints manual config instructions
  - `all` → expands to claude + cursor + continue + codex
  Config writes are idempotent (no-op if the hawp entry already exists).
- `hawp init` now also syncs `.hawp/kit/` from the latest release (same as
  `hawp update sync`) after provisioning `~/.hawp/`.

### Changed

- Registry description for `init` and `update` updated to reflect new behavior.

---

## [0.0.3] - 2026-08-24

Provider update parity, cross-platform auto-update CI, search usage docs,
and kit reference docs updated for folder-per-UUID work item layout.

### Update

- `hawp update` now syncs all provider overlays by default (equivalent to
  `--provider all`). Use `--no-providers` to update the binary and kit only.
- Help text updated: `update  update binary + kit + all providers (--no-providers for kit-only)`

### CI

- New workflow: `test-auto-update.yml` — runs on every `main` push across
  Linux (ubuntu-latest), macOS (macos-latest), and Windows (windows-latest).
  Downloads the `0.0.1` binary, runs `hawp update`, and asserts the binary
  upgraded to the latest release. Skips gracefully when only one release
  exists or baseline is already latest.

### Docs

- Added `core/.hawp/kit/usage/search.md`: full guide for `hawp search`
  including search modes, context packing (`--context`), LLM reshape
  (`--llm-reshape`), embedding backends, flag reference, and agent workflow.
- Updated `core/.hawp/kit/references/backlog-alignment.md`: documents
  folder-per-UUID work item layout (`active/{uuid}/plan.md`,
  `closed/YYYY/MM/DD/{uuid}/plan.md`) as the current standard alongside
  the legacy flat-file format.
- Updated `core/.hawp/kit/references/work-item-file-tracking.md`: aligns
  file location examples and agent usage instructions with folder-per-item
  layout; legacy `TASK-NNN-files.md` format noted as still accepted.

---

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
