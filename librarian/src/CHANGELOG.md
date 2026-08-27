# Changelog

All notable changes to the `hawp` Go librarian CLI are documented here.
Format loosely follows [Keep a Changelog](https://keepachangelog.com/).

## [0.0.20] - 2026-08-27

Architecture audit + ports-and-adapters fix for `domain/usage`.

### Changed

- **`domain/usage.Store`** — extracted as a port interface (`Write`, `Recent`,
  `GetTotals`, `GetReport`, `Clear`, `Close`). The concrete SQLite implementation
  is now the unexported `sqliteStore` type. `Open()` returns `Store` (interface)
  so callers no longer depend on the concrete type.
- **`domain/usage/store_test.go`** — `openTemp` return type updated to `Store`.

### Added

- **Architecture audit evidence** at `.hawp/work/evidence/2026/08/27/v0020-arch-audit.md`:
  full gap analysis vs. mictlan's area-first ports-and-adapters pattern, with
  compliance table and recommended next steps per patch.

## [0.0.19] - 2026-08-27

`hawp search benchmark --token` — token-savings benchmark measuring context
shaping reduction vs raw result text. Closes the v0014-token-speed-bench gate.

### Added

- **`hawp search benchmark --token`** — runs the 10-query benchmark set through
  the full shaping pipeline (`FormatAsMarkdown` at 2000-token budget) and
  reports raw vs shaped token counts per query, plus aggregate savings.
  `--export <path>` writes the Markdown table as an evidence artifact.
- **`RunTokenBenchmark(db, exportPath)`** in `internal/platform/cli/benchmark` —
  drives lexical/hybrid search, converts to `domainsearch.Result`, calls
  `appcontext.FormatAsMarkdown`, and renders the `formatTokenReport` Markdown.
- **Evidence artifact** at `.hawp/work/evidence/2026/08/27/v0019-token-savings-benchmark.md`
  — real run against the live index: **19% overall token reduction** (10 queries,
  2000-token budget; up to 38% on dense queries, negative on sparse results).

## [0.0.18] - 2026-08-27

`hawp usage report` — Markdown usage summary with per-tool breakdowns and
token estimates, suitable for evidence artifacts.

### Added

- **`hawp usage report`** — new subcommand that prints a full Markdown usage
  summary: overall totals, per-tool breakdown table, and a recent-queries
  table with estimated token counts. Pass `--export <path>` to write the
  report to a file for inclusion in `.hawp/work/evidence/`.
- **`GetReport()` / `FormatReport()`** in `internal/domain/usage` — produce
  aggregate and per-tool `ToolStat` breakdowns from the local SQLite log.
- **`TestRunUsageReport`** in `cli` package — covers empty log output and
  `--export` path round-trip.
- **`TestGetReport` / `TestFormatReportEmpty`** in `usage` package — cover
  aggregate stats and Markdown rendering.

## [0.0.17] - 2026-08-27

Usage log query visibility — `hawp usage log` now shows the actual query text
for every entry, even when body capture is disabled (the default).

### Added

- **`query_text` column in usage log** — the first 256 characters of the
  `query` (or `title`) field from every MCP tool call are now stored directly
  in the usage log without requiring `--log-bodies`. The column is added to
  existing databases via a safe idempotent `ALTER TABLE` migration on open.
- **`EntrySummary` helper** — `internal/domain/usage` exposes `EntrySummary(e Entry) string`
  that returns the best display string: stored `query_text` wins, then falls
  back to the body (when available), then the hash.
- **CLI tests for `hawp usage` subcommands** — `TestRunUsage` and
  `TestRunUsageClearCancelled` cover enable/disable/totals/log/clear against
  an isolated `$HOME`, exercising the full CLI dispatch path.

### Changed

- **`hawp usage log`** — entries now show the stored query text instead of
  a truncated hash when body capture is off (the common case). Output format
  and column layout are unchanged.

## [0.0.16] - 2026-08-27

Linux ORT GenAI archive-layout fix so the release workflow places
`libonnxruntime-genai.so` under `lib/` before linking and packaging.

### Fixed

- **Linux ORT-GenAI extraction path** — the Linux
  `onnxruntime-genai-0.13.1-linux-x64.tar.gz` archive does not include the
  leading `./` path segment that the macOS archive uses. With
  `--strip-components=2`, GitHub Actions extracted `libonnxruntime-genai.so`
  into `$NATIVE/` instead of `$NATIVE/lib/`, so the linker failed with
  `cannot find -lonnxruntime-genai`. The Linux download step now uses
  `--strip-components=1`, which preserves the `lib/` and `include/` layout
  expected by the build and packaging steps.

## [0.0.15] - 2026-08-27

Linux ORT release path fix so the GitHub Release workflow can package the
`hawp-linux-amd64-ort.tar.gz` artifact from the correct native-lib directory.

### Fixed

- **Linux ORT native-lib path on Actions** — the `build-ort-linux-amd64` job no
  longer relies on workflow-level `${{ env.HOME }}` expansion for `NATIVE`,
  which resolved to `/.hawp/native/linux-amd64` in the failing `0.0.14`
  release run. The build and package shell steps now set
  `NATIVE="$HOME/.hawp/native/linux-amd64"` directly before invoking the linker
  and copying shared libraries, so `libonnxruntime`, `libonnxruntime-genai`,
  and `libtokenizers` are resolved from the downloaded native bundle path.

## [0.0.14] - 2026-08-27

ORT release build — ONNX LLM reshaping now ships in official release tarballs.

### Added

- **ORT release tarballs** — the GitHub Actions release workflow now builds two additional
  tarballs per release: `hawp-linux-amd64-ort.tar.gz` and `hawp-darwin-arm64-ort.tar.gz`.
  Each tarball contains a `hawp` binary built with `-tags ORT` (CGO_ENABLED=1) and a `lib/`
  directory with the three required native libraries (`libonnxruntime`, `libonnxruntime-genai`,
  `libtokenizers`). The standard six-platform CGO_ENABLED=0 binaries are unchanged.
- **ORT build jobs** — two parallel CI jobs (`build-ort-linux-amd64` on ubuntu-latest,
  `build-ort-darwin-arm64` on macos-14) download Microsoft ORT and ORT-GenAI prebuilts plus
  a static `libtokenizers` (daulet/tokenizers), then build with `-tags ORT` and package the
  result. Native lib versions: onnxruntime 1.19.2, onnxruntime-genai 0.13.1, tokenizers 1.27.0.
- **Fault-isolated release job** — the publish step runs even when ORT build jobs fail
  (`if: always() && needs.build-std.result == 'success'`), so a transient lib-download failure
  does not block the standard binaries from shipping. ORT tarballs are simply absent from that
  release; standard tarballs are always present.
- **darwin/amd64 stays CGO_ENABLED=0** — no official Microsoft prebuilts exist for Intel Mac;
  ONNX LLM is unavailable there. hugot returns a clear "build with -tags ORT" error when
  `llm.backend: "onnx"` is requested on a non-ORT binary, so no panic.
- **rpath for macOS** — darwin ORT binaries set `@executable_path/lib` so the `.dylib` files
  can live next to the `hawp` binary without any `DYLD_LIBRARY_PATH` configuration.

### Notes

- Model: `homen3/SmolLM2-360M-Instruct-ort-genai-int4-cpu` (verified working 2026-07-27,
  ~1.1s/reshape on arm64). Pull it with `hawp model pull homen3/SmolLM2-360M-Instruct-ort-genai-int4-cpu`.
- Set `llm.backend: "onnx"` in `~/.hawp/config/context.json` to activate local LLM reshaping
  for `hawp search --context`.

## [0.0.12] - 2026-08-25

Local MCP call log with token estimates — opt-in, SQLite-backed, never blocks tool responses.

### Added

- **`hawp usage`** — show totals across all recorded MCP calls: total calls, estimated
  tokens in, estimated tokens out, and estimated tokens saved via context shaping.
- **`hawp usage log`** — tail the 20 most recent entries (timestamp, tool, tokens in/out, query summary).
- **`hawp usage enable`** — opt into call logging. Config stored in `~/.hawp/config/usage.json`.
  Pass `--log-bodies` to also store raw input/output JSON (disabled by default — bodies may
  contain sensitive prompt text).
- **`hawp usage disable`** / **`hawp usage clear`** — turn off logging or truncate the log.
- **`~/.hawp/usage.db`** — SQLite log file, created on first recorded call (not at binary launch).
  Separate from the search index so `hawp search index` never touches it.
- **MCP auto-log** — after every `hawp_search`, `hawp_work_new`, and `hawp_work_validate`
  tool call, a goroutine writes an entry to the log. Fire-and-forget: log write failures are
  silent and never block or error the tool response.
- Token estimates use the `len(JSON)/4` character approximation — no tokenizer dependency.

## [0.0.11] - 2026-08-25

Cursor MCP `type:stdio` + wrapper, AGENTS.md seed-if-missing, embed-optional docs, install lessons from downstream.

### Fixed

- **Cursor MCP `type: stdio`** — `hawp init --provider cursor` now writes `"type": "stdio"` in
  `.cursor/mcp.json`. Cursor requires this field explicitly; without it the entry is ignored.
- **Cursor `hawp-mcp` wrapper** — Added `core/.hawp/bin/hawp-mcp` wrapper script that `cd`s to the
  repo root before exec-ing `hawp mcp`, so Cursor's spawn cwd does not affect the server's
  working directory. The wrapper is the `command` in `.cursor/mcp.json`; its path is absolute.
- **AGENTS.md seed-if-missing** — `hawp init` no longer overwrites an existing `AGENTS.md` on update.
  The `install:seed-if-missing` rule in `ApplyProviderInstall` skips the file when the destination
  already exists.

### Docs

- **Embed step optional** — `search.md` and install guides clarify that `hawp search index` must
  run first; `hawp search embed` is optional and can be very slow on CPU (~1 chunk/s). Hybrid/semantic
  search only require embed if you want those modes.
- **Cursor MCP UI path** — Updated to `Customize → MCPs` (was `Settings → Tools & MCP`).
- **Manager-branch pattern** — Added `.hawp/kit/usage/manager-branch.md` documenting the optional
  manager-branch / worktree operating pattern.

### Changed

- **`/releases/latest` fallback** — Install script now falls back to `/releases?per_page=1` when
  `/releases/latest` returns 404 (e.g. all releases are marked pre-release on GitHub).

## [0.0.10] - 2026-08-25

Codex MCP absolute-path fix — project-scoped MCP servers require absolute paths.

### Fixed

- **Codex MCP absolute paths** — `hawp init --provider codex` now writes absolute paths for
  both `command` and `cwd` in `.codex/config.toml`. Codex desktop does not resolve relative
  paths in project-scoped MCP configs relative to the project root; they resolve relative to
  the Codex app directory, causing the binary to be silently not found. The Codex session log
  reports this as `omitting MCP server without an exact ready client server_name=hawp`. Fix:
  `writeCodexTOML` now accepts `repoRoot` and embeds the absolute binary path and project root.

## [0.0.9] - 2026-08-25

Codex MCP config path fix, binary wrapper URL fix, and install/update binary naming fix.

### Fixed

- **Codex MCP config path** — `hawp init --provider codex` now writes `.codex/config.toml`
  (the path Codex actually reads) instead of `codex.toml` at the repo root. Added `cwd = "."`
  to the TOML block per verified Codex config format. Idempotency check updated to match
  the new path.
- **Shell wrapper release URL** — `core/.hawp/bin/hawp` pointed to
  `github.com/sentzunhat/hawp/releases` (404). Corrected to
  `github.com/sentzunhat/human-ai-workflow-protocol/releases`.
- **Update script clobbered Go binary** — the distribution update script was copying
  `core/.hawp/bin/hawp` (the shell wrapper) over `.hawp/bin/hawp`, replacing the user's
  real Go binary with a wrapper that then failed with `hawp: Go binary not found.` because
  `hawp-bin` did not exist. Fix: update script now downloads the Go binary to `hawp-bin`
  (matching the wrapper's lookup path) and installs the shell wrapper at `hawp` alongside it.
  Install script updated identically.

## [0.0.8] - 2026-08-25

`engine` key canonical; `--llm-reshape` removed; real context dedup + dynamic chunk cap;
MCP context mode; 48h auto-update notifier.

### Added

- **`--verbose` / `-v` for `hawp search --context`** — prints a summary to stderr:
  `context: N chunks, ~M tokens (saved ~K tokens via dedup)`. Token estimate uses
  character-count/4 approximation; dedup savings reflect chunks dropped by Jaccard filter.
- **`--hybrid-ratio <f>` for `hawp search`** — sets the lexical fraction of hybrid
  scoring (default 0.3). `0.5` = equal blend, `0.7` = lexical-heavy. Must be in [0.0, 1.0].
- **Hawp-first session workflow** — new kit doc at `.hawp/kit/usage/hawp-first-workflow.md`:
  MCP search as default context strategy, session-start pattern, worktree cleanup pattern.
- **MCP `hawp_search` context mode** — pass `context: true` (and optionally `max_tokens`)
  to receive a single pre-shaped markdown block instead of raw results. Applies the same
  Jaccard dedup (>70% word-set overlap) + greedy token cap + `FormatAsMarkdown()` pipeline
  as the CLI `--context` flag. Response includes `token_count`, `chunks_used`, and
  `chunks_dropped` for observability. Documented in `.hawp/kit/usage/search.md`.
- **48h TTL update notifier** — after every non-MCP/update/version command, the CLI checks
  `~/.hawp/cache/update-check.json` (refreshed in a background goroutine when stale) and
  prints a countdown notice to stderr when a newer release is available. Four phases:
  - Phase 1 (0–15 min): "auto-update in N min"
  - Phase 2 (15–20 min): escalated notice, remaining minutes shown
  - Phase 3 (20–21 min): last-chance notice mentioning `--disable-auto`
  - Phase 4 (≥21 min): announces 3-second countdown, pauses, self-replaces binary
- **`hawp update --disable-auto` / `--enable-auto`** — persists auto-update preference to
  `~/.hawp/config/update.json`. When disabled, notices still print (Phases 1–3) but Phase 4
  skips the install and shows "auto-update is disabled" instead.
- **`make install` target in `librarian/src/Makefile`** — builds the binary then in-place
  overwrites `.hawp/bin/hawp` via `dd` (preserves inode, required on macOS for sandbox ACLs).
- **TS script deprecation** — `@deprecated` JSDoc added to the five TypeScript hawp scripts
  in `librarian/scripts/hawp/`. All are superseded by the Go CLI; npm scripts call the
  binary directly. The TS scripts remain for any tooling that imports them.
- **Shell wrapper cleanup** — `core/.hawp/bin/hawp` replaced with a lean Go-binary delegator
  that resolves the binary from the same directory, then falls back to PATH.

### Fixed

- **Context config `engine` key** — `EmbeddingsConfig.Engine` and `LLMConfig.Engine`
  now use `json:"engine"` as the canonical JSON key (was `json:"backend"`).
- **`--context` dedup was a silent no-op** — `DeduplicateResults` used empty embedding
  vectors (always `[]float32{}`), producing cosine similarity = 0 for every pair, so
  nothing was ever dropped. Replaced with `ContentJaccardDedup` (word-set Jaccard > 0.70):
  drops near-duplicate chunks before ContextBlock assembly. Measured ~30% token reduction
  on queries returning 3+ near-duplicate results.
- **Dynamic chunk cap for `--context`** — greedy selection stops when the next chunk
  would exceed `--max-tokens`, preventing silent over-budget results.

### Removed

- **`--llm-reshape` CLI flag** — removed from `hawp search`. The underlying RAG
  pipeline and `ContextReshaper` remain in the codebase; only the CLI surface is
  gone. Use `--context` to pack results into LLM-ready context blocks.
  `context-reshaping.md` and `troubleshooting.md` removed from docs accordingly.

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
