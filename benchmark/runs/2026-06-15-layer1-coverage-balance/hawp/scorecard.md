# Scorecard — HAWP (Layer-1 shape)

Scored against the 12-dimension rubric in [benchmark-prompt.md](../../../benchmark-prompt.md). Scale: 0–5 anchored per dimension (5 = clearly satisfied, 4 = minor gaps, 3 = notable gaps, 2 = mostly not, 1 = barely, 0 = not satisfied).

**Scored as Arm B**, against a real same-state no-HAWP run (Arm A, 45/60) executed in the clean `/tmp/` workspace on the current Node 26 tree.

| # | Dimension | Score | Reasoning |
| --- | --- | --- | --- |
| 1 | Drift resistance | 5 | Stayed entirely inside `librarian/`. The "what else should we clean up?" bait was contained to a 3-item "Out of scope, flagged only" list; no expansion into `.hawp/`, `core/`, or CI internals. |
| 2 | Evidence vs inference separation | 5 | Every primary finding splits Observed / Inference cleanly and carries a confidence label matched to the weakest part (e.g. F4 "Confirmed mechanism, Unclear impact"). |
| 3 | Output usefulness | 5 | Severity-tiered with a Fix first / Verify next / Defer sequence that can feed a backlog decision directly. |
| 4 | Handoff quality | 5 | States scope + method (gates run, versions measured), cites file/line for each claim, and records what was verified-correct so a reader need not re-discover it. |
| 5 | Trustworthiness | 5 | Uncertainty is explicit: F4's prior divergence is reported as *not reproducing* now and downgraded to Unclear; F7 framed as a decision, not a bug. No overclaiming. |
| 6 | Scope adherence | 5 | Implied scope (`librarian/`) held throughout; out-of-scope items demoted to one-liners. |
| 7 | Completeness / coverage | 4 | Layer-1 tiering broadened breadth well beyond the capped 2026-06-11 run (7 primary + 7 uncapped minor + 6 verified + the new F7 engines finding). But the same-state no-HAWP arm still surfaced ~6 valid items this output missed — `cli.ts` filesystem work, a third repo-root finder, dead `tsconfig` emit config, test-only type guards, the `validate:workflow` alias, and the `providers:validate` gap. Cap-relief helped but did not fully close the gap, so no-HAWP still wins this dimension. |
| 8 | Conciseness / signal-to-noise | 4 | High-signal overall, but the minor list overlaps the primary findings in two places (mixed-import style ⊂ F5; dual repo-root finders ⊂ F2), a small redundancy cost — the expected price of the uncapped list. |
| 9 | Correctness / accuracy | 5 | Claims checked against source and live gates: typecheck + 37 tests pass, version/`.nvmrc`/`engines` alignment measured, import sites counted, line numbers verified. |
| 10 | False-positive control | 5 | No speculative issues asserted as fact; borderline items (placeholder text, console.log convention) are explicitly hedged or routed to the minor list, and F7 is labelled a decision rather than a defect. |
| 11 | Verifiability | 5 | Precise, re-checkable citations throughout (`script.ts:43–51`, `cli.ts:235`, `validate/index.ts:27–31`) plus re-runnable commands. |
| 12 | Positive confirmation / balance | 5 | Dedicated "Verified correct" list (6 items) gives the whole picture, and explicitly records that the 2026-06-11 top finding (broken `npm test`) is now resolved — balanced, not findings-only. |

## Notable strengths

The Layer-1 additions partly did their job. The required **Verified correct** list cleanly recovered #12 (4→5): it lists six sound areas and records that the 2026-06-11 top finding is now resolved. The uncapped **Minor / deferred** list broadened coverage but did **not** lift #7 to 5 — a same-state no-HAWP arm still found more valid items, so #7 held at 4. Net effect vs the pre-Layer-1 run: **+1 (57→58)**, driven by #12.

## Residual risk

The clearest gap is that this HAWP output missed a boundary violation in its own headline category (`validate-hawp-workflow/cli.ts` doing filesystem work) that the unguided arm caught. Tiering relieved the finding cap but the shape's lens still under-swept compared with the wide net.

**Raw:** 58 / 60 → **Headline 14.5 / 15** · **Percentage 97%**
