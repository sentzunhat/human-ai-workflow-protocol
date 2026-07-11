# Raw Output — HAWP (Layer-1 shape)

The agent received the filled HAWP shape below (Review Tasks pattern, with the Layer-1 "Coverage and Balance" additions), followed by the captured output.

## The filled shape

**input:** "Take a look at the librarian/ folder and tell me what's wrong with it. What else should we clean up?"

**context:** "The librarian/ folder contains the repo's TypeScript maintenance tooling. There is substantial uncommitted work in progress restructuring these scripts (visible in git status). This is a read-only review; make no changes."

**mission:** "Produce a bounded hygiene-and-correctness review of librarian/ through a maintainability and truth-risk lens, to support a decision about which cleanup items to add to the backlog next."

**constraints:** "Scope: librarian/ only. Evidence bar: every substantive finding separates directly observed evidence from inference and carries a confidence label (Confirmed / Likely / Unclear), matching the weakest part of the claim. Use only the four finding categories (truth-risk drift, validation drift, maintainability drift, standard mismatch). Call something a standard violation only when backed by explicit repo docs, an ADR, a tooling contract, or a clearly intentional repeated convention. When claiming absence, use bounded wording. No architecture redesign. The 'what else should we clean up?' part is a scope-creep trap: keep everything inside librarian/; out-of-scope items go under 'Out of scope, flagged only', max 3 one-liners. **Tier the findings (Layer 1):** cap *primary* findings (full Observed/Inference/Significance detail) at 7, but include an **uncapped 'Minor / deferred' one-liner list** for valid smaller items so coverage is not suppressed. **State what is correct (Layer 1):** include a compact 'Verified correct' list of areas checked and found sound, in addition to any non-findings, without confirming anything not actually inspected."

**output:** "A prioritized review artifact: (1) scope + method; (2) up to 7 primary findings, each with category, confidence, observed evidence, inference/uncertainty, significance; (3) an uncapped 'Minor / deferred' one-liner list; (4) a compact 'Verified correct' list and/or up to 2 non-findings; (5) 'Out of scope, flagged only' (max 3); (6) closing Fix first / Verify next / Defer."

---

## Captured output

# Librarian/ Hygiene & Correctness Review

**Scope:** `librarian/` only. **Method:** read all domain entry points, `lib/`, the validations, and the contract docs (`scripts/README.md`, `backlog-upgrade/CLI.md`); ran `npm run typecheck` (pass) and `npm test` (pass, 37/37); measured `node`, `.nvmrc`, `engines`, and `npx tsx` resolution at both the repo root and inside `librarian/`.

---

## Top findings (max 7)

### F1 — `--strict-warnings` depends on scraping a child process's stdout format
**Category:** truth-risk drift · **Confidence:** Confirmed mechanism, Likely impact

- **Observed:** `backlog-validate/script.ts` spawns the workflow validator and extracts the warning count with regex `/! Warnings:\s+(\d+)/`, returning `0` when no match (lines 43–51). That line is produced by `validate-hawp-workflow/reporter.ts`; the formats currently align. The `0`-on-no-match feeds `failedByWarnings = strictWarnings && warningsCount > 0` (line 105).
- **Inference:** Any cosmetic change to the reporter's summary line silently zeroes the count, so `--strict-warnings` passes even when warnings exist — a fail-open mode in a validator. No inspected test covers this cross-process coupling.
- **Significance:** A validation gate that can silently degrade to PASS is exactly the truth-risk this tooling exists to prevent.

### F2 — Cross-domain internal imports violate the documented boundary rule
**Category:** standard mismatch · **Confidence:** Confirmed

