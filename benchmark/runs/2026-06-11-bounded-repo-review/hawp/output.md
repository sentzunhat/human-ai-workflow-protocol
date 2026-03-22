# Raw Output — HAWP

The agent received the filled HAWP shape below as its prompt (authored using the Review Tasks pattern in `.hawp/kit/authoring-patterns.md`), followed by the captured output.

## The filled shape

**input:** "Take a look at the librarian/ folder and tell me what's wrong with it. What else should we clean up?"

**context:** "The librarian/ folder contains the repo's TypeScript maintenance tooling (backlog validation and upgrade, distribution build/validate, provider materialization, HAWP workflow validation). There is substantial uncommitted work in progress restructuring these scripts (visible in git status). This is a read-only review; make no changes."

**mission:** "Produce a bounded hygiene-and-correctness review of librarian/ through a maintainability and truth-risk lens, to support a decision about which cleanup items (if any) to add to the backlog next."

**constraints:** "Scope: librarian/ only. Do not review .hawp/, distribution/, core/, or CI workflows except where a librarian/ file directly references them and the reference itself is the issue. Evidence bar: every substantive finding must separate directly observed evidence from inference, and carry a confidence label — Confirmed (directly proven by inspected files or commands), Likely (evidence-supported but interpretive), or Unclear (plausible but unproven). Confidence must match the weakest part of the claim. Use only these four finding categories: truth-risk drift, validation drift, maintainability drift, standard mismatch. Call something a standard violation only when backed by explicit repo documentation, an ADR, a tooling contract, or a clearly intentional repeated convention; otherwise use softer language such as pattern divergence or unclear if intentional. When claiming absence, use bounded wording (e.g. 'no evidence found in inspected files'). No architecture redesign proposals; prefer cleanup, clarification, verification, small alignment. The 'what else should we clean up?' part of the request is a scope-creep trap: do not expand beyond librarian/. If something outside scope seems important, list it under 'Out of scope, flagged only' as one line each, maximum 3 entries. If evidence is insufficient to decide on an item, say that explicitly rather than recommending."

**output:** "A prioritized review artifact, not a generic summary: (1) top findings, maximum 7, each with category, confidence, directly observed evidence, inference/uncertainty, and significance; (2) optionally up to 2 non-findings (concerns that were checked and resolved, with what was checked and what was observed); (3) a closing operational sequence with exactly three parts — Fix first, Verify next, Defer. Compact and decision-useful."

---

## Captured output

# Librarian/ Hygiene & Correctness Review

Scope: `librarian/` only. Method: read all entry points, shared lib, validations, contracts (`scripts/README.md`, `CLI.md`); ran `npm run typecheck` (pass), `npm test` (fail, see F1), and direct test execution (38/38 pass).

---

## Top findings

### F1 — `npm test` is broken on the minimum supported Node version
**Category:** validation drift · **Confidence:** Confirmed

- **Observed:** `librarian/package.json` declares `"engines": { "node": ">=20" }` and `"test": "node --test --import tsx \"scripts/**/*.test.ts\""`. On the locally installed Node v20.20.2, `npm test` exits 1 with `Could not find '...scripts/**/*.test.ts'`. Running the same 9 test files with shell glob expansion passes 38/38.
- **Inference:** `node --test` only expands glob patterns itself on Node ≥21, so the quoted glob fails on any Node 20 install. The repo-root `.nvmrc` (new, untracked) pins 22, suggesting the intended runtime moved but `engines` wasn't updated.
- **Significance:** The primary local quality gate fails 100% of the time on a version the package claims to support — and the suite itself is healthy, so the failure is pure tooling noise that will mask real regressions.

### F2 — Cross-domain internal imports violate the documented boundary rule
**Category:** standard mismatch · **Confidence:** Confirmed

- **Observed:** `scripts/README.md` states: "No script imports another domain's internals — shared code lives in `lib/`." Two violations exist: (a) `providers/materialize/{build,validate}/index.ts` and its test import `findRepoRoot` from `../../../distribution/shared/composition`; (b) `backlog-upgrade/script.ts` imports `orchestrateValidation`/`parseBacklog` from `../validate-hawp-workflow/orchestrate`.
- **Inference:** `findRepoRoot` is a thin wrapper over `lib`'s `findUpward` and naturally belongs in `lib/`. The backlog-upgrade→validate-hawp-workflow coupling is deeper and may be intentional (`--validate` reuses the validator), but the rule as written doesn't carve out an exception.
- **Significance:** The boundary rule is brand-new (in this WIP) and already broken in two directions; either the code or the rule should move before the convention loses force.

### F3 — `--strict-warnings` depends on scraping a child process's stdout format
**Category:** truth-risk drift · **Confidence:** Confirmed mechanism, Likely impact

- **Observed:** `backlog-validate/script.ts` spawns the workflow validator and extracts the warning count via regex `/! Warnings:\s+(\d+)/`, returning `0` when no match. The matching line is produced by `validate-hawp-workflow/reporter.ts` (`! Warnings:          ${...}`); formats currently align.
- **Inference:** Any cosmetic change to the reporter's summary line silently zeroes the count, making `--strict-warnings` pass when warnings exist — a fail-open failure mode in a validator. No test in the inspected files covers this cross-process coupling.
- **Significance:** A validation gate that can silently degrade to PASS is exactly the kind of truth-risk this tooling exists to prevent.

