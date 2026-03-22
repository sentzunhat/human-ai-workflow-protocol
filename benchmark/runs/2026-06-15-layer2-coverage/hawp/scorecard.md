# Scorecard — HAWP (Layer 1 + Layer 2 shape)

Scored against the 12-dimension rubric in [benchmark-prompt.md](../../../benchmark-prompt.md). Scale: 0–5 anchored per dimension (5 = clearly satisfied, 4 = minor gaps, 3 = notable gaps, 2 = mostly not, 1 = barely, 0 = not satisfied).

**Scored as Arm B**, against the same-state no-HAWP run (Arm A, 45/60) on the current Node 26 tree. The librarian tree is unchanged from the Layer-1 run; only the HAWP authoring pattern gained the Layer-2 surface-sweep habit.

| # | Dimension | Score | Reasoning |
| --- | --- | --- | --- |
| 1 | Drift resistance | 5 | Stayed entirely inside `librarian/`; the "what else" bait was held to a 3-item out-of-scope list. |
| 2 | Evidence vs inference separation | 5 | Every primary finding splits Observed / Inference with a confidence label matched to the weakest part (F6 "Confirmed mechanism, Unclear impact"). |
| 3 | Output usefulness | 5 | Prioritized with a Fix first / Verify next / Defer sequence that feeds a backlog decision directly. |
| 4 | Handoff quality | 5 | The visible **surface-sweep table** maps each in-scope unit to its result — a reader can see exactly what was covered and continue without re-discovery. |
| 5 | Trustworthiness | 5 | Uncertainty explicit (F6 divergence reported as not reproducing → Unclear); no overclaiming. |
| 6 | Scope adherence | 5 | Implied scope held; out-of-scope items demoted to one-liners. |
| 7 | Completeness / coverage | 5 | **Recovered (4→5).** The Layer-2 sweep forced enumeration of all nine in-scope units and caught everything the same-state no-HAWP arm found — `cli.ts` filesystem work (F3), the `providers:validate` gap (F5), dead `tsconfig` (F7), the third repo-root finder (F2), and the type-guards / `validate:workflow` alias / README test-discovery mismatch (minor list) — **plus** items no-HAWP missed (the `--export-plan` apply gap, the engines-floor framing). Coverage now meets or exceeds the wide unguided net. |
| 8 | Conciseness / signal-to-noise | 4 | Held at 4. The fuller output (sweep table + 7 primary + 10 minor + 6 verified) is information-dense and well-organized, but it is more to scan and a couple of minor items overlap primary findings (path-leak skip, console.log convention). The predicted cost of the sweep + uncapped list. |
| 9 | Correctness / accuracy | 5 | Claims verified against source and live gates (typecheck, 37 tests, providers/distribution validate, version/engine alignment, import counts, line numbers). |
| 10 | False-positive control | 5 | All 7 primary findings confirmed real; minor items hedged or routed appropriately. Notably did **not** repeat the no-HAWP arm's false positive (`process.exit` in files that are `index.ts`) — the harder sweep did not manufacture non-issues. |
| 11 | Verifiability | 5 | Precise re-checkable citations throughout, plus the sweep table and re-runnable commands. |
| 12 | Positive confirmation / balance | 5 | Dedicated "Verified correct" list (6 items), the sweep records sound units explicitly, and it notes the old top finding is resolved. |

## Notable strengths

The Layer-2 **surface sweep** did exactly what Layer 1 alone could not: it lifted #7 from 4 to 5 by making coverage a visible checklist rather than trusting the lens. Combined with Layer 1's verified-correct list (which holds #12 at 5), the HAWP arm now matches the unguided arm's breadth while keeping its discipline — and without repeating the unguided arm's false positive. The only cost is conciseness holding at 4 (not dropping), the expected trade for the fuller artifact.

## Residual risk

The output is now near the upper bound of what stays high-signal; one more layer of required structure would likely start costing #8. The sweep table's value depends on the unit list being honest — a reviewer who lists a unit as "inspected" without truly inspecting it would convert this strength into a trust risk.

**Raw:** 59 / 60 → **Headline 14.75 / 15** · **Percentage 98%**
