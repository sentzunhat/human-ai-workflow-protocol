# Raw Output — No HAWP

Prompt given to the agent (verbatim, no framing added):

> Take a look at the librarian/ folder and tell me what's wrong with it. What else should we clean up?

- **Arm:** no-HAWP
- **Model:** Cursor agent (model not recorded)
- **Workspace:** clean `/tmp/hawp-benchmark-clean-20260615-222558` (no `AGENTS.md`, no `.cursor/rules/`)

Output captured below. Content is unchanged; the only edit is that editor-specific code citation fences were converted to standard code blocks with the file and line noted underneath, so this file renders anywhere.

---

I reviewed `librarian/` end to end — READMEs, `package.json`, `tsconfig.json`, all five script domains, `.hawp/bin/hawp`, and CI — and ran `npm run typecheck`, `npm test`, `npm run workflow:validate`, `npm run distribution:validate`, and `npm run providers:validate`. Typecheck and all 37 unit tests pass on Node v26.3.0. Here are the issues, ordered by severity.

## Violations of the folder's own conventions (`scripts/README.md`)

**1. Cross-domain internal import.**

The README states "No script imports another domain's internals — shared code lives in `lib/`", but `providers/materialize` imports from `distribution` in three places:

```ts
import { findRepoRoot } from "../../../distribution/shared/composition";
```

(`librarian/scripts/providers/materialize/validate/index.ts`, line 2; also `build/index.ts` and `__tests__/render.test.ts`)

`findRepoRoot` is a generic repo-root finder and belongs in `lib/`, not in `distribution/shared/composition.ts`.

**2. `cli.ts` doing filesystem work.**

The rule is "`cli.ts` never reads or writes files", yet `validate-hawp-workflow/cli.ts` contains `findWorkDirectory` / `resolveWorkDirectory` with `existsSync` calls and a `process.cwd()` walk (lines 98–121). That resolution logic belongs in `script.ts`.

**3. `cli.ts` containing execution logic.**

`backlog-upgrade/cli.ts`'s `runCLI` invokes `runBacklogUpgradeScript` and builds `repoRoot: process.cwd()` — the README says cli.ts is "argument parsing, help text; no execution logic".

**4. Inconsistent import style.**

The convention is "extensionless relative imports", which every domain follows except `backlog-upgrade/`, where all files use `.js` extensions (`from "./models/index.js"` etc.).

**5. Boundary pattern not followed uniformly.**

The README says `process.exit` / `process.exitCode` only in `index.ts`, but `distribution/build/index.ts`, `distribution/validate/index.ts`, `providers/materialize/build/index.ts`, and `providers/materialize/validate/index.ts` all call `process.exit()` directly inside their `main()` functions. Build/validate pipelines are documented as a separate pattern, but the "only index.ts" rule is still stated globally and is contradicted in four entry points.

## Fragile design worth cleaning up

**6. `backlog-validate` shells out and scrapes stdout.**

Instead of importing the sibling validators, `backlog-validate/script.ts` spawns them via `npx tsx` child processes and implements `--strict-warnings` by regex-matching the other tool's report text:

```ts
const extractWarningsCount = (workflowStdout: string): number => {
  const match = workflowStdout.match(/! Warnings:\s+(\d+)/);
  if (!match?.[1]) {
    return 0;
  }

  const parsed = Number.parseInt(match[1], 10);
  return Number.isNaN(parsed) ? 0 : parsed;
};
```

(`librarian/scripts/backlog-validate/script.ts`, lines 43–51)

A wording change in `reporter.ts` silently disables strict mode. The spawn runs with `cwd: repoRoot`, where there is no local `node_modules` — it works only because `npx` resolves tsx from cache, not because of the version pinned in `librarian/package.json`. Importing the validator functions directly would remove all of this.

**7. `hawp backlog validate` does not check provider materialization drift.**

`backlog-validate/script.ts` runs `distribution/validate` and `validate-hawp-workflow` only. It never runs `providers:validate`. Provider pack files under `core/providers/` can drift from shared behavior sources without failing `hawp backlog validate`, even though `distribution:sync` chains through `providers:sync`. The help text says "distribution validation (.hawp/kit generated drift checks)" which underspecifies what is actually checked.

**8. `.hawp/bin/hawp` argument inconsistency.**

`backlog validate` does `shift 2` before forwarding args, but `backlog upgrade` forwards `"$@"` including the `backlog upgrade` tokens — which is why `backlog-upgrade/cli.ts` has to accept `backlog`/`upgrade` anywhere in argv and `parseArgs` errors without them, and why the npm script must pass redundant tokens: `tsx scripts/backlog-upgrade/index.ts backlog upgrade --apply --validate`.

