# human-ai-workflow-protocol (HAWP)

> A minimal protocol for reliable human–AI collaboration. Shape the work **before** execution begins — and stop drifting.

**Less drift, cheaper handoffs, zero lock-in.** Portable task-shaping for humans and agents.

[![Validate Distribution Generated](https://github.com/sentzunhat/human-ai-workflow-protocol/actions/workflows/sync-distribution-generated.yml/badge.svg)](https://github.com/sentzunhat/human-ai-workflow-protocol/actions/workflows/sync-distribution-generated.yml)

## Get Started

**New to HAWP?**

→ [**Install HAWP**](distribution/generated/install-main.md) | [Update](distribution/generated/update-main.md)

Then open `.hawp/kit/start-here.md` and shape your first task.

**Already installed?**

1. Edit `.hawp/kit/start-here.md` to shape a task
2. Track work in `.hawp/work/BACKLOG.md`
3. Use templates under `.hawp/kit/templates/` and patterns under `.hawp/kit/standards/`

**Other branches?** [Dev install](distribution/generated/install-dev.md) | [Dev update](distribution/generated/update-dev.md)

**Compare with other methods?** → [Benchmark](benchmark/README.md)

---

## Core Concept

HAWP is a **task-shaping protocol** — five fields that lock intent before work begins:

```ts
type Shape = {
  input: string; // the request as received
  context: string; // minimal background
  mission: string; // concrete objective
  checkpoint?: string; // optional lightweight pause
  constraints: string; // hard boundaries and quality bars
  output: string; // what done looks like
};
```

That's the full core. Templates, patterns, examples, and tooling are optional usage aids.

**HAWP is not** a runtime, memory system, validator, or orchestrator. It stays portable across projects.

## Why It Works

- **Less drift** — Shape locks intent before first tool call
- **Cheap handoffs** — Five fields work as context transfer between humans and agents
- **No lock-in** — Plain Markdown + one TypeScript type. No runtime or database
- **Optional everything** — Use just `start-here.md`, or layer in intake loop, status reports, and ADRs as needed

---

## Next Steps

After install:

1. **Shape work:** Edit `.hawp/kit/start-here.md` to fill the five fields
2. **Track it:** Add rows to `.hawp/work/BACKLOG.md`
3. **Use patterns:** Reuse templates from `.hawp/kit/templates/` for recurring shapes
4. **Handle handoffs:** Use `.hawp/kit/templates/status-report.md` to transfer context
5. **Review safety:** Check `.hawp/kit/standards/` for domain-specific guidance

See `.hawp/kit/examples/` for concrete filled shapes.

---

## More Details

## Repository Layout (Source)

**Install and update:**

- `distribution/generated/` — branch-based install/update guides (auto-generated, do not edit)
- `distribution/sources/` — authoritative sources for install/update scripts

**Reusable kit (installed to downstream `.hawp/kit/`):**

- `core/.hawp/kit/` — protocol docs, templates, patterns, standards, examples, usage guides

**This repo's working state (NOT installed downstream):**

- `.hawp/work/` — backlog, active tasks, closed work, decisions, evidence

**Tooling and reference:**

- `.github/` — Copilot instructions and prompt library
- `librarian/` — distribution validation and backlog tooling
- `benchmark/` — optional HAWP vs other methods comparison

---

## Contributing

**Distribution maintenance:**

Authoritative scripts live in `distribution/sources/install/script.md` and `distribution/sources/update/script.md`. Do not edit files under `distribution/generated/` directly.

After editing sources, regenerate and validate:

```bash
cd librarian
npm run distribution:sync
```

Both install and update are single copy/paste blocks. Safe to re-run. They never overwrite `.hawp/work/` or `.github/copilot-instructions.md`.

**Validation:**

GitHub Actions runs on `main`/`dev` pushes and fails if generated distribution guides drift from sources.

---

## Roadmap

Active development tracked in [`.hawp/work/BACKLOG.md`](.hawp/work/BACKLOG.md). Recent focus: backlog upgrade command, distribution clarity, UUID migration for parallel-safe work items.

---

## License

Apache 2.0 — see [LICENSE](./LICENSE). Installed kit ships its own `.hawp/LICENSE` copy.
