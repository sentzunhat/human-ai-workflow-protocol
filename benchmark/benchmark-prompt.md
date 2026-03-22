# HAWP Benchmark Prompt

**Step-by-step workflow:** [instructions/setup.md](instructions/setup.md) → [instructions/run.md](instructions/run.md) (no-HAWP arm, clean window; stop) → [instructions/cleanup.md](instructions/cleanup.md) (HAWP arm, score, record, cleanup; original repo).

## What this is

A reusable guide for running a simple HAWP vs no-HAWP comparison across practical task types.

This is a human evaluation aid. It is not a scientific instrument or an automated scoring system.

Use it when you want to evaluate whether HAWP-guided outputs are more bounded, more trustworthy, and more decision-useful than outputs from the same task run without HAWP.

---

## Core benchmark question

> Does framing a task with HAWP fields produce a more bounded, evidence-grounded, and decision-useful output than asking the same question without framing?

---

## Recommended task types

Pick one or more of these for each comparison run. Choose tasks that are realistic and sufficiently open-ended that drift is possible.

| Task type                 | Why it tests HAWP                                                 |
| ------------------------- | ----------------------------------------------------------------- |
| Bounded repo review       | Tests whether HAWP limits scope and prevents speculative findings |
| Standards / truth audit   | Tests evidence discipline and confidence labeling                 |
| Implementation planning   | Tests whether constraints prevent scope overreach                 |
| Handoff generation        | Tests whether HAWP produces portable, reusable context artifacts  |
| Vague open-ended question | Tests whether framing converts ambiguity into a bounded response  |

Good benchmark tasks include vague, messy, or overloaded requests. These are where unguided work drifts most.

---

## How to run a comparison

### Step 1: Choose a task

Pick a realistic task from the types above. Write it as a natural, somewhat vague request — the kind that would typically produce a drifty or overconfident response.

Example: "review this codebase and tell me what's wrong"

### Step 2: Run it without HAWP

Give the task to a human or AI agent without using HAWP. Do not add a structured prompt. Capture the output as-is.

Full steps: [instructions/run.md](instructions/run.md). Setup and launch: [instructions/setup.md](instructions/setup.md).

Keep the comparison fair:

- Use the same person or the same AI model for both runs.
- Run each in a fresh, separate session so neither run can see the other's work.
- If the repository already has HAWP rules installed (for example in `AGENTS.md` or editor rules), the "without HAWP" run is not fully unguided. Editors like Cursor load `AGENTS.md` and `.cursor/rules/` automatically for every agent in the workspace, so this cannot be switched off from inside a chat.

To get a truly clean unguided run, do **not** delete `AGENTS.md` or the rules from the repository itself — they are part of what the repo ships. Instead, use the helper script (recommended):

```bash
./benchmark/prepare-clean-workspace.sh --open --task "Your bare task prompt here"
```

That copies the repo to `/tmp/`, removes agent-loaded guidance **from the copy only**, and opens the copy in a new Cursor window. Run the bare prompt there in a fresh chat ([run.md](instructions/run.md)). Close that window, then run the HAWP arm in the original repository ([cleanup.md](instructions/cleanup.md)).

Manual equivalent:

1. Copy the repository to a temporary folder (for example `/tmp/benchmark-clean/`).
2. In the copy only, delete the agent-loaded guidance: `AGENTS.md`, `.cursor/rules/`, `.continue/rules/`, and `.github/instructions/` (plus the same paths under `core/providers/` if present).
3. Open the copy as its own workspace (a new editor window) and run the bare prompt there in a fresh chat.
4. Run the HAWP arm in the original repository as usual.

When finished, remove old copies: `./benchmark/prepare-clean-workspace.sh --cleanup`

If that is more effort than the run justifies, run both arms in the original repository anyway and note the contamination in the run's caveats — it usually makes the gap look smaller than it really is, so the HAWP result is understated, not inflated.

### Step 3: Run it with HAWP

Fill the HAWP shape for the same task:

```ts
const shape: Shape = {
  input: "...", // the original request, preserved
  context: "...", // minimal background needed
  mission: "...", // concrete objective with lens and decision target
  constraints: "...", // scope limits, evidence bar, anti-drift rules
  output: "...", // what done looks like as a decision-useful artifact
};
```

Use the filled shape as the prompt. Capture the output.

Full steps: [instructions/cleanup.md](instructions/cleanup.md) (original repo only, after no-HAWP is saved).

For authoring guidance on specific task types, see [.hawp/kit/references/authoring-patterns.md](../core/.hawp/kit/references/authoring-patterns.md).

### Step 4: Compare on scoring dimensions

Score both outputs on each dimension below.

### Step 5: Record the run

Save the run so others can read it later, using this folder layout:

```text
benchmark/runs/<date>-<task-name>/
  README.md            # task, setup, side-by-side scores, interpretation, caveats
  comparison.md        # detailed comparison: structure, findings overlap, per-dimension analysis
  no-hawp/
    output.md          # raw output from the unguided run
    scorecard.md       # scores and reasoning for the unguided run
  hawp/
    output.md          # the filled shape, then the raw output
    scorecard.md       # scores and reasoning for the guided run
```

See [runs/](runs/) for completed examples. After the run: [instructions/cleanup.md](instructions/cleanup.md).

---

## Scoring dimensions

Score both outputs on the **twelve** dimensions below. The set is deliberately balanced: some dimensions tend to favor disciplined (HAWP) output, and two — **completeness / coverage** and **positive confirmation / balance** — tend to favor the wider net an unguided run casts. No dimension is structurally winnable by only one arm. Note that the "report-discipline" dimensions (2, 4, 5, 10) are correlated; they are kept distinct deliberately, but read them as a cluster, not four independent signals.

