# Scorecard — HAWP

Scored against the 12-dimension rubric in [benchmark-prompt.md](../../../benchmark-prompt.md). Scale: 0–5 anchored per dimension (5 = clearly satisfied, 4 = minor gaps, 3 = notable gaps, 2 = mostly not, 1 = barely, 0 = not satisfied).

**Scored blind** as Arm B before the HAWP/no-HAWP mapping was revealed.

| # | Dimension | Score | Reasoning |
| --- | --- | --- | --- |
| 1 | Drift resistance | 5 | Caps respected exactly (5 blockers, 3 nice-to-haves); no feature proposals or roadmap beyond v1.0. |
| 2 | Evidence vs inference separation | 5 | Per-criterion Confirmed / Likely / Unclear with the weakest-link rule applied (C3 Likely despite a Confirmed foundation). |
| 3 | Output usefulness | 5 | Explicit criteria table, per-criterion status, a decidable go/no-go verdict, capped blockers, and Fix/Verify/Defer. |
| 4 | Handoff quality | 5 | Scope note states branch, working-tree state, commands run, and how absence claims are bounded; Verify-next names the exact checks still owed. |
| 5 | Trustworthiness | 5 | Hit the "not enough evidence" failure case correctly — public-safety gate marked Cannot assess / Unclear rather than guessed. |
| 6 | Scope adherence | 5 | Judged readiness only against documented-or-labeled-assumed criteria; no generic release checklist silently imported. |
| 7 | Completeness / coverage | 4 | Comparable coverage and found the public-safety criterion the other arm missed, but the blocker/nice-to-have caps left the process gaps (untriaged-defect backlog rows) to the other arm. |
| 8 | Conciseness / signal-to-noise | 4 | High-signal, but two tables plus per-criterion prose add length. |
| 9 | Correctness / accuracy | 5 | Claims verified via `git ls-tree`, tag list, repo-wide search, and commit counts. |
| 10 | False-positive control | 5 | Bounded wording and an explicit "cannot assess" instead of manufacturing a verdict; no over-claiming. |
| 11 | Verifiability | 5 | Highly re-checkable — every absence claim is tied to a named command (tag list, `git ls-tree`, repo-wide search) a reader can re-run. |
| 12 | Positive confirmation / balance | 4 | Records verified positives per criterion (C1/C2/C4 Pass with evidence), but woven into the criteria table rather than a broad at-a-glance positives section. |

## Notable strengths

The adopter-first lens drove the run's most important finding — `core/providers/` is absent on `main`, so the advertised Stable install path would fail — which the unguided run missed by assuming "Distribution works." Strongest on discipline, trustworthiness, correctness, and verifiability. The caps cost it some breadth on completeness and positive-confirmation. (Delivery note: the artifact needed one follow-up to be output verbatim — a round-trip cost, not an artifact-quality issue, so not scored here.)

**Raw:** 57 / 60 → **Headline 14.25 / 15** · **Percentage 95%**
