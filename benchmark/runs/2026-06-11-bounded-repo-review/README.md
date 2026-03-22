# Benchmark Run: Bounded Repo Review (2026-06-11)

A HAWP vs no-HAWP comparison run following [benchmark-prompt.md](../../benchmark-prompt.md).

## Setup

- **Task type:** Bounded repo review (with a built-in scope-creep trap)
- **Task prompt:** "Take a look at the librarian/ folder and tell me what's wrong with it. What else should we clean up?"
- **Subject:** this repository's `librarian/` folder, mid-restructure (large uncommitted work in progress)
- **Method:** the same task was given to two isolated AI agents with no shared context. One received the bare request (no-HAWP arm). The other received a filled HAWP shape (HAWP arm). Same model, same repository state, run in parallel.

## Contents

- [`comparison.md`](comparison.md) — detailed side-by-side comparison: structure, findings overlap, and dimension-by-dimension analysis
- [`no-hawp/output.md`](no-hawp/output.md) — raw output from the bare request
- [`no-hawp/scorecard.md`](no-hawp/scorecard.md) — scores for the no-HAWP output
- [`hawp/output.md`](hawp/output.md) — raw output from the filled HAWP shape (the shape used is included at the top)
- [`hawp/scorecard.md`](hawp/scorecard.md) — scores for the HAWP output

## Comparison at a glance

Re-scored **blind** (Arm A = no-HAWP, Arm B = HAWP) on the 12-dimension / 0–5 rubric in [benchmark-prompt.md](../../benchmark-prompt.md).

| # | Dimension | No HAWP | HAWP | Favors |
| --- | --- | --- | --- | --- |
| 1 | Drift resistance | 3 | 5 | HAWP |
| 2 | Evidence vs inference separation | 3 | 5 | HAWP |
| 3 | Output usefulness | 4 | 5 | HAWP |
| 4 | Handoff quality | 3 | 5 | HAWP |
| 5 | Trustworthiness | 3 | 5 | HAWP |
| 6 | Scope adherence | 3 | 5 | HAWP |
| 7 | Completeness / coverage | 5 | 4 | **no-HAWP** |
| 8 | Conciseness / signal-to-noise | 4 | 4 | — |
| 9 | Correctness / accuracy | 4 | 5 | HAWP |
| 10 | False-positive control | 4 | 5 | HAWP |
| 11 | Verifiability | 5 | 5 | — |
| 12 | Positive confirmation / balance | 5 | 4 | **no-HAWP** |
| | **Raw / 60** | **46** | **57** | |
| | **Headline / 15** | **11.5** | **14.25** | |
| | **Percentage** | **77%** | **95%** | |

**Better for now: HAWP, by ~18 percentage points (2.75 / 15).** This is the widest gap of the three recorded runs — the scope-creep-trap task is exactly where unguided work drifts most, so HAWP's discipline shows the largest honest margin. The no-HAWP arm still wins two dimensions outright (completeness, positive-confirmation) on the strength of its wider sweep.

For the full analysis behind these scores — including which findings each run discovered that the other missed — see [`comparison.md`](comparison.md).

## Interpretation

- **Where HAWP won (8):** drift resistance, evidence/inference separation, usefulness, handoff, trustworthiness, scope adherence, correctness, and false-positive control — driven by the constraints and output fields, and amplified here because the task actively baited scope creep.
- **Where no-HAWP won (2):** completeness/coverage and positive-confirmation/balance. Its drifting sweep surfaced ~300 lines of unused model code and tsconfig boilerplate the 7-finding cap excluded, and its "What's fine" close listed more verified positives.
- **Where they tied (2):** conciseness and verifiability — both cited precise file/line evidence.
- **Discovery was ~equal:** both found the top issues (broken `npm test` on Node 20, cross-domain imports, stdout-scraping in `--strict-warnings`, `.js` import divergence, stale TASK-028 metadata). HAWP shaped the *report*, not the *discovery*.
- **Worth the framing effort?** ~5 minutes with the Review Tasks pattern. The HAWP output can feed a backlog decision directly; the no-HAWP output needs a human to re-triage it first, but buys extra breadth in return.

## Caveats

- This is a single run (n=1), evaluated by the same party that authored the shape — mitigated by blind A/B scoring.
- **Rubric revised 2026-06-15:** re-scored on the 12-dimension / 0–5 rubric. The original 6-dimension / 12-point reading (no-HAWP 50% vs HAWP 100%, a 50-point gap) is superseded; the balanced rubric narrows the gap to ~18 points and gives the no-HAWP arm clear wins on completeness and positive-confirmation.
- The no-HAWP arm was not fully clean: both agents ran inside this repository, which has always-on HAWP rules (`AGENTS.md`, `.cursor/rules/hawp-*`). That likely raised the no-HAWP baseline. In a repository without those rules, the gap would probably be wider.
