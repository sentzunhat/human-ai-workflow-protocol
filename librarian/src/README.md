# hawp (librarian/src/)

Go source for the `hawp` CLI — the HAWP intelligence tool. Distribution
tooling (provider materialization, generated guides) remains in `librarian/`.

Roadmap:

- `db init`
- `work/` + `kit/` ingest
- lexical search
- later vector search, local context building, and prompt handoff

Assets are not embedded: `init` downloads the ONNX Runtime and models into
`~/.hawp/` (see work item `e98de8c4`), and the binary self-updates from
published releases (see work item `cdcf9f78`).

Current status:

- All read/mutate/composite workflow commands ported from TypeScript
  (uuid, links check, kit/work validate + normalize, check)
- `hawp init` provisions ONNX Runtime + an embedding model into `~/.hawp/`
- `hawp model pull` / `hawp embed` run real Hugging Face ONNX models
  locally via [hugot](https://github.com/knights-analytics/hugot)'s
  pure-Go backend — no cgo, no external ONNX Runtime binary needed for
  this path (see "Local models" below)
- `hawp update`/`hawp version` are implemented against real GitHub
  Releases API semantics
- no database implementation yet (`hawp db init` is still a layout-planning
  scaffold; real indexing is `fbf12a93`)

## Command taxonomy

The full `hawp` command surface, mapped across the port phases (see
`.hawp/work/notes/2026/07/20/librarian-ts-to-go-port-plan.md`) and the
intelligence lane work items. Aliases keep the existing `.hawp/bin/hawp`
wrapper surface stable.

| Command | Source today | Status | Lane / work item |
| --- | --- | --- | --- |
| `hawp db init` | Go scaffold | available (scaffold, layout-planning only) | superseded by `hawp init` for real provisioning |
| `hawp index build [--scope] [--export]` | Go, real | **available** — enriches kit/work docs with folder role + backlog metadata, exports JSON | fbf12a93 Slice 1; ported 2026-07-21 (`f93bee55`) |
| `hawp uuid [--short]` | npm `uuid` | **available** | ported 2026-07-20 (`39bc92b6`) |
| `hawp links check` | npm `check:markdown-links` | **available** (fixes dead TS link regex) | ported 2026-07-20 (`39bc92b6`) |
| `hawp kit validate` | npm `kit:validate` | **available** | ported 2026-07-20 (`39bc92b6`) |
| `hawp work validate` | npm `work:validate` | **available** (count parity verified) | ported 2026-07-20 (`39bc92b6`) |
| `hawp kit normalize [--apply]` | npm `kit:normalize` | **available** | ported 2026-07-20 (`eddd8339`) |
| `hawp work normalize [--dry-run --apply --validate]` | npm `work:normalize` | **available** (full flag set, clean-tree guard) | ported 2026-07-20 (`eddd8339`) |
| `hawp backlog upgrade` | wrapper alias | **available** | alias for `work normalize` |
| `hawp check` | npm `hawp:check` | **available** | composite kit + work + links; ported 2026-07-20 |
| `hawp backlog validate` | wrapper alias | **available** | alias for `check` |
| `hawp init` | — | **available** | provision `~/.hawp` (ONNX Runtime + embedding model); ported 2026-07-20 (`e98de8c4`); supersedes `db init` for the intelligence lane |
| `hawp search <query>` | — | planned | lexical then vector search (`fbf12a93`) |
| `hawp context <query>` | — | planned | compact context packs for prompt handoff (`fbf12a93`) |
| `hawp update [--check --provider --skip-kit]` | — | **available** | self-update binary + kit + provider overlay, one unified version; ported 2026-07-20 (`cdcf9f78`), extended 2026-07-21 (`273f3e4b`) |
| `hawp version` | — | **available** | prints the running binary's version; ported 2026-07-20 (`cdcf9f78`) |
| `hawp commands [--json]` | — | **available** | agent-facing command/usage discovery; added 2026-07-21 (`b4c8af81`) |
| `hawp model pull <hf-repo>` | — | **available** | pull any HF ONNX model via hugot; added 2026-07-21 (`748609a8`) |
| `hawp embed <text>...` | — | **available** | local embeddings via hugot's pure-Go backend; added 2026-07-21 (`748609a8`) |
| `hawp generate <prompt>` | — | planned | text generation; needs hugot's ORT/cgo backend, deferred pending a build-matrix decision (`748609a8`) |

Stays npm-only (maintainer tooling, not part of the installable CLI):
`providers:materialize|validate|sync`, `distribution:build|validate|sync`,
`typecheck`, `test`, `validate`. The Node CLI PoC was retired 2026-07-20.

## Build

```bash
cd librarian/src
make build                    # → bin/hawp, version "dev"
make build VERSION=v0.1.0     # → bin/hawp with a real version baked in
make check                    # vet + tests + build
make dist VERSION=v0.1.0      # → bin/dist/hawp-<os>-<arch> for all release platforms
./bin/hawp --help
```

Release platforms: macOS arm64/x64, Linux arm64/x64, Windows x64/arm64.
`bin/` is gitignored — distribution happens via published releases, never
by committing binaries.

## Releasing and self-update

Two ways to cut a release, via
[.github/workflows/release-librarian-go.yml](../../.github/workflows/release-librarian-go.yml):

1. **Manual dispatch (recommended)** — GitHub Actions tab → this
   workflow → "Run workflow" → enter a `version` (e.g. `0.1.0`, no `v`
   prefix). The workflow creates and pushes the plain `<version>` tag
   itself; no local `git tag`/`git push` needed.
2. **Push a tag directly** — `git tag 0.1.0 && git push origin 0.1.0`
   also triggers the workflow.

Either way, `make dist` builds all six platform binaries, the workflow
also packages `.hawp/kit/` + `core/providers/manifest.yaml` +
`core/providers/.{claude,codex,cursor,continue,github}/` into one
`hawp-kit-bundle.tar.gz` asset, and publishes all of it as a single
GitHub Release — the CLI, kit, and provider overlays now share one
version. **The release requires a `CHANGELOG.md` section for that exact
version** — the workflow hard-fails before publishing anything if one is
missing, so a release can never go out without notes.

**Publish timing**: by default the release goes live immediately
(auto-approved) — no extra step needed. To require review first, tick
the `draft` input when dispatching manually: the release lands as an
unpublished draft, invisible to `hawp update` (GitHub hides drafts from
the unauthenticated API calls `hawp update` makes — confirmed empirically
on this repo, including a real workflow-created draft that a running
binary correctly could not see). A maintainer reviews the built binaries
and changelog on GitHub and clicks "Publish release"; only then does it
become visible. Tag-push triggers have no inputs to set, so they always
auto-publish.

`hawp update --check` compares the running binary's version (`hawp
version`) against the latest *published* release for the current
platform; `hawp update` downloads and verifies that release asset's
SHA-256 (GitHub computes a digest for every uploaded asset
automatically) before atomically replacing the running binary. A
checksum mismatch or download failure leaves the current binary
untouched. Self-replace is not supported on Windows (an executable
cannot be renamed over while running); Windows users replace the binary
manually.

**`hawp update` also refreshes `.hawp/kit/` and the installed provider
overlay** from the same release — one unified version for the CLI and
the protocol content, not two separate update paths. When run inside a
HAWP repo, it downloads and verifies `hawp-kit-bundle.tar.gz`, refreshes
`.hawp/kit/` wholesale, and refreshes only the `update: refresh` files
for the detected provider (per `core/providers/manifest.yaml`) —
`update: skip` files, like a hand-customized `CLAUDE.md`, are never
touched. Auto-detection works for Claude, Cursor, and Continue (their
manifest entries have a file-pattern marker to detect against); Codex
and GitHub overlays have no such marker and need an explicit
`--provider <name>`. Pass `--skip-kit` to update the binary only, or run
outside any HAWP repo (e.g. a globally installed binary) and the sync
step skips itself gracefully rather than erroring.

```bash
hawp update                        # binary + kit + auto-detected provider
hawp update --provider codex       # also sync an overlay that can't auto-detect
hawp update --skip-kit             # binary only, old behavior
```

Current binary proof on this machine (all workflow commands ported,
including init provisioning, update/version, and local models via hugot):

- `bin/hawp` ≈ 16–18 MB stripped per platform (`make dist` output) — up
  from ~6.6 MB before adopting hugot for `hawp embed`/`model pull`. The
  increase comes from hugot's GoMLX/onnx-gomlx dependency tree, pulled
  in even though only the pure-Go feature-extraction path is used. Still
  small relative to embedding an ONNX Runtime binary and models directly
  (tens of MB more), and cross-compiles cleanly with `CGO_ENABLED=0` for
  all six platforms.

## Local models

`hawp embed` and `hawp model pull` wrap
[hugot](https://github.com/knights-analytics/hugot) (Apache 2.0) — the
closest Go equivalent to transformers.js — to run Hugging Face ONNX
models locally. Both use hugot's **pure-Go backend** (`NewGoSession`,
backed by [GoMLX](https://github.com/gomlx/gomlx)'s `simplego`), which
needs no cgo and no external ONNX Runtime binary. Verified empirically
on this repo: `CGO_ENABLED=0 go build`/`go test` succeed with the real
hugot package linked in, and cross-compile cleanly for every `make dist`
platform — adopting hugot does not change the single-Linux-runner
cross-compile setup at all.

```bash
hawp embed "The quick brown fox" "A fast auburn fox"      # default model
hawp embed "some text" --model sentence-transformers/all-mpnet-base-v2
hawp model pull Xenova/all-MiniLM-L6-v2                   # pull only, no inference
```

`hawp embed` defaults to the pinned `Xenova/all-MiniLM-L6-v2` model
(the same one `hawp init` provisions for the future vector-search work,
`fbf12a93`), pulling it into `~/.hawp/models/` on first use if not
already present. Any other repo works via `--model`; `--onnx-file`
disambiguates repos that publish more than one `.onnx` export (hugot's
downloader otherwise errors rather than guessing).

**Text generation is intentionally not wired in yet.** hugot's
text-generation pipeline currently requires its ORT backend (cgo +
libonnxruntime), which would break `make dist`'s single-runner
cross-compile the way the pure-Go embedding path does not — that needs
a real per-OS build matrix decision first, tracked in work item
`748609a8`'s follow-up scope rather than shipped silently.

## Agent usage

`hawp commands --json` is the canonical way for an agent to discover the
full command surface without parsing free-text help — one registry backs
both `--help` and this output, so they cannot drift apart. Each entry has
a `usage` line, one-line `description`, `flags`, and per-command
`exitCodes` string; `hawp commands` (no flag) prints the same data as
human-readable text.

```bash
hawp commands --json | jq '.commands[] | select(.name == "work validate")'
```

Shared exit-code convention across every validating/mutating command
(also returned as `exitCodeConvention` in the JSON output):

- `0` — success / no issues found
- `1` — issues found, or a command-specific failure (e.g. checksum
  mismatch, missing release asset)
- `2` — usage error, or a guarded refusal (e.g. `--apply` on a dirty
  git worktree)

For scripted/agent use, prefer `--format json` where a command offers it
(`work normalize`) and `hawp check`'s exit code over parsing stdout text.
