# Scripts Guidelines

How `librarian/scripts/` is organized for the retained Node maintainer
pipelines. Follows the repo-wide clean-code rules in
`.hawp/kit/instructions/clean-code-and-structure.md` (one responsibility per
file, split only with a reason).

## Domain map

The active Node script surface is now just `librarian/` for release
maintenance. The canonical user-facing workflow path is the Go CLI under
`librarian/src`. Shared helpers live in `lib/`.

| Folder | Purpose | Entry point | npm script / CLI |
|---|---|---|---|
| `librarian/distribution/` | Build/validate generated install/update guides | `build/index.ts`, `validate/index.ts` | `distribution:build`, `distribution:validate`, `distribution:sync` |
| `librarian/providers/materialize/` | Materialize shared behaviors into provider packs | `build/index.ts`, `validate/index.ts` | `providers:materialize`, `providers:validate`, `providers:sync` |
| `lib/` | Shared utilities (node builtins only) | — (library) | — |

## Current ownership

- Use the Go CLI for all user-facing HAWP workflow commands (`uuid`, `links
  check`, `kit/work validate`, `kit/work normalize`, `check`).
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
- Larger domains can add purpose folders below that when they have a clear
  reason.

Build/validate pipelines (`librarian/distribution/`, `librarian/providers/materialize/`) are simple single-use scripts: a `build/index.ts` and a `validate/index.ts` per pipeline, with shared composition logic in a sibling module (`shared/composition.ts`, `composition.ts` + `render.ts`).

## Shared `lib/`

- Keep it dependency-free (node builtins only) so any domain can use it without coupling.
- Current contents: upward repo-root finders, markdown tree walker, `toRepoRelative`, `normalizeForCompare`, legacy-cutoff constant.
- Add here only when two or more domains need the same helper — don't pre-abstract.

## Adding a new script

1. Create a new domain folder (kebab-case, single purpose) under
   `librarian/` for Node maintainer pipelines — don't grow an existing one
   past its responsibility.
2. Use the boundary pattern above; start with `index.ts` + `script.ts`, add `cli.ts` when flags appear.
3. Add an npm script in `librarian/package.json` (`tsx scripts/<group>/<domain>/index.ts`).
4. If it's a user-facing workflow command, add it to the Go CLI instead of
   creating a new TypeScript surface. Only add a TypeScript entry point when
   the command is intentionally Node-only maintainer tooling.
5. Add `__tests__/<name>.test.ts` — `npm test` picks up
   `scripts/librarian/**/*.test.ts` automatically.
6. Validate with the runtime that owns the surface:
   - Go workflow commands: `go test ./...` and the relevant `hawp` command
   - Node maintainer pipelines: `npm run typecheck && npm test`

## Conventions

- Imports: `node:` prefix for builtins; extensionless relative imports (tsconfig uses bundler resolution).
- Errors: never swallow — surface skipped files/branches with a stderr `[<domain>] warning:` line.
- Paths from untrusted content (markdown links, backlog rows) must be containment-checked before filesystem access.
- Mutating maintainer scripts should fail closed when they would overwrite
  generated outputs or drift from the declared source state.
