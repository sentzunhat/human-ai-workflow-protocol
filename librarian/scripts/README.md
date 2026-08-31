# Scripts Guidelines

How `librarian/scripts/` is organized and how to add new scripts. Follows the repo-wide clean-code rules in `.hawp/kit/instructions/clean-code-and-structure.md` (one responsibility per file, split only with a reason).

## Domain map

Scripts are grouped by owner: `hawp/` now holds deprecated legacy workflow
implementations, while `librarian/` holds the active Node-based build
pipelines that still power distribution/provider maintenance. The canonical
user-facing workflow path is the Go CLI under `librarian/src`; npm wrappers
call that binary directly. No script imports another domain's internals —
shared code lives in `lib/` (flat, used by both groups).

| Folder | Purpose | Entry point | npm script / CLI |
|---|---|---|---|
| `hawp/kit-validate/` | Deprecated legacy implementation of `.hawp/kit/` validation | `index.ts` | superseded by `hawp kit validate` |
| `hawp/kit-normalize/` | Deprecated legacy implementation of `.hawp/kit/` normalization | `index.ts` | superseded by `hawp kit normalize` |
| `hawp/work-validate/` | Deprecated legacy implementation of work-record validation | `index.ts` | superseded by `hawp work validate` |
| `hawp/work-normalize/` | Deprecated legacy implementation of work-record normalization | `index.ts` | superseded by `hawp work normalize` |
| `hawp/hawp-check/` | Deprecated legacy implementation of the composite workflow check | `index.ts` | superseded by `hawp check` |
| `librarian/distribution/` | Build/validate generated install/update guides | `build/index.ts`, `validate/index.ts` | `distribution:build`, `distribution:validate`, `distribution:sync` |
| `librarian/providers/materialize/` | Materialize shared behaviors into provider packs | `build/index.ts`, `validate/index.ts` | `providers:materialize`, `providers:validate`, `providers:sync` |
| `lib/` | Shared utilities (node builtins only) | — (library) | — |

## Current ownership

- Use the Go CLI for all user-facing HAWP workflow commands (`uuid`, `links
  check`, `kit/work validate`, `kit/work normalize`, `check`).
- Keep `librarian/scripts/hawp/` only as transitional source for import
  compatibility and deletion planning; do not add new workflow features there.
- Keep `librarian/scripts/librarian/` and `lib/` for Node-based maintainer
  pipelines that are still active release tooling.

## Boundary pattern (CLI-shaped domains)

Domains invoked as commands follow a three-file boundary so logic stays testable without process assumptions:

```text
<domain>/
├── index.ts    # executable boundary: argv in, process.exitCode out — nothing else
├── cli.ts      # adapter: argument parsing, help text; no execution logic
├── script.ts   # logic: pure-ish functions returning structured results/exit codes
└── __tests__/  # node:test files (*.test.ts), run by `npm test`
```

Rules:

- `process.exit` / `process.exitCode` only in `index.ts`.
- `cli.ts` never reads or writes files.
- `script.ts` (and deeper modules) never reads `process.argv`.
- Larger domains add purpose folders below that (e.g. `hawp/work-normalize/detection/`, `models/`, `output/`; `hawp/work-validate/validations/`; `hawp/kit-normalize/mutations/`).

Build/validate pipelines (`librarian/distribution/`, `librarian/providers/materialize/`) are simple single-use scripts: a `build/index.ts` and a `validate/index.ts` per pipeline, with shared composition logic in a sibling module (`shared/composition.ts`, `composition.ts` + `render.ts`).

## Shared `lib/`

- Keep it dependency-free (node builtins only) so any domain can use it without coupling.
- Current contents: upward repo-root finders, markdown tree walker, `toRepoRelative`, `normalizeForCompare`, legacy-cutoff constant.
- Add here only when two or more domains need the same helper — don't pre-abstract.

## Adding a new script

1. Create a new domain folder (kebab-case, single purpose) under `hawp/` (workflow tooling) or `librarian/` (build pipelines) — don't grow an existing one past its responsibility.
2. Use the boundary pattern above; start with `index.ts` + `script.ts`, add `cli.ts` when flags appear.
3. Add an npm script in `librarian/package.json` (`tsx scripts/<group>/<domain>/index.ts`).
4. If it's a user-facing workflow command, add it to the Go CLI instead of
   creating a new `hawp/` TypeScript surface. Only add a TypeScript entry point
   when the command is intentionally Node-only maintainer tooling.
5. Add `__tests__/<name>.test.ts` — `npm test` picks up `scripts/**/*.test.ts` automatically.
6. Validate with the runtime that owns the surface:
   - Go workflow commands: `go test ./...` and the relevant `hawp` command
   - Node maintainer pipelines: `npm run typecheck && npm test`

## Conventions

- Imports: `node:` prefix for builtins; extensionless relative imports (tsconfig uses bundler resolution).
- Errors: never swallow — surface skipped files/branches with a stderr `[<domain>] warning:` line.
- Paths from untrusted content (markdown links, backlog rows) must be containment-checked before filesystem access.
- Mutating scripts must guard with a clean-working-tree check that fails closed (see `hawp/work-normalize/script.ts`).
