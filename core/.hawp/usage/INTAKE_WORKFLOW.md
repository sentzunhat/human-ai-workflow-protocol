# Work Intake Workflow

This is the operating loop for bugs, tasks, and improvements in a HAWP-enabled repo.
It is built on top of the HAWP task-shaping protocol — keep it lean.

---

## The Loop

```
1. INTAKE       You drop a bug or task (natural language is fine)
2. ANALYZE      I investigate, shape it as a HAWP task, and fill the intake template
3. PLAN         I write a plan file with root cause, options, and recommended fix
4. REVIEW GATE  Auto-implement if low-risk; hold for your approval if risky
5. IMPLEMENT    Execute the fix
6. VERIFY       Confirm it works / document what is still unproven
7. CLOSE        Move backlog row to Done; save status report if context matters
```

---

## Step 1 — Intake

You can drop work in natural language, for example:

> "this feature is broken in production"
> "add X to the Y screen"
> "the output from this command is wrong"

I will:

- Assign a backlog ID and add a row to [BACKLOG.md](BACKLOG.md) with status `inbox`
- Move to analysis immediately unless you say "just log it"

---

## Step 2 — Analyze

I investigate and produce a filled-in intake form (see `../templates/intake-plan.md`).
The backlog row moves to `analyzing`.

Analysis covers:

- Where the problem lives (files, lines)
- Root cause (or most likely cause if not directly provable)
- Scope — is anything else affected?
- Confidence level — what is proven vs inferred

---

## Step 3 — Plan

I write a plan file to `usage/status/<date>-<id>-<slug>.md`.
The backlog row moves to `plan-ready` and links to the plan file.

The plan includes:

- Root cause summary
- Options considered (at least 2 when non-trivial)
- Recommended fix with rationale
- Risk level: **low / medium / high**
- What the fix touches
- What to verify after

---

## Step 4 — Review Gate

| Risk   | Default                                                               |
| ------ | --------------------------------------------------------------------- |
| Low    | I implement immediately after writing the plan, note it in the status |
| Medium | I present the plan and ask: "OK to proceed?"                          |
| High   | I hold — you must explicitly approve                                  |

**Low risk examples:** typo, wrong path, missing config value, isolated utility fix
**Medium risk examples:** changes to shared code, routing or auth-adjacent change
**High risk examples:** schema migration, deploy config, security-related change

You can override: say "just do it" or "hold for review" at any point.

---

## Step 5 — Implement

I make the changes. The backlog row moves to `in-progress`.

I use the smallest correct fix — no refactoring beyond scope unless it is the root cause.

Before editing files, confirm the plan's overlap check shows `Can implement now: yes` or explicit user approval for overlap.

---

## Step 6 — Verify

I confirm the fix is structurally sound (imports resolve, build passes if checkable).
I note what is directly verified vs what requires a live environment to confirm.

---

## Step 7 — Close

The backlog row moves to `done`.
A status report is saved if:

- The item was non-trivial
- Context should survive into the next session
- You ask for one

---

## Parallel Safety

- The backlog is the single source of truth — check it before starting a new item
- Each work item has one plan file — never two agents working the same ID
- If you start something outside this loop, add a row so it is visible

---

## Quick Reference — Risk Levels

| Scenario                           | Risk   |
| ---------------------------------- | ------ |
| Typo, copy, or label fix           | Low    |
| Isolated config or path correction | Low    |
| Standalone utility or helper fix   | Low    |
| Shared module or component change  | Medium |
| New route, page, or API endpoint   | Medium |
| Auth, env var, or secrets change   | High   |
| Deploy config or infrastructure    | High   |
| Database or schema migration       | High   |
