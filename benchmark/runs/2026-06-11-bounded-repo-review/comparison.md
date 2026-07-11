# Detailed Comparison — No HAWP vs HAWP

This file compares the two outputs of this run in detail. For the raw material, see [`no-hawp/output.md`](no-hawp/output.md) and [`hawp/output.md`](hawp/output.md); for per-arm scores, see the `scorecard.md` file in each folder.

Both runs used the same model, the same repository state, and isolated sessions. The only variable was the prompt: a bare request vs a filled HAWP shape.

---

## How the two outputs are built

The clearest difference is structural, before reading a single finding.

| Aspect | No HAWP | HAWP |
| --- | --- | --- |
| Overall form | Severity-ordered list of 12 items, grouped under informal headings ("Broken right now", "Dead weight") | Review artifact with a stated scope and method, 7 capped findings, 2 non-findings, an out-of-scope list, and a closing action sequence |
| Per-finding structure | Prose paragraph with code citations | Category + confidence label + Observed + Inference + Significance, for every finding |
| Confidence signals | None | Every finding labeled Confirmed / Likely, including split labels (e.g. "Confirmed mechanism, Likely impact") |
| What was already checked | A closing "What's fine" paragraph | 2 explicit non-findings stating what was checked, what was observed, and why the concern was dismissed |
| Next steps | None — the reader must triage | Fix first / Verify next / Defer, each tied to specific findings |
| Scope handling | Silently expanded (reviewed `.hawp/bin/hawp`, commented on repo rules) | Stayed in `librarian/`; out-of-scope items flagged in 3 one-liners as the constraints allowed |

---

## What each run found

### Found by both (independently)

These five issues were discovered by both runs, which suggests the underlying investigation ability was equal — HAWP changed the report, not the discovery.

1. `npm test` is broken on Node 20 while `engines` claims `>=20` (the top finding in both).
2. Cross-domain import of `findRepoRoot` from `distribution/shared/composition` into `providers/materialize`, violating the documented boundary rule.
3. `--strict-warnings` works by regex-scraping another tool's stdout, so a cosmetic wording change silently disables it.
4. `.js`-suffixed relative imports in `backlog-upgrade/` diverge from the documented extensionless convention.
5. Stale TASK-028 references and a hardcoded `v1.1.0` version string in user-facing CLI output.

### Found only by the no-HAWP run

Mostly products of its wider (drifting) sweep:

- ~300 lines of unused "future API" model code in `models/evidence-report.ts`, referenced only by its own tests.
- Dead tsconfig emit configuration (`outDir`, `declaration`, etc.) that nothing uses, plus `exclude` entries for folders that do not exist.
- `package.json` nits: a duplicate script alias and an empty `dependencies` block.
- `.hawp/bin/hawp` argument-forwarding inconsistency between `backlog validate` and `backlog upgrade` (note: outside the `librarian/` scope).
- `cli.ts` boundary violations: filesystem access in `validate-hawp-workflow/cli.ts` and execution logic in `backlog-upgrade/cli.ts`, both against the folder's own README rules.
- Smaller items: a deprecated legacy-guides guard and three separate "find the repo root" helpers that could consolidate.

### Found only by the HAWP run

Mostly products of deeper verification within the bounded scope:

- A second, separate cross-domain import direction: `backlog-upgrade/script.ts` importing from `validate-hawp-workflow/` (the no-HAWP run caught only the materialize→distribution direction).
- `--apply --export-plan` silently ignores the export flag — the apply branch returns before the export handling, and `CLI.md` documents no caveat.
- The distribution path-leak check silently skips missing target files (`continue` with no warning), against the folder's documented "never swallow" rule.
- Measured evidence for the tsx version risk: `npx tsx` resolves v4.22.4 at the repo root vs the lockfile-pinned v4.21.0 inside `librarian/` (the no-HAWP run asserted the cause without measuring, and stated it more confidently than the evidence supported).
- Two verified non-findings: the apply-mode dirty-tree guard fails closed, and evidence-path containment works as documented.

### The takeaway from the overlap

Neither run strictly dominated on discovery. The no-HAWP run bought extra breadth with its drift; the HAWP run bought extra depth and verification with its discipline. The HAWP run's 7-finding cap is the direct reason it has no equivalent of the unused-model-code finding — that is feedback on how the constraints field was authored, not on the protocol itself.

