This repository uses HAWP as a lightweight workflow method.

Follow the repo-local HAWP guidance in:

- `.hawp/kit/start-here.md`
- `.hawp/kit/usage/status-report.md`

Use `.hawp/kit/start-here.md` as the operating guide for how this repo applies HAWP in practice.

Use `.hawp/kit/usage/status-report.md` when the user asks for a status report, checkpoint summary, context transfer summary, or second-brain review artifact.

Saved status reports belong in `.hawp/work/status/YYYY/MM/DD/{uuid}/status.md` — use `hawp work status --title "<label>"` to generate the path. Similarly: `hawp work evidence|decision|note --title "<label>"` for other secondary document types.

For bugs/tasks, track in `.hawp/work/BACKLOG.md`. Active plan files go in `.hawp/work/active/`. Close by moving to `.hawp/work/closed/YYYY/MM/DD/`.

Cursor rules under `.cursor/rules/hawp-*.mdc` provide scoped HAWP behavior. See `.hawp/kit/references/backlog-alignment.md` for backlog compaction policy.

For changes that affect installed files or for work in a downstream consumer repo, consult `.github/skills/multi-repo-context.md` for the source → downstream → consumer propagation model.

HAWP in this repository is a practical workflow layer, not a runtime engine, compiler, validator, orchestrator, or memory system.

Do not invent per-field folders such as input/, context/, mission/, constraints/, output/, or checkpoint/.

Prefer compact, decision-useful outputs. Separate direct evidence from inference.