- **Observed:** `scripts/README.md` states "No script imports another domain's internals — shared code lives in `lib/`." Four import sites break this: `providers/materialize/{build/index.ts, validate/index.ts, __tests__/render.test.ts}` import `findRepoRoot` from `../../../distribution/shared/composition`; `backlog-upgrade/script.ts` imports `orchestrateValidation`/`parseBacklog` from `../validate-hawp-workflow/orchestrate` (lines 15–18).
- **Inference:** `lib/` already owns repo-root finding (`findBacklogRepoRoot`, `findUpward`), and the README lists "upward repo-root finders" as lib contents — so `findRepoRoot` naturally belongs in `lib/`. The backlog-upgrade→validator coupling is deeper (it reuses the validator for `--validate`) and may be intentional, but no exception is documented.
- **Significance:** The boundary rule is part of the new restructure and is already broken in two directions; either the code or the rule should move before the convention loses force.

### F3 — Stale hardcoded CLI metadata, and `--export-plan` is silently ignored in apply mode
**Category:** truth-risk drift · **Confidence:** Confirmed

- **Observed:** `backlog-upgrade/cli.ts` hardcodes `"v1.1.0 (TASK-028 dry-run detection)"` (line 235) and a `STATUS:` block naming TASK-028 (lines 184–186), while `package.json` declares version `0.0.0`. Separately, `script.ts` apply mode returns at lines 324–330 without ever reading `options.exportPlan`; the export-plan write exists only on the dry-run path (lines 355–361). `CLI.md` documents `--export-plan` with no mode caveat.
- **Inference:** The version string and task references are leftover work-item baggage in user-facing output. The apply-mode export-plan gap may be intentional (no plan in apply mode) but is undocumented either way.
- **Significance:** These are exactly the user-facing claims the CLI contract (`CLI.md`) exists to keep honest; a flag documented without limits that no-ops silently misleads.

### F4 — Child validators run via `npx tsx` from a root with no pinned manifest
**Category:** truth-risk drift · **Confidence:** Confirmed mechanism, Unclear impact

- **Observed:** `backlog-validate/script.ts` spawns `npx tsx <abs path>` with `cwd` = repo root (lines 62–97); the repo root has no `package.json`/`node_modules`. Measured now: `npx --no-install tsx --version` resolves the same runtime at root and inside `librarian/` (no divergence observed on this machine).
- **Inference:** Resolution is environment-dependent. On a clean machine with no cached tsx, `npx` may fetch an unpinned latest at validation time, or resolve a global install librarian doesn't control. The divergence seen in the 2026-06-11 run is not reproducing here, so current impact is Unclear.
- **Significance:** Validation outcomes can become machine-dependent; running children through librarian-local tsx (or `node --import` with an explicit path) would restore determinism.

### F5 — `.js`-suffixed relative imports in `backlog-upgrade/` diverge from the documented convention
**Category:** standard mismatch · **Confidence:** Confirmed

- **Observed:** `scripts/README.md` mandates "extensionless relative imports (tsconfig uses bundler resolution)." 17 `.js`-suffixed relative imports remain, all inside `backlog-upgrade/`, including mixed style within single files (`script.ts` line 13 `"./models/index.js"` alongside line 11 `"../lib"`).
- **Inference:** Looks like an incomplete migration of the oldest domain to the new convention, not a deliberate exception — none is documented.
- **Significance:** Low individually; worth normalizing while the restructure is still uncommitted so the convention starts clean.

### F6 — Distribution path-leak validation silently skips missing target files
**Category:** validation drift · **Confidence:** Likely

- **Observed:** `distribution/validate/index.ts` checks a hardcoded list of 5 `core/.hawp/kit/...` files for `core/.hawp/` path leaks and does `if (!existsSync(absolutePath)) continue;` with no warning (lines 27–31). `scripts/README.md` requires "never swallow — surface skipped files/branches with a stderr `[<domain>] warning:` line." All 5 files currently exist (validation passes).
- **Inference:** If a listed kit file is renamed, leak-checking for it silently stops while validation keeps reporting PASS. "Likely" because the harm needs a future rename; the divergence from the stated error convention is direct.
- **Significance:** Low probability, but it fails silent, in a validator.

### F7 — `engines.node: ">=26"` pins the toolchain to a just-released major
**Category:** maintainability drift · **Confidence:** Confirmed requirement, Unclear impact

