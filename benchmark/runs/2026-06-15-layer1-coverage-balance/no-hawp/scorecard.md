# Scorecard — No HAWP (same-state, current tree)

Scored against the 12-dimension rubric in [benchmark-prompt.md](../../../benchmark-prompt.md). Scale: 0–5 anchored per dimension (5 = clearly satisfied, 4 = minor gaps, 3 = notable gaps, 2 = mostly not, 1 = barely, 0 = not satisfied).

**Scored as Arm A.** This is a real no-HAWP run executed in the clean `/tmp/` workspace against the **current** tree (Node 26), not the reused 2026-06-11 baseline — see [`output.md`](output.md).

| # | Dimension | Score | Reasoning |
| --- | --- | --- | --- |
| 1 | Drift resistance | 3 | Took the "what else should we clean up?" bait — a full "repo-wide" section (provider drift, CI coverage gap, benchmark hygiene) plus `.hawp/bin/hawp` (#8), all outside `librarian/`. It labels the section as outside scope, so notable but acknowledged drift. |
| 2 | Evidence vs inference separation | 3 | Strong file/line citations, but states inference as fact: "it works only because `npx` resolves tsx from cache" (unproven cause) and "speculative scaffolding" (intent read as fact). |
| 3 | Output usefulness | 4 | Severity-ordered with concrete fixes, but a flat 14-item list plus repo-wide section the reader must triage and sequence. |
| 4 | Handoff quality | 3 | Good method line (ran all five validators) and a "What's fine" close, but no scope statement and no explicit record of what was checked and ruled out. |
| 5 | Trustworthiness | 3 | Reads uniformly confident; presents an unproven causal claim (npx cache) and a misread rule (#5) as settled fact. |
| 6 | Scope adherence | 3 | Implied scope was `librarian/`; expanded into CI, provider sync, and `.hawp/bin/` without flagging the expansion as a limit. |
| 7 | Completeness / coverage | 5 | Widest net of any arm: 14 findings + repo-wide section + "What's fine". Surfaced valid items the HAWP arm missed — `cli.ts` doing filesystem work (#2, confirmed), a third repo-root finder (#9, confirmed), dead `tsconfig` emit config (#11, confirmed), test-only type guards (#12), the `validate:workflow` alias (#13), the `providers:validate` gap in `backlog validate` (#7), and the `scripts/README.md` test-discovery doc mismatch (#14). |
| 8 | Conciseness / signal-to-noise | 4 | Mostly high-signal; #14 bundles four loosely related things into one grab-bag and the repo-wide section adds length. |
| 9 | Correctness / accuracy | 4 | Most findings verified, but #5 misreads the rule (the cited files **are** `index.ts`, where `process.exit` is allowed), and the npx-cache cause is asserted without proof. |
| 10 | False-positive control | 3 | #5 is a genuine false positive (no rule is violated — `process.exit` in `index.ts` is exactly what the contract permits); "Future Extensions speculative" is also borderline. |
| 11 | Verifiability | 5 | Precise, re-checkable citations throughout (`script.ts:43–51`, `cli.ts:184–185`, `composition.ts`). |
| 12 | Positive confirmation / balance | 5 | Dedicated "What's fine" close lists ~7 verified positives (clean typecheck, 37 tests pass on Node 26, `node_modules` gitignored, dependency-free `lib/`, sound domain split). |

## Notable strengths

Cast the widest net of any recorded arm and surfaced six-plus valid findings the HAWP arm did not. The cost is discipline: it drifted past `librarian/`, blended unproven inference into confident prose, and shipped one outright false positive (#5).

**Raw:** 45 / 60 → **Headline 11.25 / 15** · **Percentage 75%**
