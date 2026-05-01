# HAWP Benchmark Prompt

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

For authoring guidance on specific task types, see [.hawp/kit/authoring-patterns.md](../core/.hawp/kit/authoring-patterns.md).

### Step 4: Compare on scoring dimensions

Score both outputs on each dimension below.

---

## Scoring dimensions

Use **pass / mixed / fail** for each dimension, or a 1–5 scale if you prefer finer resolution.

| Dimension                            | What to look for                                                                                            |
| ------------------------------------ | ----------------------------------------------------------------------------------------------------------- |
| **Drift resistance**                 | Did the output stay within the stated scope, or did it expand into unrequested areas?                       |
| **Evidence vs inference separation** | Are direct observations and interpretations clearly separated, or collapsed into confident-sounding claims? |
| **Output usefulness**                | Is the output decision-useful — does it help a reader decide what to do next — or is it a generic summary?  |
| **Handoff quality**                  | Could this output be given to another human or agent to continue the work without significant re-discovery? |
| **Trustworthiness**                  | Does the output acknowledge uncertainty, or does it overstate confidence?                                   |
| **Constraint discipline**            | Were the stated constraints respected throughout, or did the output drift past them?                        |

### Scoring key

| Score     | Meaning                                                 |
| --------- | ------------------------------------------------------- |
| **pass**  | This dimension was clearly satisfied                    |
| **mixed** | Partially satisfied, with notable gaps or inconsistency |
| **fail**  | This dimension was not satisfied                        |

### 1–5 scale (optional)

If you want finer resolution:

- 5 — clearly satisfied, no meaningful issues
- 4 — satisfied with minor gaps
- 3 — partially satisfied
- 2 — mostly not satisfied
- 1 — not satisfied at all

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

If HAWP consistently improves drift resistance and evidence discipline but not output usefulness, that is a meaningful finding — and a direction for improving authoring guidance.

---

## Usage note

This benchmark is a human evaluation aid.

It is not a proof, a scientific claim, or a formal evaluation framework. It is a structured way to compare outputs and develop a practical sense of whether HAWP helps.

Run it with real tasks from your own work context. The more realistic the task, the more meaningful the comparison.

---

## Related docs

- [.hawp/kit/authoring-patterns.md](../core/.hawp/kit/authoring-patterns.md) — how to fill the shape for specific task types
- [.hawp/kit/examples/](../core/.hawp/kit/examples/) — concrete filled-shape examples
- [.hawp/kit/spec.md](../core/.hawp/kit/spec.md) — field semantics and evidence discipline rules
