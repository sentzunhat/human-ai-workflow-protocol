This repository uses HAWP as a lightweight workflow method.

Follow the repo-local HAWP guidance in:

- human-ai-workflow-protocol/hawp/usage/INIT.md
- human-ai-workflow-protocol/hawp/usage/STATUS_REPORT.md

Use human-ai-workflow-protocol/hawp/usage/INIT.md as the operating guide for how this repo applies HAWP in practice.

Use human-ai-workflow-protocol/hawp/usage/STATUS_REPORT.md when the user asks for a:

- status report
- checkpoint summary
- context transfer summary
- second-brain review artifact

Saved status reports belong in:

- human-ai-workflow-protocol/hawp/usage/status/

Keep the repo-local HAWP layer lean.

HAWP in this repository is a practical workflow layer, not a runtime engine, compiler, validator, orchestrator, or memory system.

Do not:

- invent per-field folders such as input/, context/, mission/, constraints/, output/, or checkpoint/
- imply a runtime engine, compiler, validator, orchestrator, persistence layer, or memory system
- overstate certainty; keep direct evidence separate from inference

Prefer compact, decision-useful outputs.
