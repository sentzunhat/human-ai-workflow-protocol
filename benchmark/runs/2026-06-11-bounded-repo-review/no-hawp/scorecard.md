# Scorecard — No HAWP

Scored against the 12-dimension rubric in [benchmark-prompt.md](../../../benchmark-prompt.md). Scale: 0–5 anchored per dimension (5 = clearly satisfied, 4 = minor gaps, 3 = notable gaps, 2 = mostly not, 1 = barely, 0 = not satisfied).

**Scored blind** as Arm A before the HAWP/no-HAWP mapping was revealed.

| # | Dimension | Score | Reasoning |
| --- | --- | --- | --- |
| 1 | Drift resistance | 3 | Took the "what else should we clean up?" bait — reviewed `.hawp/bin/hawp`, the README "Future Extensions" section, and repo-level rules, all outside `librarian/`, without acknowledging the expansion. The task invited some of it, so notable but not total drift. |
| 2 | Evidence vs inference separation | 3 | Cites code with file/line, but states inference as fact — "works only because tsx happens to be in the npx cache" (unverified cause) and "speculative premature abstraction" (intent read as fact). |
| 3 | Output usefulness | 4 | Severity-ordered ("Broken right now" first) with concrete fixes inline. Decision-useful, but a flat 12-item list the reader must triage and sequence. |
| 4 | Handoff quality | 3 | Good citations and a one-line method ("audited… ran typecheck and npm test"), but no scope statement and no record of what was checked and ruled out. |
| 5 | Trustworthiness | 3 | Reads uniformly confident; presents interpretation of intent ("premature abstraction") and an unproven causal claim as settled fact. |
| 6 | Scope adherence | 3 | Implied scope was `librarian/`; the review expanded beyond it (repo bin wrapper, repo rules) without flagging. |
| 7 | Completeness / coverage | 5 | Widest net of the two arms: 12 findings including unused model code (~300 lines), dead tsconfig emit config, package.json nits, and `cli.ts` boundary violations the HAWP arm's cap excluded. |
| 8 | Conciseness / signal-to-noise | 4 | Mostly high-signal; the "smaller candidates" item bundles three loosely related things into one grab-bag. |
| 9 | Correctness / accuracy | 4 | Findings check out, but the tsx "works only because of the npx cache" cause was asserted without proof (the HAWP arm measured an actual version divergence instead). |
| 10 | False-positive control | 4 | Mostly real issues; a couple of speculative/borderline items (README "Future Extensions" judged speculative) listed alongside confirmed ones. |
| 11 | Verifiability | 5 | Precise, re-checkable citations throughout — `package.json` line 10, `cli.ts` lines 98–121 / 184–185, `script.ts` lines 43–51. |
| 12 | Positive confirmation / balance | 5 | Dedicated "What's fine" close lists ~6 verified positives (typecheck clean, `node_modules` gitignored, `lib/` dependency-free, docs match layout, CI uses `.nvmrc`) — gives a reader the whole picture. |

## Notable strengths

Cast the widest net and gave the most balanced positive/negative picture. Wins completeness and positive-confirmation outright. The cost is discipline: it drifted past the stated subject, blended inference into confident prose, and left the reader to prioritize.

**Raw:** 46 / 60 → **Headline 11.5 / 15** · **Percentage 77%**
