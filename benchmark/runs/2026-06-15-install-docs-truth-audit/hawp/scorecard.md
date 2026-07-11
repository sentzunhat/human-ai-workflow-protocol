# Scorecard — HAWP

Scored against the 12-dimension rubric in [benchmark-prompt.md](../../../benchmark-prompt.md). Scale: 0–5 anchored per dimension (5 = clearly satisfied, 4 = minor gaps, 3 = notable gaps, 2 = mostly not, 1 = barely, 0 = not satisfied).

**Scored blind** as Arm B before the HAWP/no-HAWP mapping was revealed.

| # | Dimension | Score | Reasoning |
| --- | --- | --- | --- |
| 1 | Drift resistance | 5 | Install-only scope held. Update guides referenced only where seed-vs-refresh required it. |
| 2 | Evidence vs inference separation | 5 | Every finding pairs doc claim with script observation; confidence labels applied; inference sections explicit (F1 user impact marked Likely). |
| 3 | Output usefulness | 5 | Opens with verdict ("partially trustworthy"), ends with Fix first / Verify next / Defer tied to specific findings. Directly supports a ship-or-fix decision. |
| 4 | Handoff quality | 5 | Scope, method, and validation command stated upfront. Non-findings record what was checked and ruled out. |
| 5 | Trustworthiness | 5 | Downgrades where appropriate (F4 harm Likely not Confirmed). Verdict is nuanced, not binary. |
| 6 | Scope adherence | 5 | Stayed within the implied install scope; update referenced only to explain seed-vs-refresh. Scored on the same basis as the no-HAWP arm — no credit for merely having explicit constraints. |
| 7 | Completeness / coverage | 4 | Found the same five substantive mismatches, but the deliberate 5-finding cap excluded the minor kit-subfolder under-count the other arm surfaced. Discipline traded for a little breadth. |
| 8 | Conciseness / signal-to-noise | 4 | Structured and high-signal; the per-finding template (Observed/Inference/Significance ×5) adds some repetition. |
| 9 | Correctness / accuracy | 5 | Factual claims check out against the sources; no incorrect findings identified. |
| 10 | False-positive control | 5 | Bounded "no evidence found in inspected prose" wording; two explicit non-findings record dismissed concerns rather than padding the finding list. |
| 11 | Verifiability | 5 | Pairs Observed (doc) with Observed (script) for every finding, with exact files/lines and the validate command — a reader can independently re-check each claim. |
| 12 | Positive confirmation / balance | 4 | Two non-findings record verified alignment, but the 2-non-finding cap gives a narrower positive picture than the other arm's broad "Matches" inventory. |

## Notable strengths

The shape's evidence bar produced side-by-side doc/script citations with confidence labels for every mismatch, an explicit trust verdict, and an operational Fix/Verify/Defer sequence. Strongest on the report-discipline and verifiability dimensions. The cost is breadth: the finding and non-finding caps left some coverage and positive confirmation to the other arm.

**Raw:** 57 / 60 → **Headline 14.25 / 15** · **Percentage 95%**
