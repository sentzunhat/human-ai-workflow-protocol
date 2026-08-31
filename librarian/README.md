# librarian

HAWP maintenance and documentation generation tooling.

## Purpose

The `librarian` folder now holds two clear lanes:

- **Go CLI** ([src/](src/README.md)) — the canonical `hawp` runtime for user-facing workflow commands.
- **npm/TypeScript maintainer tooling** (this level) — retained only for distribution guide generation, provider materialization, and their supporting tests.

The former TypeScript workflow command tree under `scripts/hawp/` was retired on 2026-08-31 after the Go CLI reached parity for `uuid`, `links check`, `kit/work validate`, `kit/work normalize`, and `check`.

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

### Workflow Commands

HAWP workflow validation and normalization now live in the Go CLI under
`librarian/src`. Use [src/README.md](src/README.md) for the command taxonomy
and verification contract.

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

Dry-run the Go normalize path:

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

Inside this source repository, the npm wrappers call `go run ./cmd/hawp` from
`librarian/src` so maintainer checks always exercise the current Go source
instead of depending on a separately installed binary. The external
`./.hawp/bin/hawp` entrypoint remains the user-facing wrapper.

```bash
./.hawp/bin/hawp kit validate [options]
./.hawp/bin/hawp kit normalize [options]
./.hawp/bin/hawp work validate [options]
./.hawp/bin/hawp work normalize [options]
./.hawp/bin/hawp check                       # Go composite: kit + work + links
./.hawp/bin/hawp links check                 # Go markdown link check
./.hawp/bin/hawp backlog upgrade [options]   # alias for work normalize
./.hawp/bin/hawp backlog validate [options]  # alias for check
./.hawp/bin/hawp uuid [--short]              # generate a work item UUID
```

Use either `./.hawp/bin/hawp work normalize ...` or
`npm --prefix librarian run work:normalize:cli -- ...` for modes such as
`--dry-run`, `--format json`, `--export-plan`, or `--export-research-queue`.

## Future Extensions

The librarian architecture supports:

- **CI/automation integration** via structured JSON output
- **UUID-based work item migration** (extensible id-parser)
- **Multi-root validation** for monorepo support

## Go binary

The installable CLI is the Go implementation in [src/](src/README.md).
`make build` in `librarian/src` produces `src/bin/hawp`; `make dist`
cross-compiles all release
platforms. The former Node/oclif/SEA PoC (`scripts/hawp-cli-poc/`) and its
CI workflow were retired 2026-07-20 after the Go port reached mutation
parity (see `.hawp/work/closed/2026/07/20/`).
