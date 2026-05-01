---
name: intake
description: Run the full work intake loop — analyze, plan, implement or hold based on risk
---

You are running the work intake loop for this repository.

Read `kit/usage/INTAKE_WORKFLOW.md` for the full operating loop before proceeding.
Read `work/BACKLOG.md` to find the next available backlog ID.

Adapt paths to the local layout: `.hawp/work/` in downstream projects, root `.work/` in the HAWP source repo.

## Your job

Given the bug or task the user described, execute steps 1–4 of the workflow:

### 1. Assign & log

- Pick the next backlog ID (`TASK-XXX`)
- Add a row to `work/BACKLOG.md` with status `analyzing`

### 2. Investigate

- Search the codebase for the affected area
- Identify root cause or most likely cause
- Note scope — what else could be affected
- Be explicit about what is directly verified vs inferred

### 3. Write the plan

- Use `kit/templates/intake-plan.md` as the template
- Save to `work/active/<ID>.md`
- Update the BACKLOG.md row to `plan-ready` with a link to the plan file

### 4. Apply the review gate

- **Low risk** (wrong path, missing asset, typo, isolated config fix): implement immediately, note it in the plan
- **Medium risk** (shared module, routing, layout change): present the plan and ask "OK to proceed?"
- **High risk** (auth, API, env vars, deploy config, security): hold — do not implement until user explicitly approves

## Output format

After analysis, respond with:

1. A one-paragraph plain-English summary of what you found
2. The risk level and gate decision
3. Either the fix applied (low risk) or the plan for review (medium/high)

Keep it compact and decision-useful.
