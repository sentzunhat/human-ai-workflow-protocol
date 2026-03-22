# Scripts Guidelines

How `librarian/scripts/` is organized and how to add new scripts. Follows the repo-wide clean-code rules in `.hawp/kit/instructions/clean-code-and-structure.md` (one responsibility per file, split only with a reason).

## Domain map

Scripts are grouped by owner: `hawp/` holds HAWP-workflow tooling, `librarian/` holds librarian build pipelines. Each domain folder has a single purpose. No script imports another domain's internals — shared code lives in `lib/` (flat, used by both groups).

| Folder | Purpose | Entry point | npm script / CLI |
|---|---|---|---|
| `hawp/kit-validate/` | Validate `.hawp/kit/` structure (naming, required files, links) | `index.ts` | `kit:validate` · `kit validate` |
| `hawp/kit-normalize/` | Detect and normalize `.hawp/kit/` drift | `index.ts` | `kit:normalize` · `kit normalize` |
| `hawp/work-validate/` | Backlog/plan/evidence integrity checks | `index.ts` | `work:validate` · `work validate` |
| `hawp/work-normalize/` | Detect and normalize work record drift | `index.ts` | `work:normalize` · `work normalize` · `backlog upgrade` |
| `hawp/hawp-check/` | Combined distribution + work:validate in one command | `index.ts` | `hawp:check` · `backlog validate` |
| `librarian/distribution/` | Build/validate generated install/update guides | `build/index.ts`, `validate/index.ts` | `distribution:build`, `distribution:validate`, `distribution:sync` |
| `librarian/providers/materialize/` | Materialize shared behaviors into provider packs | `build/index.ts`, `validate/index.ts` | `providers:materialize`, `providers:validate`, `providers:sync` |
| `lib/` | Shared utilities (node builtins only) | — (library) | — |

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
4. If it's a user-facing command, wire it as a subcommand in `.hawp/bin/hawp`.
5. Add `__tests__/<name>.test.ts` — `npm test` picks up `scripts/**/*.test.ts` automatically.
6. Validate: `npm run typecheck && npm test`.

## Conventions

- Imports: `node:` prefix for builtins; extensionless relative imports (tsconfig uses bundler resolution).
- Errors: never swallow — surface skipped files/branches with a stderr `[<domain>] warning:` line.
- Paths from untrusted content (markdown links, backlog rows) must be containment-checked before filesystem access.
- Mutating scripts must guard with a clean-working-tree check that fails closed (see `hawp/work-normalize/script.ts`).
