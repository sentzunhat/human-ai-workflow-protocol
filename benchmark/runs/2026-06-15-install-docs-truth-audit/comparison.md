# Detailed Comparison — No HAWP vs HAWP

This file compares the two outputs of this run in detail. For the raw material, see [`no-hawp/output.md`](no-hawp/output.md) and [`hawp/output.md`](hawp/output.md); for per-arm scores, see the `scorecard.md` file in each folder.

Both runs used the same model, the same repository state, and isolated sessions. The only variable was the prompt: a bare request vs a filled HAWP shape.

**Task prompt (both arms):** "Check whether our install docs match what the scripts actually do."

---

## How the two outputs are built

| Aspect | No HAWP | HAWP |
| --- | --- | --- |
| Overall form | Narrative audit: verdict paragraph, "Matches" section, numbered mismatches, "What looks fine," suggested fixes | Structured artifact: scope/method header, explicit verdict, 5 capped findings with category + confidence + Observed/Inference/Significance, 2 non-findings, Fix/Verify/Defer |
| Per-issue structure | Prose with inline file references | Fixed template per finding; doc claim and script behavior always paired |
| Confidence signals | None | Confirmed / Likely on every finding |
| Positive confirmation | Dedicated "Matches" and "What looks fine" sections | Two explicit non-findings (provider refresh alignment; CI validate pass) |
| Decision framing | "Mostly match, but…" in opening; fixes listed at end | "Partially trustworthy" verdict with one-sentence rationale; operational sequence at end |
| Scope handling | Implicit install focus; mentioned update only where relevant | Explicit install-only scope in shape; update referenced only for seed-vs-refresh |

---

## What each run found

### Found by both (independently)

These five issues were discovered by both runs, suggesting equal underlying investigation ability — HAWP changed the report format, not core discovery.

1. **Cursor install contract vs script for `AGENTS.md`** — install-contract says refresh on every install; boundaries/safety/script say seed-if-missing.
2. **Global no-clobber claim in `safety.md` overgeneralizes** — kit uses `rm -rf .hawp/kit`; overlays use plain `cp`.
3. **GitHub legacy filename cleanup undocumented** — script deletes `human-ai-workflow-protocol-*` files; no guide prose mentions it.
4. **Orphan retirement trigger misdescribed** — docs say "data row"; script gates on Status Key table rows too.
5. **Re-run install vs update for seed files** — install continues with seed semantics when `.hawp/` exists; docs say use update for full refresh.

### Found only by the no-HAWP run

- **Kit `instructions/` and `standards/` not named in prose summaries** — minor completeness gap under "What Was Added" bullets (script copies them; `.hawp/kit/**` covers it).
- **Auto-dispatch four-backtick fence pattern** — noted as intentional design in "What looks fine."
- **Broader positive inventory** — dedicated "Matches" section listing kit refresh, work scaffold, GitHub copilot seeding, Continue overlay, legacy migrations, HAWP_LOCAL_CORE, README mapping.

### Found only by the HAWP run

- **Explicit `distribution:validate` as validation evidence** — ran and recorded as Non-Finding 2 with CI implication.
- **Formal non-findings** for GitHub/Continue refresh alignment (confirmed doc-script match where no-HAWP listed these under "Matches" without a "checked and ruled out" frame).
- **Verdict taxonomy** — "partially trustworthy" with explicit guidance: trust copy/paste blocks, don't trust all prose without checking boundaries.

### The takeaway from the overlap

Discovery was nearly identical on the substantive mismatches. The no-HAWP run added minor positive detail and one small prose-completeness nit; the HAWP run added formal validation evidence and explicit non-findings. Neither run missed the headline Cursor contract bug.

---

## Dimension-by-dimension (0–5 anchored, scored blind; no-HAWP / HAWP)

**Drift resistance (5 / 5).** Both stayed on install docs vs scripts. Neither drifted into release readiness, librarian review, or benchmark meta-commentary. Tie.

**Evidence vs inference separation (3 / 5).** The no-HAWP run cites real files and script lines but does not systematically label confidence or separate observed from inferred ("intentional and consistent" stated as fact). The HAWP run does both on every finding. HAWP.

**Output usefulness (4 / 5).** Both are actionable and both carry a verdict. The no-HAWP fix list maps 1:1 to findings but has no ship-or-fix taxonomy or prioritization; the HAWP run adds Fix first / Verify next / Defer. HAWP edge.

**Handoff quality (3 / 5).** The no-HAWP run has a method sentence but no explicit scope boundary or ruled-out frame. The HAWP run opens with scope + method and uses non-findings to record dismissed concerns. HAWP.

**Trustworthiness (4 / 5).** The no-HAWP run is generally measured and hedged, with a few interpretations stated as fact. The HAWP run uses Likely where user impact is inferred (F1, F4). HAWP edge.

**Scope adherence (5 / 5).** Scored symmetrically: both respected the implied install scope and referenced update behavior only where seed-vs-refresh required it. Full marks both — this dimension no longer rewards merely *having* explicit constraints. Tie.

**Completeness / coverage (5 / 4).** The no-HAWP run cast the widest net — all five substantive mismatches plus a minor kit-subfolder under-count. The HAWP run's 5-finding cap deliberately excluded the minor item. **no-HAWP wins.**

**Conciseness / signal-to-noise (4 / 4).** Neither padded meaningfully. The no-HAWP "Matches"/"What looks fine" sections overlap slightly; the HAWP per-finding template repeats its structure. Tie.

**Correctness / accuracy (5 / 5).** Both arms' factual claims check out against the sources. Tie.

**False-positive control (4 / 5).** The no-HAWP run lists a borderline non-issue (#6) as a numbered point; the HAWP run uses bounded wording and two explicit non-findings to record dismissed concerns. HAWP edge.

**Verifiability (4 / 5).** Both cite the `install-contract.md` line and the validate command. The HAWP run pairs Observed (doc) with Observed (script) for every finding, making each claim slightly easier to re-check; a few no-HAWP claims describe script behavior without an exact location. HAWP edge.

**Positive confirmation / balance (5 / 4).** The no-HAWP "Matches" section (8 items) plus "What looks fine" gives the fuller picture of what is correct; the HAWP run's 2-non-finding cap records less. **no-HAWP wins.**

---

## Which one is better, for now

**HAWP — on this run, by a modest margin.**

| Arm | Raw / 60 | Headline / 15 | Percentage |
| --- | --- | --- | --- |
| No HAWP | 51 | 12.75 | 85% |
| HAWP | 57 | 14.25 | 95% |

The HAWP output scored **~10 percentage points higher (1.5 / 15)** — HAWP wins 6 dimensions, ties 4, and loses 2. This is much narrower than the legacy-rubric reading (25 points) because the rebalanced rubric dropped the structurally HAWP-only "constraint discipline" dimension (now a symmetric "scope adherence" both arms win) and added two counterweights — completeness/coverage and positive-confirmation/balance — that the **no-HAWP arm wins outright**. The remaining gap is concentrated in the correlated report-discipline cluster (evidence separation, handoff, trustworthiness, false-positive control) plus verifiability, not in discovery — the two arms found the same five substantive mismatches.

## Bottom line

Same model, same repository, same task. Both found the same substantive doc-script drifts. The unguided run produced a competent narrative audit; the guided run produced a bounded, labeled artifact with an explicit trust verdict and operational sequence. For standards/truth audits where the answer is "mostly yes, but fix these before trusting prose," HAWP's constraints added clear value without suppressing findings.
