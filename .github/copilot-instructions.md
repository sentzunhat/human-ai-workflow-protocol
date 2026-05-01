This repository uses HAWP as a lightweight workflow method.

Follow the repo-local HAWP guidance in:

- core/.hawp/kit/usage/INIT.md
- core/.hawp/kit/usage/STATUS_REPORT.md

Use core/.hawp/kit/usage/INIT.md as the operating guide for how this repo applies HAWP in practice.

Use core/.hawp/kit/usage/STATUS_REPORT.md when the user asks for a:

- status report
- checkpoint summary
- context transfer summary
- second-brain review artifact

Saved status reports belong in:

- .work/active/

For bug reports, tasks, and improvement work in this repo:

- Follow the intake loop in core/.hawp/kit/usage/INTAKE_WORKFLOW.md
- Track all open and completed work in .work/BACKLOG.md
- Write plan files using core/.hawp/kit/templates/intake-plan.md
- Active plan files go in .work/active/
- Deferred items can go in .work/parked/
- Close by moving to .work/closed/YYYY/MM/DD/
- ADRs and decisions go in .work/decisions/YYYY/MM/DD/

Keep the repo-local HAWP layer lean.

HAWP in this repository is a practical workflow layer, not a runtime engine, compiler, validator, orchestrator, or memory system.

Do not:

- invent per-field folders such as input/, context/, mission/, constraints/, output/, or checkpoint/
- imply a runtime engine, compiler, validator, orchestrator, persistence layer, or memory system
- overstate certainty; keep direct evidence separate from inference

Prefer compact, decision-useful outputs.
