# HAWP — Human-AI Workflow Protocol

> Shape the work before execution begins. Stop re-explaining, stop drifting, stop losing context.

[![Validate Distribution Generated](https://github.com/sentzunhat/human-ai-workflow-protocol/actions/workflows/sync-distribution-generated.yml/badge.svg)](https://github.com/sentzunhat/human-ai-workflow-protocol/actions/workflows/sync-distribution-generated.yml)

---

## The problem

Every AI session starts from zero. You re-explain the goal, the constraints, what done looks like. Midway through, the agent drifts. At handoff, context evaporates and the next session re-derives everything again.

The root cause: **intent is never locked before execution begins.**

---

## What HAWP does

HAWP is a **task-shaping protocol** — five fields you fill once, before any tool call:

| Field | What goes here |
|-------|---------------|
| `input` | The request as received |
| `context` | Minimal background the agent needs |
| `mission` | The concrete objective — one sentence |
| `constraints` | Hard limits and quality bars |
| `output` | What done looks like |

That shape travels with the work. Agents read it at session start. Handoffs are a copy-paste. Re-explanation drops to near zero.

**HAWP is not** a runtime, a memory system, a framework, or an orchestrator. It is plain Markdown. No database. No API calls. No lock-in.

---

## By the numbers

From benchmark runs against the kit itself:

| Search mode | Latency | Quality |
|-------------|---------|---------|
| Lexical | <1 ms | 10 / 10 |
| Hybrid (lexical + vector) | 72 ms | 10 / 10 |
| Semantic (vector-only) | 479 ms | 9 / 10 |

- **~30% token reduction** on context packing — Jaccard dedup drops near-duplicate chunks before they reach the LLM
- **33 ms / chunk** embedding with Ollama (nomic-embed-text, warm batch, 3 100+ chunks)
- **MCP server** — any connected agent (Claude Code, Cursor, Continue) can search the kit and work index directly via `hawp_search`

Full benchmark details: [benchmark/README.md](benchmark/README.md)

---

## What you get after install

```
.hawp/
  bin/hawp          ← CLI: search, index, embed, init, update, mcp, work
  kit/              ← Protocol docs, templates, patterns, standards, examples
  work/             ← Your backlog, active plans, evidence, status reports
```

The `hawp` CLI ships with:

- **`hawp search <query>`** — hybrid lexical+vector search over kit and work docs
- **`hawp search --context`** — packs results into a single LLM-ready block with token cap and dedup
- **`hawp mcp`** — stdio MCP server; wire it into Claude Code, Cursor, or Continue in one command
- **`hawp init --provider <name>`** — provisions `~/.hawp/` and writes the MCP config for your agent
- **`hawp update`** — self-updates the binary and kit from the latest release (48h auto-update notifier built in; Windows uses manual binary replacement)
- **`hawp work new`** — scaffolds a new work item with UUID, plan file, and BACKLOG row

Maintainer and workflow-repair commands run from [`librarian/src`](./librarian/src/README.md),
including `work validate`, `work normalize`, `kit validate`, `distribution sync`,
and `providers sync`.

---

## Install — pick your agent

Each guide is a single copy-paste block. Safe to re-run. Never overwrites `.hawp/work/`.

| Provider | Install | Update |
|----------|---------|--------|
| Claude Code | [claude/install/main.md](distribution/generated/claude/install/main.md) | [claude/update/main.md](distribution/generated/claude/update/main.md) |
| GitHub Copilot | [github/install/main.md](distribution/generated/github/install/main.md) | [github/update/main.md](distribution/generated/github/update/main.md) |
| Cursor | [cursor/install/main.md](distribution/generated/cursor/install/main.md) | [cursor/update/main.md](distribution/generated/cursor/update/main.md) |
| Codex | [codex/install/main.md](distribution/generated/codex/install/main.md) | [codex/update/main.md](distribution/generated/codex/update/main.md) |
| Continue | [continue/install/main.md](distribution/generated/continue/install/main.md) | [continue/update/main.md](distribution/generated/continue/update/main.md) |

Dev channel guides: [distribution/generated/README.md](distribution/generated/README.md)

---

## First five minutes

```bash
# 1. Open your project, install for your provider (example: Claude Code)
#    → run the Install command from the table above

# 2. Index and embed the kit
hawp search index
hawp search embed --backend ollama   # or --backend onnx for offline

# 3. Wire up MCP (one-time per provider)
hawp init --provider claude   # writes .mcp.json
#    or: hawp init --provider codex    # writes .codex/config.toml

# 4. Shape your first task
open .hawp/kit/start-here.md
```

After that: use `hawp_search` from your agent instead of reading kit files directly. Ranked chunks at the token budget you set — no skimming.

If you use Codex in this repository, run `scripts/setup-codex-mcp.sh` to write a
local `.codex/config.toml` for this machine. That file stays gitignored because
Codex requires machine-local absolute paths.

---

## Why it works

**Before HAWP:** agent starts cold, re-asks for context, drifts mid-task, produces output that doesn't match the original intent.

**With HAWP:** the shape is in the repo. The agent reads it once. Constraints are visible before the first tool call. Handoffs are the shape file plus a status report — both already written.

The protocol is five fields. The kit adds templates and patterns. The CLI adds search, indexing, and MCP access. All three layers are optional: use just the shape, or use all of it.

---

## Roadmap

Active development in [`.hawp/work/BACKLOG.md`](.hawp/work/BACKLOG.md). Next focus: smart context sizing and dynamic chunk cap to push measurable token reduction beyond the current 30% dedup baseline.

---

## Contributing

Shared agent behaviors: `core/providers/shared/behaviors/` → materialize into provider packs with:

```bash
cd librarian/src && go run ./cmd/hawp distribution sync
```

See [librarian/README.md](librarian/README.md) for tooling details and validation commands.

---

## License

Apache 2.0 — see [LICENSE](./LICENSE).
