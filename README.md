# human-ai-workflow-protocol (HAWP)

> A minimal protocol for reliable human–AI collaboration. Shape the work **before** execution begins — and stop drifting.

## Get started

| You want to…                              | Open this                                    |
| ----------------------------------------- | -------------------------------------------- |
| **Install HAWP** into a repo (copy/paste) | → [core/install.md](core/install.md)         |
| **Update HAWP** to latest `main`          | → [core/update.md](core/update.md)           |
| Shape your first task                     | → `.hawp/kit/START_HERE.md` (after install)  |
| Run the HAWP vs no-HAWP comparison        | → [benchmark/README.md](benchmark/README.md) |

Both install and update are single copy/paste shell blocks. Safe to re-run. They never overwrite your `.hawp/work/` content or your `.github/copilot-instructions.md`.

---

## What HAWP is

A compact task-shaping protocol — five required fields, one optional checkpoint:

```ts
type Shape = {
  input: string; // the request as received
  context: string; // minimal background needed to interpret it
  mission: string; // the concrete objective for this run
  checkpoint?: string; // optional lightweight pause marker
  constraints: string; // hard boundaries and quality bars
  output: string; // what done looks like
};
```

That's the whole core. Templates, patterns, examples, and the GitHub Copilot overlay are optional usage aids around it.

HAWP is **not** a runtime engine, memory system, validator, or orchestration framework. It is a lean shaping protocol intended to stay portable across projects. Public examples in this repo are kept fictional or clearly generic on purpose.

## Why use it

- **Less drift.** A filled shape locks intent before the first tool call.
- **Cheap handoffs.** The same five fields work as a context-transfer artifact between humans, between agents, and across sessions.
- **No lock-in.** Plain Markdown + one TypeScript type. No runtime, no service, no database.
- **Optional everything.** Use just `START_HERE.md`, or layer in the intake loop, status reports, ADRs, and Copilot overlay as you need them.

## Using HAWP after install

1. Read `.hawp/kit/README.md` to understand the protocol.
2. Use `.hawp/kit/START_HERE.md` to shape a task.
3. Use `.hawp/kit/AUTHORING_PATTERNS.md` for recurring task-shape recipes.
4. For bugs/tasks, follow `.hawp/kit/usage/INTAKE_WORKFLOW.md` and track in `.hawp/work/BACKLOG.md`.
5. For checkpoints / context handoffs, use `.hawp/kit/templates/status-report.md`.
6. See `.hawp/kit/examples/` for concrete filled shapes.

Publication-safety guidance lives under `.hawp/kit/reviews/`.

## Repo layout (source)

- `core/install.md` — install script + boundary notes
- `core/update.md` — update script + migration notes
- `core/.hawp/kit/` — the reusable kit (docs, templates, patterns, reviews, examples, types, usage guides)
- `core/.hawp/work/` — clean install scaffold seeded into downstream `.hawp/work/` (READMEs + starter `BACKLOG.md`)
- `.work/` — this repo's own operating state: backlog, ADRs, status plans, evidence (not installed downstream)
- `.github/` — Copilot instructions + prompt pack (intake, status, handoff, docs alignment)
- `benchmark/` — optional HAWP vs no-HAWP comparison harness

Downstream installs flatten `core/.hawp/` to repo-root `.hawp/`. The install allowlist seeds a clean `.hawp/work/` from `core/.hawp/work/` (READMEs + starter `BACKLOG.md`); this repo's own operating state lives at root `.work/` and never ships.

## License

Apache 2.0 — see [LICENSE](./LICENSE). The installed kit ships its own `.hawp/LICENSE` copy.
