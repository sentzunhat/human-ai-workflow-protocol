---
applyTo: "**"
---

# Work Intake — Ambient Rules

These rules apply whenever the user reports a bug, describes a broken behavior, or requests a fix or task.

## Trigger Recognition

Treat any of these as a work intake:

- "I have a bug", "there's a bug", "something is broken"
- "X is not working / not loading / not showing"
- "can you fix / look at / check X"
- "X looks wrong", "X is missing", "X shows an error"
- "add X", "update X", "change X"

## Operating Rules

1. **Check the backlog first.** Read `BACKLOG.md` (in the active `.hawp/usage/` or `core/.hawp/usage/` path) before starting any work.
   Assign the next sequential ID. Never start on an ID already marked `in-progress`.

2. **Write a plan file** for all non-trivial work before implementing.
   Use `templates/intake-plan.md`. Save to `usage/status/<date>-<id>-<slug>.md`.

3. **Respect the review gate** defined in `usage/INTAKE_WORKFLOW.md`.
   Do not implement medium or high risk changes without explicit user approval.

4. **Keep scope tight.** Fix the reported issue only. Do not refactor adjacent code unless it is the direct root cause.

5. **Update BACKLOG.md** at every status transition: `inbox → analyzing → plan-ready → in-progress → done`.

6. **Separate evidence from inference** in every plan. Never present an inference as a confirmed fact.

## Source Resolution

For full loop detail, read in order:

1. `usage/INTAKE_WORKFLOW.md`
2. `usage/BACKLOG.md`
3. `templates/intake-plan.md`

Adapt paths to match the local install: `core/.hawp/` in the HAWP source repo, `.hawp/` in downstream projects.
