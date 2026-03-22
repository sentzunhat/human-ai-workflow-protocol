# Benchmark Run: v1.0 Release Readiness (2026-06-12)

A HAWP vs no-HAWP comparison run following [benchmark-prompt.md](../../benchmark-prompt.md). Second recorded run; see [2026-06-11-bounded-repo-review](../2026-06-11-bounded-repo-review/README.md) for the first.

## Setup

- **Task type:** Vague open-ended question (with two built-in failure cases: an unlimited-wishlist trap in "what do we still need to do?", and at least one criterion whose correct answer is "not enough evidence to decide")
- **Task prompt:** "Is HAWP ready for a v1.0 release? What do we still need to do?"
- **Subject:** this repository itself, mid-restructure (38 uncommitted paths, `dev` 68 commits ahead of `main`)
- **Method:** the same task was given to two isolated AI agents with no shared context. One received the bare request (no-HAWP arm). The other received a filled HAWP shape (HAWP arm). Same model, same repository state, run in parallel.

## Contents

- [`comparison.md`](comparison.md) — detailed side-by-side comparison: structure, findings overlap, and dimension-by-dimension analysis
- [`no-hawp/output.md`](no-hawp/output.md) — raw output from the bare request
- [`no-hawp/scorecard.md`](no-hawp/scorecard.md) — scores for the no-HAWP output
- [`hawp/output.md`](hawp/output.md) — the filled shape, then the raw output
- [`hawp/scorecard.md`](hawp/scorecard.md) — scores for the HAWP output

## Comparison at a glance

Re-scored **blind** (Arm A = no-HAWP, Arm B = HAWP) on the 12-dimension / 0–5 rubric in [benchmark-prompt.md](../../benchmark-prompt.md).

| # | Dimension | No HAWP | HAWP | Favors |
| --- | --- | --- | --- | --- |
| 1 | Drift resistance | 4 | 5 | HAWP |
| 2 | Evidence vs inference separation | 3 | 5 | HAWP |
| 3 | Output usefulness | 5 | 5 | — |
| 4 | Handoff quality | 3 | 5 | HAWP |
| 5 | Trustworthiness | 4 | 5 | HAWP |
| 6 | Scope adherence | 4 | 5 | HAWP |
| 7 | Completeness / coverage | 5 | 4 | **no-HAWP** |
| 8 | Conciseness / signal-to-noise | 4 | 4 | — |
| 9 | Correctness / accuracy | 5 | 5 | — |
| 10 | False-positive control | 4 | 5 | HAWP |
| 11 | Verifiability | 4 | 5 | HAWP |
| 12 | Positive confirmation / balance | 5 | 4 | **no-HAWP** |
| | **Raw / 60** | **50** | **57** | |
| | **Headline / 15** | **12.5** | **14.25** | |
| | **Percentage** | **83%** | **95%** | |

**Better for now: HAWP, by ~12 percentage points (1.75 / 15).** A clear yes/no question self-bounds the unguided run well (it resisted the wishlist trap and ran the gates), so HAWP's margin is smaller than the scope-creep-trap run. HAWP wins 7 dimensions, ties 3, and loses 2 (completeness, positive-confirmation) to the no-HAWP arm's wider net.

For the full analysis behind these scores, see [`comparison.md`](comparison.md).

## Interpretation

- **The gap was narrower than in run 1.** The unguided run was strong: it reached the same overall verdict (not ready), ran the repo's own validators, and resisted the wishlist trap on its own. A clear yes/no question seems to bound unguided work better than an open "what's wrong?" request does.
- **The shape's lens produced the run's most important unique finding.** The mission field framed readiness as "what a first-time adopter would actually receive." That led the HAWP run to inspect `main` directly and discover that the advertised Stable install path points at a branch missing `core/providers/` entirely — so the published install guides would error. The unguided run said "Distribution works" and treated the dev-to-main merge as routine mechanics.
- **The "not enough evidence" failure case worked as designed.** The HAWP run marked the public-safety-checklist criterion "Cannot assess / Unclear" — the correct answer. The unguided run never mentioned that gate.
- **The unguided run still found things the guided run missed:** the untriaged defects from the previous benchmark run having no backlog rows, and a stale installation note in the benchmark README.

## Caveats

- This is a single run (n=1), evaluated by the same party that authored the shape — mitigated by blind A/B scoring.
- **Rubric revised 2026-06-15:** re-scored on the 12-dimension / 0–5 rubric. The original 6-dimension / 12-point reading (no-HAWP 75% vs HAWP 100%, a 25-point gap) is superseded; the balanced rubric narrows the gap to ~12 points and gives the no-HAWP arm clear wins on completeness and positive-confirmation.
- The no-HAWP arm was not fully clean: both agents ran inside this repository, which has always-on HAWP rules (`AGENTS.md`, `.cursor/rules/hawp-*`). That likely raised the no-HAWP baseline.
- **Delivery hiccup in the HAWP arm:** the agent produced the full artifact internally but its first final response only referred to it instead of including it. One follow-up message ("output the full artifact verbatim") retrieved it. The artifact content was scored as delivered after that retry; the hiccup is a real cost worth counting against round-trip efficiency, though it concerns the delivery channel rather than the artifact's quality.
