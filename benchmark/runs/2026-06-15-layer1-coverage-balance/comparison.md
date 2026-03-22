# Comparison — Layer-1 Coverage & Balance (2026-06-15)

This run tests one thing: does the **Layer-1 authoring change** (tiered findings + a positive-confirmation element in the Review Tasks pattern) recover the two dimensions HAWP normally loses — #7 completeness/coverage and #12 positive-confirmation/balance — without regressing the dimensions it already wins?

Unlike the projection originally recorded here, both arms were run **same-state** against the current tree (Node 26): a fresh no-HAWP arm in the clean `/tmp/` workspace (Arm A) and the Layer-1 HAWP shape in the repo (Arm B). The pre-Layer-1 HAWP column (2026-06-11) is shown for the before/after read on the authoring change.

## Scores

| # | Dimension | No-HAWP (Arm A, same-state) | HAWP pre-Layer-1 (2026-06-11) | HAWP Layer-1 (Arm B) |
| --- | --- | --- | --- | --- |
| 1 | Drift resistance | 3 | 5 | 5 |
| 2 | Evidence vs inference separation | 3 | 5 | 5 |
| 3 | Output usefulness | 4 | 5 | 5 |
| 4 | Handoff quality | 3 | 5 | 5 |
| 5 | Trustworthiness | 3 | 5 | 5 |
| 6 | Scope adherence | 3 | 5 | 5 |
| 7 | Completeness / coverage | **5** | 4 | 4 |
| 8 | Conciseness / signal-to-noise | 4 | 4 | 4 |
| 9 | Correctness / accuracy | 4 | 5 | 5 |
| 10 | False-positive control | 3 | 5 | 5 |
| 11 | Verifiability | 5 | 5 | 5 |
| 12 | Positive confirmation / balance | **5** | 4 | **5** |
| | **Raw / 60** | **45** | **57** | **58** |
| | **Headline / 15** | **11.25** | **14.25** | **14.5** |
| | **Percentage** | **75%** | **95%** | **97%** |

## What the Layer-1 change did

- **#12 positive-confirmation/balance: 4 → 5 (recovered).** The Layer-1 shape requires a **Verified correct** list. The HAWP arm listed six verified-sound areas (clean typecheck, 37/37 tests, dependency-free `lib/`, fail-closed dirty guard, evidence containment, boundary pattern) and explicitly noted the old top finding is resolved — matching the balanced picture the no-HAWP arm wins on. Clean recovery.
- **#7 completeness/coverage: 4 → 4 (not recovered).** The uncapped **Minor / deferred** list broadened the HAWP output well beyond the pre-Layer-1 capped run, but the same-state no-HAWP arm still surfaced ~6 valid items the HAWP arm missed (see below). Cap-relief helped breadth in absolute terms but did not close the gap to the wide unguided net, so no-HAWP still wins #7.
- **No regression elsewhere.** Discipline dimensions (1–6, 9–11) held at 5; conciseness held at 4. Net Layer-1 effect: **+1 raw (57 → 58)**, entirely from #12.

## What each arm found that the other missed

**No-HAWP found, HAWP missed (the #7 gap):**
- `validate-hawp-workflow/cli.ts` does filesystem work (`findWorkDirectory`/`resolveWorkDirectory` with `existsSync`), violating "cli.ts never reads or writes files" — a boundary violation in HAWP's own headline category.
- A **third** repo-root finder (`findWorkDirectory`) — HAWP counted only two.
- Dead `tsconfig` emit config (`declaration`, `sourceMap`, `outDir: build-src-esm`, `noEmit: false`) that never emits.
- Test-only type guards / speculative model scaffolding.
- `validate:workflow` is a pure alias of `workflow:validate`.
- `hawp backlog validate` never runs `providers:validate` (provider-pack drift can pass).
- `scripts/README.md` still documents test discovery as a `**/*.test.ts` glob while `package.json` uses `find`.

**HAWP found, no-HAWP missed (or framed weaker):**
- `--apply --export-plan` silently ignores the flag (apply mode returns before export handling) — a concrete contract gap.
- `engines.node: ">=26"` pins the toolchain to a just-released major with no documented rationale.
- Sharper truth-risk framing of the stdout-scraping fail-open path, with confidence labels.

**Both found:** cross-domain imports, stdout-scraping for `--strict-warnings`, stale `v1.1.0`/TASK-028 CLI metadata, `.js` imports in `backlog-upgrade/`. Discovery of the top issues was roughly equal; HAWP shaped the *report* (tiers, confidence, verified-correct, operational sequence) while no-HAWP cast a wider but flatter net.

**No-HAWP's discipline costs (the 13-point gap):** one outright false positive (#5 — `process.exit` flagged in files that *are* `index.ts`, where it is permitted), an unproven causal claim stated as fact (npx cache), and drift into CI / provider-sync / `.hawp/bin/` outside the stated subject.

## How the subject changed since 2026-06-11

Same task and folder, but the tree advanced, which both arms reflected:
- **Resolved:** the broken `npm test` on Node 20 (old F1). `engines.node` is `>=26`, `.nvmrc` and runtime are Node 26, the test script uses `find`. Both arms logged this as fine.
- **New:** `engines.node: ">=26"` floor (HAWP F7).
- **No longer reproducing:** the `npx tsx` version divergence (old F4) — HAWP honestly downgraded it to "Unclear"; no-HAWP still asserted the npx-cache cause.

## Which is better, for now

**HAWP (Layer-1), 97% vs no-HAWP 75% — a +22-point margin (13 raw points) on a clean same-state comparison.** The gap is wider than the original 2026-06-11 reading (+18) largely because the unguided arm shipped a false positive and more drift this time, not because HAWP improved much. The Layer-1 change itself was a **modest, honest win: +1 raw point**, recovering #12 cleanly but leaving #7 to the no-HAWP arm.

## Caveats

- **Author = scorer**, mitigated by the anchored rubric but not blind in the strict sense. The most consequential judgment was holding HAWP #7 at 4 (not 5): the same-state no-HAWP arm demonstrably found more valid items, so awarding HAWP a perfect coverage score would have been self-serving.
- **The original projection (59/60, "recovers both") was optimistic.** The measured result is 58/60: #12 recovered, #7 did not. The honest lesson is that a positive-confirmation requirement is easy to satisfy on demand, but cap-relief alone does not make a bounded lens sweep as wide as an unbounded one.
- **n=1.** Single run; the +1 Layer-1 effect is real but small and not independently replicated.
