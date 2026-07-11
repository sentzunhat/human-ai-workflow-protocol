# Raw Output — No HAWP

Prompt given to the agent (verbatim, no framing added):

> Repository: /home/beltrd/Desktop/projects/human-ai-workflow-protocol
>
> Take a look at the librarian/ folder and tell me what's wrong with it. What else should we clean up?
>
> Return your complete findings as your final response.

Output captured below. Content is unchanged; the only edit is that editor-specific code citation fences were converted to standard code blocks with the file and line noted underneath, so this file renders anywhere.

---

I audited `librarian/` (structure, configs, conventions doc, all five script domains) and ran `typecheck` (passes) and `npm test` (fails). Here are the findings, ordered by severity.

## Broken right now

**1. `npm test` fails on the documented minimum Node version — and on this machine.**

```json
    "test": "node --test --import tsx \"scripts/**/*.test.ts\"",
```

(`librarian/package.json`, line 10)

The quoted glob relies on `node --test` glob expansion, which only exists in Node ≥22. But:

- `package.json` `engines` says `">=20"`
- `.nvmrc` says `22`
- The locally installed/default Node is **v20.20.2** (Node 22 isn't installed at all)

Result: `npm test` exits 1 locally with `Could not find '.../scripts/**/*.test.ts'`. CI passes only because `librarian-quality.yml` reads `.nvmrc`. Fix: bump `engines` to `>=22` (and install Node 22 locally), or make the test command Node-20-compatible.

## Violations of the folder's own conventions (`scripts/README.md`)

**2. Cross-domain internal import.** The README states "No script imports another domain's internals — shared code lives in `lib/`", but `providers/materialize` imports from `distribution` in three places:

```ts
import { findRepoRoot } from "../../../distribution/shared/composition";
```

(`librarian/scripts/providers/materialize/validate/index.ts`, line 2; also `build/index.ts` and `__tests__/render.test.ts`). `findRepoRoot` should move to `lib/`.

**3. `cli.ts` doing filesystem work.** The rule is "`cli.ts` never reads or writes files", yet `validate-hawp-workflow/cli.ts` contains `findWorkDirectory`/`resolveWorkDirectory` with `existsSync` calls and a `process.cwd()` read (lines 98–121). That resolution logic belongs in `script.ts`.

**4. `cli.ts` containing execution logic.** `backlog-upgrade/cli.ts`'s `runCLI` invokes `runBacklogUpgradeScript` and builds `repoRoot: process.cwd()` — the README says cli.ts is "argument parsing, help text; no execution logic".

**5. Inconsistent import style.** The convention is "extensionless relative imports", which every domain follows except `backlog-upgrade/`, where all ~15 files use `.js` extensions (`from "./models/index.js"` etc.).

## Fragile design worth cleaning up

**6. `backlog-validate` shells out and scrapes stdout.** Instead of importing the sibling validators, `backlog-validate/script.ts` spawns them via `npx tsx` child processes and implements `--strict-warnings` by regex-matching the other tool's report text:

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

A wording change in `reporter.ts` silently disables strict mode. Also, the `npx tsx` spawn runs with `cwd: repoRoot`, where there's no `node_modules` — it currently works only because tsx happens to be in the npx cache, not because of the version pinned in `librarian/package.json`. Importing the validator functions directly would remove all of this.

**7. `.hawp/bin/hawp` argument inconsistency.** `backlog validate` does `shift 2` before forwarding args, but `backlog upgrade` forwards `"$@"` including the `backlog upgrade` tokens — which is why `backlog-upgrade/cli.ts` has to accept `backlog`/`upgrade` anywhere in argv and `parseArgs` errors without them, and why the npm script must pass redundant tokens: `tsx scripts/backlog-upgrade/index.ts backlog upgrade --apply --validate`.

## Dead weight

**8. ~300 lines of unused "future API" model code.** `models/evidence-report.ts` (264 lines — `EvidenceReport`, `createEvidenceReport`, `createFileHashRecord`, `createValidatorStateSnapshot`, `assessValidatorImprovement`, `isEvidenceReport`) is referenced by **no production code**, only its own tests. The `is*` type guards in the other model files are similarly test-only. This is speculative "JSON-first / future APIs/UIs/agents" scaffolding — the kind of premature abstraction the repo's own HAWP rules discourage.

**9. Stale task-ID droppings in user-facing CLI output.**

```text
STATUS:
  dry-run detection pipeline enabled (TASK-028)
```

(`librarian/scripts/backlog-upgrade/cli.ts`, lines 184–185)

plus a hardcoded version string `"hawp backlog upgrade v1.1.0 (TASK-028 dry-run detection)"` that contradicts `package.json`'s `0.0.0`. Internal task IDs shouldn't leak into help/version text.

**10. Dead tsconfig boilerplate.** The config is fully set up to emit (`declaration`, `declarationMap`, `sourceMap`, `inlineSources`, `outDir: "build-src-esm"`, `removeComments`, `noEmitOnError`) but nothing ever emits — only `tsx` and `tsc --noEmit` are used, and no `build-src-esm/` exists. The `exclude` list references nonexistent folders (`build-*-cjs`, `coverage`, `test`, `examples`). Could be trimmed to ~15 lines.

**11. package.json nits.** `validate:workflow` is a pure alias of `workflow:validate` (pick one naming scheme — currently `validate:backlog` vs `workflow:validate` coexist); empty `"dependencies": {}` block.

**12. Smaller candidates.** `LEGACY_ROOT_GUIDES` in `composition.ts` guards against root-level guide files that were already removed (self-described `@deprecated` — retire once confidence is there); three separate upward root-finders (`findBacklogRepoRoot` in lib, `findRepoRoot` in distribution, `findWorkDirectory` in workflow cli) use different markers but invite drift — worth consolidating in `lib/`; the README's "Future Extensions" section is speculative.

## What's fine

Typecheck is clean, `node_modules` is gitignored, no stray build artifacts, `lib/` is genuinely dependency-free, the docs (`README.md`, `scripts/README.md`, `CLI.md`) match the actual layout, and the new CI workflow correctly uses `.nvmrc` with npm caching. The structure itself is sound — most of the cleanup is about the code actually honoring the rules `scripts/README.md` declares, plus deleting the unused model layer and fixing the Node-version mismatch that breaks `npm test` locally.
