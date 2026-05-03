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

1. **Check the backlog first.** Read `BACKLOG.md` in `.hawp/work/` before starting any work.
   Assign the next sequential ID. Never start on an ID already marked `in-progress`.

2. **Write a plan file** for all non-trivial work before implementing.
   Use `kit/templates/intake-plan.md`. Save to `work/active/<id>.md`.

3. **Respect the review gate** defined in `kit/usage/intake-workflow.md`.
   Do not implement medium or high risk changes without explicit user approval.

4. **Keep scope tight.** Fix the reported issue only. Do not refactor adjacent code unless it is the direct root cause.

5. **Update BACKLOG.md** at every status transition: `inbox → analyzing → plan-ready → in-progress → done`.

6. **Separate evidence from inference** in every plan. Never present an inference as a confirmed fact.

7. **For this repo, write project work to `.hawp/work/`.** Do not write repo work into `core/.hawp/work/`.

8. **Treat `core/.hawp/work/` as downstream scaffold source only.**

9. **Treat `core/.hawp/kit/` as reusable package source only.**

10. **Keep `BACKLOG.md` compact.** Do not append completed work endlessly; move closed work to `.hawp/work/closed/YYYY/MM/DD/`.

11. **When Done rows grow large, create a compaction item first.** Add a work item titled `Compact BACKLOG.md and archive closed work.` before adding more Done rows.

12. **Apply backlog alignment guardrails.** Follow `.github/instructions/hawp-backlog-alignment.instructions.md`; use `.github/prompts/hawp-backlog-alignment.prompt.md` for explicit review/compaction requests.

## Source Resolution

For full loop detail, read in order:

1. `kit/usage/intake-workflow.md`
2. `work/BACKLOG.md`
3. `kit/templates/intake-plan.md`

Adapt paths to match the local layout: `.hawp/work/` in downstream projects, repo-root `.hawp/work/` in the HAWP source repo.