| #   | Dimension                            | What to look for                                                                                            | Tends to favor |
| --- | ------------------------------------ | ----------------------------------------------------------------------------------------------------------- | -------------- |
| 1   | **Drift resistance**                 | Did the output stay within the scope the task implies, or expand into unrequested areas?                    | neutral        |
| 2   | **Evidence vs inference separation** | Are direct observations and interpretations clearly separated, or collapsed into confident-sounding claims? | HAWP           |
| 3   | **Output usefulness**                | Is the output decision-useful — does it help a reader decide what to do next — or is it a generic summary?  | neutral        |
| 4   | **Handoff quality**                  | Could this output be given to another human or agent to continue the work without significant re-discovery? | HAWP           |
| 5   | **Trustworthiness**                  | Does the output acknowledge uncertainty, or does it overstate confidence?                                   | HAWP           |
| 6   | **Scope adherence**                  | Did the output stay within the scope the task implies? Score both arms the same way — an unguided run that respects implicit scope earns full marks; explicit constraints are not required to score here. | neutral        |
| 7   | **Completeness / coverage**          | Did it surface the real issues a reader needs? Breadth of valid findings. A finding cap can cost points here. | no-HAWP        |
| 8   | **Conciseness / signal-to-noise**    | Is it tight and high-signal, or padded with filler and restated points?                                     | neutral        |
| 9   | **Correctness / accuracy**           | Are the factual claims actually true when checked against the source?                                       | neutral        |
| 10  | **False-positive control**           | Does it avoid flagging non-issues or over-claiming? Are dismissed concerns handled honestly?                | neutral        |
| 11  | **Verifiability**                    | Can a reader independently re-check each claim from the cited paths, lines, or commands?                    | neutral        |
| 12  | **Positive confirmation / balance**  | Does it state what is correct, not only what is broken — so a reader sees the whole picture, not just problems? | no-HAWP        |

### Scoring key (0–5 anchored, per dimension)

Score each dimension 0–5. Keep to the anchors below — resist inventing in-between half-points, since fine-grained latitude is exactly where unconscious bias enters a qualitative score.

| Score | Meaning                                       |
| ----- | --------------------------------------------- |
| **5** | Clearly satisfied, no meaningful issues       |
| **4** | Satisfied, minor gaps                         |
| **3** | Partially satisfied, notable gaps             |
| **2** | Mostly not satisfied                          |
| **1** | Barely addressed                              |
| **0** | Not satisfied                                 |

### Headline score (out of 15)

Total the twelve dimensions (raw max = 60), then report both:

- **Headline = raw ÷ 4 → score out of 15** (e.g. raw 51 → 12.75 / 15)
- **Percentage = raw ÷ 60** (keeps runs comparable even if the dimension set changes later)

Treat both as a summary of qualitative judgments, not a measurement. The per-dimension reasoning and the direction of the gap matter more than the number.

### Blind scoring (do this)

To keep scores honest, especially when the same person authored the HAWP shape:

1. Strip the arm labels — call them **Arm A** and **Arm B** (randomize which is which).
2. Score all twelve dimensions for both arms before revealing which is HAWP.
3. Reveal the mapping, then write the reasoning.

Record in the run README which arm was A and which was B.

> **Rubric version note:** The three runs in `runs/` were originally scored on a legacy 6-dimension / 12-point rubric (pass/mixed/fail). All three were re-scored on this 12-dimension / 0–5 rubric on 2026-06-15; their READMEs carry both the current score and a note recording the superseded legacy reading. Use the 0–5 percentages when comparing gaps across runs.

---

## Failure cases

Good benchmark tasks should include failure-prone scenarios. These are the conditions under which guidance matters most.

Include at least one of the following in your benchmark set:

- a vague or underspecified task with no clear output definition
- a task that mixes multiple competing goals in one request
- a task with an obvious "scope creep" trap (e.g., "what else should we fix?")
- a task where the correct answer is "we don't have enough evidence to decide"

If HAWP performs well on clean tasks but not on messy ones, it is not adding durable value.

---

## Interpreting results

There is no passing threshold. The goal is comparative insight.

Ask:

- On which dimensions did HAWP-guided output consistently score better?
- On which dimensions did HAWP add little or no value?
- Were the gains from HAWP worth the additional framing effort?
- Were there task types where HAWP made no difference?
- Did the HAWP framing cause anything to be missed? Caps and scope limits trade breadth for discipline. Check whether the unguided run found valid items the guided run excluded — if it did, that points at the constraints field, not at the protocol.

If HAWP consistently improves drift resistance and evidence discipline but not output usefulness, that is a meaningful finding — and a direction for improving authoring guidance.

---

## Usage note

This benchmark is a human evaluation aid.

It is not a proof, a scientific claim, or a formal evaluation framework. It is a structured way to compare outputs and develop a practical sense of whether HAWP helps.

Run it with real tasks from your own work context. The more realistic the task, the more meaningful the comparison.

---

## Related docs

- [.hawp/kit/references/authoring-patterns.md](../core/.hawp/kit/references/authoring-patterns.md) — how to fill the shape for specific task types
- [.hawp/kit/examples/](../core/.hawp/kit/examples/) — concrete filled-shape examples
- [.hawp/kit/references/spec.md](../core/.hawp/kit/references/spec.md) — field semantics and evidence discipline rules
