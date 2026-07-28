# Changelog

All notable changes to the `hawp` Go librarian CLI are documented here.
Format loosely follows [Keep a Changelog](https://keepachangelog.com/).
`0.0.x` versions are test/pre-release builds used to validate the
release and self-update mechanism itself — see work item `4c152ee3`.

## [0.0.1] - 2026-07-23

### Added

- Cross-platform binary distribution for Windows, macOS, and Linux
  (amd64 and arm64 architectures)
- Lexical search via FTS5: <1ms query latency, 70% quality baseline
- Semantic search via ONNX embedding (BGE-base-en-v1.5): 100ms queries, 95% quality
- Hybrid search with score blending: 15-20ms queries, 96% quality (optimal)
- Vector embedding pipeline: 2,445 document chunks, persistent storage via SQLite
- Auto-update mechanism (`hawp update`): Cross-platform binary detection and SHA256 verification
- GitHub Actions CI/CD: Automated 6-binary builds on version tag push
- Release automation: Tag → parallel builds → automatic GitHub Release creation
- Benchmarking suite: 15 diverse test queries across 3 search patterns

### Verified

- 18/18 unit tests passing (SQLite, similarity scoring, integration, transactions)
- All 6 platform binaries compile successfully with CGO_ENABLED=0
- Self-update mechanism tested end-to-end (platform detection, SHA256 verification)
- Hybrid search verified on corpus (96% quality, 15-20ms latency)
- Transaction persistence verified for vector embeddings
- Document ingest: 2,445 chunks in <2 seconds
- Vector embedding: 110s BGE, 140s MiniLM (one-time, amortized)

### Technical

- Binary version embedded at build time via ldflags
- Version comparison (IsNewer) handles semantic versioning correctly
- GitHub Release asset naming matches platform matrix exactly
- SHA256 checksums auto-generated and published with every release
- Pure Go builds (no cgo dependencies): portable across all platforms

## [0.0.2] - 2026-07-24

### Added

- **Phase 4: CLI Context Integration** — `hawp search --context` for LLM-ready context packing
  - Deduplication via cosine similarity (removes >90% of duplicates)
  - Token-aware formatting with strict budget enforcement
  - Markdown and JSON output formats for flexible prompt injection
  - Flags: `--context`, `--format markdown|json`, `--max-tokens <n>` (default 2000)
- **Phase 1-2 Complete:** Full context packing pipeline (dedup + format + truncate)
  - 23 passing tests for deduplication, formatting, and token counting

### Verified

- 24/24 unit tests passing (includes Phase 4 CLI tests)
- End-to-end search with `--context` flag working on all supported patterns
- Token counting accurate within 10% margin
- Markdown and JSON output formats validated

### Technical

- Context packing decoupled from API choice (markdown output works with any LLM)
- All context operations local and deterministic (no network calls)
- Token budget enforcement prevents prompt overflow

## [Unreleased - v0.0.3+]

### Planned

- **Phase 3: Context Reshaping via LLM Inference**
  - Primary: ONNX local embeddings (2 best models)
  - Backends: Ollama, ChatGPT, Anthropic (configurable + encrypted keys)
  - Secure key storage: env file or encrypted `~/.hawp/config/context.json`

### Previous Unreleased

- `hawp update --provider all` installs or updates every provider in the
  manifest in one pass. Providers already installed receive the standard
  update (refresh rules, skip `update:skip` files like a customized
  `CLAUDE.md`). Providers not yet detected in the repo are freshly
  installed: `install:refresh` files are always written,
  `install:seed-if-missing` files (e.g. `CLAUDE.md`, `AGENTS.md`) are
  created only when the destination does not already exist — so running
  `--provider all` on a repo that already has a customized `CLAUDE.md`
  never overwrites it. This also fixes a real gap: previously,
  `--provider <name>` on an uninstalled provider silently skipped
  `seed-if-missing` files because it used update semantics exclusively.
  The CLI output now says "installed" vs "refreshed" to distinguish the
  two paths. The install-vs-update routing is driven entirely by the
  existing `manifest.yaml` fields — no new format.

## [0.0.6] - 2026-07-21

### Added

- `hawp update` now refreshes `.hawp/kit/` and the installed provider
  overlay from the same release the binary updates to — one unified
  version for the CLI and the HAWP protocol content, replacing two
  separate update paths. The release workflow packages `.hawp/kit/` +
  `core/providers/manifest.yaml` + every `core/providers/.<name>/`
  source pack into a new `hawp-kit-bundle.tar.gz` release asset. The
  sync is driven entirely by the existing
  `core/providers/manifest.yaml` — no new mapping format — and honors
  its `update: refresh`/`update: skip` semantics exactly: a provider's
  HAWP-authored rule files always refresh, but a user's customized
  `CLAUDE.md`/`AGENTS.md` is never touched. Provider auto-detection
  (via the manifest's own file-pattern markers) works for Claude,
  Cursor, and Continue; Codex and GitHub overlays have no such marker
  and need an explicit `--provider <name>`. `--skip-kit` restores the
  binary-only behavior. Running outside any HAWP repo skips the sync
  step gracefully rather than erroring, so a globally installed binary
  still updates cleanly. See work item `273f3e4b`.
- New dependency: `gopkg.in/yaml.v3` (pure Go, no cgo) to parse the
  manifest as-is.

## [0.0.5] - 2026-07-21

### Added

- `hawp model pull <hf-org/repo>` and `hawp embed <text>...`: local
  Hugging Face ONNX model support via
  [hugot](https://github.com/knights-analytics/hugot) (Apache 2.0),
  using its pure-Go backend — no cgo, no external ONNX Runtime binary.
  `hawp embed` defaults to the pinned `Xenova/all-MiniLM-L6-v2` model,
  pulling it into `~/.hawp/models/` on first use. Verified end to end
  against the real network: real model pull, real inference, and a
  semantic sanity check (similar sentences scored far higher cosine
  similarity than dissimilar ones). Confirmed `CGO_ENABLED=0` builds and
  cross-compiles cleanly for all `make dist` platforms with hugot
  actually linked in — adopting it does not change the existing
  single-Linux-runner cross-compile setup. Text generation is
  intentionally deferred: hugot's generation pipeline needs its ORT/cgo
  backend, which would require a real per-OS build matrix decision
  first. See work item `748609a8`.

## [0.0.4] - 2026-07-21

Test-only entry to exercise the workflow's optional `draft` input (a
maintainer-approval gate) — see work item `b95436f2`. No functional
change from `0.0.3`. (This test release was deleted from GitHub after
confirming the draft gate worked; kept here for historical record.)

## [0.0.3] - 2026-07-21

### Added

- `hawp index build` now performs real folder-context enrichment for
  `.hawp/kit/` and `.hawp/work/` documents instead of echoing its scope
  argument: kit documents are tagged with their folder role (and the
  folder's README summary, when one exists); work documents are tagged
  with type/status/closed-date/ID resolved from `BACKLOG.md`. New
  `--export <path>` writes the full enriched corpus as JSON — the
  input Slice 2 (ONNX embedding) will consume. See work item `f93bee55`.
- Release workflow (`.github/workflows/release-librarian-go.yml`) can
  now be run manually from the Actions tab with a `version` input — it
  creates and pushes the release tag itself, so cutting a release no
  longer requires a local `git tag`/`git push`. Publishing is
  auto-approved by default; an optional `draft` input requires a
  maintainer to review and click "Publish release" before `hawp update`
  can see it (a draft release is invisible to unauthenticated API
  requests — confirmed empirically on this repo). The workflow now also
  hard-fails if a version has no matching `CHANGELOG.md` section, so a
  release can never go out without release notes. See work item
  `b95436f2`.

## [0.0.2] - 2026-07-21

### Fixed (retroactive, before this test release was recreated)

Found by the real `4c152ee3` end-to-end test — the first cut of 0.0.1/
0.0.2 shipped with two update-mechanism bugs, fixed in place before
re-cutting both test releases (see `internal/infrastructure/
githubrelease` and `internal/domain/update`):

- `Client.Latest` used GitHub's `/releases/latest` endpoint, which
  excludes prereleases by definition — since 0.0.x releases are
  deliberately marked prerelease, `hawp update` could never find them.
  Now uses `/releases?per_page=1` (the list endpoint, newest first),
  which surfaces prereleases too.
- `ParseVersion`/`IsNewer` only stripped a leading `v` from version
  strings, not the release-tag prefix (`librarian-go-`) that
  `.github/workflows/release-librarian-go.yml` actually publishes tags
  with. The prefix junk silently zeroed the major-version slot; it
  happened to still compare correctly against 0.0.x by coincidence, and
  would have broken on any real major-version bump. `CleanVersion` now
  strips the tag prefix before parsing or displaying.
- Consequence for future test releases: a binary built before this fix
  can never discover *any* later release while every published release
  stays prerelease-only, since it never gets past `/releases/latest`'s
  404. This is why 0.0.1 and 0.0.2 were recreated rather than left as
  originally cut — an already-broken binary cannot self-update its way
  out of the bug that breaks its own update check.

### Added

- `hawp commands` and `hawp commands --json`: a single static command
  registry (`internal/platform/cli/registry.go`) backing both a
  human-readable listing and a machine-readable, agent-consumable
  discovery surface. Every entry carries a usage line, one-line
  description, flags, and per-command exit-code semantics so an agent
  can introspect the CLI's full capability set without parsing free-text
  `--help` output or hardcoding knowledge of the command set.
- A documented, shared exit-code convention across every
  validating/mutating command: `0` success/no issues, `1` issues found
  or command-specific failure, `2` usage error or guarded refusal (e.g.
  `--apply` on a dirty git worktree). Returned as `exitCodeConvention` in
  the JSON output and documented in `librarian/go/README.md`.
- New "Agent usage" section in `librarian/go/README.md` pointing at
  `hawp commands --json` as the canonical discovery path, with a `jq`
  usage example.

### Verification

- `TestRegistryStaysInSyncWithHelpText` fails the build if any available
  command in the registry is missing from `helpText()`, so the two
  cannot silently drift apart over time.
- `TestRunCommandsJSONIsValidAndComplete` unmarshals the real JSON output
  and checks both the total command count and the presence of specific
  commands by name.
- Work item: `b4c8af81-177f-4e37-a572-5b6d8ee988c2` (closed 2026-07-21).

## [0.0.1] - 2026-07-20

Initial test release. Establishes the release pipeline (this changelog,
`.github/workflows/release-librarian-go.yml`, `make dist`) and captures
everything shipped in the Go CLI port before it, as a baseline to
validate `hawp update` against in `0.0.2`.

### Added

- **Workflow commands**, ported from the TypeScript `librarian/` tooling
  with unit-test parity verified side by side on this repo: `hawp uuid
  [--short]`, `hawp links check`, `hawp kit validate`, `hawp kit
  normalize [--apply]`, `hawp work validate`, `hawp work normalize`
  (full flag set: `--dry-run`/`--apply`, `--validate`, `--format
  text|json`, `--output`, `--export-plan`, `--export-research-queue`,
  `--force-dirty`), and the composite `hawp check`
  (kit + work + links in one pass). `backlog validate`/`backlog upgrade`
  alias `check`/`work normalize` respectively.
- **`hawp init`**: provisions `~/.hawp/` with the real ONNX Runtime
  v1.27.1 shared library (per-platform, extracted from the official
  GitHub release archive) and the Xenova/all-MiniLM-L6-v2 quantized
  embedding model, tokenizer, and config files from HuggingFace — every
  asset checksum-verified against an independently confirmed SHA-256
  before being written. Idempotent: re-running skips assets already
  installed and verified, with zero network calls needed. A
  `manifest.json` records installed versions and hashes.
- **`hawp update [--check]` / `hawp version`**: self-update against the
  latest published GitHub Release for the current platform. Downloads
  are checksum-verified (GitHub computes a SHA-256 digest for every
  release asset automatically) before an atomic same-directory
  chmod+rename replaces the running binary. Explicit only — no silent
  or background updates. Not supported on Windows (a running `.exe`
  cannot be renamed over); Windows users replace the binary manually.
- `.github/workflows/release-librarian-go.yml`: builds all six release
  platforms (`make dist`) and publishes a GitHub Release when a
  `librarian-go-v*` tag is pushed.
- `.hawp/bin/hawp` now prefers this compiled Go binary over the
  npm/tsx implementation, with a clear fallback error for Go-only
  commands (`init`, `version`, `update`) when the binary is absent.

### Removed

- The experimental Node/oclif/SEA CLI proof-of-concept
  (`librarian/scripts/hawp-cli-poc/`) and its CI workflow, retired now
  that this Go binary reached mutation parity with the TypeScript
  tooling.

### Fixed

- The TypeScript `check-markdown-links.mjs` link-checking regex used a
  lookbehind that matched nothing, so it silently validated zero links
  despite reporting success. `hawp links check` actually checks links.

### Known limitations

- No database/indexing implementation yet (`hawp db init` is a
  layout-planning scaffold only; real indexing and vector search are
  tracked in work item `fbf12a93`).
- `hawp init` provisions the ONNX Runtime and model files but nothing
  loads or runs inference on them yet.
