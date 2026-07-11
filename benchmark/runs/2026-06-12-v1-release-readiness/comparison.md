# Detailed Comparison — No HAWP vs HAWP

This file compares the two outputs of this run in detail. For the raw material, see [`no-hawp/output.md`](no-hawp/output.md) and [`hawp/output.md`](hawp/output.md); for per-arm scores, see the `scorecard.md` file in each folder.

Both runs used the same model, the same repository state, and isolated sessions. The only variable was the prompt: a bare request vs a filled HAWP shape.

---

## How the two outputs are built

| Aspect | No HAWP | HAWP |
| --- | --- | --- |
| Overall form | Verdict-first essay: short answer, what's good, five numbered needs, ordered checklist | Decision artifact: scope note, explicit criteria table (documented vs assumed), per-criterion status table, verdict, capped blockers/nice-to-haves, Fix/Verify/Defer |
| What "ready" means | Never defined — implied by the items chosen | Seven explicit criteria, each labeled documented (with the source) or assumed |
| Confidence signals | None | Confirmed / Likely / Unclear per criterion, with the weakest-link rule applied (C3 is Likely despite a Confirmed foundation) |
| Handling of unknowables | Not addressed | One criterion (public-safety gate) explicitly marked "Cannot assess / Unclear" |
| Verdict | "Not yet — but it's close on substance" | "No-go today," explicitly framed as a decidable verdict rather than an evidence gap |
| Next steps | Ordered four-step checklist | Fix first / Verify next / Defer tied to specific criteria |

Both verdicts agree. The difference is that the HAWP artifact shows its reasoning structure — a reader can challenge any single criterion without re-deriving the rest.

---

## What each run found

### Found by both (independently)

1. No git tags, no CHANGELOG, no defined versioning — "v1.0" is currently undefined.
2. The protocol self-describes as v0.1 (`spec.md` title, `start-here.md` "locked to v0.1"), which a v1.0 label would contradict unless explained.
3. `dev` is 68 commits ahead of `main`, while the README points "Stable" users at `main`.
4. The large uncommitted working tree (kit `references/` restructure plus the benchmark run) must land before any release cut.
5. All of the repo's own quality gates pass on the current tree (both runs executed them rather than trusting the README).
6. The Node 20 / `engines: >=20` mismatch in `librarian/` (both flagged it as non-blocking polish).

### Found only by the no-HAWP run

- The previous benchmark run's defects have **no backlog rows**, against the repo's own rule that work started outside the loop still gets a row — a genuine process gap.
- The benchmark README's installation note reads stale/inverted relative to the actual install flow.
- Live re-verification of two earlier findings (reproduced the Node 20 test failure; located the hardcoded `v1.1.0` string at its current line).

### Found only by the HAWP run

- **The run's most important finding:** `core/providers/` does not exist on `main` (verified via `git ls-tree`), yet the generated Stable install guides pin `REF="main"` and error when the provider pack is missing. The advertised stable channel would fail for a first-time adopter — and GitHub's default branch shows the old product entirely. The unguided run said "Distribution works" and treated the merge as routine mechanics.
- The public-safety checklist (`reviews/public-safety-checklist.md`) exists as a named gate but has no recorded run against the current state — marked honestly as unassessable.
- The explicit observation that the repo defines no release criteria at all, verified by repo-wide search rather than asserted.

### The takeaway from the overlap

The shared core (six findings) confirms equal investigation ability. The unique finds split along predictable lines: the unguided run noticed *process* details near the edges of the question, while the guided run's adopter-first lens drove it to verify the one claim that mattered most — what the stable channel actually delivers — instead of assuming it.

---

## Dimension-by-dimension (0–5 anchored, scored blind; no-HAWP / HAWP)

**Drift resistance (4 / 5).** Unlike run 1, the unguided output resisted the trap on its own — its extras are marked optional or deferrable. "Is X ready?" has a natural boundary that "tell me what's wrong" lacks. The HAWP run's caps were respected exactly (5 blockers, 3 nice-to-haves). HAWP edge.

**Evidence vs inference separation (3 / 5).** The same fact viewed two ways. The unguided run inferred "Distribution works" from validators passing — collapsing "sources match generated outputs" into "an adopter can install." The guided run separated them: drift checks Confirmed pass, install path Likely fail, end-to-end test still missing. HAWP.

**Output usefulness (5 / 5).** Both are strongly actionable; the unguided ordered checklist is genuinely good and the HAWP artifact's explicit criteria are equally decision-useful. Tie at full marks.

**Handoff quality (3 / 5).** The unguided run names the commands it ran but not what it skipped — a continuing agent cannot tell that `main`'s contents were never inspected. The HAWP scope note states branch, working-tree state, what was executed, and what is still owed. HAWP.

**Trustworthiness (4 / 5).** The built-in failure case was a criterion whose correct answer is "not enough evidence." The HAWP run hit it (public-safety gate → Cannot assess / Unclear). The unguided run never mentioned that gate, and its one overconfident claim ("Distribution works") is the one a release decision leans on hardest. HAWP edge.

**Scope adherence (4 / 5).** Scored symmetrically: both stayed within the implied "is it ready / what to do" scope and proposed no new features. The unguided run's polish items edge toward a wishlist; the HAWP run respected explicit caps. HAWP edge.

**Completeness / coverage (5 / 4).** The unguided run noticed *process* details near the edges — the previous benchmark's defects have no backlog rows, a stale install note — that the HAWP caps left out. **no-HAWP wins.**

**Conciseness / signal-to-noise (4 / 4).** The unguided "smaller polish" section and the HAWP two-table structure roughly offset. Tie.

**Correctness / accuracy (5 / 5).** Both verified their claims by running gates and counting commits; no factual errors. Tie.

**False-positive control (4 / 5).** The unguided run flags some optional polish as needs; the HAWP run uses bounded wording and an explicit "cannot assess" rather than manufacturing a verdict. HAWP edge.

**Verifiability (4 / 5).** Both cite re-runnable checks. The HAWP run ties every absence claim to a named command (`git ls-tree`, tag list, repo-wide search) a reader can reproduce; a few unguided figures lack an exact path. HAWP edge.

**Positive confirmation / balance (5 / 4).** The unguided "What's already in good shape" section gives a broad at-a-glance positive picture; the HAWP positives are woven into the criteria table. **no-HAWP wins.**

---

## Which one is better, for now

**HAWP — but by a smaller margin than run 1.**

| Arm | Raw / 60 | Headline / 15 | Percentage |
| --- | --- | --- | --- |
| No HAWP | 50 | 12.5 | 83% |
| HAWP | 57 | 14.25 | 95% |

The HAWP output scored **~12 percentage points higher (1.75 / 15)** — HAWP wins 7 dimensions, ties 3, and loses 2 (completeness and positive-confirmation, both to the no-HAWP arm's wider net). The gap is smaller than the scope-creep-trap run (~18 points) and concentrated in evidence separation, handoff, and verifiability — the dimensions a release decision depends on most. Treat the number as a summary of qualitative judgments, not a measurement.

## Bottom line

Same verdict from both arms — not ready — so on this task the framing did not change the answer. It changed what the answer rests on. The unguided run produced a strong essay with one load-bearing overconfident claim; the guided run produced a criteria-backed decision record that caught the broken stable channel, admitted the one thing it could not assess, and named the verification still owed. The gap was narrower than in the first benchmark run, which itself is useful data: clear yes/no questions partially self-bound, while open-ended "what's wrong?" requests benefit most from the shape. One process cost on the HAWP side is recorded honestly: the agent needed a follow-up message to deliver the artifact it had already produced.