---

## Dimension-by-dimension (0–5 anchored, scored blind; no-HAWP / HAWP)

**Drift resistance (3 / 5).** The trap in the task was "what else should we clean up?" The no-HAWP run took it: it reviewed `.hawp/bin/hawp`, judged the README's "Future Extensions" section, and invoked repo-level HAWP rules — all outside `librarian/`, with the expansion never acknowledged. The HAWP run routed the same impulse into a bounded "Out of scope, flagged only" list of exactly three one-liners. HAWP.

**Evidence vs inference separation (3 / 5).** Both runs cite real evidence. The difference is what happens at the edge of the evidence. Both noticed the `npx tsx` resolution risk; the no-HAWP run wrote "it currently works only because tsx happens to be in the npx cache" — an unverified cause stated as fact, while the HAWP run measured the version divergence and explicitly noted the cause "was not proven." HAWP.

**Output usefulness (4 / 5).** The no-HAWP output is genuinely useful, severity-ordered raw material, but a flat list the reader must sequence. The HAWP output ends with Fix first / Verify next / Defer and ties each finding to the stated decision (what to add to the backlog). HAWP edge.

**Handoff quality (3 / 5).** The no-HAWP run gives good citations and a one-line method but never states its scope or what it ruled out. The HAWP run opens with scope, method, and commands run, and its two non-findings record dismissed concerns. HAWP.

**Trustworthiness (3 / 5).** The no-HAWP run reads uniformly confident; its judgment that unused code is "speculative premature abstraction" presents intent as fact. The HAWP run downgrades thin-evidence findings (F7 Likely) and uses bounded absence wording. HAWP.

**Scope adherence (3 / 5).** Scored symmetrically against the implied `librarian/` scope: the no-HAWP run expanded past it (repo bin wrapper, repo rules) without flagging; the HAWP run stayed in and flagged out-of-scope items separately. HAWP.

**Completeness / coverage (5 / 4).** The no-HAWP run's wider sweep surfaced ~300 lines of unused model code, dead tsconfig emit config, and package.json nits the HAWP 7-finding cap excluded. **no-HAWP wins.**

**Conciseness / signal-to-noise (4 / 4).** Neither padded meaningfully; the no-HAWP "smaller candidates" grab-bag and the HAWP per-finding template repetition roughly offset. Tie.

**Correctness / accuracy (4 / 5).** Both arms' findings check out, but the no-HAWP tsx causal claim was asserted unproven where the HAWP run measured an actual version divergence. HAWP edge.

**False-positive control (4 / 5).** The no-HAWP run lists a couple of speculative/borderline items alongside confirmed ones; the HAWP run uses bounded wording and non-findings to avoid over-claiming. HAWP edge.

**Verifiability (5 / 5).** Both cite precise, re-checkable evidence — the no-HAWP run with exact file/line numbers, the HAWP run with files plus measured commands and version numbers. Tie.

**Positive confirmation / balance (5 / 4).** The no-HAWP "What's fine" close lists ~6 verified positives; the HAWP run records two deeper non-findings but a narrower positive picture. **no-HAWP wins.**

---

## Which one is better, for now

**HAWP — clearly, on this run.**

| Arm | Raw / 60 | Headline / 15 | Percentage |
| --- | --- | --- | --- |
| No HAWP | 46 | 11.5 | 77% |
| HAWP | 57 | 14.25 | 95% |

The HAWP output scored **~18 percentage points higher (2.75 / 15)** — HAWP wins 8 dimensions, ties 2, and loses 2 (completeness and positive-confirmation, both to the no-HAWP arm's wider sweep). This is the widest gap of the three recorded runs because the scope-creep-trap task is exactly where unguided work drifts most. Treat the number as a summary of qualitative judgments, not a measurement.

## Bottom line

Same model, same repository, same task. The unguided run produced a competent but unbounded brain dump that a human must triage; the guided run produced a bounded, evidence-labeled artifact that can feed a decision directly. The cost was about 5 minutes of shape authoring and a small loss of breadth at the finding cap. For this task type, the framing paid for itself — with the caveat that this is one run, scored by the same party that authored the shape.