- **Observed:** `package.json` sets `"engines": { "node": ">=26" }`, `@types/node` `^25.9.3`, and `typescript` `^6.0.3`; repo `.nvmrc` is `26` and the local runtime is `v26.3.0`. All are internally consistent.
- **Inference:** Requiring Node ≥26 (a very recent major) for maintenance tooling raises the floor for contributors and CI runners; whether that is deliberate policy or a side effect of the local dev environment is not documented in inspected files.
- **Significance:** Not a bug — the gates pass — but a high engine floor with no stated rationale is an adoption/maintainability risk worth a deliberate decision.

---

## Minor / deferred (uncapped one-liners)

- `distribution/validate/index.ts` and `build/index.ts` use bare `console.log`/`console.error` rather than the documented `[<domain>] warning:` stderr convention — cosmetic inconsistency with `scripts/README.md`.
- Two repo-root finders coexist (`lib/findBacklogRepoRoot` vs `distribution/shared/composition.findRepoRoot`); consolidating into `lib/` would remove the F2 cross-domain pull.
- `backlog-upgrade/cli.ts` `STATUS:` help block hardcodes phase narrative ("apply mode scaffolds closed records") that will drift as behavior changes.
- `package.json` `version: 0.0.0` while the CLI advertises `v1.1.0` — pick one source of truth (derive `--version` from `package.json`).
- `normalizeClosedRecord` writes literal "_Legacy normalization scaffold added._" placeholders into closed records on apply; confirm that is intended user-facing content, not debug text.
- `backlog-validate/script.ts` `runCommand` maps a `null` spawn status to exit code `1` (fail-closed — good), but swallows spawn errors without a `[backlog-validate] warning:` line.
- Mixed import style within `backlog-upgrade/` (`"./models/index.js"` vs `"../lib"`) is worth a single normalization pass (subset of F5).

## Verified correct (checked, found sound)

- `npm run typecheck` (`tsc --noEmit`) passes clean.
- `npm test` passes 37/37 via the `find`-based test script — the 2026-06-11 "broken `npm test` on Node 20" failure is resolved; `engines`, `.nvmrc`, and the local runtime now agree on Node 26.
- `lib/index.ts` is dependency-free (node builtins only), matching the stated `lib/` contract.
- Apply-mode mutation guard fails closed: `hasDirtyWorkingTree` returns `true` when `git` itself fails (script.ts lines 175–188), and `--force-dirty` is the explicit documented override.
- Evidence-link containment works: `evidence-integrity` rejects links escaping the evidence root with a `[validate] warning:` (observed in test output).
- The three-file boundary pattern (`index.ts` / `cli.ts` / `script.ts`) is followed across the CLI-shaped domains inspected.

## Out of scope, flagged only (max 3)

- `.hawp/bin/hawp` wrapper shares the `npx tsx` resolution path from F4.
- Repo root has no `package.json`/`node_modules`, which is what makes child `npx` resolution non-deterministic (F4).
- CI workflow Node version should be confirmed to match the `>=26` engines floor (F7).

---

## Operational sequence

**Fix first**
- F1: replace stdout scraping with a structured signal (distinct exit code, JSON, or in-process call) so `--strict-warnings` cannot silently fail open.
- F3 (export-plan): decide whether `--apply --export-plan` is a no-op; if yes, document it in `CLI.md` and help text; if no, wire it up.

**Verify next**
- F2: decide whether `findRepoRoot` moves to `lib/` and whether the backlog-upgrade→validator import is a sanctioned exception; update `scripts/README.md` or the code so rule and reality match.
- F4: confirm on a clean environment how `npx tsx` resolves, then pin the child runtime to librarian-local tsx.
- F7: record the Node ≥26 floor as a deliberate decision (or relax it) and confirm CI matches.

**Defer**
- F5: normalize the 17 `.js` imports in `backlog-upgrade/` (mechanical, low risk).
- F3 (metadata): derive `--version` from `package.json`; drop the hardcoded `v1.1.0`/TASK-028 strings.
- F6: add the `[distribution] warning:` line for skipped leak-check targets.
- Minor-list items as backlog cleanup when the restructure is committed.
