# Raw Output — HAWP (Layer 1 + Layer 2 shape)

The agent received the filled HAWP shape below (Review Tasks pattern, with the Layer-1 "Coverage and Balance" additions **and** the new Layer-2 "Sweep the scoped surface" habit), followed by the captured output.

## The filled shape

**input:** "Take a look at the librarian/ folder and tell me what's wrong with it. What else should we clean up?"

**context:** "The librarian/ folder contains the repo's TypeScript maintenance tooling. There is substantial uncommitted work in progress restructuring these scripts (visible in git status). This is a read-only review; make no changes."

**mission:** "Produce a bounded hygiene-and-correctness review of librarian/ through a maintainability and truth-risk lens, to support a decision about which cleanup items to add to the backlog next."

**constraints:** "Scope: librarian/ only. Evidence bar: every substantive finding separates directly observed evidence from inference and carries a confidence label (Confirmed / Likely / Unclear), matching the weakest part of the claim. Use only the four finding categories (truth-risk drift, validation drift, maintainability drift, standard mismatch). Call something a standard violation only when backed by explicit repo docs, an ADR, a tooling contract, or a clearly intentional repeated convention. When claiming absence, use bounded wording. No architecture redesign. The 'what else should we clean up?' part is a scope-creep trap: keep everything inside librarian/; out-of-scope items go under 'Out of scope, flagged only', max 3 one-liners. **Tier the findings (Layer 1):** cap *primary* findings at 7, but include an uncapped 'Minor / deferred' one-liner list. **State what is correct (Layer 1):** include a compact 'Verified correct' list. **Sweep the scoped surface (Layer 2):** before finalizing, enumerate every in-scope unit — each domain folder and its index/cli/script files, lib/, the config files (tsconfig.json, package.json), and the contract docs (scripts/README.md, CLI.md) — and confirm each was inspected. Make the sweep visible as a coverage checklist; use bounded-absence wording for anything not inspected. A unit producing no finding is recorded as verified-correct or a non-finding, not as silence."

**output:** "A prioritized review artifact: (1) scope + method + a visible surface-sweep checklist of the in-scope units; (2) up to 7 primary findings, each with category, confidence, observed evidence, inference/uncertainty, significance; (3) an uncapped 'Minor / deferred' one-liner list; (4) a compact 'Verified correct' list and/or non-findings; (5) 'Out of scope, flagged only' (max 3); (6) closing Fix first / Verify next / Defer."

---

## Captured output

# Librarian/ Hygiene & Correctness Review

**Scope:** `librarian/` only. **Method:** read every domain entry point, `lib/`, the validations, and the contract docs; ran `npm run typecheck` (pass), `npm test` (37/37 pass), `npm run providers:validate` and `npm run distribution:validate` (both pass); measured `node`/`.nvmrc`/`engines` and `npx tsx` resolution at the repo root and inside `librarian/`.

### Surface sweep (coverage checklist)

| In-scope unit | Inspected | Result |
| --- | --- | --- |
| `backlog-upgrade/` (index, cli, script, detection/, models/, output/) | yes | F4, F5(minor); `.js` imports (minor) |
| `backlog-validate/` (index, cli, script) | yes | F1, F6 |
| `validate-hawp-workflow/` (index, cli, orchestrate, reporter, validations/) | yes | F3 |
| `distribution/` (build, validate, shared/composition) | yes | F2, path-leak skip (minor) |
| `providers/materialize/` (build, validate, composition, render) | yes | F2 |
| `lib/` | yes | verified correct |
| `tsconfig.json` | yes | F7 |
| `package.json` | yes | F4, alias (minor), engines (minor) |
| `scripts/README.md`, `backlog-upgrade/CLI.md` | yes | F4, doc mismatch (minor) |

---

## Top findings (max 7)

### F1 — `--strict-warnings` depends on scraping a child process's stdout format
**Category:** truth-risk drift · **Confidence:** Confirmed mechanism, Likely impact

- **Observed:** `backlog-validate/script.ts` extracts the warning count with `/! Warnings:\s+(\d+)/`, returning `0` on no match (lines 43–51); that feeds `failedByWarnings = strictWarnings && warningsCount > 0` (line 105). The matching line is produced by `validate-hawp-workflow/reporter.ts`.
- **Inference:** Any cosmetic change to the reporter's summary line silently zeroes the count, so `--strict-warnings` passes when warnings exist — fail-open. No inspected test covers this cross-process coupling.
- **Significance:** A validation gate that can silently degrade to PASS is the exact truth-risk this tooling exists to prevent.

