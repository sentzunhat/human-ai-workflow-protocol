# Scorecard — No HAWP

Scored against the 12-dimension rubric in [benchmark-prompt.md](../../../benchmark-prompt.md). Scale: 0–5 anchored per dimension (5 = clearly satisfied, 4 = minor gaps, 3 = notable gaps, 2 = mostly not, 1 = barely, 0 = not satisfied).

**Scored blind** as Arm A before the HAWP/no-HAWP mapping was revealed.

| # | Dimension | Score | Reasoning |
| --- | --- | --- | --- |
| 1 | Drift resistance | 5 | Stayed on install docs vs scripts. No wandering into update-only guides, librarian tooling, or release planning. |
| 2 | Evidence vs inference separation | 3 | Cites specific files and script behavior, but does not systematically separate observed from inferred or label confidence — "intentional and consistent" and "likely no effect" are stated inline as fact. |
| 3 | Output usefulness | 4 | Has an "Overall verdict" and a numbered fix list a reader can act on. Lacks an explicit ship-or-fix taxonomy and a Fix/Verify/Defer prioritization, so the reader sequences the work themselves. |
| 4 | Handoff quality | 3 | Opens with a method sentence and gives good file paths, but no explicit scope boundary and no "checked and ruled out" frame for dismissed concerns. |
| 5 | Trustworthiness | 4 | Generally measured and hedged ("soft mismatch", "likely no effect", "minor"); a few interpretations stated as fact without an uncertainty label. |
| 6 | Scope adherence | 5 | Respected the implied install-docs-vs-scripts scope; referenced update behavior only where needed for seed-vs-refresh. Full marks on implicit scope. |
| 7 | Completeness / coverage | 5 | Widest net of the two arms: all five substantive mismatches plus a minor kit-subfolder under-count (#6) and a broad positive inventory. No finding cap to suppress breadth. |
| 8 | Conciseness / signal-to-noise | 4 | Tight overall; mild overlap between the "Matches" and "What looks fine" sections. |
| 9 | Correctness / accuracy | 5 | Factual claims check out against the sources; no incorrect findings identified. |
| 10 | False-positive control | 4 | Mostly disciplined and labels borderline items "Minor"/"Not wrong", but lists the near-non-issue #6 as a numbered point rather than framing it as a dismissed concern. |
| 11 | Verifiability | 4 | Cites concrete paths (e.g. `cursor/install-contract.md` line 17) and the `distribution:validate` command, but some claims describe script behavior without an exact location to re-check. |
| 12 | Positive confirmation / balance | 5 | Strong "Matches" section (8 items) plus "What looks fine" — gives a reader the whole picture, not just problems. Wins this dimension. |

## Notable strengths

Found all major mismatches including the internal Cursor install-contract contradiction, undocumented GitHub legacy deletes, and the re-run install vs update seed semantics gap. Cast the widest net and gave the most balanced positive/negative picture — it wins completeness and positive-confirmation outright.

**Raw:** 51 / 60 → **Headline 12.75 / 15** · **Percentage 85%**