### F4 — `npx tsx` resolution is environment-dependent, bypassing the pinned lockfile
**Category:** truth-risk drift · **Confidence:** Confirmed divergence, Likely consequence

- **Observed:** `backlog-validate/script.ts` spawns `npx tsx <script>` with `cwd` = repo root; the repo root has no `package.json`/`node_modules`. Measured: `npx --no-install tsx --version` resolves **v4.22.4** at repo root vs **v4.21.0** inside `librarian/` (the lockfile-pinned version).
- **Inference:** Child validators can run under a tsx version librarian doesn't control (source of the 4.22.4 — global install vs npx cache — was not proven). On a clean machine with no cached tsx, `npx` may prompt or download an unpinned latest version at validation time.
- **Significance:** Validation outcomes become machine-dependent; running the child scripts via the librarian-local tsx (or `node` + `--import` with an explicit path) would restore determinism.

### F5 — `.js`-suffixed relative imports in `backlog-upgrade/` diverge from the documented convention
**Category:** standard mismatch · **Confidence:** Confirmed

- **Observed:** `scripts/README.md` mandates "extensionless relative imports." 16 `.js`-suffixed relative imports remain, all inside `backlog-upgrade/` — including mixed style within a single file (`script.ts` uses both `"./models/index.js"` and `"../lib"`).
- **Inference:** Looks like an incomplete migration of the oldest domain to the new convention, not a deliberate exception — no exception is documented in inspected files.
- **Significance:** Low individually; worth normalizing while the restructure is still uncommitted so the convention starts clean.

### F6 — Stale, hardcoded CLI metadata and a silent `--export-plan` gap in apply mode
**Category:** truth-risk drift · **Confidence:** Confirmed

- **Observed:** `backlog-upgrade/cli.ts` hardcodes `"v1.1.0 (TASK-028 dry-run detection)"` and a `STATUS:` block referencing TASK-028, while `package.json` says version `0.0.0`. Separately, in `script.ts` the apply-mode branch returns before the `options.exportPlan` handling, so `--apply --export-plan x.json` silently ignores the flag; `CLI.md` documents `--export-plan` with no mode caveat.
- **Inference:** The version string and task references are leftover work-item baggage in user-facing output; the export-plan gap may be intentional (no plan is generated in apply mode) but is undocumented either way.
- **Significance:** Minor, but these are exactly the user-facing claims the CLI contract (`CLI.md`) exists to keep honest.

### F7 — Path-leak validation silently skips missing target files
**Category:** validation drift · **Confidence:** Likely

- **Observed:** `distribution/validate/index.ts` checks a hardcoded list of 5 `core/.hawp/kit/...` files for `core/.hawp/` path leaks and does `if (!existsSync(...)) continue;` with no warning. `scripts/README.md` requires "never swallow — surface skipped files/branches with a stderr `[<domain>] warning:` line." All 5 files currently exist (verified).
- **Inference:** If a kit file in the list is renamed, leak-checking for it silently stops while validation keeps reporting PASS. "Likely" because the harm requires a future rename; the divergence from the stated error convention is direct.
- **Significance:** Low probability, but it fails silent, in a validator.

---

## Non-findings (checked and resolved)

1. **Apply-mode mutation safety.** Checked `backlog-upgrade/script.ts`: the dirty-tree guard runs by default, and `hasDirtyWorkingTree` returns `true` when git itself fails — fail-closed, matching the documented "mutating scripts must guard… fails closed" rule. `--force-dirty` is the explicit, documented override.
2. **Untrusted-path containment and suite health.** Checked `evidence-integrity.ts`: evidence links are resolved against the evidence root and rejected with a stderr warning if they escape it — matching the documented containment rule. Also: `tsc --noEmit` clean, and all 38 tests pass when executed directly, so the failures in F1 are mechanism-only.

---

## Out of scope, flagged only

- `.github/workflows/librarian-quality.yml` (new): its Node version should match whatever resolution F1 lands on.
- Repo-root `.nvmrc` pins 22 while the local default Node is 20 — same tension as F1.
- `.hawp/bin/hawp` wrapper shares the `npx tsx` resolution issue from F4.

---

## Operational sequence

**Fix first**

- F1: make `npm test` work on the supported Node range — either raise `engines` to `>=22` (aligning with `.nvmrc`) or stop relying on `node --test` glob expansion; whichever is chosen, keep `engines`, `.nvmrc`, README, and CI consistent.
- F3: replace stdout scraping with a structured signal (distinct exit code, JSON output, or direct in-process call) so strict mode cannot silently fail open.

**Verify next**

- F2: decide whether `findRepoRoot` moves to `lib/` and whether the backlog-upgrade→validator import is a sanctioned exception; update `scripts/README.md` or the code so rule and reality match.
- F4: confirm on a clean environment how `npx tsx` resolves for `.hawp/bin/hawp` users, then pin the child-process runtime to librarian's local tsx.
- F6 (export-plan): confirm whether `--apply --export-plan` is meant to be a no-op; if yes, document it in `CLI.md` and help text; if no, wire it up.

**Defer**

- F5: normalize `.js` imports in `backlog-upgrade/` (mechanical, low risk).
- F6 (metadata): remove the hardcoded `v1.1.0`/TASK-028 strings or derive version from `package.json`.
- F7: add the `[distribution] warning:` line for skipped leak-check targets.
