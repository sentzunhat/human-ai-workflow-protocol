# Benchmark setup

Use this once before your first run, or whenever you start a new comparison session.

This benchmark compares two outputs from the same task: one without HAWP framing (no-HAWP arm) and one with a filled HAWP shape (HAWP arm). They run in **two separate Cursor windows** with **no shared chat history**.

| Phase | Where | Instruction file |
| --- | --- | --- |
| Setup + launch clean workspace | **Original repo** (this window) | This file |
| No-HAWP arm | **Clean copy** under `/tmp/` | [run.md](run.md) |
| HAWP arm + scoring + recording + cleanup | **Original repo** (return here) | [cleanup.md](cleanup.md) |

---

## What you need

- This repository cloned locally
- Cursor (or another agent-capable editor) with the same model available for both runs
- About 30–60 minutes for one full comparison (task choice, two runs, scoring, recording)
- Optional: Node 22 if you want to run librarian checks yourself (see repo `.nvmrc`)

---

## One-time check

From the repository root:

```bash
chmod +x benchmark/prepare-clean-workspace.sh
./benchmark/prepare-clean-workspace.sh --help
```

You should see the script help text. The script never modifies this repository — it only prepares a copy under `/tmp/`.

---

## Read before you run

Skim these (you do not need to memorize them):

| Doc | Purpose |
| --- | --- |
| [benchmark-prompt.md](../benchmark-prompt.md) | Task types, scoring dimensions, how to interpret results |
| [runs/](../runs/) | Completed examples you can copy as templates |
| [.hawp/kit/references/authoring-patterns.md](../../core/.hawp/kit/references/authoring-patterns.md) | How to fill the HAWP shape for your task type |

---

## Step 1 — Choose the task

Pick one task type and write a natural, somewhat vague prompt.

| Task type | Example prompt |
| --- | --- |
| Bounded repo review | "Take a look at the librarian/ folder and tell me what's wrong with it. What else should we clean up?" |
| Vague open-ended question | "Is HAWP ready for a v1.0 release? What do we still need to do?" |
| Standards / truth audit | "Check whether our install docs match what the scripts actually do." |
| Implementation planning | "We need to merge dev to main — what's the plan?" |
| Handoff generation | "I'm stopping for today — what should the next person know?" |

Write your chosen prompt down. You will use the **exact same words** for both arms.

---

## Step 2 — Create a run folder

Pick a date and short task name, then create the folder skeleton in **this** repository:

```bash
DATE=$(date +%Y-%m-%d)
NAME="my-task-name"   # e.g. bounded-repo-review, v1-release-readiness
mkdir -p "benchmark/runs/${DATE}-${NAME}/no-hawp"
mkdir -p "benchmark/runs/${DATE}-${NAME}/hawp"
```

You will fill these folders during [run.md](run.md) and [cleanup.md](cleanup.md).

---

## Step 3 — Launch the clean workspace (opens a new Cursor window)

From the repository root, replace the quoted text with your task prompt:

```bash
./benchmark/prepare-clean-workspace.sh --open --task "YOUR BARE TASK PROMPT HERE"
```

This will:

1. Copy the repo to `/tmp/hawp-benchmark-clean-<timestamp>/`
2. Remove agent-loaded guidance from the copy only (`AGENTS.md`, rules, instructions)
3. Write `BENCHMARK-TASK.txt` in the copy
4. Open the copy in a **new Cursor window** (if `cursor` CLI is available)

If the window did not open, open the printed path manually as its own workspace.

**Do not delete `AGENTS.md` or rules from the real repository.**

---

## Step 4 — Switch to the new window

In the **new** Cursor window (the `/tmp/` copy):

1. Open **[run.md](run.md)** and follow it through the STOP step.
2. Save `no-hawp/output.md` in the real repo.
3. Close the clean workspace window.

**Stop there.** Do not run the HAWP arm in the clean window.

---

## Step 5 — Return here for the HAWP arm

Back in **this** original repository window:

1. Open a **brand-new agent chat**.
2. Follow **[cleanup.md](cleanup.md)** for the HAWP arm, scoring, recording, and cleanup.

---

## Fairness rule (important)

In this repo, Cursor loads `AGENTS.md` and `.cursor/rules/` automatically. A bare prompt in the main workspace is **not** fully unguided.

For a fair no-HAWP arm, always use `prepare-clean-workspace.sh` as in Step 3. If you run the no-HAWP arm in the main repo anyway, note contamination in your run README (HAWP's advantage will look smaller than it really is).

---

## Setup checklist

- [ ] Script is executable and `--help` works
- [ ] You have read `benchmark-prompt.md` (at least the task types and scoring table)
- [ ] Task prompt chosen and run folder created under `benchmark/runs/`
- [ ] Clean workspace launched with `--open --task "..."`
- [ ] You will use a **fresh chat** in each window (no shared context)

**Next (clean workspace window):** [run.md](run.md)  
**After no-HAWP is saved (original repo):** [cleanup.md](cleanup.md)
