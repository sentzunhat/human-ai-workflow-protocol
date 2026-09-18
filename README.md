# HAWP — Human-AI Workflow Protocol

> The open workflow layer for humans and AI agents.
> Shape intent once. Carry it across agents, sessions, and handoffs.

**Your model can change. Your intent shouldn't.**

[![Validate Distribution Generated](https://github.com/sentzunhat/human-ai-workflow-protocol/actions/workflows/sync-distribution-generated.yml/badge.svg)](https://github.com/sentzunhat/human-ai-workflow-protocol/actions/workflows/sync-distribution-generated.yml)

---

HAWP keeps the shape of the work inside your project instead of leaving it
trapped inside an AI conversation.

Define what you are building, the context that matters, the constraints that
cannot drift, and what done looks like. Then hand that same intent to Claude
Code, Codex, Cursor, GitHub Copilot, Continue, or the next agent you decide to
use.

No hosted platform. No proprietary memory layer. No workflow lock-in. Just a
small open protocol, your repository, and optional local tooling when you want
more.

---

## Why HAWP?

Every AI session starts from zero. You re-explain the goal, the constraints, what done looks like. Midway through, the agent drifts. At handoff, context evaporates and the next session re-derives everything again.

The root cause: **intent is never locked before execution begins.**

AI agents are getting faster. The hard part is no longer getting an agent to
write code. The hard part is keeping the work aligned.

```text
human intent
     |
HAWP shape
     |
repo context + constraints
     |
Claude / Codex / Cursor / Copilot / Continue
     |
work + evidence + handoff
```

The agent can change. The shape stays.

---

## The protocol

At the center of HAWP is a deliberately small shape:

```text
input: |
  Build authentication for the API.

context: |
  Fastify service using PostgreSQL.
  Existing services use JWT access tokens.

mission: |
  Implement the authentication flow.

constraints: |
  Preserve the existing service architecture.
  Do not introduce a second authentication framework.
  Include tests.

output: |
  Working implementation with tests and a concise handoff.
```

The canonical v0.1 shape has five required fields and one optional checkpoint:

| Field | What goes here |
|-------|---------------|
| `input` | The request as received |
| `context` | Minimal background the agent needs |
| `mission` | The concrete objective — one sentence |
| `checkpoint` | Optional pause point or handoff marker |
| `constraints` | Hard limits and quality bars |
| `output` | What done looks like |

That is the protocol. Everything else in HAWP exists to make that shape easier
to create, retrieve, carry, validate, and reuse.

**HAWP is not** an agent runtime, an orchestration engine, a proprietary memory
system, a hosted AI platform, or another abstraction around your model. It
answers a smaller question: what does an AI agent need to understand before it
starts doing the work?

---

## One protocol. Three layers.

Use as much or as little as you need.

### 1. Shape

Plain text. Human-readable. Agent-readable. Portable.

Use HAWP with nothing more than Markdown when that is enough.

### 2. Kit + Work

```
.hawp/
  kit/     reusable guidance, templates, standards, and patterns
  work/    plans, backlog, evidence, decisions, and handoffs
```

The knowledge stays close to the code and evolves with the project.

### 3. Librarian

The `hawp` Go CLI turns repository knowledge into searchable agent context.

The `hawp` CLI ships with:

- **`hawp search <query>`** — hybrid lexical+vector search over kit and work docs
- **`hawp search --context`** — packs results into a single LLM-ready block with token cap and dedup
- **`hawp mcp`** — stdio MCP server; wire it into Claude Code, Cursor, or Continue in one command
- **`hawp init --provider <name>|all`** — provisions `~/.hawp/`, syncs the kit, and configures Claude, Cursor, Codex, GitHub/Copilot, plus Continue guidance
- **`hawp update`** — self-updates the binary and kit from the latest release (48h auto-update notifier built in; Windows uses manual binary replacement)
- **`hawp work new`** — scaffolds a new work item with UUID, plan file, and BACKLOG row

Maintainer and workflow-repair commands run from [`librarian/src`](./librarian/src/README.md),
including `work validate`, `work normalize`, `kit validate`, `distribution sync`,
and `providers sync`.

---

## Start in five minutes

Each guide is a single copy-paste block. Safe to re-run. Never overwrites `.hawp/work/`.

| Provider | Install | Update |
|----------|---------|--------|
| Claude Code | [claude/install/main.md](distribution/generated/claude/install/main.md) | [claude/update/main.md](distribution/generated/claude/update/main.md) |
| GitHub Copilot | [github/install/main.md](distribution/generated/github/install/main.md) | [github/update/main.md](distribution/generated/github/update/main.md) |
| Cursor | [cursor/install/main.md](distribution/generated/cursor/install/main.md) | [cursor/update/main.md](distribution/generated/cursor/update/main.md) |
| Codex | [codex/install/main.md](distribution/generated/codex/install/main.md) | [codex/update/main.md](distribution/generated/codex/update/main.md) |
| Continue | [continue/install/main.md](distribution/generated/continue/install/main.md) | [continue/update/main.md](distribution/generated/continue/update/main.md) |

Development channel guides: [distribution/generated/README.md](distribution/generated/README.md)

```bash
# Build the local search index
hawp search index

# Optional semantic search
hawp search embed --backend ollama

# Or completely local
hawp search embed --backend onnx

# Connect all supported providers through MCP
hawp init --provider all

# Or configure only one provider without downloads or kit sync
hawp mcp configure --provider continue --repo-root .

# Shape the first task
open .hawp/kit/start-here.md
```

After that: use `hawp_search` from your agent instead of reading kit files directly. Ranked chunks at the token budget you set — no skimming.

If you use Codex in this repository, run
`hawp mcp configure --provider codex --repo-root .` from the repository root.
The CLI merges supported existing settings into the local `.codex/config.toml`;
it refuses unsupported configurations instead of replacing them. That file stays
gitignored because the local server configuration contains machine-specific paths.

---

## Built for the multi-agent world

HAWP does not try to become another AI IDE. It does not require your team to
choose one model vendor. It gives different agents a common contract for the
work.

```text
                     Claude Code
                     Codex
human -> HAWP -> repo -> Cursor
                     GitHub Copilot
                     Continue
```

Today those are provider adapters. Tomorrow there will be different models,
agents, editors, and runtimes. The protocol should survive all of them.

---

## Local by design

Your workflow should not require sending your entire project history to another
platform just to remember what you are building.

HAWP keeps protocol and work artifacts in the repository. Its tooling supports
local indexing, local retrieval, SQLite-backed usage data, Ollama, and
ONNX-based models. Cloud integrations can exist; they are not the foundation.

---

## Proof, not vibes

HAWP includes benchmarks because workflow infrastructure should be measurable.

Current search benchmarks against the kit include:

| Search mode | Latency | Quality |
|-------------|---------|---------|
| Lexical | <1 ms | 10 / 10 |
| Hybrid (lexical + vector) | 72 ms | 10 / 10 |
| Semantic (vector-only) | 479 ms | 9 / 10 |

Context shaping removes near-duplicate retrieval results before they reach the
model. The current evidence records **19% total token reduction** across the
10-query benchmark run, with dense query cases saving up to **38%**.

Full benchmark details: [benchmark/README.md](benchmark/README.md). Release
evidence lives under [`.hawp/work/evidence/`](.hawp/work/evidence/).

---

## Open source vision

Builders will not work with one AI agent. Different models will plan, research,
implement, review, test, and hand work to one another.

The missing primitive is not another chatbot. It is a durable way to express
intent.

HAWP is an experiment in making that primitive open: a workflow contract that
belongs to the builder, is readable by humans and agents, stays with the
project, and remains simple enough to use without the tooling.

The goal is not to make humans adapt to agents. The goal is to give humans a
durable way to tell agents what matters.

---

## Project structure

```text
.hawp/         HAWP running against this repository
core/          canonical protocol, kit, and provider sources
librarian/     Go CLI and local intelligence tooling
distribution/  generated install and update packages
benchmark/     evaluation and benchmark tooling
```

---

## Roadmap

Active development lives in [`.hawp/work/BACKLOG.md`](.hawp/work/BACKLOG.md).
The roadmap is deliberately incremental: prove each layer before making the
protocol itself more complicated.

---

## Contributing

Shared agent behaviors: `core/providers/shared/behaviors/` → materialize into provider packs with:

```bash
cd librarian/src && go run ./cmd/hawp distribution sync
```

See [librarian/README.md](librarian/README.md) for tooling details and validation commands.

Try it. Break it. Benchmark it. Bring an agent HAWP does not support yet.
Challenge the protocol.

---

## License

Apache 2.0 — see [LICENSE](./LICENSE).
