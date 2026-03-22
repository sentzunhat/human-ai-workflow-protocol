# HAWP arm, scoring, and cleanup

Run this in the **original repository** — the workspace where `AGENTS.md` and `.cursor/rules/` exist.

By now the no-HAWP arm should be done: `benchmark/runs/<date>-<task-name>/no-hawp/output.md` is saved and the clean `/tmp/` window is closed. If not, finish **[run.md](run.md)** first.

This file does three things in order: (1) run the HAWP arm, (2) score and record both arms, (3) remove the temp workspace.

For scoring rules and task-type ideas, see [benchmark-prompt.md](../benchmark-prompt.md).

---

## Before you start

- [ ] `benchmark/runs/<date>-<task-name>/no-hawp/output.md` exists in this repo
- [ ] The clean-workspace Cursor window is closed
- [ ] You will use a **fresh agent chat** here (not the no-HAWP chat, not the setup chat)

---

## Part 1 — Run the HAWP arm

### 1a. Fill the HAWP shape

Use the **same task** as `input`. Fill the other fields using [.hawp/kit/references/authoring-patterns.md](../../core/.hawp/kit/references/authoring-patterns.md):

```ts
const shape: Shape = {
  input: "...",      // your bare prompt, unchanged
  context: "...",    // minimal background
  mission: "...",    // objective, lens, decision target
  constraints: "...", // scope, evidence bar, anti-drift rules
  output: "...",     // what done looks like
};
```

If your run folder already has a filled shape (for example [runs/2026-06-15-layer1-coverage-balance/README.md](../runs/2026-06-15-layer1-coverage-balance/README.md)), reuse it.

### 1b. Run the shaped prompt

1. Start a **new agent chat** in this window.
2. Paste the filled shape and ask the agent to treat each field as binding.
3. Add: `Return the complete artifact as your final response.`
4. If the agent only summarizes instead of including the full artifact, follow up once: `Output the full artifact verbatim.`

### 1c. Save the HAWP output

```text
benchmark/runs/<date>-<task-name>/hawp/output.md
```

Include the filled shape at the top, then the captured output. See [runs/2026-06-12-v1-release-readiness/hawp/output.md](../runs/2026-06-12-v1-release-readiness/hawp/output.md) for an example.

---

## Part 2 — Score and record

### 2a. Score both outputs

Score both arms **blind**: label them Arm A / Arm B (randomized), score before revealing which is HAWP. For each arm, score all twelve dimensions on the 0–5 anchored scale (5 = clearly satisfied, 4 = minor gaps, 3 = notable gaps, 2 = mostly not, 1 = barely, 0 = not satisfied):

| # | Dimension | What to look for | Tends to favor |
| --- | --- | --- | --- |
| 1 | Drift resistance | Stayed in scope vs expanded into unrequested areas | neutral |
| 2 | Evidence vs inference separation | Observations separated from interpretation | HAWP |
| 3 | Output usefulness | Decision-useful artifact vs generic summary | neutral |
| 4 | Handoff quality | Could someone continue without re-discovery? | HAWP |
| 5 | Trustworthiness | Uncertainty acknowledged vs overstated confidence | HAWP |
| 6 | Scope adherence | Stayed within implied scope — scored the same for both arms | neutral |
| 7 | Completeness / coverage | Breadth of valid findings (a finding cap can cost points here) | no-HAWP |
| 8 | Conciseness / signal-to-noise | Tight and high-signal vs padded | neutral |
| 9 | Correctness / accuracy | Claims are true when checked against the source | neutral |
| 10 | False-positive control | Avoids flagging non-issues / over-claiming | neutral |
| 11 | Verifiability | Reader can re-check each claim from cited paths/lines/commands | neutral |
| 12 | Positive confirmation / balance | States what is correct, not only what is broken | no-HAWP |

Save scorecards:

```text
benchmark/runs/<date>-<task-name>/no-hawp/scorecard.md
benchmark/runs/<date>-<task-name>/hawp/scorecard.md
```

Headline score: total the twelve dimensions (raw max 60), then report **raw ÷ 4 → out of 15** and **raw ÷ 60 → percentage**. See [benchmark-prompt.md](../benchmark-prompt.md#headline-score-out-of-15).

### 2b. Write the comparison and README

Create or update these files in the run folder:

```text
benchmark/runs/<date>-<task-name>/
  README.md         # setup, score table, short interpretation, caveats
  comparison.md     # detailed side-by-side analysis + "which is better for now"
  no-hawp/
    output.md
    scorecard.md
  hawp/
    output.md
    scorecard.md
```

In `comparison.md`, include:

- How the two outputs differ in structure
- What each run found that the other missed
- Dimension-by-dimension notes
- **Which one is better, for now** — with the simple score percentage if you used it

Copy structure from [runs/2026-06-11-bounded-repo-review/comparison.md](../runs/2026-06-11-bounded-repo-review/comparison.md).

In `README.md` caveats, note:

- Same model and isolated sessions used?
- Was the no-HAWP arm run in a clean `/tmp/` copy or in the main repo (contamination)?
- Any delivery issues (e.g. agent summarized instead of outputting full artifact)?

---

## Part 3 — Clean up

### 3a. Remove temporary clean workspaces

The no-HAWP arm used a copy under `/tmp/`. Remove all benchmark copies:

```bash
./benchmark/prepare-clean-workspace.sh --cleanup
```

This deletes folders matching `/tmp/hawp-benchmark-clean-*`. Your real repository is not touched.

To remove one specific copy instead:

```bash
rm -rf /tmp/hawp-benchmark-clean-YYYYMMDD-HHMMSS
```

### 3b. Close extra editor windows

If the clean copy window is still open, close it now. Keep working in your normal repository window. Do not commit or push anything from the `/tmp/` copy — it is throwaway.

### 3c. Verify the real repo is unchanged

From the repository root:

```bash
test -f AGENTS.md && echo "AGENTS.md OK"
test -d .cursor/rules && echo ".cursor/rules OK"
```

Both should print OK. If either is missing, restore from git — you should never have deleted them from the real repo.

```bash
git checkout -- AGENTS.md .cursor/rules .continue/rules .github/instructions
```

(Only if something was accidentally removed.)

### 3d. Optional: commit your run

```bash
git add benchmark/runs/<date>-<task-name>/
git status
```

Commit only when you intend to — the benchmark folder does not require a commit after every run.

---

## Checklist

- [ ] HAWP arm: original repo + fresh chat + filled shape saved to `hawp/output.md`
- [ ] Both scorecards written with reasoning
- [ ] `comparison.md` and `README.md` written
- [ ] `/tmp/hawp-benchmark-clean-*` removed
- [ ] Real repo still has `AGENTS.md` and rules

You are done. For the next comparison, start again with **[setup.md](setup.md)**.
