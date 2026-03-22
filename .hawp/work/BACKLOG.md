# Backlog

Active coordination index for open work. Closed history is archived under `.hawp/work/closed/`.

---

## Status Key

| Status        | Meaning                             |
| ------------- | ----------------------------------- |
| `inbox`       | Received, not yet analyzed          |
| `analyzing`   | Under investigation                 |
| `plan-ready`  | Plan written, awaiting review       |
| `approved`    | Plan approved, ready to implement   |
| `in-progress` | Being implemented                   |
| `parked`      | Deferred without closing            |
| `done`        | Implemented and verified            |
| `blocked`     | Blocked — reason noted in plan file |
| `wont-fix`    | Decided not to fix — reason noted   |

---

## Active Work

| UUID | Type | Title | Status | Owner | Plan File | Updated |
| ---- | ---- | ----- | ------ | ----- | --------- | ------- |
| `t5l7m9n1` | planning | Benchmarking infrastructure for v0.0.2+: harness scaffold + test data | plan-ready | unassigned | [plan](active/t5l7m9n1/plan.md) | 2026-07-26 |
| `j7h9e4l0` | test | Auto-update verification: test `hawp update` from 0.0.1 binary on all platforms | plan-ready | unassigned | [plan](active/j7h9e4l0/plan.md) | 2026-08-22 |

---

## Blocked / Parked

| ID  | Type | Title | Reason | Detail | Updated |
| --- | ---- | ----- | ------ | ------ | ------- |
| `bee15107` | improvement | Defer CLI participant adapters for Codex, Claude, and GitHub | Not needed yet; provider packs are enough | [plan](parked/bee15107/plan.md) | 2026-07-06 |
| `v010-3-3a` | feature | ONNX LLM Text2Text models (FLAN-T5-small) | No working ONNX text2text models available; deferred | [detail](parked/v010-3-3a/plan.md) | 2026-07-27 |
| `v010-3-2c` | feature | OpenAI Embeddings (text-embedding-3-small/large) | Cloud-API tracks parked; see incremental-only strategy | [detail](parked/v010-3-2c/plan.md) | 2026-07-27 |
| `v010-3-3c` | feature | OpenAI LLM (gpt-3.5-turbo, gpt-4-turbo) | Cloud-API tracks parked; see incremental-only strategy | [detail](parked/v010-3-3c/plan.md) | 2026-07-27 |
| `v010-3-2d` | feature | Anthropic Embeddings Stub | Cloud-API tracks parked; see incremental-only strategy | [detail](parked/v010-3-2d/plan.md) | 2026-07-27 |
| `v010-3-3d` | feature | Anthropic LLM (claude-3-sonnet, claude-3-opus) | Cloud-API tracks parked; see incremental-only strategy | [detail](parked/v010-3-3d/plan.md) | 2026-07-27 |
| `v010-cost` | feature | Cost Tracking + Rate Limiting (cross-cutting for cloud backends) | Cloud-API tracks parked; see incremental-only strategy | [detail](parked/v010-cost/plan.md) | 2026-07-27 |

---

## Recently Closed

Limited to the last 10 items.

| ID | Type | Title | Closed | Detail |
| -- | ---- | ----- | ------ | ------ |
| `b7e2a4f9` | refactor | Rename `librarian/go/` → `src/` + retire TS validators | 2026-08-22 | [plan](closed/2026/08/22/b7e2a4f9/plan.md) |
| `b8d3e1f0` | batch-close | All v0.0.1/0.0.2/0.0.3 planned work — shipped in 0.0.1 release (2026-08-21) | 2026-08-22 | [index](closed/2026/08/22/b8d3e1f0/plan.md) |
| `w4x6y8z0` | fix | `--llm-reshape` wired into the CLI | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md) |
| `n4o6p8q0` | fix | `ReshapingConfig` now honors configured Ollama URL | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md) |
| `j2k4l6m8` | fix | `DefaultConfig()` no longer defaults to nonexistent ONNX LLM model | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md) |
| `a1b3c5d7` | fix | 8 failing Ollama LLM tests repaired | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md) |
| `e9f1g3h5` | fix | Token budget `1995/0` display bug fixed | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md) |
| `i7j9k1l3` | fix | `aes256` silent downgrade now rejected at config validation | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md) |
| `r5s7t9u1` | fix | `Reshape()` empty-block guard fixed | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md) |
| `e2c4f9g5` | infrastructure | Update command v2 (platform detect, download, SHA256 verify) | 2026-07-23 | [plan](closed/2026/07/23/e2c4f9g5/plan.md) |
| `4b114c0c` | feature | Document index ingest pipeline (SQLite + FTS5 lexical search) | 2026-07-22 | [plan](closed/2026/07/22/4b114c0c/plan.md) |

---

## Archive

- Closed work: `closed/`
- Status reports: `status/`
- Evidence: `evidence/`
- Decisions: `decisions/`

## Notes

- Check this file before starting any new item.
- Active plan files go in `active/{uuid}/plan.md`. Close by moving to `closed/YYYY/MM/DD/{uuid}/plan.md`.
- Deferred items live in `parked/{id}/plan.md` without being closed.
- ADRs and decisions go in `decisions/YYYY/MM/DD/`.
- Each item gets one plan folder — no two agents on the same ID.
- Keep `Recently Closed` capped at 10; archive history lives in `closed/`.

## Future Improvements

- **Optional MCP server overlay for librarian tooling** (deferred, much later).
- **Evidence recovery for legacy unproven claims** (optional, low priority): 85 legacy claims carry explicit unproven annotations (`a7ebe68a`, 2026-07-20).
