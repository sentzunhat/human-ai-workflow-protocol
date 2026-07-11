# Benchmark Run: Install Docs Audit × Layer 1+2 (2026-06-15) — EXECUTED

Two purposes: (1) check that the **Layer 1 + Layer 2** authoring changes generalize to a *different task type* (truth audit, not the scope-creep-trap review they were tuned on); (2) confirm the changes did **not** move the **no-HAWP** score — the integrity concern raised from the audit 85% vs review 75% gap.

> **Status: executed 2026-06-15.** HAWP arm re-run with the Audit pattern carrying Layer 1+2; no-HAWP arm carried over same-state (subject `distribution/` is byte-unchanged). Result: **HAWP 59/60 (98%) vs no-HAWP 51/60 (85%, unchanged)**. The layers recovered #7 and #12 here too (+2), and no-HAWP did not move. See [`comparison.md`](comparison.md), [`hawp/scorecard.md`](hawp/scorecard.md), [`no-hawp/README.md`](no-hawp/README.md).

## Setup

- **Task type:** Standards / truth audit (same as [2026-06-15-install-docs-truth-audit](../2026-06-15-install-docs-truth-audit/README.md)).
- **Task prompt (both arms):** "Check whether our install docs match what the scripts actually do."
- **Subject:** `distribution/` install sources + generated guides — **unchanged** (`git status --short distribution/` clean; `distribution:validate` passes).
- **What changed:** only the HAWP authoring pattern (Layer 1 tiering + verified-correct, Layer 2 surface sweep). No-HAWP arm cannot be affected by construction.

## Headline

| Arm | Raw / 60 | Percentage | vs original audit |
| --- | --- | --- | --- |
| No-HAWP (same-state, unchanged) | 51 | **85%** | **no change** |
| HAWP original (caps) | 57 | 95% | — |
| **HAWP Layer 1+2** | **59** | **98%** | **+2 (recovered #7, #12)** |

## The two answers

1. **Layers generalize:** on the audit task, Layer 1+2 recovered #7 (4→5, via the surface sweep + uncapped minor list) and #12 (4→5, via the verified-correct list), with no regression — the same +2 pattern as the review-task runs.
2. **No-HAWP not harmed:** it held at 85%. The earlier "10-point drop" was a **task-type** difference (bounded audit 85% vs scope-creep-trap review 75%), not a scoring regression — and HAWP authoring changes can never touch the unguided arm.

## Caveats

- Author = scorer (anchored rubric; no-HAWP is a fixed carry-over with no re-scoring latitude).
- No-HAWP reused, not re-run — valid because the subject is unchanged.
- n=1 per task type; two task types now show the same layer effect.
