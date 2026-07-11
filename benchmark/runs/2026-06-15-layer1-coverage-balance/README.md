# Benchmark Run: Layer-1 Coverage & Balance (2026-06-15) — EXECUTED

Purpose: measure whether the **Layer-1 authoring change** (the new "Coverage and Balance" section in `authoring-patterns.md`: tiered findings + a positive-confirmation element) recovers the two dimensions HAWP loses — **#7 completeness/coverage** and **#12 positive-confirmation/balance** — without regressing the dimensions it already wins.

> **Status: executed 2026-06-15.** Both arms were run **same-state** against the current Node 26 tree (a fresh no-HAWP arm in the clean `/tmp/` workspace, plus the Layer-1 HAWP shape) and scored on the 12-dimension rubric. Result: **HAWP 58/60 (97%) vs no-HAWP 45/60 (75%)**. The Layer-1 change recovered **#12 (4→5)** cleanly but **did not** recover **#7** (held at 4) — the same-state no-HAWP arm still found more valid items. Net Layer-1 effect: **+1 raw (57→58)**, below the original 59/60 projection. See [`comparison.md`](comparison.md), [`hawp/scorecard.md`](hawp/scorecard.md), and [`no-hawp/scorecard.md`](no-hawp/scorecard.md).

## Setup

- **Task type:** Bounded repo review (scope-creep trap) — same task/folder as [2026-06-11-bounded-repo-review](../2026-06-11-bounded-repo-review/README.md), re-run same-state.
- **Task prompt (both arms):** "Take a look at the librarian/ folder and tell me what's wrong with it. What else should we clean up?"
- **What changed:** the HAWP shape's `constraints`/`output` gained the Layer-1 tiering + verified-correct requirements. Both arms were run against the **current** Node 26 tree (the no-HAWP arm freshly, in the clean workspace), so this is a same-state head-to-head — not a reused baseline.
- **Clean workspace used for the no-HAWP arm:** `/tmp/hawp-benchmark-clean-20260615-222558` (since removed during cleanup).
- **Model:** Cursor agent (model not recorded for either arm).

## Layer-1 HAWP shape (Review Tasks pattern, with Coverage & Balance)

**input:** "Take a look at the librarian/ folder and tell me what's wrong with it. What else should we clean up?"

**context:** "The librarian/ folder contains the repo's TypeScript maintenance tooling. There is substantial uncommitted work in progress restructuring these scripts (visible in git status). This is a read-only review; make no changes."

**mission:** "Produce a bounded hygiene-and-correctness review of librarian/ through a maintainability and truth-risk lens, to support a decision about which cleanup items to add to the backlog next."

**constraints:** "Scope: librarian/ only. Evidence bar: every substantive finding separates directly observed evidence from inference and carries a confidence label (Confirmed / Likely / Unclear), matching the weakest part of the claim. Use only the four finding categories (truth-risk drift, validation drift, maintainability drift, standard mismatch). Call something a standard violation only when backed by explicit repo docs, an ADR, a tooling contract, or a clearly intentional repeated convention. When claiming absence, use bounded wording. No architecture redesign. The 'what else should we clean up?' part is a scope-creep trap: keep everything inside librarian/; out-of-scope items go under 'Out of scope, flagged only', max 3 one-liners. **Tier the findings (Layer 1):** cap *primary* findings (full Observed/Inference/Significance detail) at 7, but include an **uncapped 'Minor / deferred' one-liner list** for valid smaller items so coverage is not suppressed. **State what is correct (Layer 1):** include a compact 'Verified correct' list of areas checked and found sound, in addition to any non-findings, without confirming anything not actually inspected."

**output:** "A prioritized review artifact: (1) scope + method; (2) up to 7 primary findings, each with category, confidence, observed evidence, inference/uncertainty, significance; (3) an uncapped 'Minor / deferred' one-liner list; (4) a compact 'Verified correct' list and/or up to 2 non-findings; (5) 'Out of scope, flagged only' (max 3); (6) closing Fix first / Verify next / Defer."

## How this run was executed

1. **No-HAWP arm** (clean `/tmp/` workspace, per [instructions/run.md](../../instructions/run.md)): bare prompt, current tree → [`no-hawp/output.md`](no-hawp/output.md), scored in [`no-hawp/scorecard.md`](no-hawp/scorecard.md).
2. **HAWP arm** (this repo, per [instructions/cleanup.md](../../instructions/cleanup.md)): the Layer-1 shape above → [`hawp/output.md`](hawp/output.md), scored in [`hawp/scorecard.md`](hawp/scorecard.md).
3. Both scored on the 12-dimension / 0–5 rubric; see [`comparison.md`](comparison.md) for the side-by-side and which findings each arm uniquely caught.

## Measured result

| Arm | Raw / 60 | Headline / 15 | Percentage |
| --- | --- | --- | --- |
| No-HAWP (Arm A, same-state current tree) | 45 | 11.25 | 75% |
| HAWP pre-Layer-1 (2026-06-11) | 57 | 14.25 | 95% |
| **HAWP Layer-1 (Arm B, current tree)** | **58** | **14.5** | **97%** |

Effect of the Layer-1 change on its two target dimensions:

| # | Dimension | Pre-Layer-1 | Layer-1 | Δ |
| --- | --- | --- | --- | --- |
| 12 | Positive confirmation / balance | 4 | 5 | **+1 (recovered)** |
| 7 | Completeness / coverage | 4 | 4 | 0 (not recovered — no-HAWP still wins it) |
| 8 | Conciseness / signal-to-noise | 4 | 4 | 0 (held) |

**Outcome vs projection:** the original 59/60 "recovers both" projection was optimistic. Measured at **58/60**: the required **Verified correct** list recovered #12 cleanly, but the uncapped **Minor / deferred** list did not lift #7 — a same-state no-HAWP arm still surfaced ~6 valid items the HAWP arm missed (`cli.ts` filesystem work, a third repo-root finder, dead `tsconfig` config, type-guard scaffolding, the `validate:workflow` alias, the `providers:validate` gap). Net Layer-1 gain: **+1 raw point**, entirely from #12.

## Caveats

- **Author = scorer**, mitigated by the anchored rubric. The key judgment was holding HAWP #7 at 4 rather than 5: the same-state no-HAWP arm found more valid items, so a perfect coverage score would have been self-serving. See [`comparison.md`](comparison.md).
- **The Layer-1 win is modest and honest (+1), not the projected +2.** A positive-confirmation requirement is easy to satisfy on demand (#12); cap-relief alone does not make a bounded lens sweep as wide as an unbounded one (#7).
- **n=1.** Single same-state run; not independently replicated.