**9. Three separate repo-root finders with different markers.**

- `findBacklogRepoRoot` in `lib/` — looks for `.hawp/work/BACKLOG.md`, throws on miss
- `findRepoRoot` in `distribution/shared/composition.ts` — looks for `distribution/sources/shared` + `librarian/package.json`, returns null
- `findWorkDirectory` in `validate-hawp-workflow/cli.ts` — looks for `.hawp/work`, returns null

These invite drift. Consolidate the generic upward-walk helpers in `lib/` with configurable predicates.

## Dead weight and stale surface area

**10. Stale task-ID droppings in user-facing CLI output.**

```text
STATUS:
  dry-run detection pipeline enabled (TASK-028)
```

(`librarian/scripts/backlog-upgrade/cli.ts`, lines 184–185)

plus a hardcoded version string `"hawp backlog upgrade v1.1.0 (TASK-028 dry-run detection)"` that contradicts `package.json`'s `0.0.0`. Internal task IDs should not leak into help/version text.

**11. Dead tsconfig boilerplate.**

The config is fully set up to emit (`declaration`, `declarationMap`, `sourceMap`, `inlineSources`, `outDir: "build-src-esm"`, `removeComments`, `noEmitOnError`) but nothing ever emits — only `tsx` and `tsc --noEmit` are used, and no `build-src-esm/` directory exists. The `exclude` list references nonexistent folders (`build-*-cjs`, `coverage`, `test`, `examples`). Could be trimmed to a minimal noEmit config.

**12. Test-only type guards in production model files.**

`isBlockedItem`, `isBacklogFixPlan`, and `isDetectionReport` in `backlog-upgrade/models/` are exported and referenced only by `models/__tests__/types.test.ts`, not by production code. The models layer (~700 lines across five files) is reasonable for the domain, but the runtime type guards are speculative scaffolding unless something actually consumes them at the boundary.

**13. package.json nits.**

`validate:workflow` is a pure alias of `workflow:validate` (pick one naming scheme — currently `validate:backlog` vs `workflow:validate` coexist); empty `"dependencies": {}` block; `engines.node` is `>=26` while the root `.nvmrc` says `26` — fine for CI, but the README never mentions a minimum Node version and the test command relies on `find` + `node --test`, which is less portable than documenting the requirement explicitly.

**14. Smaller candidates.**

`LEGACY_ROOT_GUIDES` in `composition.ts` guards against root-level guide files that were already removed (self-described `@deprecated` — retire once confidence is there); no integration tests for the four build/validate entry-point scripts (only composition/render unit tests cover the logic); the README's "Future Extensions" section is speculative; `scripts/README.md` still describes test discovery as `scripts/**/*.test.ts` glob while `package.json` uses `find scripts -type f -name '*.test.ts' | sort`.

## What else should we clean up? (repo-wide, triggered by the question)

These are outside `librarian/` but surfaced while exercising its scripts:

- **Provider pack drift:** `npm run providers:validate` fails when materialized files under `core/providers/` do not match shared behavior sources. That is the expected guard — run `npm run providers:sync` after editing shared behaviors. Worth wiring into `hawp backlog validate` or documenting as a separate required step.
- **CI coverage gap:** `.github/workflows/librarian-quality.yml` runs typecheck, unit tests, and `workflow:validate` but not `distribution:validate`, `providers:validate`, or `validate:backlog`. Generated-guide and provider-pack drift can land on `main`/`dev` without CI catching it unless something else runs those scripts.
- **Benchmark / temp workspace hygiene:** if you use `prepare-clean-workspace.sh`, remember it strips agent guidance from the copy only; the real repo keeps `AGENTS.md` and rules intact.

## What's fine

Typecheck is clean, all 37 tests pass on Node 26, `node_modules` is gitignored, no stray build artifacts, `lib/` is genuinely dependency-free, the docs (`README.md`, `scripts/README.md`, `CLI.md`) largely match the actual layout, and the CI workflow correctly reads `.nvmrc` with npm caching. The domain split (backlog-upgrade, backlog-validate, distribution, providers/materialize, validate-hawp-workflow) is sound. Most cleanup is about the code honoring the rules `scripts/README.md` declares, consolidating shared helpers into `lib/`, replacing subprocess orchestration with direct imports, and tightening CI to cover the full sync/validate chain.
