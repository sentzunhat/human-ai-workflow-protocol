# Benchmark Run: Layer-2 Surface Sweep (2026-06-15) — EXECUTED

Purpose: measure whether the **Layer-2 authoring change** — a "Sweep the scoped surface" habit in the Review/Audit patterns (enumerate every in-scope unit and confirm each was inspected) — recovers **#7 completeness/coverage**, the one dimension the Layer-1 change failed to lift, without regressing the others.

> **Status: executed 2026-06-15.** HAWP arm run against the current Node 26 tree with the Layer 1+2 shape; no-HAWP arm reused same-state (subject unchanged). Result: **HAWP 59/60 (98%) vs no-HAWP 45/60 (75%)**. Layer 2 recovered **#7 (4→5)** with no regression, and the unguided arm no longer wins any dimension outright. See [`comparison.md`](comparison.md), [`hawp/scorecard.md`](hawp/scorecard.md), and [`no-hawp/README.md`](no-hawp/README.md).

## Setup

- **Task type:** Bounded repo review (scope-creep trap) — same task/folder/state as the Layer-1 run, so the no-HAWP baseline carries over.
- **Task prompt (both arms):** "Take a look at the librarian/ folder and tell me what's wrong with it. What else should we clean up?"
- **What changed:** the HAWP authoring pattern gained the Layer-2 surface-sweep habit (`core/.hawp/kit/references/authoring-patterns.md` and the mirrored root `.hawp/` copy). Provider/distribution artifacts were regenerated (no content change — authoring-patterns does not flow into them) and all gates pass.
- **Model:** Cursor agent (model not recorded).

## Layer-2 addition to the shape

On top of the Layer-1 tiering + verified-correct requirements, the shape's `constraints`/`output` now require:

> **Sweep the scoped surface (Layer 2):** before finalizing, enumerate every in-scope unit — each domain folder and its index/cli/script files, lib/, the config files (tsconfig.json, package.json), and the contract docs (scripts/README.md, CLI.md) — and confirm each was inspected. Make the sweep visible as a coverage checklist; use bounded-absence wording for anything not inspected.

## Measured result

| Arm | Raw / 60 | Headline / 15 | Percentage |
| --- | --- | --- | --- |
| No-HAWP (Arm A, same-state) | 45 | 11.25 | 75% |
| HAWP pre-Layer-1 (2026-06-11) | 57 | 14.25 | 95% |
| HAWP Layer-1 (2026-06-15) | 58 | 14.5 | 97% |
| **HAWP Layer-1+2 (this run)** | **59** | **14.75** | **98%** |

Cumulative effect of the two authoring layers:

| Change | Dimension recovered | HAWP raw |
| --- | --- | --- |
| + Layer 1 (tiering + verified-correct) | #12 (4→5) | 57 → 58 |
| + Layer 2 (surface sweep) | #7 (4→5) | 58 → 59 |

The surface sweep forced coverage of the units the Layer-1 arm missed (`cli.ts` filesystem work, the `providers:validate` gap, dead `tsconfig`, the third repo-root finder, type guards, the `validate:workflow` alias, the README test-discovery mismatch) — closing the coverage gap to the unguided arm without manufacturing false positives.

## Caveats

- **Author = scorer**, mitigated by the anchored rubric. #7 = 5 is defensible (the output covers everything the unguided arm found, plus more, via a visible sweep); #8 was honestly held at 4, not raised, to reflect the fuller artifact.
- **No-HAWP reused same-state**, valid because `librarian/` is unchanged — see [`no-hawp/README.md`](no-hawp/README.md).
- **n=1**, and the sweep is only as honest as the unit list. Diminishing headroom: at 59/60, conciseness is the lone non-5 and further structure would likely cost #8.

## Contents

- [`comparison.md`](comparison.md) — four-column side-by-side (no-HAWP, pre-Layer-1, Layer-1, Layer-1+2)
- [`hawp/output.md`](hawp/output.md) — filled Layer 1+2 shape + the review (with the visible surface sweep)
- [`hawp/scorecard.md`](hawp/scorecard.md) — HAWP scores
- [`no-hawp/README.md`](no-hawp/README.md) — reused same-state baseline pointer
