# Scorecard — HAWP (Audit pattern, Layer 1 + Layer 2)

Scored against the 12-dimension rubric in [benchmark-prompt.md](../../../benchmark-prompt.md). Scale: 0–5 anchored per dimension.

**Scored as Arm B**, against the unchanged same-state no-HAWP audit run (Arm A, 51/60). Subject `distribution/` is byte-identical to the 2026-06-15 install-docs run.

| # | Dimension | Score | Δ vs 95% run | Reasoning |
| --- | --- | --- | --- | --- |
| 1 | Drift resistance | 5 | — | Install-only scope held; update referenced only for seed-vs-refresh. |
| 2 | Evidence vs inference separation | 5 | — | Every finding pairs Observed (doc) with Observed (script) + confidence label. |
| 3 | Output usefulness | 5 | — | Verdict + Fix/Verify/Defer tied to findings; supports a ship-or-fix call. |
| 4 | Handoff quality | 5 | — | Scope, method, validate command, and the visible sweep table make coverage auditable. |
| 5 | Trustworthiness | 5 | — | Confidence downgrades where appropriate (F4 harm Likely); nuanced verdict. |
| 6 | Scope adherence | 5 | — | Stayed within install scope. |
| 7 | Completeness / coverage | 5 | **+1** | The Layer-2 sweep enumerated all install units; the uncapped minor list carries the breadth (README branch mapping, the `cp -Rn` example inaccuracy, REF=dev default, benchmark exclusion) that the prior 7-finding/2-non-finding caps suppressed. Now matches the unguided arm's net. |
| 8 | Conciseness / signal-to-noise | 4 | — | Held. Information-dense; a couple of shallow-checked minor items add slight noise — the expected cost of the sweep + uncapped list. |
| 9 | Correctness / accuracy | 5 | — | Claims verified against sources and `distribution:validate`. |
| 10 | False-positive control | 5 | — | Minor items hedged ("not deep-checked", "confirm"); no non-issues asserted as facts. |
| 11 | Verifiability | 5 | — | Doc+script citations with files/lines, plus the sweep table and validate command. |
| 12 | Positive confirmation / balance | 5 | **+1** | The Layer-1 "Verified correct" list (5 items) replaces the prior 2-non-finding cap, giving the broad positive picture the unguided arm used to win. |

## Notable strengths

The two authoring layers lifted exactly the two dimensions the capped audit lost in the 95% run: #7 (sweep + uncapped minor list) and #12 (verified-correct list). The discipline cluster stayed at 5 and no false positives crept in. Net **+2 (57 → 59)**, the same recovery pattern seen on the review-task runs — confirming the layers generalize beyond the scope-creep-trap task.

**Raw:** 59 / 60 → **Headline 14.75 / 15** · **Percentage 98%**
