# Librarian TypeScript → Go Port Plan

Related work: `.hawp/work/active/56067dc8-1cc6-4dbe-8b6d-dd7ffdfea80c.md`
Builds on: `.hawp/work/notes/2026/07/08/librarian-cli-zacatl-architecture.md`
(command taxonomy, bin contract, layering) — the Node/oclif/SEA packaging
sections of that note are superseded by the Go-first decision; its command
surface and boundary discipline carry forward unchanged.

## Direct Evidence — current TypeScript inventory

From `librarian/package.json` and `librarian/scripts/` (58 source files,
12 test files, all `node:test`):

| TS domain | Size | Kind | Mutating | Tests |
| --- | --- | --- | --- | --- |
| `scripts/hawp/kit-validate/` | 28K | validator | no | yes |
| `scripts/hawp/kit-normalize/` | 24K | normalizer | yes (clean-tree guard) | yes |
| `scripts/hawp/work-validate/` | 88K, 6 validation modules | validator | no | yes |
| `scripts/hawp/work-normalize/` | 144K, largest (detection/rules 516-line evaluator) | normalizer | yes (clean-tree guard) | yes |
| `scripts/hawp/hawp-check/` | 16K | composite (distribution + work validate) | no | yes |
| `scripts/check-markdown-links.mjs` | 4K | validator | no | no |
| `uuid` npm script | one-liner | utility | no | no |
| `scripts/librarian/providers/materialize/` | 24K | build pipeline | yes (writes provider packs) | yes |
| `scripts/librarian/distribution/` | 24K | build pipeline | yes (writes generated docs) | yes |
| `scripts/hawp-cli-poc/` | 40K | Node/oclif/SEA PoC | no | yes |
| `scripts/lib/` | 1 file | shared helpers (repo-root finder, markdown walker, `toRepoRelative`, `normalizeForCompare`, legacy-cutoff constant) | no | — |

Go scaffold today (`librarian/go/`): `db init` and
`index build [--scope]` routes only, 168 lines total, **zero test files**,
module path still `github.com/sentzunhat/hawp/golang`.

## Command mapping — where each TS command goes

Target Go layout follows the existing Zacatl layers
(`application` / `domain` / `infrastructure` / `platform`).

| Today (npm / wrapper) | Go CLI command | Go package home | Phase |
| --- | --- | --- | --- |
| `uuid` | `hawp uuid` | `internal/application/uuidgen` | 1 |
| `check:markdown-links` | `hawp links check` | `internal/application/links` + `internal/infrastructure/markdown` | 1 |
| `kit:validate` / `hawp kit validate` | `hawp kit validate` | `internal/application/kit` + `internal/domain/kit` | 1 |
| `work:validate` / `hawp work validate` | `hawp work validate` | `internal/application/work` + `internal/domain/work` (one file per validation, mirroring the 6 TS validation modules) | 1 |
| `kit:normalize` / `hawp kit normalize` | `hawp kit normalize [--apply]` | `internal/application/kit` (mutations under `internal/domain/kit`) | 2 |
| `work:normalize` / `backlog upgrade` | `hawp work normalize [--apply --validate]` + `backlog upgrade` alias | `internal/application/work` (detection rules under `internal/domain/work/rules`) | 2 |
| `hawp:check` / `backlog validate` | `hawp check` | `internal/application/check` composing kit+work+links | 3 |
| `providers:*` | stays npm (maintainer-only) | — | deferred |
| `distribution:*` | stays npm (maintainer-only); `hawp check` in Go covers only what it can run natively, npm `hawp:check` remains canonical until then | — | deferred |
| `cli:poc*` (oclif + Node SEA) | not ported — superseded by the Go binary; park/remove after Phase 2 parity | — | retired |
| `scripts/lib/` helpers | — | `internal/infrastructure/repo` (root finding, repo-relative paths), `internal/infrastructure/markdown` (tree walk, link parse) | 0 |
| `.hawp/bin/hawp` wrapper | unchanged surface; switches from `tsx` scripts to the Go binary per-subcommand as phases land | — | 1–3 |

Existing Go routes `db init` and `index build` are the future
intelligence lane (`e98de8c4`, `fbf12a93`) and are untouched by this port.

## Inference — port order and rationale

Phases ordered read-only → mutating → composite, so each phase can be
verified against the TS implementation's output before the next lands:

- **Phase 0 — foundations.** Fix module path to
  `github.com/sentzunhat/hawp/librarian/go`; add `go test ./...` harness;
  port `scripts/lib/` helpers with unit tests; add `testdata/` fixture
  repos (a minimal valid `.hawp/` tree and known-broken variants).
- **Phase 1 — read-only commands.** `uuid`, `links check`,
  `kit validate`, `work validate`. Parity check: run TS and Go against
  this repo and diff exit codes + finding counts.
- **Phase 2 — mutating commands.** `kit normalize`, `work normalize`
  with the clean-worktree fail-closed guard and `--dry-run`/`--apply`
  semantics preserved exactly. `work-normalize` is the largest module
  (144K, rule evaluator) — port rules table-driven, one Go test per rule.
- **Phase 3 — composite + switchover.** `hawp check`; point
  `.hawp/bin/hawp` subcommands at the Go binary; retire the Node CLI PoC
  and its CI workflow; npm scripts stay for `providers:*`/`distribution:*`.

## Go unit test standard (applies to every phase)

- Standard library `testing` only; table-driven tests; no assertion deps.
- Every ported validation/mutation rule gets at least one passing and one
  failing fixture case — mirror what the 12 TS test files cover before
  deleting any TS.
- Fixtures live in `testdata/` per package (Go convention, excluded from
  builds); golden files for validator report output.
- CLI routing tests at `internal/platform/cli`: help, aliases, unknown
  command, exit codes (mirrors the architecture note's test list).
- CI: extend validation to run `go vet ./...` and `go test ./...` in
  `librarian/go/` alongside the npm suite; TS tests are deleted only when
  the Go equivalent covers the same behavior (per-phase, not big-bang).

## Verification plan (for the implementing items)

- `go build ./cmd/hawp && go vet ./... && go test ./...` in `librarian/go/`
- Per-phase parity diff: TS vs Go exit codes and finding counts on this repo
- `npm --prefix librarian run validate` stays green throughout (TS remains
  canonical until its phase completes)

## Open decisions

- Whether `hawp check` in Go eventually shells out to npm for
  distribution checks or waits until providers/distribution are ported.
- When to remove `@oclif/core` and the PoC packaging scripts (proposed:
  end of Phase 2).
- Whether `links check` output format copies the TS one-line summary or
  adopts the richer validator report style.
