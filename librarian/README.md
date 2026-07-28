# librarian

HAWP maintenance and documentation generation tooling.

## Purpose

The `librarian` folder is the single librarian home. It holds two lanes:

- **npm/TypeScript tooling** (this level) — the current working implementation of packaging, validation, and workflow integrity checks.
- **Go implementation** ([go/](go/README.md)) — the future installable `hawp` CLI. The TypeScript commands are being ported there phase by phase (see `.hawp/work/notes/2026/07/20/librarian-ts-to-go-port-plan.md`); `make build` in `go/` produces the distributable `go/bin/hawp` binary. The full command taxonomy lives in [go/README.md](go/README.md).

The npm scripts remain canonical until their Go equivalent lands with matching unit tests.

## Runtime

The librarian tooling targets Node 26 and npm 11.17.0. Use `nvm` before running commands here:

```bash
nvm install
nvm use
```

The repository also supports other Node 26 installs as long as they satisfy `.nvmrc`, `engines.node`, and `engines.npm`.

## Script Organization

Layout rules, the CLI boundary pattern (`index.ts` / `cli.ts` / `script.ts`), and how to add new scripts are documented in [scripts/README.md](scripts/README.md).

### Distribution Tools

Scripts for building and validating user-facing installation and update guides.

**Location:** `scripts/librarian/distribution/`

- **`build/index.ts`** — Composes distribution guide variants by combining source fragments with bash script extraction
  - Reads source fragments from `distribution/sources/` (shared + `providers/<name>/`)
  - Composes bash from `script-core` + `providers/<name>/script-install|update` + `script-footer`
  - Applies branch-specific (`main` or `dev`) and provider-specific (`PROVIDER`) substitution
  - Outputs twenty guides: `claude`, `codex`, `github`, `cursor`, and `continue` × install/update × main/dev

- **`validate/index.ts`** — Verifies generated outputs match expected composition
  - Compares actual files against source fragment lists
  - Reports missing or stale generated files
  - Exit code: 0 (success), 1 (validation failed)

### Provider Materialization

Materializes shared behaviors into provider pack files.

**Location:** `scripts/librarian/providers/materialize/`

- **`build/index.ts`** — Renders `core/providers/shared/behaviors/*.md` into generated provider overlay files
- **`validate/index.ts`** — Fails when materialized files drift from shared sources
- **`composition.ts`** — Emit map (frontmatter + output paths per provider)

Run `npm run providers:sync` before `distribution:sync` (wired automatically).

### Workflow Validation Tools

Scripts for HAWP backlog and work item integrity verification.

**Location:** `scripts/hawp/work-validate/`

- **`index.ts`** — CLI entry point and orchestration
- **`cli.ts`** — Command-line argument parsing
- **`reporter.ts`** — Structured output formatting (plain text)
- **`orchestrate.ts`** — Validation pipeline coordination
- **`types.ts`** — Shared type definitions

**Validations** (`validations/`):

- `backlog-consistency.ts` — Ensures BACKLOG.md rows match actual plan files
- `id-parser.ts` — Extracts and validates work item IDs (extensible for UUID support)
- `closed-task-completeness.ts` — Checks closed plans include required sections (with legacy tolerance)
- `evidence-integrity.ts` — Validates evidence file references and completeness
- `verification-clarity.ts` — Ensures verification sections avoid ambiguous or unproven claims

## Common Tasks

### Regenerate distribution guides

When shared behaviors, provider sources, or distribution fragments change:

```bash
cd librarian
npm run providers:sync
npm run distribution:sync
```

(`distribution:sync` runs `providers:sync` first.)

### Add or update a provider

Provider work has three source layers plus generated outputs:

1. Add or update the source pack under `core/providers/.<provider>/`.
2. Register install targets in `core/providers/manifest.yaml` and `core/providers/README.md`.
3. Add source fragments under `distribution/sources/providers/<provider>/`.
4. Add the provider to `ACTIVE_PROVIDERS` in `scripts/librarian/distribution/shared/composition.ts`.
5. If shared behaviors should materialize into provider files, add targets in `scripts/librarian/providers/materialize/composition.ts`.
6. Update indexes such as `distribution/sources/providers/README.md` and `distribution/generated/README.md`.
7. Run `npm --prefix librarian run distribution:sync` from the repo root.

Provider guides install the kit plus exactly one provider overlay. They do not install `librarian/` into target repos.

### Validate workflow records

Check backlog consistency and task completeness:

```bash
npm --prefix librarian run work:validate
```

### Validate the kit

Check `.hawp/kit/` naming, required files, and internal links:

```bash
npm --prefix librarian run kit:validate
```

### Normalize workflow records

Dry-run the raw normalize CLI:

```bash
npm --prefix librarian run work:normalize:cli -- --dry-run --validate
```

Apply closed-record normalization and run validation:

```bash
npm --prefix librarian run work:normalize
```

### Normalize the kit

Apply kit file renames and internal-link updates:

```bash
npm --prefix librarian run kit:normalize -- --apply
```

### Type checking

Verify TypeScript code:

```bash
npm --prefix librarian run typecheck
```

### Run unit tests

```bash
cd librarian
npm test
```

Tests use Node's built-in `node:test` runner via `tsx`; no extra test framework is required.

### Full validation

Run the common maintenance checks in one pass:

```bash
cd librarian
npm run validate
```

## CLI

`./.hawp/bin/hawp` is a thin bash wrapper. It prefers the compiled Go binary
(`librarian/go/bin/hawp`, build with `cd librarian/go && make build`) and
falls back to these npm scripts when the binary is absent:

```bash
./.hawp/bin/hawp kit validate [options]      # Go; fallback scripts/hawp/kit-validate/
./.hawp/bin/hawp kit normalize [options]     # Go; fallback scripts/hawp/kit-normalize/
./.hawp/bin/hawp work validate [options]     # Go; fallback scripts/hawp/work-validate/
./.hawp/bin/hawp work normalize [options]    # Go; fallback scripts/hawp/work-normalize/ raw CLI
./.hawp/bin/hawp check                       # Go composite: kit + work + links
./.hawp/bin/hawp links check                 # Go markdown link check
./.hawp/bin/hawp backlog upgrade [options]   # alias for work normalize
./.hawp/bin/hawp backlog validate [options]  # alias for check (npm fallback: work:validate)
./.hawp/bin/hawp uuid [--short]              # generate a work item UUID
```

`npm --prefix librarian run work:normalize` is the apply-and-validate shortcut.
Use `./.hawp/bin/hawp work normalize ...` or
`npm --prefix librarian run work:normalize:cli -- ...` when you need the raw
CLI modes such as `--dry-run`, `--format json`, `--export-plan`, or
`--export-research-queue`.

See `scripts/hawp/work-normalize/CLI.md` for the full work-normalize command contract. The wrapper resolves paths relative to this repository, so it only works in the HAWP source repo (not in downstream installs).

## Future Extensions

The librarian architecture supports:

- **CI/automation integration** via structured JSON output
- **UUID-based work item migration** (extensible id-parser)
- **Multi-root validation** for monorepo support

## Go binary

The installable CLI is the Go implementation in [go/](go/README.md).
`make build` produces `go/bin/hawp`; `make dist` cross-compiles all release
platforms. The former Node/oclif/SEA PoC (`scripts/hawp-cli-poc/`) and its
CI workflow were retired 2026-07-20 after the Go port reached mutation
parity (see `.hawp/work/closed/2026/07/20/`).
