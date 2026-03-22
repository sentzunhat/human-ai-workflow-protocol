# Benchmark Run: Install Docs Truth Audit (2026-06-15)

A HAWP vs no-HAWP comparison run following [benchmark-prompt.md](../../benchmark-prompt.md).

## Setup

- **Task type:** Standards / truth audit
- **Task prompt:** "Check whether our install docs match what the scripts actually do."
- **Subject:** `distribution/generated/*/install/` guides vs composed bash in `distribution/sources/`
- **Method:** Same task given to two isolated AI sessions. No-HAWP arm ran in a clean `/tmp/` workspace copy (agent rules stripped). HAWP arm ran with a filled audit shape. Same model. **Scored blind** on the 12-dimension / 0–5 rubric: arms labeled A/B, all dimensions scored before the HAWP mapping was revealed.

## Contents

- [`comparison.md`](comparison.md) — detailed side-by-side comparison
- [`no-hawp/output.md`](no-hawp/output.md) — raw output from the bare request
- [`no-hawp/scorecard.md`](no-hawp/scorecard.md) — scores for the no-HAWP output
- [`hawp/output.md`](hawp/output.md) — filled shape + raw output from the guided run
- [`hawp/scorecard.md`](hawp/scorecard.md) — scores for the HAWP output

## Comparison at a glance

Scored **blind** (Arm A = no-HAWP, Arm B = HAWP, revealed after scoring) on the 12-dimension / 0–5 rubric in [benchmark-prompt.md](../../benchmark-prompt.md).

| # | Dimension | No HAWP | HAWP | Favors |
| --- | --- | --- | --- | --- |
| 1 | Drift resistance | 5 | 5 | — |
| 2 | Evidence vs inference separation | 3 | 5 | HAWP |
| 3 | Output usefulness | 4 | 5 | HAWP |
| 4 | Handoff quality | 3 | 5 | HAWP |
| 5 | Trustworthiness | 4 | 5 | HAWP |
| 6 | Scope adherence | 5 | 5 | — |
| 7 | Completeness / coverage | 5 | 4 | **no-HAWP** |
| 8 | Conciseness / signal-to-noise | 4 | 4 | — |
| 9 | Correctness / accuracy | 5 | 5 | — |
| 10 | False-positive control | 4 | 5 | HAWP |
| 11 | Verifiability | 4 | 5 | HAWP |
| 12 | Positive confirmation / balance | 5 | 4 | **no-HAWP** |
| | **Raw / 60** | **51** | **57** | |
| | **Headline / 15** | **12.75** | **14.25** | |
| | **Percentage** | **85%** | **95%** | |

**Better for now: HAWP, by ~10 percentage points (1.5 / 15).** HAWP wins 6 dimensions (the report-discipline cluster plus verifiability), ties 4, and **loses 2** — completeness/coverage and positive-confirmation/balance, both won by the no-HAWP arm's wider net. The legacy-rubric reading (25 points) is superseded; rebalancing dropped the structurally HAWP-only dimension and added counterweights the unguided arm wins.

For the full analysis, see [`comparison.md`](comparison.md).

## Interpretation

- **Where HAWP won (6):** evidence/inference separation, output usefulness, handoff quality, trustworthiness, false-positive control, and verifiability — the disciplined-report cluster plus re-checkability, driven by the audit constraints and output shape.
- **Where no-HAWP won (2):** completeness/coverage and positive-confirmation/balance. The unguided arm cast a wider net (minor kit-subfolder under-count, an 8-item "Matches" inventory); HAWP's 5-finding and 2-non-finding caps deliberately traded breadth for focus.
- **Where they tied (4, full or equal marks):** drift resistance, scope adherence, conciseness, correctness. Discovery of the five substantive mismatches was identical — HAWP shaped the *report*, not the *findings*.
- **Task-type note:** Standards audits are less drift-prone than open-ended reviews, so the no-HAWP baseline was high (85%). HAWP's value here is report discipline and decision framing, not finding things the other run missed.
- **Worth the framing effort?** ~5 minutes to fill the Audit Tasks pattern. The HAWP output can go directly to a "fix docs before ship" decision; the no-HAWP output is already usable but needs a reader to extract the verdict, and it accepts a couple of extra low-value items in exchange for broader coverage.

## Caveats

- Single run (n=1). Scored by the same party that authored the HAWP shape, mitigated by blind A/B scoring (labels hidden until after scores were assigned).
- **Rubric revised 2026-06-15:** re-scored on the rebalanced 12-dimension / 0–5 rubric (out of 15 + percentage). The earlier 6-dimension / 12-point score (no-HAWP 75% vs HAWP 100%, a 25-point gap) is superseded; the balanced rubric narrows the gap to ~10 points and gives the no-HAWP arm clear wins on completeness/coverage and positive-confirmation/balance.
- **No-HAWP arm was clean:** ran in `/tmp/hawp-benchmark-clean-20260615-211005` with `AGENTS.md`, `.cursor/rules/`, `.continue/rules/`, and `.github/instructions/` stripped from the copy only. Source repo untouched.
- **HAWP arm:** executed in the same session as the benchmark orchestrator (same model, same repo state). Ideally a fully separate fresh chat in the original repo window; content should be equivalent since both read the same tree.
- **No delivery issues:** both arms returned full artifacts in the first response (no "output verbatim" follow-up needed).
- Neither arm executed install scripts end-to-end in a scratch repo; comparison is doc-vs-source static analysis plus `distribution:validate`.

## Run checklist

- [x] Same task prompt used for both arms
- [x] No-HAWP arm: clean workspace + fresh chat
- [x] HAWP arm: filled shape saved in output.md
- [x] Both outputs saved under `benchmark/runs/2026-06-15-install-docs-truth-audit/`
- [x] Both scorecards written with reasoning
- [x] `comparison.md` and `README.md` written
