# Comparison — Install Docs Audit, re-run with Layer 1 + Layer 2 (2026-06-15)

This run has two purposes:

1. **Generalization check** — do the Layer 1 (tiering + verified-correct) and Layer 2 (surface sweep) authoring changes help on a *different task type* (truth audit), not just the scope-creep-trap review they were tuned on?
2. **No-HAWP integrity check** — confirm the HAWP authoring changes did **not** move the no-HAWP score (the concern raised after the no-HAWP audit 85% vs review 75% gap).

The subject (`distribution/`) is unchanged, so the no-HAWP arm is the same-state 85% baseline.

## Scores

| # | Dimension | No-HAWP (Arm A, unchanged) | HAWP original (95% run) | HAWP Layer 1+2 (this run) |
| --- | --- | --- | --- | --- |
| 1 | Drift resistance | 5 | 5 | 5 |
| 2 | Evidence vs inference separation | 3 | 5 | 5 |
| 3 | Output usefulness | 4 | 5 | 5 |
| 4 | Handoff quality | 3 | 5 | 5 |
| 5 | Trustworthiness | 4 | 5 | 5 |
| 6 | Scope adherence | 5 | 5 | 5 |
| 7 | Completeness / coverage | 5 | 4 | **5** |
| 8 | Conciseness / signal-to-noise | 4 | 4 | 4 |
| 9 | Correctness / accuracy | 5 | 5 | 5 |
| 10 | False-positive control | 4 | 5 | 5 |
| 11 | Verifiability | 4 | 5 | 5 |
| 12 | Positive confirmation / balance | 5 | 4 | **5** |
| | **Raw / 60** | **51** | **57** | **59** |
| | **Headline / 15** | **12.75** | **14.25** | **14.75** |
| | **Percentage** | **85%** | **95%** | **98%** |

## Purpose 1 — the layers generalize

The original install-docs HAWP arm (95%) lost #7 and #12 to its hard caps (7 findings, 2 non-findings) — the same two dimensions HAWP loses everywhere. With Layer 1 + Layer 2:

- **#7 4→5:** the surface sweep enumerated every install unit and the uncapped minor list carried the smaller items (README branch mapping, the inaccurate `cp -Rn` example, the REF=dev default, the benchmark-exclusion check) the caps had suppressed.
- **#12 4→5:** the "Verified correct" list (5 items) replaced the 2-non-finding cap with the broad positive picture the unguided arm used to win.
- **No regression:** discipline dimensions held at 5; conciseness held at 4; no false positives.

Net **+2 (57 → 59)** — the same recovery pattern as the review-task layer runs. The authoring layers are not task-specific; they help any capped review/audit pattern.

## Purpose 2 — no-HAWP was not harmed

**No-HAWP held at exactly 51/60 (85%).** It could not have moved: the no-HAWP arm never sees the HAWP shape, and the audited subject is byte-identical. The "10-point drop" the question flagged was between *different task types*, not a regression:

| Run | Task type | No-HAWP |
| --- | --- | --- |
| install-docs truth audit | standards/truth audit | 85% |
| install-docs **this re-run** | standards/truth audit | 85% (unchanged) |
| bounded review / layer1 / layer2 | scope-creep-trap review | 75% |

Audits are naturally bounded (little room to drift → high unguided baseline); the scope-creep-trap review actively baits drift (→ lower unguided baseline). Both are correct for their task and are not comparable to each other.

## Which is better, for now

**HAWP (Layer 1+2) 98% vs no-HAWP 85% — a +13-point margin**, up from +10 in the original audit, entirely because HAWP rose (no-HAWP was flat). As in the layer-2 review run, the unguided arm no longer wins any dimension outright: HAWP now ties it on #7 and #12 while keeping the discipline-cluster wins.

## Caveats

- **Author = scorer**, mitigated by the anchored rubric and by the no-HAWP arm being a fixed, unchanged carry-over (no re-scoring latitude there).
- **No-HAWP reused, not re-run** — valid and in fact the point: the subject is unchanged, so an honest re-run would reproduce 51/60.
- **n=1** per task; two task types now show the same +2 layer effect, which is corroborating but still small-sample.
