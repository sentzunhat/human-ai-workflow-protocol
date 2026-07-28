# human-ai-workflow-protocol (HAWP)

> A minimal protocol for reliable human–AI collaboration. Shape the work **before** execution begins — and stop drifting.

**Less drift, cheaper handoffs, zero lock-in.** Portable task-shaping for humans and agents.

[![Validate Distribution Generated](https://github.com/sentzunhat/human-ai-workflow-protocol/actions/workflows/sync-distribution-generated.yml/badge.svg)](https://github.com/sentzunhat/human-ai-workflow-protocol/actions/workflows/sync-distribution-generated.yml)

## Get Started

Pick your agent provider. Each guide installs **HAWP kit** (`.hawp/kit/**`) plus that provider's overlay.

**Recommended:** open the linked guide and run the visible **Install Command (Copy/Paste)** block after review. Optional guide-fetch helpers write a script to `/tmp` for inspection — they do not auto-execute remote content.

### GitHub / Copilot

| | Stable (`main`) | Dev (`dev`) |
|--|-----------------|-------------|
| **Install** | [github/install/main.md](distribution/generated/github/install/main.md) | [github/install/dev.md](distribution/generated/github/install/dev.md) |
| **Update** | [github/update/main.md](distribution/generated/github/update/main.md) | [github/update/dev.md](distribution/generated/github/update/dev.md) |

Installs `core/providers/.github/` → `.github/instructions/`, `.github/prompts/`, `.github/copilot-instructions.md`

### Claude Code

| | Stable (`main`) | Dev (`dev`) |
|--|-----------------|-------------|
| **Install** | [claude/install/main.md](distribution/generated/claude/install/main.md) | [claude/install/dev.md](distribution/generated/claude/install/dev.md) |
| **Update** | [claude/update/main.md](distribution/generated/claude/update/main.md) | [claude/update/dev.md](distribution/generated/claude/update/dev.md) |

Installs `core/providers/.claude/` → `.claude/rules/hawp-*.md`, `CLAUDE.md`

### Codex

| | Stable (`main`) | Dev (`dev`) |
|--|-----------------|-------------|
| **Install** | [codex/install/main.md](distribution/generated/codex/install/main.md) | [codex/install/dev.md](distribution/generated/codex/install/dev.md) |
| **Update** | [codex/update/main.md](distribution/generated/codex/update/main.md) | [codex/update/dev.md](distribution/generated/codex/update/dev.md) |

Installs `core/providers/.codex/` → `AGENTS.md`

### Cursor

| | Stable (`main`) | Dev (`dev`) |
|--|-----------------|-------------|
| **Install** | [cursor/install/main.md](distribution/generated/cursor/install/main.md) | [cursor/install/dev.md](distribution/generated/cursor/install/dev.md) |
| **Update** | [cursor/update/main.md](distribution/generated/cursor/update/main.md) | [cursor/update/dev.md](distribution/generated/cursor/update/dev.md) |

Installs `core/providers/.cursor/` → `.cursor/rules/hawp-*.mdc`, `AGENTS.md`

### Continue

| | Stable (`main`) | Dev (`dev`) |
|--|-----------------|-------------|
| **Install** | [continue/install/main.md](distribution/generated/continue/install/main.md) | [continue/install/dev.md](distribution/generated/continue/install/dev.md) |
| **Update** | [continue/update/main.md](distribution/generated/continue/update/main.md) | [continue/update/dev.md](distribution/generated/continue/update/dev.md) |

Installs `core/providers/.continue/` → `.continue/rules/hawp-*.md`

**Index:** [distribution/generated/README.md](distribution/generated/README.md)

Then open `.hawp/kit/start-here.md` and shape your first task.

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

## Repository Layout (Source)

**Install and update:**

- `distribution/generated/` — per-provider install/update guides by branch (`github/install/main.md`, etc.; auto-generated)
- `distribution/sources/` — authoritative sources for install/update scripts and provider fragments

**Reusable kit (installed to downstream `.hawp/kit/`):**

- `core/.hawp/kit/` — protocol docs, templates, patterns, standards, examples, usage guides

**Agent provider packs (installed downstream per guide):**

- `core/providers/.claude/` → `.claude/rules/`, `CLAUDE.md`
- `core/providers/.codex/` → `AGENTS.md`
- `core/providers/.github/` → `.github/`
- `core/providers/.cursor/` → `.cursor/rules/`, `AGENTS.md`
- `core/providers/.continue/` → `.continue/rules/`

**This repo's working state (NOT installed downstream):**

- `.hawp/work/` — backlog, active tasks, closed work, decisions, evidence

**Tooling and reference:**

- `.github/` — Copilot instructions and prompt library
- `librarian/` — distribution validation and backlog tooling (see [librarian/README.md](librarian/README.md)); `librarian/go/` holds the Go scaffold for the future small native librarian product
- `.hawp/bin/hawp` — repo-local CLI wrapper (`kit validate`, `kit normalize`, `work validate`, `work normalize`, plus `backlog` aliases)
- `benchmark/` — optional HAWP vs other methods comparison

---

## Contributing

**Distribution maintenance:**

Shared agent behaviors live in `core/providers/shared/behaviors/` and materialize into provider packs.

Install scripts are composed per provider:

- `distribution/sources/install/script-core.md` + `providers/<provider>/script-install.md` + `install/script-footer.md`
- `distribution/sources/update/script-core.md` + `providers/<provider>/script-update.md` + `update/script-footer.md`

After editing shared behaviors or distribution sources:

```bash
npm --prefix librarian run providers:sync
npm --prefix librarian run distribution:sync
```

`librarian/` is a repo-local source tree for HAWP maintenance. You do not install it separately; the commands above run against the checked-in workspace and materialize the provider files used by the generated install/update guides.

(`distribution:sync` runs `providers:sync` first — materializes generated provider overlays from shared behaviors.)

Both install and update are single copy/paste blocks per provider guide (`PROVIDER` is set in each generated bash block). Safe to re-run. They never overwrite `.hawp/work/`. Provider overlays install only the paths for that provider (see [distribution/generated/README.md](distribution/generated/README.md)).

**Validation:**

GitHub Actions runs on `main`/`dev` pushes and pull requests:

- `Validate Distribution Generated` — fails if generated guides or materialized provider overlays drift from sources.
- `Librarian Quality` — typechecks, runs the librarian unit tests, and validates `.hawp/work/` workflow state.

---

## Roadmap

Active development tracked in [`.hawp/work/BACKLOG.md`](.hawp/work/BACKLOG.md). Recent focus: provider packs and shared behavior materialization (Claude Code, Codex, GitHub, Cursor, Continue), public/private standards boundaries, and tooling hardening. UUID migration for parallel-safe work items remains on the roadmap.

## Local Setup

The repo targets Node 26 and npm 11.17.0. Use `nvm` to select the version before running librarian commands:

```bash
nvm install
nvm use
```

If you do not use `nvm`, make sure your active Node runtime satisfies `librarian/package.json` and `.nvmrc`.

## Git Publishing

For publish flows in this repository, prefer GitKraken CLI (`gk`) or the GitHub connector/tools when they are available.

- Use `gk` for terminal-based auth, provider sync, and GitKraken-managed publishing.
- Use GitHub tooling when you need repo or PR operations through the GitHub connector.
- Fall back to plain `git` only when the GitKraken or GitHub paths are unavailable.

---

## License

Apache 2.0 — see [LICENSE](./LICENSE). Installed kit ships its own `.hawp/LICENSE` copy.
