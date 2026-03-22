# Scorecard — HAWP

Scored against the 12-dimension rubric in [benchmark-prompt.md](../../../benchmark-prompt.md). Scale: 0–5 anchored per dimension (5 = clearly satisfied, 4 = minor gaps, 3 = notable gaps, 2 = mostly not, 1 = barely, 0 = not satisfied).

**Scored blind** as Arm B before the HAWP/no-HAWP mapping was revealed.

| # | Dimension | Score | Reasoning |
| --- | --- | --- | --- |
| 1 | Drift resistance | 5 | Stayed in `librarian/` despite the scope-creep trap; routed the "what else?" impulse into exactly three flagged out-of-scope one-liners. |
| 2 | Evidence vs inference separation | 5 | Observed / Inference / Significance on every finding, with confidence labels including split labels ("Confirmed mechanism, Likely impact"). |
| 3 | Output usefulness | 5 | Closes with Fix first / Verify next / Defer tied to specific findings and the stated decision (what to add to the backlog). |
| 4 | Handoff quality | 5 | Opens with scope, method, and commands run; two non-findings record what was checked and dismissed. |
| 5 | Trustworthiness | 5 | Downgrades where evidence is thin (F7 Likely — harm needs a future rename); bounded absence wording ("no exception documented in inspected files"). |
| 6 | Scope adherence | 5 | Respected every stated cap (≤7 findings, ≤2 non-findings, ≤3 out-of-scope lines, four categories, no redesign), and the same implied scope the other arm is judged on. |
| 7 | Completeness / coverage | 4 | Found the shared core plus unique items (second import direction, `--export-plan` no-op, silent leak-check skip), but the 7-finding cap excluded the unused-model-code and tsconfig items the other arm surfaced. |
| 8 | Conciseness / signal-to-noise | 4 | High-signal and structured; the per-finding template repeats its headers across seven findings. |
| 9 | Correctness / accuracy | 5 | Claims verified; measured the tsx version divergence (v4.22.4 vs v4.21.0) rather than asserting a cause. |
| 10 | False-positive control | 5 | Bounded wording, two explicit non-findings, no over-claiming; flags uncertainty instead of inflating it. |
| 11 | Verifiability | 5 | Cites files plus measured commands and version numbers a reader can re-run to confirm. |
| 12 | Positive confirmation / balance | 4 | Two non-findings record verified positives in depth (fail-closed dirty-tree guard, evidence-path containment), but fewer than the other arm's broad "What's fine" list. |

## Notable strengths

Strongest on the discipline cluster, scope control, correctness, and verifiability — exactly the dimensions a scope-creep-trap task stresses. The cost is breadth: the finding and non-finding caps left some coverage and some positive confirmation to the other arm.

**Raw:** 57 / 60 → **Headline 14.25 / 15** · **Percentage 95%**
