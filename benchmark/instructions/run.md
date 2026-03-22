# Run the no-HAWP arm (clean workspace)

This is the **no-HAWP arm only**. Run it in the clean Cursor window that `setup.md` opened — the one whose path looks like `/tmp/hawp-benchmark-clean-*`.

The HAWP arm, scoring, and recording happen later, back in the original repo, in **[cleanup.md](cleanup.md)**.

---

## Confirm you are in the right place

Before starting, check all of these:

- [ ] Workspace path is under `/tmp/hawp-benchmark-clean-*` (not your normal clone path)
- [ ] `AGENTS.md` is **missing** from the workspace root (removed in the copy on purpose)
- [ ] `.cursor/rules/` is **missing** (same reason)
- [ ] You see `BENCHMARK-TASK.txt` in the workspace root (if setup used `--task`)
- [ ] This is a **brand-new agent chat** — not a continuation from the original repo window

If any check fails, stop. Do not run the no-HAWP arm in the original repository.

---

## Step 1 — Run the bare prompt

1. Open `BENCHMARK-TASK.txt` and copy the task prompt.
2. Paste it into this chat. Do not add HAWP fields or extra structure.
3. If the file is missing, paste the exact same words you wrote down during setup.
4. Wait for the full output. Do not ask for summaries — you want the complete answer.

Optional nudge if the agent stops early:

```text
Return your complete answer as your final response.
```

---

## Step 2 — Save the output in the real repo

Switch to your **original repository window** (or use a terminal pointed at the real clone — not this `/tmp/` copy).

Save the captured agent output here:

```text
benchmark/runs/<date>-<task-name>/no-hawp/output.md
```

At the top of `output.md`, include:

- Which arm this is (no-HAWP)
- The exact prompt used
- Model name (if you know it)
- That the run used the clean `/tmp/` workspace

See [runs/2026-06-11-bounded-repo-review/no-hawp/output.md](../runs/2026-06-11-bounded-repo-review/no-hawp/output.md) for an example.

---

## Step 3 — STOP and hand back to the original repo

**Do not start the HAWP arm in this window.**

1. Confirm `no-hawp/output.md` is saved in the **real** repo.
2. **Close this Cursor window** (the `/tmp/` clean workspace).
3. Return to the **original repository** Cursor window.
4. Continue with **[cleanup.md](cleanup.md)** in a **fresh chat** there — it runs the HAWP arm, scores both, records the run, then removes the temp copy.

The two arms must never share a workspace or chat history. Skipping this step is a contamination caveat to record in the run README.

---

## No-HAWP checklist

- [ ] Ran in clean `/tmp/` workspace with no `AGENTS.md` / rules
- [ ] Fresh chat; bare prompt only (no HAWP shape)
- [ ] Full output saved to `benchmark/runs/<date>-<task-name>/no-hawp/output.md` in the real repo
- [ ] Clean workspace window closed before continuing

**Next (original repo, fresh chat):** [cleanup.md](cleanup.md)
