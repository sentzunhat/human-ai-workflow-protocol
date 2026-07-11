# Comparison — Layer-2 Surface Sweep (2026-06-15)

This run tests the **Layer-2 authoring change**: a "Sweep the scoped surface" habit added to the Review/Audit patterns, requiring the reviewer to enumerate every in-scope unit and confirm each was inspected. It targets the one dimension Layer 1 failed to recover — **#7 completeness/coverage** — without regressing the rest.

The no-HAWP arm is the same-state baseline (subject `librarian/` is unchanged); the comparison is read against both that baseline and the prior HAWP runs.

## Scores

| # | Dimension | No-HAWP (Arm A, same-state) | HAWP pre-Layer-1 (06-11) | HAWP Layer-1 (06-15) | HAWP Layer-1+2 (this run) |
| --- | --- | --- | --- | --- | --- |
| 1 | Drift resistance | 3 | 5 | 5 | 5 |
| 2 | Evidence vs inference separation | 3 | 5 | 5 | 5 |
| 3 | Output usefulness | 4 | 5 | 5 | 5 |
| 4 | Handoff quality | 3 | 5 | 5 | 5 |
| 5 | Trustworthiness | 3 | 5 | 5 | 5 |
| 6 | Scope adherence | 3 | 5 | 5 | 5 |
| 7 | Completeness / coverage | 5 | 4 | 4 | **5** |
| 8 | Conciseness / signal-to-noise | 4 | 4 | 4 | 4 |
| 9 | Correctness / accuracy | 4 | 5 | 5 | 5 |
| 10 | False-positive control | 3 | 5 | 5 | 5 |
| 11 | Verifiability | 5 | 5 | 5 | 5 |
| 12 | Positive confirmation / balance | 5 | 4 | 5 | 5 |
| | **Raw / 60** | **45** | **57** | **58** | **59** |
| | **Headline / 15** | **11.25** | **14.25** | **14.5** | **14.75** |
| | **Percentage** | **75%** | **95%** | **97%** | **98%** |

## What the Layer-2 change did

- **#7 completeness/coverage: 4 → 5 (recovered).** This is the dimension Layer 1's tiering could not lift. The Layer-2 **surface sweep** forced enumeration of all nine in-scope units as a visible checklist, which caught every item the same-state no-HAWP arm had found and the Layer-1 HAWP arm had missed:
  - `validate-hawp-workflow/cli.ts` doing filesystem work (now F3)
  - `hawp backlog validate` omitting `providers:validate` (now F5)
  - dead `tsconfig` emit config (now F7)
  - the third repo-root finder (folded into F2)
  - test-only type guards, the `validate:workflow` alias, and the `scripts/README.md` test-discovery doc mismatch (minor list)
- **No regression.** Discipline dimensions held at 5; #12 stayed 5 (Layer 1's verified-correct list); conciseness held at 4 despite the fuller output.
- **It did not manufacture false positives.** Sweeping harder is the obvious risk for #10, but every added finding is confirmed real, and the arm avoided the no-HAWP arm's `process.exit`-in-`index.ts` false positive. #10 stayed 5.

## Cumulative effect of the two layers

| Change | Dimension recovered | HAWP raw | HAWP % |
| --- | --- | --- | --- |
| Pre-layer baseline (06-11) | — | 57 | 95% |
| + Layer 1 (tiering + verified-correct) | #12 (4→5) | 58 | 97% |
| + Layer 2 (surface sweep) | #7 (4→5) | 59 | 98% |

The two authoring habits recovered exactly the two dimensions the unguided arm historically wins, one each, for a combined **+2 raw points** — landing on the original Layer-1 projection of 59/60, but only after Layer 2 supplied the missing half.

## HAWP vs no-HAWP, same-state

**HAWP (Layer 1+2) 98% vs no-HAWP 75% — a +23-point margin.** HAWP now wins or ties every dimension: it matches no-HAWP's coverage (#7 tied at 5) and balance (#12 tied at 5) while keeping its decisive wins on the discipline cluster (drift, evidence separation, handoff, trustworthiness, scope, false-positive control). The unguided arm no longer wins a single dimension outright — the first run in the set where that is true.

## Caveats

- **Author = scorer**, mitigated by the anchored rubric. The #7 = 5 award is defensible because the HAWP output demonstrably covers everything the unguided arm found plus more, with a visible sweep; the judgment to hold #8 at 4 (not 5) reflects the real length growth.
- **Same-state no-HAWP is reused, not re-run** — valid because the subject `librarian/` is unchanged; see [`no-hawp/README.md`](no-hawp/README.md).
- **n=1**, and the sweep's coverage guarantee is only as honest as the unit list. The result shows the habit *can* close the gap, not that it always will.
- **Diminishing headroom:** at 59/60 with conciseness the lone non-5, further required structure would likely start costing #8 rather than adding net score.