### F2 — Cross-domain imports and three separate repo-root finders break the boundary rule
**Category:** standard mismatch · **Confidence:** Confirmed

- **Observed:** `scripts/README.md`: "No script imports another domain's internals — shared code lives in `lib/`." `providers/materialize/{build,validate}/index.ts` + the render test import `findRepoRoot` from `../../../distribution/shared/composition`; `backlog-upgrade/script.ts` imports `orchestrateValidation`/`parseBacklog` from `../validate-hawp-workflow/orchestrate`. Separately there are **three** upward repo-root finders: `lib/findBacklogRepoRoot`, `distribution/shared/composition.findRepoRoot`, and `validate-hawp-workflow/cli.ts findWorkDirectory`.
- **Inference:** Repo-root finding is exactly the kind of helper `lib/` is meant to own (the README lists "upward repo-root finders" there). Consolidating into `lib/` would remove the cross-domain pull and the triplication.
- **Significance:** The boundary rule is new in this restructure and already broken in two directions, with duplicated logic that will drift.

### F3 — `validate-hawp-workflow/cli.ts` performs filesystem work, violating the cli.ts contract
**Category:** standard mismatch · **Confidence:** Confirmed

- **Observed:** `scripts/README.md`: "`cli.ts` never reads or writes files." But `validate-hawp-workflow/cli.ts` defines `findWorkDirectory` (lines 98–103, `existsSync` + `process.cwd()` walk) and `resolveWorkDirectory` (lines 108–121, `existsSync`/`resolve`).
- **Inference:** Path resolution and existence checks are logic that belongs in `script.ts` (or `lib/`), not the argument adapter. No exception is documented.
- **Significance:** The boundary the restructure is establishing is contradicted in the validator's own adapter — the layer most likely to be copied as a template.

### F4 — Stale hardcoded CLI metadata, version mismatch, and a silent `--export-plan` apply gap
**Category:** truth-risk drift · **Confidence:** Confirmed

- **Observed:** `backlog-upgrade/cli.ts` hardcodes `"v1.1.0 (TASK-028 dry-run detection)"` (line 235) and a `STATUS:` block naming TASK-028 (lines 184–186); `package.json` version is `0.0.0`. `script.ts` apply mode returns at lines 324–330 without reading `options.exportPlan` (export only happens on the dry-run path, lines 355–361); `CLI.md` documents `--export-plan` with no mode caveat.
- **Inference:** Version string and task IDs are leftover work-item baggage in user-facing output; the apply-mode export-plan gap may be intentional but is undocumented.
- **Significance:** These are exactly the user-facing claims the CLI contract exists to keep honest; a flag documented without limits that no-ops silently misleads.

### F5 — `hawp backlog validate` never runs provider-pack drift checks
**Category:** validation drift · **Confidence:** Confirmed

- **Observed:** `backlog-validate/script.ts` spawns only `distribution/validate` and `validate-hawp-workflow` (lines 62–97). It never runs `providers:validate`, even though `distribution:sync` chains through `providers:sync` and provider packs under `core/providers/` are generated from shared behaviors.
- **Inference:** Provider-pack drift (materialized rules out of sync with shared sources) can pass `hawp backlog validate` while failing `providers:validate`. The combined-validate command's coverage is narrower than its name implies.
- **Significance:** A validator advertised as "combined kit + work validation" that omits a whole generated surface gives false assurance.

### F6 — Child validators run via `npx tsx` from a root with no pinned manifest
**Category:** truth-risk drift · **Confidence:** Confirmed mechanism, Unclear impact

- **Observed:** `backlog-validate/script.ts` spawns `npx tsx <abs path>` with `cwd` = repo root (lines 62–97); the repo root has no `package.json`/`node_modules`. Measured now: `npx --no-install tsx` resolves the same runtime at root and inside `librarian/` (no divergence on this machine).
- **Inference:** Resolution is environment-dependent; a clean machine with no cached tsx may fetch an unpinned latest at validation time. The 2026-06-11 divergence is not reproducing, so current impact is Unclear.
- **Significance:** Validation outcomes can become machine-dependent; running children through librarian-local tsx would restore determinism.

### F7 — Dead `tsconfig` emit configuration
**Category:** maintainability drift · **Confidence:** Confirmed

