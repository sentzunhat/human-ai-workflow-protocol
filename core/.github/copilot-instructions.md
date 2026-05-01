This repository uses HAWP as a lightweight workflow method.

Follow the repo-local HAWP guidance in:

- .hawp/kit/START_HERE.md
- .hawp/kit/templates/status-report.md

Use .hawp/kit/START_HERE.md as the operating guide for how this repo applies HAWP in practice.

Use .hawp/kit/templates/status-report.md when the user asks for a:

- status report
- checkpoint summary
- context transfer summary
- second-brain review artifact

Saved status reports belong in:

- .hawp/work/active/

For bugs/tasks, track in .hawp/work/BACKLOG.md. Active plan files go in .hawp/work/active/. Deferred items can live in .hawp/work/parked/. Close by moving to .hawp/work/closed/YYYY/MM/DD/.

Keep the repo-local HAWP layer lean.

HAWP in this repository is a practical workflow layer, not a runtime engine, compiler, validator, orchestrator, or memory system.

Do not:

- invent per-field folders such as input/, context/, mission/, constraints/, output/, or checkpoint/
- imply a runtime engine, compiler, validator, orchestrator, persistence layer, or memory system
- overstate certainty; keep direct evidence separate from inference

Prefer compact, decision-useful outputs.
