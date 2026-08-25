---
work-item: 8672216a
type: planning
title: "Backlog audit — parked review + recently-closed compaction"
status: done
date: 2026-08-24
branch: feature/v008-backlog-audit
---

# Backlog Audit — 8672216a

## Changes Made

### Recently Closed: 15 → 10 items

Removed the 5 oldest entries (all from 2026-07-26, pointing to `closed/2026/07/26/v003-ship-audit/plan.md`):

- `j2k4l6m8` — DefaultConfig() nonexistent ONNX LLM model fix
- `a1b3c5d7` — 8 failing Ollama LLM tests repaired
- `e9f1g3h5` — Token budget 1995/0 display bug fixed
- `i7j9k1l3` — aes256 silent downgrade rejection fix
- `r5s7t9u1` — Reshape() empty-block guard fixed

These moved to a new `### Archived Recently-Closed` subsection in the Archive block of BACKLOG.md. Plan files at `closed/2026/07/26/v003-ship-audit/plan.md` remain untouched.

### Parked Item: v010-3-3a reason corrected

The backlog reason said "No working ONNX text2text models available; deferred" but the plan file (`parked/v010-3-3a/plan.md`) was written to say FLAN-T5-small is NOW FEASIBLE and the plan is ready. The reason was stale. Updated to: "Plan ready (FLAN-T5-small feasible via ONNX); deferred per incremental-only patch strategy".

## Parked Items — Decisions

| ID | Decision | Rationale |
|----|----------|-----------|
| `bee15107` | **Keep parked** | hawp mcp server (v0.0.4) covers MCP search/work access; this item is specifically about runtime CLI orchestration adapters (codex-cli, claude-cli, github-cli) for a workflow-loop runner. That scope remains out of bounds. Resume trigger: user asks for workflow-loop automation. |
| `v010-3-3a` | **Keep parked** (reason corrected) | Plan is ready and technically feasible (FLAN-T5-small via ONNX). Deferred because current strategy is incremental patches (v0.0.x); v0.1.0 gate requires proven token reduction first. Not a blocker. |
| `v010-3-2c` | **Keep parked** | Cloud API (OpenAI embeddings) track. Incremental-only strategy still in effect. |
| `v010-3-3c` | **Keep parked** | Cloud API (OpenAI LLM) track. Incremental-only strategy still in effect. |
| `v010-3-2d` | **Keep parked** | Cloud API (Anthropic embeddings stub) track. Incremental-only strategy still in effect. |
| `v010-3-3d` | **Keep parked** | Cloud API (Anthropic LLM) track. Incremental-only strategy still in effect. |
| `v010-cost` | **Keep parked** | Cross-cutting cost tracking for cloud backends. Depends on cloud tracks; same park rationale. |

## Plan File Link Audit

All plan file links verified present:
- `parked/bee15107/plan.md` — OK
- `parked/v010-3-3a/plan.md` — OK
- `parked/v010-3-2c/plan.md` — OK
- `parked/v010-3-3c/plan.md` — OK
- `parked/v010-3-2d/plan.md` — OK
- `parked/v010-3-3d/plan.md` — OK
- `parked/v010-cost/plan.md` — OK
- `closed/2026/08/24/c98518bb/plan.md` — OK
- `closed/2026/08/24/6de28bdd/plan.md` — OK
- `closed/2026/08/24/t5l7m9n1/plan.md` — OK
- `closed/2026/08/24/f2d8a5c1/plan.md` — OK
- `closed/2026/08/24/3b51837c/plan.md` — OK
- `closed/2026/08/24/j7h9e4l0/plan.md` — OK
- `closed/2026/08/22/b7e2a4f9/plan.md` — OK
- `closed/2026/08/22/b8d3e1f0/plan.md` — OK
- `closed/2026/07/26/v003-ship-audit/plan.md` — OK

No broken links found.

## Active Work — Not Changed

Active items (`4c88f451`, `1c743447`, `0ca7cf49`, `8672216a`) managed by other agents; not touched per task constraint.