- **Observed:** `tsconfig.json` sets a full emit profile — `declaration`, `declarationMap`, `sourceMap`, `inlineSources`, `outDir: "build-src-esm"`, `removeComments`, `noEmit: false`, `noEmitOnError` — but nothing emits: only `tsx` runs the scripts and `typecheck` is `tsc --noEmit` (which overrides `noEmit: false`). No `build-src-esm/` exists; `exclude` lists `build-*-cjs`, `coverage`, `test`, `examples`, none of which exist.
- **Inference:** Leftover boilerplate from a build setup that was never used; trimming to a minimal `noEmit` config removes the dead surface and the misleading `exclude` entries.
- **Significance:** Low risk, but it is config that claims a build pipeline that does not exist.

---

## Minor / deferred (uncapped one-liners)

- `backlog-upgrade/` uses 17 `.js`-suffixed relative imports against the documented "extensionless" convention, including mixed style within `script.ts` (subset of the migration that left this domain behind).
- `engines.node: ">=26"` pins the toolchain to a just-released major; consistent with `.nvmrc`/runtime but undocumented as a deliberate floor — adoption risk.
- `validate:workflow` is a pure alias of `workflow:validate`; `validate:backlog` vs `workflow:validate` naming schemes coexist — pick one.
- `backlog-upgrade/models/` exports test-only type guards (`isBlockedItem`, `isBacklogFixPlan`, `isDetectionReport`) referenced only by `models/__tests__/types.test.ts`, not by production code.
- `distribution/{build,validate}/index.ts` use bare `console.log`/`console.error` rather than the documented `[<domain>] warning:` stderr convention.
- `distribution/validate/index.ts` silently `continue`s past missing path-leak target files (lines 27–31) with no `[distribution] warning:`, against the "never swallow" rule — fails silent on a future rename.
- `scripts/README.md` documents test discovery as `scripts/**/*.test.ts` (line 51) while `package.json` uses `find scripts -type f -name '*.test.ts'` — doc/reality mismatch.
- `normalizeClosedRecord` writes literal `_Legacy normalization scaffold added._` placeholder text into closed records on apply — confirm intended user-facing content.
- `backlog-validate/script.ts` `runCommand` maps a `null` spawn status to exit `1` (fail-closed, good) but does not surface spawn errors with a `[backlog-validate] warning:` line.
- `composition.ts` `LEGACY_ROOT_GUIDES` guards against root-level guides already removed (self-described legacy) — retire once confidence holds.

## Verified correct (checked, found sound)

- `tsc --noEmit` passes; `npm test` passes 37/37 via the `find`-based script — the 2026-06-11 broken-`npm test`-on-Node-20 finding is resolved (`engines`, `.nvmrc`, runtime now agree on Node 26).
- `npm run providers:validate` and `npm run distribution:validate` both pass — materialized packs and generated guides are current.
- `lib/index.ts` is dependency-free (node builtins only), matching the `lib/` contract.
- Apply-mode mutation guard fails closed: `hasDirtyWorkingTree` returns `true` when `git` fails (script.ts 175–188); `--force-dirty` is the documented override.
- Evidence-link containment works: `evidence-integrity` rejects escaping links with a `[validate] warning:` (observed in test output).
- The three-file boundary pattern (`index`/`cli`/`script`) is followed structurally across the CLI-shaped domains (the F3 exception aside).

## Out of scope, flagged only (max 3)

- `.hawp/bin/hawp` wrapper shares the `npx tsx` resolution path from F6.
- Repo root has no `package.json`/`node_modules` — the root cause of F6's non-determinism.
- CI workflow Node version should be confirmed to match the `>=26` engines floor.

---

## Operational sequence

**Fix first**
- F1: replace stdout scraping with a structured signal so `--strict-warnings` cannot fail open.
- F5: add `providers:validate` to `hawp backlog validate` (or rename/clarify the command's documented scope).

**Verify next**
- F2/F3: decide whether the repo-root finders consolidate into `lib/` and move `findWorkDirectory`/`resolveWorkDirectory` out of `cli.ts`; update `scripts/README.md` or the code so rule and reality match.
- F4 (export-plan): confirm whether `--apply --export-plan` is a no-op; document in `CLI.md` and help, or wire it up.
- F6: confirm clean-environment `npx tsx` resolution, then pin the child runtime to librarian-local tsx.

**Defer**
- F7: trim the dead `tsconfig` emit config to a minimal `noEmit` profile.
- F4 (metadata): derive `--version` from `package.json`; drop the hardcoded `v1.1.0`/TASK-028 strings.
- Minor-list items as backlog cleanup once the restructure is committed.
